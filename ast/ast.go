package ast

import "golox/token"

// Expr is the root class of expression nodes
type Expr interface {
}

// Binary is used for binary operators
type Binary struct {
	Expr
	Left     Expr
	Operator token.Token
	Right    Expr
}

// Grouping is used for parenthesized expressions
type Grouping struct {
	Expr
	Expression Expr
}

// Literal values
type Literal struct {
	Expr
	Value interface{}
}

// Unary is used for unary operators
type Unary struct {
	Expr
	Operator token.Token
	Right    Expr
}
