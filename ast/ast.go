package ast

import (
	"fmt"
	"golox/token"
	"strings"
)

// Node is the root class of AST nodes
type Node interface {
	String() string
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

// String pretty prints the operator
func (b *Binary) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(b.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(b.Left.String())
	sb.WriteString(" ")
	sb.WriteString(b.Right.String())
	sb.WriteString(")")
	return sb.String()
}

// Grouping is used for parenthesized expressions
type Grouping struct {
	Expr
	Expression Expr
}

// String pretty prints the expression grouping
func (g *Grouping) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(g.Expression.String())
	sb.WriteString(")")
	return sb.String()
}

// Literal values
type Literal struct {
	Expr
	Value interface{}
}

// String pretty prints the literal
func (l *Literal) String() string {
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

// String pretty prints the unary operator
func (u *Unary) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(u.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(u.Right.String())
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

// String pretty prints the unary operator
func (t *Ternary) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(t.Condition.String())
	sb.WriteString(" ")
	sb.WriteString(t.QMark.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(t.Then.String())
	sb.WriteString(" ")
	sb.WriteString(t.Colon.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(t.Else.String())
	sb.WriteString(")")
	return sb.String()
}

// Statements and state

// Assign is used for variable assignment
// name = value
type Assign struct {
	Expr
	Name     token.Token
	Value    Expr
	EnvIndex int
	EnvDepth int
}

// String pretty prints the assignment statement
func (a *Assign) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("=")
	sb.WriteString(" ")
	sb.WriteString(a.Name.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(a.Value.String())
	sb.WriteString(")")
	return sb.String()
}

// Variable access expression
// print x
type Variable struct {
	Expr
	Name     token.Token
	EnvIndex int
	EnvDepth int
}

// String pretty prints the assignment expression
func (v *Variable) String() string {
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
	EnvSize    int
}

// String pretty prints the block statement
func (b *Block) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	for _, stmt := range b.Statements {
		sb.WriteString(stmt.String())
	}
	sb.WriteString(")")
	return sb.String()
}

// Expression statement
type Expression struct {
	Stmt
	Expression Expr
}

// String pretty prints the expression statement
func (e *Expression) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(e.Expression.String())
	sb.WriteString(")")
	return sb.String()
}

// Print statement
// print 1 + 2
type Print struct {
	Stmt
	Expression Expr
}

// String pretty prints the print statement
func (p *Print) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("print")
	sb.WriteString(" ")
	sb.WriteString(p.Expression.String())
	sb.WriteString(")")
	return sb.String()
}

// Var is the variable declaration statement
// var <name> = <initializer>
type Var struct {
	Stmt
	Name        token.Token
	Initializer Expr
	EnvIndex    int
}

// String pretty prints the var declaration
func (v *Var) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("var")
	sb.WriteString(" ")
	sb.WriteString(v.Name.Lexeme)
	sb.WriteString(" ")
	if v.Initializer != nil {
		sb.WriteString(v.Initializer.String())
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

// String pretty prints the if statement
func (i *If) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("if")
	sb.WriteString(" ")
	sb.WriteString(i.Condition.String())
	sb.WriteString(" ")
	sb.WriteString(i.ThenBranch.String())
	sb.WriteString(" ")
	sb.WriteString(i.ElseBranch.String())
	sb.WriteString(")")
	return sb.String()
}

// For ...
type For struct {
	Stmt
	Initializer Expr
	Condition   Expr
	Increment   Expr
	Statement   Stmt
}

// String pretty prints the for statement
func (f *For) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("for")
	sb.WriteString(" ")
	sb.WriteString("(")
	sb.WriteString(f.Initializer.String())
	sb.WriteString(")")
	sb.WriteString(" ")
	sb.WriteString("(")
	sb.WriteString(f.Condition.String())
	sb.WriteString(")")
	sb.WriteString(" ")
	sb.WriteString("(")
	sb.WriteString(f.Increment.String())
	sb.WriteString(")")
	sb.WriteString(" ")
	sb.WriteString(f.Statement.String())
	sb.WriteString(")")
	return sb.String()
}

// While is the classic while statement
type While struct {
	Stmt
	Condition Expr
	Statement Stmt
}

// String pretty prints the while statement
func (w *While) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("while")
	sb.WriteString(" ")
	sb.WriteString(w.Condition.String())
	sb.WriteString(" ")
	sb.WriteString(w.Statement.String())
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

// String pretty prints the unary operator
func (l *Logical) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(l.Operator.Lexeme)
	sb.WriteString(" ")
	sb.WriteString(l.Left.String())
	sb.WriteString(" ")
	sb.WriteString(l.Right.String())
	sb.WriteString(")")
	return sb.String()
}

// Call is the node of a function call
type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

// String pretty prints the call operator
func (c *Call) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("call")
	sb.WriteString(" ")
	sb.WriteString(c.Callee.String())
	sb.WriteString(" ")
	for _, e := range c.Arguments {
		sb.WriteString(e.String())
		sb.WriteString(" ")
	}
	sb.WriteString(")")
	return sb.String()
}

// Function is the function definition node
type Function struct {
	Name     token.Token
	Params   []token.Token
	Body     []Stmt
	EnvSize  int
	EnvIndex int
}

// String pretty prints the function
func (f *Function) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("fun")
	sb.WriteString(" ")
	sb.WriteString(f.Name.Lexeme)
	sb.WriteString(" ")
	sb.WriteString("(")
	for _, p := range f.Params {
		sb.WriteString(p.Lexeme)
		sb.WriteString(" ")
	}
	sb.WriteString(")")
	sb.WriteString(" ")
	sb.WriteString("(")
	for _, stmt := range f.Body {
		sb.WriteString(stmt.String())
		sb.WriteString(" ")
	}
	sb.WriteString(")")
	sb.WriteString(")")
	return sb.String()
}

// Return is used to return from a function
type Return struct {
	Stmt
	Keyword token.Token
	Value   Expr
}

// String pretty prints the function
func (r *Return) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("return")
	sb.WriteString(" ")
	sb.WriteString(r.Value.String())
	sb.WriteString(" ")
	sb.WriteString(")")
	return sb.String()
}

// Break is used to return from a function
type Break struct {
	Stmt
	Token token.Token
}

// String pretty prints the function
func (b *Break) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("break")
	sb.WriteString(")")
	return sb.String()
}

// Continue is used to return from a function
type Continue struct {
	Stmt
	Token token.Token
}

// String pretty prints the function
func (c *Continue) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString("continue")
	sb.WriteString(")")
	return sb.String()
}
