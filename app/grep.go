package grep

import (
	"fmt"
)

// Grep -----------------------------------------------------------------------
type Grep struct {
	mustMatch bool
	hadError  bool
}

func NewGrep() *Grep {
	return &Grep{}
}

func (g *Grep) Run(line string, pattern string) (bool, error) {
	scanner := NewScanner(pattern)
	tokens := scanner.Scan(g)
	parser := NewParser(tokens)
	expr := parser.Parse()
	// printer := NewAstPrinter()
	// printer.Print(expr)
	// return false, nil
	for i := range line {
		interpreter := NewInterpreter(line[i:], g)
		ok, n := interpreter.Interpret(expr)
		if ok {
			fmt.Println("found", line[i:i+n])
			return true, nil
		}
		if g.mustMatch || g.hadError {
			break
		}
	}
	fmt.Println("not found")
	return false, nil
}

func (g *Grep) Error(err string) {
	g.hadError = true
	fmt.Println(err)
}
