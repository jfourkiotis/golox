package scanner

import (
	"golox/token"
	"testing"
)

func TestScanTokens(t *testing.T) {
	input := `(`
	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LEFTPAREN, "("},
	}

	scanner := New(input)
	for i, test := range tests {
		tokens := scanner.ScanTokens()
		if len(tokens) != 2 {
			t.Fatalf("tests[%d] - number of tokens is wrong. expected=%q, got=%q", i, len(tokens), 2)
		}
		if test.expectedType != tokens[0].Type {
			t.Fatalf("tests[%d] - token type is wrong. expected=%q, got=%q", i, test.expectedType, tokens[0].Type)
		}
		if test.expectedLiteral != tokens[0].Literal {
			t.Fatalf("tests[%d] - token literal is wrong. expected=%q, got=%q", i, test.expectedLiteral, tokens[0].Literal)
		}
	}
}
