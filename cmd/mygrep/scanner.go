package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	LOWER_CASE_D  = 'd'
	LOWER_CASE_W  = 'w'
	OPEN_BRACKET  = '['
	CLOSE_BRACKET = ']'
	OPEN_PAREN    = '('
	CLOSE_PAREN   = ')'
	CARET         = '^'
	DOLLAR_SIGN   = '$'
	PLUS          = '+'
	QUESTION      = '?'
	DOT           = '.'
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
	mustMatch   bool
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
	prev := s.prevToken()
	c := s.nextToken()

	switch c {
	case BACKSLASH:
		s.isEscaping = true
	case LOWER_CASE_D:
		if s.isEscaping {
			s.isEscaping = false
			s.matchNum()
		} else {
			s.matchChar(c)
		}
	case LOWER_CASE_W:
		if s.isEscaping {
			s.isEscaping = false
			s.matchAlphaNum()
		} else {
			s.matchChar(c)
		}
	case OPEN_BRACKET:
		s.matchCharGroup()
	case CLOSE_BRACKET:
		s.nextToken()
		return
	case CARET:
		s.mustBeStart()
	case DOLLAR_SIGN:
		s.mustBeEnd()
	case PLUS:
		s.matchAtleastOne(prev)
	case QUESTION:
		s.ok = true
	case DOT:
		s.matchAny()
	case OPEN_PAREN:
		s.matchExprGroup()
	case CLOSE_PAREN:
		s.nextToken()
		return
	default:
		s.matchChar(c)
	}

	s.lookForEOSAnchor()
	s.lookForOptional()

	if s.ok {
		s.mustMatch = true
	} else if s.mustMatch {
		s.resetRegex()
	}
}

// Matchers -------------------------------------------------------------------
func (s *Scanner) matchNum() {
	for !s.isAtLineEnd() {
		if isDigit(s.nextLineChar()) {
			s.ok = true
			break
		}
		if s.startAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
			break
		}
	}
}

func (s *Scanner) matchChar(expected byte) {
	for !s.isAtLineEnd() {
		if c := s.nextLineChar(); c == expected {
			s.ok = true
			break
		}
		if s.startAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
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
		if s.startAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
			break
		}
	}
}

func (s *Scanner) matchAtleastOne(prev byte) {
	s.ok = true
	for !s.isAtLineEnd() && s.peekLine() == prev {
		s.nextLineChar()
	}
}

func (s *Scanner) matchAny() {
	if !s.isAtLineEnd() {
		s.nextLineChar()
	}
	s.ok = s.mustMatch
}

func (s *Scanner) matchCharGroup() {
	if s.isAtLineEnd() {
		s.hasError = true
		s.err = errors.New("unexpected end of regex")
	}

	c := s.peekRegex()
	group := s.seekInRegex(CLOSE_BRACKET)

	if c == CARET {
		s.negativeCharGroup(group[1:])
	} else {
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

func (s *Scanner) matchExprGroup() {
	group := s.seekInRegex(')')
	exprs := strings.Split(string(group), "|")
	newLine := s.line[s.lineCurrent:]

	for _, expr := range exprs {
		newScanner := NewScanner(newLine, expr)
		newScanner.ScanTokens()
		if newScanner.ok && newScanner.err == nil {
			s.ok = true
			s.lineCurrent += newScanner.lineCurrent
			break
		}
	}
}

func (s *Scanner) mustBeStart() {
	s.ok = s.lineCurrent == 0
	s.startAnchor = true
	if !s.ok {
		s.done = true
	}
}

func (s *Scanner) mustBeEnd() {
	s.ok = s.ok && s.isAtLineEnd()
	s.done = true
	if !s.isAtRegexEnd() {
		s.ok = false
		s.hasError = true
		s.err = errors.New("regex not terminated by $")
	}
}

// ----------------------------------------------------------------------------

// Look aheads ----------------------------------------------------------------
func (s *Scanner) lookForEOSAnchor() {
	next := s.peekRegex()
	if next == DOLLAR_SIGN {
		s.scanToken()
	}
}

func (s *Scanner) lookForOptional() {
	next := s.peekNextRegex()
	if next == QUESTION {
		if s.ok {
			if s.peekLine() == s.peekRegex() {
				s.nextLineChar()
			}
			s.nextToken()
		}
	}
}

// ----------------------------------------------------------------------------

// Match helpers --------------------------------------------------------------
func (s *Scanner) seekInRegex(expected byte) []byte {
	group := make([]byte, 0)
	for {
		if s.isAtRegexEnd() {
			s.hasError = true
			str := fmt.Sprintf("unexpected end of regex, expected %v ", expected)
			s.err = errors.New(str)
		}
		c := s.nextToken()
		if c == expected {
			break
		}
		group = append(group, c)
	}
	return group
}

func (s *Scanner) resetRegex() {
	if s.isEscaping || s.startAnchor {
		return
	}
	s.lineCurrent--
	s.regexCurrent = 0
	s.mustMatch = false
}

func (s *Scanner) peekRegex() byte {
	if s.isAtRegexEnd() {
		return '\000'
	}
	return s.regex[s.regexCurrent]
}

func (s *Scanner) peekNextRegex() byte {
	if s.regexCurrent+1 >= len(s.regex) {
		return '\000'
	}
	return s.regex[s.regexCurrent+1]
}

func (s *Scanner) peekLine() byte {
	if s.isAtLineEnd() {
		return '\000'
	}
	return s.line[s.lineCurrent]
}

func (s *Scanner) prevToken() byte {
	if s.regexCurrent == 0 {
		return '\000'
	}
	return s.regex[s.regexCurrent-1]
}

func (s *Scanner) currentToken() byte {
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
		s.lineCurrent >= len(s.line) ||
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
