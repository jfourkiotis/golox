package semantic

import (
	"fmt"
	"golox/ast"
	"golox/semanticerror"
	"golox/token"
)

// Locals is the output of the resolve phase
type Locals = map[ast.Expr]int

// Scope represents a Lox scope
type Scope = map[string]bool

// Resolve performs name resolution to the given statements
func Resolve(statements []ast.Stmt) (Locals, error) {
	locals := make(Locals)
	resolver := &Resolver{scopes: make([]Scope, 0)}
	err := resolver.resolveStatements(statements, locals)
	return locals, err
}

// Resolver performs variable resolution on an AST
type Resolver struct {
	scopes []Scope
}

func (r *Resolver) resolve(node ast.Node, locals Locals) error {
	switch n := node.(type) {
	case *ast.Block:
		r.pushScope()
		defer r.popScope()
		for _, stmt := range n.Statements {
			if err := r.resolve(stmt, locals); err != nil {
				return err
			}
		}
	case *ast.Var:
		if err := r.declare(n.Name); err != nil {
			return nil
		}
		if n.Initializer != nil {
			if err := r.resolve(n.Initializer, locals); err != nil {
				return err
			}
		}
		r.define(n.Name)
	case *ast.Variable:
		if len(r.scopes) != 0 {
			if b, ok := r.scopes[len(r.scopes)-1][n.Name.Lexeme]; ok && !b {
				return semanticerror.MakeSemanticError("Cannot read local variable in its own initializer.")
			}
		}
		r.resolveLocal(n, n.Name, locals)
	case *ast.Assign:
		if err := r.resolve(n.Value, locals); err != nil {
			return err
		}
		r.resolveLocal(n, n.Name, locals)
	case *ast.Function:
		if err := r.declare(n.Name); err != nil {
			return err
		}
		r.define(n.Name)
		if err := r.resolveFunction(n, locals); err != nil {
			return err
		}
	case *ast.Expression:
		if err := r.resolve(n.Expression, locals); err != nil {
			return err
		}
	case *ast.If:
		if err := r.resolve(n.Condition, locals); err != nil {
			return err
		}
		if err := r.resolve(n.ThenBranch, locals); err != nil {
			return err
		}
		if n.ElseBranch != nil {
			if err := r.resolve(n.ElseBranch, locals); err != nil {
				return err
			}
		}
	case *ast.Print:
		if err := r.resolve(n.Expression, locals); err != nil {
			return err
		}
	case *ast.Return:
		if n.Value != nil {
			if err := r.resolve(n.Value, locals); err != nil {
				return err
			}
		}
	case *ast.While:
		if err := r.resolve(n.Condition, locals); err != nil {
			return err
		}
		if err := r.resolve(n.Statement, locals); err != nil {
			return err
		}
	case *ast.Binary:
		if err := r.resolve(n.Left, locals); err != nil {
			return err
		}
		if err := r.resolve(n.Right, locals); err != nil {
			return err
		}
	case *ast.Call:
		if err := r.resolve(n.Callee, locals); err != nil {
			return err
		}

		for _, e := range n.Arguments {
			if err := r.resolve(e, locals); err != nil {
				return err
			}
		}
	case *ast.Grouping:
		if err := r.resolve(n.Expression, locals); err != nil {
			return err
		}
	case *ast.Logical:
		if err := r.resolve(n.Left, locals); err != nil {
			return err
		}
		if err := r.resolve(n.Right, locals); err != nil {
			return err
		}
	case *ast.Unary:
		if err := r.resolve(n.Right, locals); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveStatements(statements []ast.Stmt, locals Locals) error {
	for _, stmt := range statements {
		if err := r.resolve(stmt, locals); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(function *ast.Function, locals Locals) error {
	r.pushScope()
	defer r.popScope()

	for _, param := range function.Params {
		if err := r.declare(param); err != nil {
			return err
		}
		r.define(param)
	}
	return r.resolveStatements(function.Body, locals)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token, locals Locals) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			locals[expr] = len(r.scopes) - i - 1
			return
		}
	}
}

func (r *Resolver) pushScope() {
	r.scopes = append(r.scopes, make(Scope))
}

func (r *Resolver) popScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name token.Token) error {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		if _, ok := scope[name.Lexeme]; ok {
			return semanticerror.MakeSemanticError(
				fmt.Sprintf("Variable '%s' already declared in this scope.", name.Lexeme))
		}
		scope[name.Lexeme] = false
	}
	return nil
}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		scope[name.Lexeme] = true
	}
}
