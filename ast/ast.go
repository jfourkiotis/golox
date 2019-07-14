package ast

import (
	"fmt"
	"golox/token"
	"strings"
)

// Node is the root class of AST nodes
type Node interface {
	ToString() string
}

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

// ToString pretty prints the operator
func (b *Binary) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(b.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(b.Left.ToString())
	sb.WriteString(" ")
	sb.WriteString(b.Right.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Grouping is used for parenthesized expressions
type Grouping struct {
	Expr
	Expression Expr
}

// ToString pretty prints the expression grouping
func (g *Grouping) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(g.Expression.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Literal values
type Literal struct {
	Expr
	Value interface{}
}

// ToString pretty prints the literal
func (l *Literal) ToString() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%v", l.Value))
	return sb.String()
}

// Unary is used for unary operators
type Unary struct {
	Expr
	Operator token.Token
	Right    Expr
}

// ToString pretty prints the unary operator
func (u *Unary) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(u.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(u.Right.ToString())
	sb.WriteString(")")
	return sb.String()
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

// ToString pretty prints the unary operator
func (t *Ternary) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(t.Condition.ToString())
	sb.WriteString(" ")
	sb.WriteString(t.QMark.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(t.Then.ToString())
	sb.WriteString(" ")
	sb.WriteString(t.Colon.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(t.Else.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Statements and state

// Assign is used for variable assignment
// name = value
type Assign struct {
	Expr
	Name  token.Token
	Value Expr
}

// ToString pretty prints the assignment statement
func (a *Assign) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("=")
	sb.WriteString(" ")
	sb.WriteString(a.Name.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(a.Value.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Variable access expression
// print x
type Variable struct {
	Expr
	Name token.Token
}

// ToString pretty prints the assignment expression
func (v *Variable) ToString() string {
	var sb strings.Builder
	sb.WriteString(v.Name.Lexeme)
	return sb.String()
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

// ToString pretty prints the block statement
func (b *Block) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	for _, stmt := range b.Statements {
		sb.WriteString(stmt.ToString())
	}
	sb.WriteString(")")
	return sb.String()
}

// Expression statement
type Expression struct {
	Stmt
	Expression Expr
}

// ToString pretty prints the expression statement
func (e *Expression) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(e.Expression.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Print statement
// print 1 + 2
type Print struct {
	Stmt
	Expression Expr
}

// ToString pretty prints the print statement
func (p *Print) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("print")
	sb.WriteString(" ")
	sb.WriteString(p.Expression.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Var is the variable declaration statement
// var <name> = <initializer>
type Var struct {
	Stmt
	Name        token.Token
	Initializer Expr
}

// ToString pretty prints the var declaration
func (v *Var) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("var")
	sb.WriteString(" ")
	sb.WriteString(v.Name.Lexeme)
	sb.WriteString(" ")
	if v.Initializer != nil {
		sb.WriteString(v.Initializer.ToString())
	} else {
		sb.WriteString("nil")
	}
	sb.WriteString(")")
	return sb.String()
}

// If is the classic if statement
type If struct {
	Stmt
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

// ToString pretty prints the if statement
func (i *If) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("if")
	sb.WriteString(" ")
	sb.WriteString(i.Condition.ToString())
	sb.WriteString(" ")
	sb.WriteString(i.ThenBranch.ToString())
	sb.WriteString(" ")
	sb.WriteString(i.ElseBranch.ToString())
	sb.WriteString(")")
	return sb.String()
}

// While is the classic while statement
type While struct {
	Stmt
	Condition Expr
	Statement Stmt
}

// ToString pretty prints the while statement
func (w *While) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("while")
	sb.WriteString(" ")
	sb.WriteString(w.Condition.ToString())
	sb.WriteString(" ")
	sb.WriteString(w.Statement.ToString())
	sb.WriteString(")")
	return sb.String()
}

// Logical is used for the "or" and "and" operators.
type Logical struct {
	Expr
	Left     Expr
	Operator token.Token
	Right    Expr
}

// ToString pretty prints the unary operator
func (l *Logical) ToString() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(l.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(l.Left.ToString())
	sb.WriteString(" ")
	sb.WriteString(l.Right.ToString())
	sb.WriteString(")")
	return sb.String()
}
