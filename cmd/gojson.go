package main

import (
	"fmt"
	"io"
	"os"

	"github.com/oabrivard/gojson/linter"
)

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}

func main() {
	var f *os.File

	if isInputFromPipe() {
		f = os.Stdin

	} else {
		fileName := ""

		if len(os.Args) != 2 {
			fmt.Fprintf(os.Stderr, "gojson filename\n")
			os.Exit(1)
		} else {
			fileName = os.Args[1]
		}

		var err error
		f, err = os.Open(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	jl := linter.NewJsonLinter(string(bytes))
	result, err := jl.Lint()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(result)
}
