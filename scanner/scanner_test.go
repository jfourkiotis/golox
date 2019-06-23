package scanner

import (
	"golox/token"
	"testing"
)

func TestScanTokens(t *testing.T) {
	input := `48( ){     }, .-      +; * =	!=><<=>===!/
				// a comment
			"some string"    
		=


		8.66 and class 
		else false	for	fun			if	nil



		or
		print
		return
					super
					this
							true var while whilevar ifnil
	`
	tests := []struct {
		expectedType   token.Type
		expectedLexeme string
	}{
		{token.NUMBER, "48"},
		{token.LEFTPAREN, "("},
		{token.RIGHTPAREN, ")"},
		{token.LEFTBRACE, "{"},
		{token.RIGHTBRACE, "}"},
		{token.COMMA, ","},
		{token.DOT, "."},
		{token.MINUS, "-"},
		{token.PLUS, "+"},
		{token.SEMICOLON, ";"},
		{token.STAR, "*"},
		{token.EQUAL, "="},
		{token.BANGEQUAL, "!="},
		{token.GREATER, ">"},
		{token.LESS, "<"},
		{token.LESSEQUAL, "<="},
		{token.GREATEREQUAL, ">="},
		{token.EQUALEQUAL, "=="},
		{token.BANG, "!"},
		{token.SLASH, "/"},
		{token.STRING, "\"some string\""},
		{token.EQUAL, "="},
		{token.NUMBER, "8.66"},
		{token.AND, "and"},
		{token.CLASS, "class"},
		{token.ELSE, "else"},
		{token.FALSE, "false"},
		{token.FOR, "for"},
		{token.FUN, "fun"},
		{token.IF, "if"},
		{token.NIL, "nil"},
		{token.OR, "or"},
		{token.PRINT, "print"},
		{token.RETURN, "return"},
		{token.SUPER, "super"},
		{token.THIS, "this"},
		{token.TRUE, "true"},
		{token.VAR, "var"},
		{token.WHILE, "while"},
		{token.IDENTIFIER, "whilevar"},
		{token.IDENTIFIER, "ifnil"},
	}

	scanner := New(input)
	tokens := scanner.ScanTokens()

	if len(tests) != len(tokens)-1 {
		t.Fatalf("tests - number of token is wrong. expected=%d, got=%d", len(tests), len(tokens)-1)
	}

	for i, test := range tests {
		if test.expectedType != tokens[i].Type {
			t.Fatalf("tests[%d] - token type is wrong. expected=%q, got=%q", i, test.expectedType, tokens[i].Type)
		}
		if test.expectedLexeme != tokens[i].Lexeme {
			t.Fatalf("tests[%d] - token literal is wrong. expected=%q, got=%q", i, test.expectedLexeme, tokens[i].Lexeme)
		}
	}

	if tokens[len(tokens)-1].Type != token.EOF {
		t.Fatalf("tests - the last token is not EOF")
	}

}
