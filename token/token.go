package token

import "fmt"

// Type is the type of the token given as a string
type Type string

//
const (
	// single-character tokens
	LEFTPAREN  = "("
	RIGHTBRACE = ")"
	LEFTBRACE  = "{"
	RIGHTPAREN = "}"
	COMMA      = ","
	DOT        = "."
	MINUS      = "-"
	PLUS       = "+"
	SEMICOLON  = ";"
	SLASH      = "/"
	STAR       = "*"
	// one or two character tokens
	BANG         = "!"
	BANGEQUAL    = "!="
	EQUAL        = "="
	EQUALEQUAL   = "=="
	GREATER      = ">"
	GREATEREQUAL = ">="
	LESS         = "<"
	LESSEQUAL    = "<="
	// literals
	IDENTIFIER = "IDENT"
	STRING     = "STRING"
	NUMBER     = "NUMBER"
	// keywords
	AND     = "and"
	CLASS   = "class"
	ELSE    = "else"
	FALSE   = "false"
	FUN     = "fun"
	FOR     = "for"
	IF      = "if"
	NIL     = "nil"
	OR      = "or"
	PRINT   = "print"
	RETURN  = "return"
	SUPER   = "super"
	THIS    = "this"
	TRUE    = "true"
	VAR     = "var"
	WHILE   = "while"
	EOF     = "eof"
	INVALID = "__INVALID__"
)

//Token contains the lexeme read by the scanner
type Token struct {
	Type    Type
	Lexeme  string
	Literal interface{}
	Line    int
}

func (token *Token) String() string {
	return fmt.Sprintf("%s %s %v", token.Type, token.Lexeme, token.Literal)
}
