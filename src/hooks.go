package main

import (
	"bufio"
	"fmt"
	"github.com/codegangsta/cli"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

const (
	// Assume that hook names contain only alphanumeric characters, _, - and . - all other hook names will not be matched
	PLUGIN_HOOK_REGEX  string = `\$plugins->run_hooks\(['|"](?P<hook_name>[\w .-]+)['|"](, ?(?P<arg>\$[\w]+))*\);`
	PHP_FILE_REGEX     string = `.*\.php$`
	TEMPLATE_FILE_NAME string = "./templates/main.html"
)

type HookInfo struct {
	File string
	Line int
	Name string
	Args []string
}

type HookList struct {
	sync.RWMutex
	Hooks map[string][]*HookInfo
}

func NewHookList() *HookList {
	list := new(HookList)
	list.Hooks = make(map[string][]*HookInfo)
	return list
}

func readHooks(c *cli.Context) {
	pluginHookRegex := regexp.MustCompile(PLUGIN_HOOK_REGEX)
	phpFileRegex := regexp.MustCompile(PHP_FILE_REGEX)

	var readFilesWaitGroup sync.WaitGroup
	hooks := NewHookList()

	fmt.Printf("Scanning directory %s for hooks.\n", c.String("input"))

	err := filepath.Walk(c.String("input"), func(path string, fileInfo os.FileInfo, err error) error {
		if phpFileRegex.MatchString(path) {
			readFilesWaitGroup.Add(1)
			go readFile(path, c.String("input"), pluginHookRegex, &readFilesWaitGroup, hooks)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	readFilesWaitGroup.Wait()

	fmt.Printf("Writing hooks from %d files to output file %s\n", len(hooks.Hooks), c.String("output"))

	writeOutputFile(hooks, c.String("output"))
}

func readFile(path, mybbRootPath string, pluginHookRegex *regexp.Regexp, readFilesWaitGroup *sync.WaitGroup, hooks *HookList) {
	defer readFilesWaitGroup.Done()

	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", path, err)
		return
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

			hooks.Add(shortPath, hookInfo)
		}

		lineNo++
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading files: %s\n", err)
		os.Exit(1)
	}
}

func (l *HookList) Add(shortPath string, hook *HookInfo) {
	l.Lock()
	defer l.Unlock()

	l.Hooks[shortPath] = append(l.Hooks[shortPath], hook)
}

func writeOutputFile(hooks *HookList, outputFileName string) {
	outputFile, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening output file: %s\n", err)
		os.Exit(1)
	}

	defer outputFile.Close()

	var funcMap = template.FuncMap{
		"join": strings.Join,
	}

	hooks.RLock()
	defer hooks.RUnlock()
	template, err := template.New(TEMPLATE_FILE_NAME).Funcs(funcMap).ParseFiles(TEMPLATE_FILE_NAME)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initialising output to file: %s\n", err)
		os.Exit(1)
	}

	_, templateName := path.Split(TEMPLATE_FILE_NAME)

	err = template.ExecuteTemplate(outputFile, templateName, hooks)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output to file: %s\n", err)
		os.Exit(1)
	}
}
