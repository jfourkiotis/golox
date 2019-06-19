package scanner

import (
	"golox/parseerror"
	"golox/token"
)

var singleCharTokens = map[byte]token.Type{
	'(': token.LEFTPAREN,
	')': token.RIGHTPAREN,
	'{': token.LEFTBRACE,
	'}': token.RIGHTBRACE,
	',': token.COMMA,
	'.': token.DOT,
	'-': token.MINUS,
	'+': token.PLUS,
	';': token.SEMICOLON,
	'*': token.STAR,
}

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

func (sc *Scanner) makeToken(tp token.Type) token.Token {
	literal := sc.source[sc.start:sc.current]
	return token.Token{Type: tp, Literal: literal, Line: sc.line}
}

func (sc *Scanner) scanToken() token.Token {
	c := sc.advance()
	tp, ok := singleCharTokens[c]
	if ok {
		return sc.makeToken(tp)
	}
	parseerror.Error(sc.line, "Parse error")
	return token.Token{Type: token.INVALID}
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

func (sc *Scanner) advance() byte {
	sc.current++
	return sc.source[sc.current-1]
}
