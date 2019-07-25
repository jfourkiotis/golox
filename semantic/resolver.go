package semantic

import (
	"fmt"
	"golox/ast"
	"golox/semanticerror"
	"golox/token"
)

// Unused local variables found by variable resolution
type Unused = map[ast.Stmt]bool // value is always false

// Locals is the output of the resolve phase
type Locals = map[ast.Expr]int

const (
	vDeclared = iota
	vDefined  = iota
)

type vInfo struct {
	status int
	isUsed bool
	stmt   ast.Stmt
}

// rScope represents a Lox scope
type rScope = map[string]*vInfo

// Resolution keeps important information about local variables and functions
type Resolution struct {
	Locals Locals
	Unused Unused
}

// Resolve performs name resolution to the given statements
func Resolve(statements []ast.Stmt) (Resolution, error) {
	resolution := Resolution{Locals: make(Locals), Unused: make(Unused)}
	resolver := &Resolver{scopes: make([]rScope, 0), currentFunction: ftNone}
	err := resolver.resolveStatements(statements, resolution.Locals, resolution.Unused)
	return resolution, err
}

const (
	ftNone     = iota
	ftFunction = iota
)

// Resolver performs variable resolution on an AST
type Resolver struct {
	scopes          []rScope
	currentFunction int
}

func (r *Resolver) resolve(node ast.Node, locals Locals, unused Unused) error {
	switch n := node.(type) {
	case *ast.Block:
		r.pushScope()
		defer r.popScope(unused)
		for _, stmt := range n.Statements {
			if err := r.resolve(stmt, locals, unused); err != nil {
				return err
			}
		}
	case *ast.Var:
		if err := r.declare(n.Name, n); err != nil {
			return nil
		}
		if n.Initializer != nil {
			if err := r.resolve(n.Initializer, locals, unused); err != nil {
				return err
			}
		}
		r.define(n.Name, n)
	case *ast.Variable:
		if len(r.scopes) != 0 {
			if v, ok := r.scopes[len(r.scopes)-1][n.Name.Lexeme]; ok && v.status == vDeclared {
				return semanticerror.MakeSemanticError("Cannot read local variable in its own initializer.")
			}
		}
		r.resolveLocal(n, n.Name, locals, unused)
	case *ast.Assign:
		if err := r.resolve(n.Value, locals, unused); err != nil {
			return err
		}
		r.resolveLocal(n, n.Name, locals, unused)
	case *ast.Function:
		if err := r.declare(n.Name, n); err != nil {
			return err
		}
		r.define(n.Name, n)
		if err := r.resolveFunction(n, locals, unused, ftFunction); err != nil {
			return err
		}
	case *ast.Expression:
		if err := r.resolve(n.Expression, locals, unused); err != nil {
			return err
		}
	case *ast.If:
		if err := r.resolve(n.Condition, locals, unused); err != nil {
			return err
		}
		if err := r.resolve(n.ThenBranch, locals, unused); err != nil {
			return err
		}
		if n.ElseBranch != nil {
			if err := r.resolve(n.ElseBranch, locals, unused); err != nil {
				return err
			}
		}
	case *ast.Print:
		if err := r.resolve(n.Expression, locals, unused); err != nil {
			return err
		}
	case *ast.Return:
		if r.currentFunction == ftNone {
			return semanticerror.MakeSemanticError("Cannot return from top-level code.")
		}
		if n.Value != nil {
			if err := r.resolve(n.Value, locals, unused); err != nil {
				return err
			}
		}
	case *ast.While:
		if err := r.resolve(n.Condition, locals, unused); err != nil {
			return err
		}
		if err := r.resolve(n.Statement, locals, unused); err != nil {
			return err
		}
	case *ast.Binary:
		if err := r.resolve(n.Left, locals, unused); err != nil {
			return err
		}
		if err := r.resolve(n.Right, locals, unused); err != nil {
			return err
		}
	case *ast.Call:
		if err := r.resolve(n.Callee, locals, unused); err != nil {
			return err
		}

		for _, e := range n.Arguments {
			if err := r.resolve(e, locals, unused); err != nil {
				return err
			}
		}
	case *ast.Grouping:
		if err := r.resolve(n.Expression, locals, unused); err != nil {
			return err
		}
	case *ast.Logical:
		if err := r.resolve(n.Left, locals, unused); err != nil {
			return err
		}
		if err := r.resolve(n.Right, locals, unused); err != nil {
			return err
		}
	case *ast.Unary:
		if err := r.resolve(n.Right, locals, unused); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStatements(statements []ast.Stmt, locals Locals, unused Unused) error {
	for _, stmt := range statements {
		if err := r.resolve(stmt, locals, unused); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(function *ast.Function, locals Locals, unused Unused, ftype int) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = ftype

	resetCurrentFunction := func() {
		r.currentFunction = enclosingFunction
	}

	defer resetCurrentFunction()

	r.pushScope()
	defer r.popScope(unused)

	for _, param := range function.Params {
		if err := r.declare(param, nil); err != nil {
			return err
		}
		r.define(param, nil)
	}
	return r.resolveStatements(function.Body, locals, unused)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token, locals Locals, unused Unused) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if v, ok := r.scopes[i][name.Lexeme]; ok {
			locals[expr] = len(r.scopes) - i - 1
			v.isUsed = true
			return
		}
	}
}

func (r *Resolver) pushScope() {
	r.scopes = append(r.scopes, make(rScope))
}

func (r *Resolver) popScope(unused Unused) {
	top := r.scopes[len(r.scopes)-1]
	r.scopes = r.scopes[:len(r.scopes)-1]
	for _, info := range top {
		// info.node is nil for function parameters
		if !info.isUsed && info.stmt != nil {
			unused[info.stmt] = true
		}
	}
}

// node is nil for function params
func (r *Resolver) declare(name token.Token, node ast.Node) error {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		if _, ok := scope[name.Lexeme]; ok {
			return semanticerror.MakeSemanticError(
				fmt.Sprintf("Variable '%s' already declared in this scope.", name.Lexeme))
		}
		scope[name.Lexeme] = &vInfo{status: vDeclared, isUsed: false, stmt: node}
	}
	return nil
}

func (r *Resolver) define(name token.Token, node ast.Node) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		scope[name.Lexeme].status = vDefined
	}
}
