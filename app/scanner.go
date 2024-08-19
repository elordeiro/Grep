package app

type Scanner struct {
	regex   string
	current int
	tokens  []*Token
}

func NewScanner(regex string) *Scanner {
	return &Scanner{
		tokens: make([]*Token, 0),
		regex:  regex,
	}
}

// Scanner --------------------------------------------------------------------
func (s *Scanner) Scan(grep *Grep) []*Token {
	for !s.isAtEnd() {
		s.scanToken(grep)
	}
	s.tokens = append(s.tokens, &Token{EOF, '\000'})
	return s.tokens
}

func (s *Scanner) scanToken(grep *Grep) {
	c := s.next()

	switch c {
	case '\\':
		s.escape(grep)
	case '.':
		s.addToken(DOT, c)
	case '*':
		s.addToken(STAR, c)
	case '-':
		s.addToken(DASH, c)
	case '+':
		s.addToken(PLUS, c)
	case '|':
		s.addToken(PIPE, c)
	case ',':
		s.addToken(COMMA, c)
	case '^':
		s.addToken(CARET, c)
	case '$':
		s.addToken(DOLLAR, c)
	case '?':
		s.addToken(QUESTION, c)
	case '[':
		s.addToken(LEFT_BRACKET, c)
	case ']':
		s.addToken(RIGHT_BRACKET, c)
	case '(':
		s.addToken(LEFT_PAREN, c)
	case ')':
		s.addToken(RIGHT_PAREN, c)
	case '{':
		s.addToken(LEFT_BRACE, c)
	case '}':
		s.addToken(RIGHT_BRACE, c)
	case ' ':
		s.addToken(CHAR_LITERAL, c)
	default:
		s.addToken(CHAR_LITERAL, c)
	}
}

// Asserters ------------------------------------------------------------------
func (s *Scanner) escape(grep *Grep) {
	c := s.next()
	switch c {
	case 'd':
		s.addToken(DIGIT, c)
	case 'w':
		s.addToken(ALPHA_NUM, c)
	case '\\':
		s.addToken(CHAR_LITERAL, c)
	default:
		if isDigit(c) {
			s.addToken(BACK_REF, c)
		} else {
			grep.Error("unexpected character: " + string(c))
		}
	}
}

// ----------------------------------------------------------------------------

// Locators -------------------------------------------------------------------
func (s *Scanner) next() byte {
	s.current++
	return s.regex[s.current-1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.regex)
}

// ----------------------------------------------------------------------------

// Helpers --------------------------------------------------------------------
func (s *Scanner) addToken(typ TokenType, literal byte) {
	s.tokens = append(s.tokens, &Token{typ, literal})
}

// ----------------------------------------------------------------------------
