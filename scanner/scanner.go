package scanner

import (
	"fmt"
	"golox/parseerror"
	"golox/token"
)

// Scanner transforms the source into tokens
type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []token.Token
}

// New creates a new scanner
func New(source string) Scanner {
	scanner := Scanner{source: source, line: 1, tokens: make([]token.Token, 0)}
	return scanner
}

// ScanTokens transforms the source into an array of tokens. The last token
// is always an token.EOF
func (sc *Scanner) ScanTokens() []token.Token {
	for !sc.isAtEnd() {
		// we're at the beginning of the next lexeme
		sc.start = sc.current
		sc.scanToken()
	}
	sc.tokens = append(sc.tokens, token.Token{Type: token.EOF})
	return sc.tokens
}

func (sc *Scanner) makeToken(tp token.Type) token.Token {
	lexeme := sc.source[sc.start:sc.current]
	return token.Token{Type: tp, Lexeme: lexeme, Line: sc.line}
}

func (sc *Scanner) addToken(tp token.Type) {
	sc.addTokenWithLiteral(tp, nil)
}

func (sc *Scanner) addTokenWithLiteral(tp token.Type, literal interface{}) {
	text := sc.source[sc.start:sc.current]
	sc.tokens = append(sc.tokens, token.Token{Type: tp, Lexeme: text, Literal: literal, Line: sc.line})
}

func (sc *Scanner) scanString() {
	for sc.peek() != '"' && !sc.isAtEnd() {
		if sc.peek() == '\n' {
			sc.line++
		}
		sc.advance()
	}

	// unterminated string
	if sc.isAtEnd() {
		parseerror.Error(sc.line, "Unterminated string.")
		return
	}

	// the closing ".
	sc.advance()

	// trim the surrounding quotes
	value := sc.source[sc.start+1 : sc.current-1]
	sc.addTokenWithLiteral(token.STRING, value)
}

func (sc *Scanner) scanToken() {
	c := sc.advance()

	switch c {
	case '(':
		sc.addToken(token.LEFTPAREN)
	case ')':
		sc.addToken(token.RIGHTPAREN)
	case '{':
		sc.addToken(token.LEFTBRACE)
	case '}':
		sc.addToken(token.RIGHTBRACE)
	case ',':
		sc.addToken(token.COMMA)
	case '.':
		sc.addToken(token.DOT)
	case '-':
		sc.addToken(token.MINUS)
	case '+':
		sc.addToken(token.PLUS)
	case ';':
		sc.addToken(token.SEMICOLON)
	case '*':
		sc.addToken(token.STAR)
	case '!':
		if sc.match('=') {
			sc.addToken(token.BANGEQUAL)
		} else {
			sc.addToken(token.BANG)
		}
	case '=':
		if sc.match('=') {
			sc.addToken(token.EQUALEQUAL)
		} else {
			sc.addToken(token.EQUAL)
		}
	case '<':
		if sc.match('=') {
			sc.addToken(token.LESSEQUAL)
		} else {
			sc.addToken(token.LESS)
		}
	case '>':
		if sc.match('=') {
			sc.addToken(token.GREATEREQUAL)
		} else {
			sc.addToken(token.GREATER)
		}
	case '/':
		if sc.match('/') {
			// A comment goes until the end of the line
			for sc.peek() != '\n' && !sc.isAtEnd() {
				sc.advance()
			}
		} else {
			sc.addToken(token.SLASH)
		}
	case '\n':
		sc.line++
	case ' ', '\r', '\t':
		// do nothing
	case '"':
		sc.scanString()
	default:
		parseerror.Error(sc.line, fmt.Sprintf("Unexpected character: %c", c))
	}
}

func (sc *Scanner) isAtEnd() bool {
	return sc.current >= len(sc.source)
}

// advance returns the current character and advances to the next
func (sc *Scanner) advance() byte {
	sc.current++
	return sc.source[sc.current-1]
}

func (sc *Scanner) match(expected byte) bool {
	if sc.isAtEnd() {
		return false
	}
	if sc.source[sc.current] != expected {
		return false
	}
	sc.current++
	return true
}

func (sc *Scanner) peek() byte {
	if sc.isAtEnd() {
		return 0
	}
	return sc.source[sc.current]
}
