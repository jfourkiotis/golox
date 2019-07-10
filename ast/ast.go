package ast

import "golox/token"

// Node is the root class of AST nodes
type Node interface{}

// Expr is the root class of expression nodes
type Expr interface {
	Node
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

// Ternary is the famous ?: operator
type Ternary struct {
	Expr
	Condition Expr
	QMark     token.Token
	Then      Expr
	Colon     token.Token
	Else      Expr
}

// Statements and state

// Assign is used for variable assignment
// name = value
type Assign struct {
	Expr
	Name  token.Token
	Value Expr
}

// Variable access expression
// print x
type Variable struct {
	Expr
	Name token.Token
}

// Stmt form a second hierarchy of syntax nodes independent of expressions
type Stmt interface {
	Node
}

// Block is a curly-braced block statement that defines a local scope
// {
//   ...
// }
type Block struct {
	Stmt
	Statements []Stmt
}

// Expression statement
type Expression struct {
	Stmt
	Expression Expr
}

// Print statement
// print 1 + 2
type Print struct {
	Stmt
	Expression Expr
}

// Var is the variable declaration statement
// var <name> = <initializer>
type Var struct {
	Stmt
	Name        token.Token
	Initializer Expr
}

// If is the classic if statement
type If struct {
	Stmt
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}
