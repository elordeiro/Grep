package main

type Expr interface {
	accept(v Visitor) bool
}

type Visitor interface {
	visitAlternation(e Alternation) bool
	visitQuantified(e Quantified) bool
	visitConcatenation(e Concatenation) bool
	visitGroup(e Group) bool
	visitCharRange(e CharRange) bool
	visitCharRangePair(e CharRangePair) bool
	visitCharLiteral(e CharLiteral) bool
	visitBackRef(e BackRef) bool
	visitAny(e Any) bool
	visitStartOfLine(e StartOfLine) bool
	visitEndOfLine(e EndOfLine) bool
}

type Alternation struct {
	left  Expr
	right Expr
}

type Concatenation struct {
	left  Expr
	right Expr
}

type Quantified struct {
	left     Expr
	min, max int
}

type Group struct {
	refIdx int
	expr   Expr
}

type BackRef struct {
	backref int
}

type CharRange struct {
	inverted bool
	ranges   []CharRangePair
}

type CharRangePair struct {
	start byte
	end   byte
}

type CharLiteral struct {
	value byte
}

type Any struct{}

type StartOfLine struct{}

type EndOfLine struct{}

func (e Alternation) accept(v Visitor) bool {
	return v.visitAlternation(e)
}

func (e Concatenation) accept(v Visitor) bool {
	return v.visitConcatenation(e)
}

func (e Quantified) accept(v Visitor) bool {
	return v.visitQuantified(e)
}

func (e Group) accept(v Visitor) bool {
	return v.visitGroup(e)
}

func (e CharRange) accept(v Visitor) bool {
	return v.visitCharRange(e)
}

func (e CharRangePair) accept(v Visitor) bool {
	return v.visitCharRangePair(e)
}

func (e Any) accept(v Visitor) bool {
	return v.visitAny(e)
}

func (e BackRef) accept(v Visitor) bool {
	return v.visitBackRef(e)
}

func (e CharLiteral) accept(v Visitor) bool {
	return v.visitCharLiteral(e)
}

func (e StartOfLine) accept(v Visitor) bool {
	return v.visitStartOfLine(e)
}

func (e EndOfLine) accept(v Visitor) bool {
	return v.visitEndOfLine(e)
}
