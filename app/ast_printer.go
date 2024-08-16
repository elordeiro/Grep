package grep

import (
	"fmt"
)

// Ast printer ----------------------------------------------------------------
type AstPrinter struct {
	tab int
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (a *AstPrinter) Print(expr Expr) {
	expr.accept(a)
}

func (a *AstPrinter) visitAlternation(e Alternation) bool {
	a.Print(e.left)
	fmt.Printf("%*sAlternation\n", a.tab, " ")
	a.Print(e.right)
	return true
}

func (a *AstPrinter) visitConcatenation(e Concatenation) bool {
	a.Print(e.left)
	fmt.Printf("%*sConcatenation\n", a.tab, " ")
	a.Print(e.right)
	return true
}

func (a *AstPrinter) visitQuantified(e Quantified) bool {
	a.Print(e.left)
	fmt.Printf("%*sQuantified\n min: %v max: %v", a.tab, " ", e.min, e.max)
	return true
}

func (a *AstPrinter) visitGroup(e Group) bool {
	fmt.Printf("%-*sGroup\n", a.tab, " ")
	a.tab += 5
	a.Print(e.expr)
	a.tab -= 5
	return true
}

func (a *AstPrinter) visitBackRef(e BackRef) bool {
	fmt.Printf("%*sBackRef ref: %v\n", a.tab, " ", e.backref)
	return true
}

func (a *AstPrinter) visitCharRange(e CharRange) bool {
	fmt.Printf("%*sCharRange\n", a.tab, " ")
	a.tab += 5
	for _, pair := range e.ranges {
		a.Print(pair)
	}
	a.tab -= 5
	return true
}

func (a *AstPrinter) visitCharLiteral(e CharLiteral) bool {
	fmt.Printf("%*sCharLiteral value: %s\n", a.tab, " ", string(e.value))
	return true
}

func (a *AstPrinter) visitAny(e Any) bool {
	fmt.Printf("%*sAny\n", a.tab, " ")
	return true
}

func (a *AstPrinter) visitStartOfLine(e StartOfLine) bool {
	fmt.Printf("%*sStartOfLine\n", a.tab, " ")
	return true
}

func (a *AstPrinter) visitEndOfLine(e EndOfLine) bool {
	fmt.Printf("%*sEndOfLine\n", a.tab, " ")
	return true
}

func (a *AstPrinter) visitCharRangePair(e CharRangePair) bool {
	fmt.Printf("%*sCharRangePair start: %c end: %c\n", a.tab, " ", e.start, e.end)
	return true
}
