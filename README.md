# Hook Finder

A tool to easily locate all hooks within a MyBB package - either the core of MyBB or a custom plugin.

This tool will traverse all PHP files within a given directory, looking for plugin hooks. All plugin hooks will then be written to a HTML file of your choosing.

## Installation

You can download pre-built binaries for Windows, OS X and Linux from the Downloads page, or install the project using [go](http://golang.org):

`go get github.com/mybb/hook-finder`

## Usage

The hook finder requires two parameters:

- The path to the project to scan for hooks
- The output file name

You run the tool as follows:

```
hook-finder -i /pat/to/mybb/root -o /path/to/hooks.html
```

The tool will then scan all of the PHP files in the specified directory and locate all available hooks which will then be written to the specified file.
