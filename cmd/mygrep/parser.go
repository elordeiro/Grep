package main

import (
	"strconv"
)

/*
// ----------------------------------------------------------------------------
	Grammar
// ----------------------------------------------------------------------------
	regex ::= alternation
	alternation ::= concatenation ("|" concatenation)*
	concatenation ::= quantified_atom+
	quantified_atom ::= atom quantifier?
	atom ::= group | character_class | special_char | literal
	group ::= "(" regex ")"
	character_class ::= "[" char_range "]" | "[^" char_range "]"
	char_range ::= char "-" char | char char_range | char
	quantifier ::= "?" | "*" | "+" | "{" number "}" | "{" number "," number "}" | "{" number "," "}"
	special_char ::= "^" | "$" | "." | "\"
	literal ::= char | "\"
	char ::= letter | digit | escaped_char | special_char
	escaped_char ::= "\" special_char
	number ::= digit+
	backreference ::= "\" digit
	letter ::= "a" | "b" | ... | "z" | "A" | "B" | ... | "Z"
	digit ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
// ----------------------------------------------------------------------------
*/

type Parser struct {
	tokens []*Token
	pos    int
}

func NewParser(tokens []*Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() Expr {
	return p.regex()
}

func (p *Parser) regex() Expr {
	return p.alternation()
}

func (p *Parser) alternation() Expr {
	expr := p.concatenation()
	for p.match(PIPE) {
		right := p.concatenation()
		expr = &Alternation{right, expr}
	}
	return expr
}

func (p *Parser) concatenation() Expr {
	expr := p.quantifiedAtom()
	for !p.isAtEnd() && p.peek().typ != PIPE && p.peek().typ != RIGHT_PAREN {
		right := p.quantifiedAtom()
		expr = &Concatenation{expr, right}
	}
	return expr
}

func (p *Parser) quantifiedAtom() Expr {
	atom := p.atom()
	if p.match(PLUS) {
		return &Quantified{atom, 1, -1}
	}
	if p.match(QUESTION) {
		return &Quantified{atom, 0, 1}
	}
	if p.match(STAR) {
		return &Quantified{atom, 0, -1}
	}
	if p.match(LEFT_BRACE) {
		min := p.number()
		if p.match(COMMA) {
			if p.match(RIGHT_BRACE) {
				return &Quantified{atom, min, -1}
			}
			max := p.number()
			if p.match(RIGHT_BRACE) {
				return &Quantified{atom, min, max}
			}
		}
	}
	return atom
}

func (p *Parser) atom() Expr {
	if p.match(LEFT_PAREN) {
		expr := p.regex()
		p.match(RIGHT_PAREN)
		return &Group{0, expr}
	}
	if p.match(LEFT_BRACKET) {
		inverted := false
		if p.match(CARET) {
			inverted = true
		}
		ranges := p.charRange()
		p.match(RIGHT_BRACKET)
		return &CharRange{inverted, ranges}
	}
	if p.match(CHAR_LITERAL) {
		return &CharLiteral{p.prev().literal}
	}
	if p.match(DOT) {
		return &Any{}
	}
	if p.match(BACK_REF) {
		return p.backreference()
	}
	if p.match(DIGIT) {
		return &CharRange{false, []CharRangePair{{'0', '9'}}}
	}
	if p.match(ALPHA_NUM) {
		return &CharRange{false, []CharRangePair{{'a', 'z'}, {'A', 'Z'}, {'0', '9'}}}
	}
	if p.match(CARET) {
		return &StartOfLine{}
	}
	if p.match(DOLLAR) {
		return &EndOfLine{}
	}
	if p.match(COMMA) {
		return &CharLiteral{','}
	}
	if p.match(DASH) {
		return &CharLiteral{'-'}
	}
	return nil
}

func (p *Parser) charRange() []CharRangePair {
	var ranges []CharRangePair
	for !p.isAtEnd() && p.peek().typ != RIGHT_BRACKET {
		if p.match(CHAR_LITERAL) {
			start := p.prev().literal
			if p.match(DASH) {
				if p.match(CHAR_LITERAL) {
					end := p.prev().literal
					ranges = append(ranges, CharRangePair{start, end})
				}
			} else {
				ranges = append(ranges, CharRangePair{start, start})
			}
		}
	}
	return ranges
}

func (p *Parser) backreference() Expr {
	idx, _ := strconv.Atoi(string(p.prev().literal))
	return &BackRef{idx}
}

func (p *Parser) number() int {
	if p.match(CHAR_LITERAL) {
		return int(p.prev().literal - '0')
	}
	return 0
}

func (p *Parser) match(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	if p.tokens[p.pos].typ == t {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) prev() *Token {
	return p.tokens[p.pos-1]
}

func (p *Parser) isAtEnd() bool {
	return p.tokens[p.pos].typ == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.pos]
}
