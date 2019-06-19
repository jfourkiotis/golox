package scanner

import (
	"golox/parseerror"
	"golox/token"
)

// Scanner transforms the source into tokens
type Scanner struct {
	source  string
	start   int
	current int
	line    int
}

// New creates a new scanner
func New(source string) Scanner {
	scanner := Scanner{source: source, line: 1}
	return scanner
}

// ScanTokens transforms the source into an array of tokens. The last token
// is always an token.EOF
func (sc *Scanner) ScanTokens() []token.Token {
	tokens := make([]token.Token, 0)

	for !sc.isAtEnd() {
		// we're at the beginning of the next lexeme
		sc.start = sc.current
		tok := sc.scanToken()
		if tok.Type != token.INVALID {
			tokens = append(tokens, tok)
		}
	}

	tokens = append(tokens, token.Token{Type: token.EOF})
	return tokens
}

func (sc *Scanner) scanToken() token.Token {
	c := sc.advance()
	switch c {
	case '(':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.LEFTPAREN, Literal: literal, Line: sc.line}
	case ')':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.RIGHTPAREN, Literal: literal, Line: sc.line}
	case '{':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.LEFTBRACE, Literal: literal, Line: sc.line}
	case '}':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.RIGHTBRACE, Literal: literal, Line: sc.line}
	case ',':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.COMMA, Literal: literal, Line: sc.line}
	case '.':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.DOT, Literal: literal, Line: sc.line}
	case '-':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.MINUS, Literal: literal, Line: sc.line}
	case '+':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.PLUS, Literal: literal, Line: sc.line}
	case ';':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.SEMICOLON, Literal: literal, Line: sc.line}
	case '*':
		literal := sc.source[sc.start:sc.current]
		return token.Token{Type: token.STAR, Literal: literal, Line: sc.line}
	default:
		parseerror.Error(sc.line, "Parse error")
		return token.Token{Type: token.INVALID}
	}
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) advance() byte {
	sc.current++
	return sc.source[sc.current-1]
}
