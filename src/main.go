package main

import (
	"bufio"
	"flag"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	// Assume that hook names contain only alphanumeric characters, _, - and . - all other hook names will not be matched
	PLUGIN_HOOK_REGEX  string = `\$plugins->run_hooks\(['|"](?P<hook_name>[\w .-]+)['|"](, ?(?P<arg>\$[\w]+))*\);`
	PHP_FILE_REGEX     string = `.*\.php$`
	TEMPLATE_FILE_NAME string = "templates/main.html"
)

type HookInfo struct {
	File string
	Line int
	Name string
	Args []string
}

func main() {
	pluginHookRegex, err := regexp.Compile(PLUGIN_HOOK_REGEX)
	phpFileRegex, err := regexp.Compile(PHP_FILE_REGEX)

	if err != nil {
		panic(err)
	}

	mybbRootPtr := flag.String("path", "./", "Specify a path to the MyBB root. Defaults to current directory.")
	outputFilePtr := flag.String("output", "./output.html", "Specify an output file. Defaults to output.html.")

	flag.Parse()

	var readFilesWaitGroup sync.WaitGroup
	hooks := make(map[string][]*HookInfo)

	err = filepath.Walk(*mybbRootPtr, func(path string, fileInfo os.FileInfo, err error) error {
		isPhpFile := phpFileRegex.MatchString(path)

		if isPhpFile {
			readFilesWaitGroup.Add(1)
			go readFile(path, *mybbRootPtr, pluginHookRegex, &readFilesWaitGroup, hooks)
		}

		return nil
	})

	readFilesWaitGroup.Wait()

	writeOutputFile(hooks, outputFilePtr)
}

func readFile(path, mybbRootPath string, pluginHookRegex *regexp.Regexp, readFilesWaitGroup *sync.WaitGroup, hooks map[string][]*HookInfo) {
	defer readFilesWaitGroup.Done()

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	lineNo := 1
	for scanner.Scan() {
		line = scanner.Text()
		if pluginHookRegex.MatchString(line) {
			subMatches := pluginHookRegex.FindStringSubmatch(line)

			var hookName string
			var args []string

			for i, name := range pluginHookRegex.SubexpNames() {
				if i == 0 || name == "" {
					continue
				}

				if name == "hook_name" {
					hookName = subMatches[i]
				} else if name == "arg" {
					args = append(args, subMatches[i])
				}
			}

			hookInfo := &HookInfo{
				File: path,
				Line: lineNo,
				Name: hookName,
				Args: args,
			}

			shortPath := strings.Replace(path, mybbRootPath+"/", "", 1)

			hooks[shortPath] = append(hooks[shortPath], hookInfo)
		}

		lineNo++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func writeOutputFile(hooks map[string][]*HookInfo, outputFileNamePtr *string) {
	outputFile, err := os.OpenFile(*outputFileNamePtr, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	var funcMap = template.FuncMap{
		"join": strings.Join,
	}

	tmp := template.Must(template.New(TEMPLATE_FILE_NAME).Funcs(funcMap).ParseFiles(TEMPLATE_FILE_NAME))
	tmp.Execute(outputFile, hooks)

	outputFile.Sync()
}
