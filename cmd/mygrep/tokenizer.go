package main

import (
	"errors"
	"fmt"
	"strconv"
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
	SPACE         = ' '
	BACKSLASH     = '\\'
)

const (
	DIGIT          = iota // -> 0
	ALPHA_NUM             // -> 1
	CHAR_LITERAL          // -> 2
	POSITIVE_GROUP        // -> 3
	NEGATIVE_GROUP        // -> 4
	EXPR_GROUP            // -> 5
	SOS_ANCHOR            // -> 6 start of string anchor
	EOS_ANCHOR            // -> 7 end of string anchor
	ONE_OR_MORE           // -> 8
	OPTIONAL              // -> 9
	ANY                   // -> 10
	BACK_REF              // -> 11
)

type Tokenizer struct {
	regex        string
	regexCurrent int
	err          error
	hasError     bool
	tokens       []*Token
	tokensBackUp []*Token
	tokenCurrent int
}

type Token struct {
	typ     int
	literal any
}

func NewTokenizer(regex string) *Tokenizer {
	t := &Tokenizer{
		tokens: make([]*Token, 0),
		regex:  regex,
	}

	for !t.isRegexAtEnd() {
		t.tokenize()
	}
	return t
}

// Tokenizer ------------------------------------------------------------------
func (t *Tokenizer) tokenize() {
	c := t.next()

	switch c {
	case BACKSLASH:
		// Escaped chars
		switch t.peek() {
		case LOWER_CASE_D:
			t.addToken(DIGIT, c)
			t.next()
		case LOWER_CASE_W:
			t.addToken(ALPHA_NUM, c)
			t.next()
		case BACKSLASH:
		default:
			t.mustBeDigit(BACK_REF)
		}
	case OPEN_BRACKET:
		t.charGroup()
	case CLOSE_BRACKET:
		t.Error(fmt.Sprintf("unexpected char '%v'", string(c)))
	case CARET:
		t.addToken(SOS_ANCHOR, c)
	case DOLLAR_SIGN:
		t.addToken(EOS_ANCHOR, c)
	case PLUS:
		t.addToken(ONE_OR_MORE, c)
	case QUESTION:
		t.addToken(OPTIONAL, c)
	case DOT:
		t.addToken(ANY, c)
	case OPEN_PAREN:
		t.exprGroup()
	case CLOSE_PAREN:
		t.Error(fmt.Sprintf("unexpected char '%v'", string(c)))
	default:
		t.addToken(CHAR_LITERAL, c)
	}

}

// Accessors ------------------------------------------------------------------
func (t *Tokenizer) NextToken() *Token {
	if t.hasError || t.tokenCurrent >= len(t.tokens) {
		return nil
	}
	t.tokenCurrent++
	return t.tokens[t.tokenCurrent-1]
}

func (t *Tokenizer) PeekToken() *Token {
	if t.hasError || t.tokenCurrent >= len(t.tokens) {
		return nil
	}
	return t.tokens[t.tokenCurrent]
}

func (t *Tokenizer) PeekNextToken() *Token {
	if t.hasError || t.tokenCurrent+1 >= len(t.tokens) {
		return nil
	}
	return t.tokens[t.tokenCurrent+1]
}

// ----------------------------------------------------------------------------

// Asserters ------------------------------------------------------------------
func (t *Tokenizer) mustMatch(c byte) {
	if n := t.next(); n != c {
		t.Error(fmt.Sprintf("wanted %v, got %v", string(c), n))
	}
}

func (t *Tokenizer) mustBeDigit(typ int) {
	c := t.next()
	if isDigit(c) {
		i, _ := strconv.Atoi(string(c))
		t.addToken(typ, i)
		return
	}
	t.Error(fmt.Sprintf("wanted digit, got %v", string(c)))
}

func (t *Tokenizer) charGroup() {
	var typ int
	if t.peek() == CARET {
		typ = NEGATIVE_GROUP
		t.next()
	} else {
		typ = POSITIVE_GROUP
	}

	group := make([]byte, 0)
	for t.peek() != CLOSE_BRACKET {
		group = append(group, t.next())
	}
	t.addToken(typ, group)
	t.mustMatch(']')
}

func (t *Tokenizer) exprGroup() {
	group := make([]byte, 0)
	for t.peek() != CLOSE_PAREN {
		group = append(group, t.next())
	}
	exprs := strings.Split(string(group), "|")
	t.addToken(EXPR_GROUP, exprs)
	t.mustMatch(')')
}

// ----------------------------------------------------------------------------

// Locators -------------------------------------------------------------------
func (t *Tokenizer) next() byte {
	t.regexCurrent++
	return t.regex[t.regexCurrent-1]
}

func (t *Tokenizer) peek() byte {
	if t.isRegexAtEnd() {
		return '\000'
	}
	return t.regex[t.regexCurrent]
}

func (t *Tokenizer) peekNext() byte {
	if t.regexCurrent+1 >= len(t.regex) {
		return '\000'
	}
	return t.regex[t.regexCurrent+1]
}

func (t *Tokenizer) isRegexAtEnd() bool {
	return t.regexCurrent >= len(t.regex)
}

func (t *Tokenizer) isAtEnd() bool {
	return t.tokenCurrent >= len(t.tokens)
}

// ----------------------------------------------------------------------------

// Helpers --------------------------------------------------------------------
func (t *Tokenizer) addToken(typ int, literal any) {
	t.tokens = append(t.tokens, &Token{typ, literal})
	t.tokensBackUp = append(t.tokensBackUp, &Token{typ, literal})
}

func (t *Tokenizer) Error(err string) {
	t.hasError = true
	t.err = errors.New(err)
}

func (t *Tokenizer) resetTokens() {
	copy(t.tokens, t.tokensBackUp)
	t.tokenCurrent = 0
}

// ----------------------------------------------------------------------------
