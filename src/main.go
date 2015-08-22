package main

import (
	"github.com/codegangsta/cli"
	"os"
)

const (
	NAME         string = "hook-finder"
	VERSION      string = "1.0.0"
	DESCRIPTION  string = "A tool to easily locate all MyBB plugin hooks within a project."
	AUTHOR       string = "MyBB Group"
	AUTHOR_EMAIL string = "euantor@mybb.com"
)

func main() {
	app := cli.NewApp()
	app.Name = NAME
	app.Version = VERSION
	app.Usage = DESCRIPTION
	app.Authors = []cli.Author{
		cli.Author{
			Name:  AUTHOR,
			Email: AUTHOR_EMAIL,
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "input,i",
			Value: "./",
			Usage: "The root path to the project to scan for hooks.",
		},
		cli.StringFlag{
			Name:  "output,o",
			Value: "./hooks.html",
			Usage: "The path to the file to save the hook information to.",
		},
	}
	app.Action = readHooks

	app.Run(os.Args)
}
