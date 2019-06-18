package token

import (
	"testing"
)

func TestTokenString(t *testing.T) {
	tok := Token{Type: NUMBER, Lexeme: "3", Literal: 3, Line: 40}
	if tok.String() != "NUMBER 3 3" {
		t.Fatalf("expected=NUMBER 3 3, got=%q", tok.String())
	}
}
