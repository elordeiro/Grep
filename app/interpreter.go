package app

type Interpreter struct {
	line     string
	pos      int
	grep     *Grep
	backRefs []string
}

func NewInterpreter(line string, grep *Grep) *Interpreter {
	return &Interpreter{line, 0, grep, make([]string, 0)}
}

func (i *Interpreter) Interpret(expr Expr) (bool, int) {
	return i.evaluate(expr), i.pos
}

func (i *Interpreter) evaluate(expr Expr) bool {
	return expr.accept(i)
}

func (i *Interpreter) visitAlternation(e Alternation) bool {
	return i.evaluate(e.left) || i.evaluate(e.right)
}

func (i *Interpreter) visitQuantified(e Quantified) bool {
	for x := 0; x < e.min; x++ {
		if !i.evaluate(e.left) {
			return false
		}
	}
	if e.max == -1 {
		for {
			if !i.evaluate(e.left) {
				break
			}
		}
	}
	for x := e.min; x < e.max; x++ {
		if !i.evaluate(e.left) {
			break
		}
	}
	return true
}

func (i *Interpreter) visitConcatenation(e Concatenation) bool {
	return i.evaluate(e.left) && i.evaluate(e.right)
}

func (i *Interpreter) visitGroup(e Group) bool {
	strIdx := i.pos
	e.refIdx = len(i.backRefs)
	i.backRefs = append(i.backRefs, "")
	ok := i.evaluate(e.expr)
	if ok {
		i.backRefs[e.refIdx] = i.line[strIdx:i.pos]
	}
	return ok
}

func (i *Interpreter) visitCharRange(e CharRange) bool {
	if i.pos >= len(i.line) || !isAlphaNumeric(i.line[i.pos]) {
		return false
	}
	for _, pair := range e.ranges {
		if !e.inverted {
			if i.visitCharRangePair(pair) {
				return true
			}
		} else {
			if !i.visitCharRangePair(pair) {
				continue
			} else {
				return false
			}
		}
	}
	if e.inverted {
		i.pos++
	}
	return e.inverted
}

func (i *Interpreter) visitCharRangePair(e CharRangePair) bool {
	if i.pos >= len(i.line) {
		return false
	}
	if i.line[i.pos] >= e.start && i.line[i.pos] <= e.end {
		i.pos++
		return true
	}
	return false
}

func (i *Interpreter) visitCharLiteral(e CharLiteral) bool {
	if i.pos >= len(i.line) {
		return false
	}
	if i.line[i.pos] == e.value {
		i.pos++
		return true
	}
	return false
}

func (i *Interpreter) visitBackRef(e BackRef) bool {
	if e.backref > len(i.backRefs) {
		i.Error("Invalid backreference")
		return false
	}
	if i.backRefs[e.backref-1] == i.line[i.pos:i.pos+len(i.backRefs[e.backref-1])] {
		i.pos += len(i.backRefs[e.backref-1])
		return true
	}
	return false
}

func (i *Interpreter) visitAny(e Any) bool {
	i.pos++
	return true
}

func (i *Interpreter) visitStartOfLine(e StartOfLine) bool {
	i.grep.mustMatch = true
	return i.pos == 0
}

func (i *Interpreter) visitEndOfLine(e EndOfLine) bool {
	return i.pos == len(i.line)
}

func (i *Interpreter) Error(err string) {
	i.grep.hadError = true
	i.grep.Error(err)
}
