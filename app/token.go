package app

// ASCII characters
const ()

// Tokens
const (
	// Single character tokens
	DOT           = iota // -> 0 .
	STAR                 // -> 1 *
	DASH                 // -> 2 -
	PLUS                 // -> 3 +
	PIPE                 // -> 4 |
	COMMA                // -> 5 ,
	CARET                // -> 6 ^
	DOLLAR               // -> 7 $
	QUESTION             // -> 8 ?
	LEFT_BRACKET         // -> 9 [
	RIGHT_BRACKET        // -> 10 ]
	LEFT_PAREN           // -> 11 (
	RIGHT_PAREN          // -> 12 )
	LEFT_BRACE           // -> 13 {
	RIGHT_BRACE          // -> 14 }
	CHAR_LITERAL         // -> 15 [a-zA-Z0-9]

	// One or two character tokens
	DIGIT     // -> 16 \d
	BACK_REF  // -> 17 \n
	ALPHA_NUM // -> 18 \w

	// Special tokens
	EOF // -> 19 end of file
)

type Token struct {
	typ     TokenType
	literal byte
}

type TokenType int

func (t Token) String() string {
	switch t.typ {
	case DIGIT, ALPHA_NUM, BACK_REF:
		if t.literal == ' ' {
			return t.typ.String() + " " + "SPACE"
		}
		return t.typ.String() + " " + string(t.literal)
	default:
		return t.typ.String()
	}
}

func (t TokenType) String() string {
	return [...]string{
		"DOT",
		"STAR",
		"DASH",
		"PLUS",
		"PIPE",
		"COMMA",
		"CARET",
		"DOLLAR",
		"QUESTION",
		"LEFT_BRACKET",
		"RIGHT_BRACKET",
		"LEFT_PAREN",
		"RIGHT_PAREN",
		"LEFT_BRACE",
		"RIGHT_BRACE",
		"CHAR_LITERAL",
		"DIGIT",
		"BACK_REF",
		"ALPHA_NUM",
		"EOF",
	}[t]
}
