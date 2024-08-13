package main

import (
	"errors"
	"slices"
)

type Scanner struct {
	tokenizer   *Tokenizer
	prevToken   *Token
	line        string
	lineCurrent int
	matchStart  int
	err         error
	backRefs    []string

	// flags
	ok           bool
	done         bool
	hasError     bool
	hasSOSAnchor bool // has start of string anchor
	mustMatch    bool
}

func NewScanner(line, regex string) *Scanner {
	return &Scanner{
		line:      line,
		tokenizer: NewTokenizer(regex),
	}
}

// Scanner --------------------------------------------------------------------
func (s *Scanner) ScanTokens() {
	if !s.tokenizerOk() {
		return
	}

	for !s.isAtEnd() {
		s.ok = false

		if !s.mustMatch {
			s.matchStart = s.lineCurrent
		}

		token := s.tokenizer.NextToken()
		s.scanToken(token)

		s.checkScannerState()

		s.prevToken = token
		s.lookForEOSAnchor()
		s.lookForOptional()
	}
	s.ok = s.ok && s.tokenizer.isAtEnd()
}

func (s *Scanner) scanToken(token *Token) {

	switch token.typ {
	case DIGIT:
		s.matchDigit()
	case ALPHA_NUM:
		s.matchAlphaNum()
	case CHAR_LITERAL:
		s.matchCharLiteral(token.literal.(byte))
	case POSITIVE_GROUP:
		s.matchPositiveGroup(token.literal.([]byte))
	case NEGATIVE_GROUP:
		s.matchNegativeGroup(token.literal.([]byte))
	case EXPR_GROUP:
		s.matchExprGroup(token.literal.([]string))
	case SOS_ANCHOR:
		s.mustBeStart()
	case EOS_ANCHOR:
		s.ok = false
		s.done = true
	case ONE_OR_MORE:
		s.matchOneOrMore()
	case OPTIONAL:
		return
	case ANY:
		s.matchAny()
	case BACK_REF:
		s.matchBackRef(token.literal.(int))
	}
}

// ----------------------------------------------------------------------------

// Asserters ------------------------------------------------------------------
func (s *Scanner) matchDigit() {
	for !s.isAtLineEnd() {
		if isDigit(s.next()) {
			s.ok = true
			break
		}
		if s.hasSOSAnchor {
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
		if isAlphaNumeric(s.next()) {
			s.ok = true
			break
		}
		if s.hasSOSAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
			break
		}
	}
}

func (s *Scanner) matchCharLiteral(expected byte) {
	for !s.isAtLineEnd() {
		if c := s.next(); c == expected {
			s.ok = true
			break
		}
		if s.hasSOSAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
			break
		}
	}
}

func (s *Scanner) matchPositiveGroup(group []byte) {
	for !s.isAtLineEnd() {
		c := s.next()
		if slices.Contains(group, c) {
			s.ok = true
			break
		}
		if s.hasSOSAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
			break
		}
	}
}

func (s *Scanner) matchNegativeGroup(group []byte) {
	for !s.isAtLineEnd() {
		if s.peek() == SPACE {
			s.ok = false
			break
		}
		c := s.next()
		if !slices.Contains(group, c) {
			s.ok = true
			break
		}
		if s.hasSOSAnchor {
			s.done = true
			break
		}
		if s.mustMatch {
			break
		}
	}
}

func (s *Scanner) matchExprGroup(exprs []string) {
	newLine := s.line[s.lineCurrent:]

	for _, expr := range exprs {
		scanner := NewScanner(newLine, expr)
		scanner.mustMatch = s.mustMatch
		scanner.ScanTokens()
		if scanner.ok {
			s.ok = true
			s.backRefs = append(s.backRefs, scanner.line[scanner.matchStart:scanner.lineCurrent])
			s.lineCurrent += scanner.lineCurrent
			return
		}
	}
	s.next()
}

func (s *Scanner) mustBeStart() {
	s.ok = s.lineCurrent == 0
	s.hasSOSAnchor = true
	if !s.ok {
		s.done = true
	}
}

func (s *Scanner) mustBeEnd() {
	s.ok = s.ok && s.isAtLineEnd()
	s.done = true
}

func (s *Scanner) matchOneOrMore() {
	oldAnchor := s.hasSOSAnchor
	s.hasSOSAnchor = false
	for !s.isAtLineEnd() {
		s.scanToken(s.prevToken)
		if !s.ok {
			s.lineCurrent--
			break
		}
		s.ok = false
	}

	s.ok = true
	s.hasSOSAnchor = oldAnchor
}

func (s *Scanner) matchOptional() {
	s.ok = false
	s.scanToken(s.tokenizer.NextToken())
	if !s.ok {
		s.lineCurrent--
	}
	s.ok = true
}

func (s *Scanner) matchAny() {
	s.next()
	s.ok = s.mustMatch
}

func (s *Scanner) matchBackRef(idx int) {
	if idx-1 >= len(s.backRefs) {
		return
	}
	expr := s.backRefs[idx-1]
	line := s.line[s.lineCurrent:]

	scanner := NewScanner(line, expr)
	scanner.hasSOSAnchor = s.mustMatch

	scanner.ScanTokens()
	s.ok = scanner.ok
	s.lineCurrent += scanner.lineCurrent
}

// ----------------------------------------------------------------------------

// Look aheads ----------------------------------------------------------------
func (s *Scanner) lookForEOSAnchor() {
	if s.tokenizer.PeekToken() == nil {
		return
	}
	next := s.tokenizer.PeekToken()
	if next.typ == EOS_ANCHOR {
		s.mustBeEnd()
		s.tokenizer.NextToken()
	}
}

func (s *Scanner) lookForOptional() {
	if s.tokenizer.PeekNextToken() == nil {
		return
	}
	next := s.tokenizer.PeekNextToken()
	if next.typ == OPTIONAL {
		if s.ok {
			s.matchOptional()
			s.tokenizer.NextToken()
		}
	}
}

// ----------------------------------------------------------------------------

// Locators -------------------------------------------------------------------
func (s *Scanner) peek() byte {
	if s.isAtLineEnd() {
		return '\000'
	}
	return s.line[s.lineCurrent]
}

func (s *Scanner) next() byte {
	s.lineCurrent++
	return s.line[s.lineCurrent-1]
}

func (s *Scanner) isAtEnd() bool {
	return s.lineCurrent >= len(s.line) ||
		s.tokenizer.isAtEnd() ||
		s.done ||
		s.hasError
}

func (s *Scanner) isAtLineEnd() bool {
	return s.lineCurrent >= len(s.line)
}

// ----------------------------------------------------------------------------

// Helpers --------------------------------------------------------------------
func (s *Scanner) resetRegex() {
	if s.hasSOSAnchor {
		return
	}
	s.matchStart = s.lineCurrent
	s.tokenizer.resetTokens()
}

func (s *Scanner) checkScannerState() {
	if s.ok {
		s.mustMatch = true
	} else if s.mustMatch {
		s.resetRegex()
	}
}

func (s *Scanner) Error(err string) {
	s.hasError = true
	s.err = errors.New(err)
}

func (s *Scanner) tokenizerOk() bool {
	if !s.tokenizer.hasError {
		return true
	}
	s.hasError = true
	s.err = s.tokenizer.err
	return false
}

// ----------------------------------------------------------------------------
