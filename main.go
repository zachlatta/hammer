package main

import (
	"flag"
	"io"
	"os"
	"pkg/text/template"
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}
}

const usageTemplate = `Hammer is a utility for finding short usernames.

Usage:

  hammer [network] [flags]

`

func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	tmpl(w, usageTemplate, nil)
}

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}
