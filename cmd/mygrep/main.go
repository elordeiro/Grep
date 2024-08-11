package main

import (
	"fmt"
	"io"
	"os"

	"unicode/utf8"
)

/*
Usage: echo <input_text> | your_program.sh -E <pattern>
Exit codes:

	0 means success
	1 means no lines were selected, >1 means error
*/
func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: mygrep -E <pattern>\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin) // assume we're only dealing with a single line
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(line, pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		fmt.Println("not found")
		os.Exit(1)
	}

	fmt.Println("found")
}

func matchLine(line []byte, pattern string) (bool, error) {
	if utf8.RuneCountInString(pattern) < 1 {
		return false, fmt.Errorf("unsupported pattern: %q", pattern)
	}

	scanner := NewScanner(string(line), pattern)

	scanner.ScanTokens()

	return scanner.ok, scanner.err
}
