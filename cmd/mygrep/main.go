package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
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

	for _, c := range line {
		ok, restOfPattern, err := matchRune(rune(c), pattern)
		pattern = restOfPattern
		if err != nil {
			return false, err
		}

		if ok && pattern == "" {
			return true, nil
		}
	}
	return false, nil
}

func matchRune(r rune, pattern string) (bool, string, error) {
	var ok bool
	var patIdx int

	switch pattern[patIdx] {

	// Match single digit or letter
	case '\\':
		patIdx++
		switch pattern[patIdx] {
		// Match single digit
		case 'd':
			patIdx++
			ok = unicode.IsDigit(r)
			// Match single letter
		case 'w':
			patIdx++
			ok = unicode.IsDigit(r) || unicode.IsLetter(r) || r == rune('_')
		}

	case '[':
		patIdx++
		endIdx := strings.Index(pattern, "]")
		if endIdx == -1 {
			return false, "", errors.New("missing closing bracket")
		}

		include := true
		insideExpr := pattern[patIdx:endIdx]
		patIdx = endIdx + 1

		if insideExpr[0] == '^' {
			include = false
			insideExpr = insideExpr[1:]
		}

		for _, w := range insideExpr {
			if include {
				ok = r == w
				if ok {
					break
				}
			} else {
				ok = r != w
				if !ok {
					break
				}
			}
		}

	default:
		ok = r == rune(pattern[patIdx])
		patIdx++
	}

	if ok {
		pattern = pattern[patIdx:]
	}

	return ok, pattern, nil
}
