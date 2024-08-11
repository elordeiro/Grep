package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
)

const (
	LOWER_CASE_D  = 'd'
	LOWER_CASE_W  = 'w'
	OPEN_BRACKET  = '['
	CLOSE_BRACKET = ']'
	CARET         = '^'
	DOLLAR_SIGN   = '$'
	BACKSLASH     = '\\'
)

type Scanner struct {
	line         string
	regex        string
	lineCurrent  int
	regexCurrent int
	err          error

	// flags
	ok          bool
	done        bool
	hasError    bool
	isEscaping  bool
	startAnchor bool
}

func NewScanner(line, regex string) *Scanner {
	return &Scanner{
		line:  line,
		regex: regex,
	}
}

func (s *Scanner) ScanTokens() {
	for !s.isAtEnd() {
		s.ok = false
		s.scanToken()
	}
	s.ok = s.ok && s.isAtRegexEnd()
}

func (s *Scanner) scanToken() {
	c := s.nextToken()

	switch c {
	case BACKSLASH:
		s.isEscaping = true
	case LOWER_CASE_D:
		if s.isEscaping {
			s.isEscaping = false
			s.matchNum()
		} else {
			s.matchAlpha(c)
		}
	case LOWER_CASE_W:
		if s.isEscaping {
			s.isEscaping = false
			s.matchAlphaNum()
		} else {
			s.matchAlpha(c)
		}
	case OPEN_BRACKET:
		s.matchCharGroup()
	case CLOSE_BRACKET:
		s.nextToken()
	case CARET:
		s.mustBeStart()
	case DOLLAR_SIGN:
		s.mustBeEnd()
	default:
		s.matchAlpha(c)
	}

	if next := s.peek(); next == DOLLAR_SIGN {
		s.scanToken()
	}
}

// Matchers -------------------------------------------------------------------
func (s *Scanner) matchNum() {
	for !s.isAtLineEnd() {
		if isDigit(s.nextLineChar()) {
			s.ok = true
			break
		}
	}
}

func (s *Scanner) matchAlpha(expected byte) {
	for !s.isAtLineEnd() {
		if c := s.nextLineChar(); c == expected {
			s.ok = true
			break
		}
		if s.startAnchor {
			s.done = true
			break
		}
	}
}

func (s *Scanner) matchAlphaNum() {
	for !s.isAtLineEnd() {
		if isAlphaNumeric(s.nextLineChar()) {
			s.ok = true
			break
		}
	}
}

func (s *Scanner) matchCharGroup() {
	if s.isAtLineEnd() {
		s.hasError = true
		s.err = errors.New("unexpected end of regex")
	}

	c := s.peek()
	if c == CARET {
		s.nextToken()
		group := s.seek(CLOSE_BRACKET)
		s.negativeCharGroup(group)
	} else {
		group := s.seek(CLOSE_BRACKET)
		s.positiveCharGroup(group)
	}
}

func (s *Scanner) negativeCharGroup(group []byte) {
	s.ok = true
	for !s.isAtLineEnd() {
		c := s.nextLineChar()
		if slices.Contains(group, c) {
			s.ok = false
			s.done = true
			break
		}
	}
}

func (s *Scanner) positiveCharGroup(group []byte) {
	for !s.isAtLineEnd() {
		c := s.nextLineChar()
		if slices.Contains(group, c) {
			s.ok = true
			break
		}
	}
}

func (s *Scanner) mustBeStart() {
	s.ok = s.lineCurrent == 0
	if !s.ok {
		s.done = true
	}
	s.startAnchor = true
}

func (s *Scanner) mustBeEnd() {
	s.ok = s.ok && s.isAtLineEnd()
	s.done = true
}

// ----------------------------------------------------------------------------

// Match helpers --------------------------------------------------------------
func (s *Scanner) seek(expected byte) []byte {
	group := make([]byte, 0)
	for {
		if s.isAtLineEnd() {
			s.hasError = true
			fmt.Fprintf(os.Stderr, "unexpected line end, expected %v ", expected)
		}
		c := s.nextToken()
		if c == expected {
			break
		}
		group = append(group, c)
	}
	return group
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.regex[s.regexCurrent]
}

func (s *Scanner) nextToken() byte {
	s.regexCurrent++
	return s.regex[s.regexCurrent-1]
}

func (s *Scanner) nextLineChar() byte {
	s.lineCurrent++
	return s.line[s.lineCurrent-1]
}

func (s *Scanner) isAtEnd() bool {
	return s.regexCurrent >= len(s.regex) ||
		s.done ||
		s.hasError
}

func (s *Scanner) isAtLineEnd() bool {
	return s.lineCurrent >= len(s.line)
}

func (s *Scanner) isAtRegexEnd() bool {
	return s.regexCurrent >= len(s.regex)
}

// ----------------------------------------------------------------------------

// Helpers --------------------------------------------------------------------
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

// ----------------------------------------------------------------------------
