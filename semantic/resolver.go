package semantic

import (
	"fmt"
	"golox/ast"
	"golox/semanticerror"
	"golox/token"
)

// Unused local variables found by variable resolution
type Unused = map[ast.Stmt]bool // value is always false

// EnvSize keeps the environment size of each ast.Block & ast.Function
// nodes
type EnvSize = map[ast.Stmt]int

const (
	vDeclared = iota
	vDefined  = iota
)

type vInfo struct {
	name   string
	status int
	isUsed bool
	stmt   ast.Stmt
}

// rScope represents a Lox scope
type rScope = []vInfo

func scopeLookup(name string, scope rScope) int {
	for i := len(scope) - 1; i >= 0; i-- {
		if name == scope[i].name {
			return i
		}
	}
	return -1
}

// Resolution keeps important information about local variables and functions
type Resolution struct {
	Unused Unused
}

// NewResolution creates an empty resolution object
func NewResolution() Resolution {
	return Resolution{Unused: make(Unused)}
}

// Resolve performs name resolution to the given statements
func Resolve(statements []ast.Stmt) (Resolution, error) {
	resolution := NewResolution()
	resolver := &Resolver{scopes: make([]rScope, 0), currentFunction: ftNone, currentClass: ctNone}
	err := resolver.resolveStatements(statements, resolution)
	return resolution, err
}

const (
	ftNone        = iota
	ftFunction    = iota
	ftMethod      = iota
	ftInitializer = iota
)

const (
	ctNone  = iota
	ctClass = iota
)

// Resolver performs variable resolution on an AST
type Resolver struct {
	scopes          []rScope
	currentFunction int
	currentClass    int
}

func (r *Resolver) resolve(node ast.Node, res Resolution) error {
	switch n := node.(type) {
	case *ast.Block:
		r.pushScope()
		defer r.popScope(n, res)
		for _, stmt := range n.Statements {
			if err := r.resolve(stmt, res); err != nil {
				return err
			}
		}
	case *ast.Var:
		index, err := r.declare(n.Name, n)
		if err != nil {
			return nil
		}
		n.EnvIndex = index
		if n.Initializer != nil {
			if err := r.resolve(n.Initializer, res); err != nil {
				return err
			}
		}
		r.define(n.Name, n)
	case *ast.Variable:
		if len(r.scopes) != 0 {
			top := r.scopes[len(r.scopes)-1]
			index := scopeLookup(n.Name.Lexeme, top)
			if index >= 0 && top[index].status == vDeclared {
				return semanticerror.Make("Cannot read local variable in its own initializer.")
			}
		}
		index, depth := r.resolveLocal(n, n.Name, res)
		n.EnvIndex = index
		n.EnvDepth = depth
	case *ast.Assign:
		if err := r.resolve(n.Value, res); err != nil {
			return err
		}
		index, depth := r.resolveLocal(n, n.Name, res)
		n.EnvIndex = index
		n.EnvDepth = depth
	case *ast.Function:
		index, err := r.declare(n.Name, n)
		if err != nil {
			return err
		}
		n.EnvIndex = index
		r.define(n.Name, n)
		if err := r.resolveFunction(n, res, ftFunction); err != nil {
			return err
		}
	case *ast.Expression:
		if err := r.resolve(n.Expression, res); err != nil {
			return err
		}
	case *ast.If:
		if err := r.resolve(n.Condition, res); err != nil {
			return err
		}
		if err := r.resolve(n.ThenBranch, res); err != nil {
			return err
		}
		if n.ElseBranch != nil {
			if err := r.resolve(n.ElseBranch, res); err != nil {
				return err
			}
		}
	case *ast.Print:
		if err := r.resolve(n.Expression, res); err != nil {
			return err
		}
	case *ast.Return:
		if r.currentFunction == ftNone {
			return semanticerror.Make("Cannot return from top-level code.")
		}
		if n.Value != nil {
			if r.currentFunction == ftInitializer {
				return semanticerror.Make("Cannot return a value from an initializer.")
			}
			if err := r.resolve(n.Value, res); err != nil {
				return err
			}
		}
	case *ast.For:
		if err := r.resolve(n.Increment, res); err != nil {
			return err
		}
		if err := r.resolve(n.Condition, res); err != nil {
			return err
		}
		if err := r.resolve(n.Increment, res); err != nil {
			return err
		}
		if err := r.resolve(n.Statement, res); err != nil {
			return err
		}
	case *ast.While:
		if err := r.resolve(n.Condition, res); err != nil {
			return err
		}
		if err := r.resolve(n.Statement, res); err != nil {
			return err
		}
	case *ast.Binary:
		if err := r.resolve(n.Left, res); err != nil {
			return err
		}
		if err := r.resolve(n.Right, res); err != nil {
			return err
		}
	case *ast.Call:
		if err := r.resolve(n.Callee, res); err != nil {
			return err
		}

		for _, e := range n.Arguments {
			if err := r.resolve(e, res); err != nil {
				return err
			}
		}
	case *ast.Grouping:
		if err := r.resolve(n.Expression, res); err != nil {
			return err
		}
	case *ast.Logical:
		if err := r.resolve(n.Left, res); err != nil {
			return err
		}
		if err := r.resolve(n.Right, res); err != nil {
			return err
		}
	case *ast.Unary:
		if err := r.resolve(n.Right, res); err != nil {
			return err
		}
	case *ast.Class:
		enclosingClass := r.currentClass
		r.currentClass = ctClass

		resetCurrentClass := func() {
			r.currentClass = enclosingClass
		}

		defer resetCurrentClass()

		index, err := r.declare(n.Name, n)
		if err != nil {
			return err
		}
		n.EnvIndex = index
		r.define(n.Name, n)

		r.pushScope()
		defer r.popScope(n, res)

		top := r.scopes[len(r.scopes)-1]
		top = append(top, vInfo{name: "this", status: vDefined, isUsed: true})
		r.scopes[len(r.scopes)-1] = top // FIXME: is this needed

		for _, method := range n.Methods {
			declaration := ftMethod
			if method.Name.Lexeme == "init" {
				declaration = ftInitializer
			}
			if err := r.resolveFunction(method, res, declaration); err != nil {
				return err
			}
		}
	case *ast.Get:
		if err := r.resolve(n.Expression, res); err != nil {
			return err
		}
	case *ast.Set:
		if err := r.resolve(n.Value, res); err != nil {
			return err
		}
		if err := r.resolve(n.Object, res); err != nil {
			return err
		}
	case *ast.This:
		if r.currentClass == ctNone {
			return semanticerror.Make("Cannot use 'this' outside of a class.")
		}
		index, depth := r.resolveLocal(n, n.Keyword, res)
		n.EnvIndex = index
		n.EnvDepth = depth
	}
	return nil
}

func (r *Resolver) resolveStatements(statements []ast.Stmt, res Resolution) error {
	for _, stmt := range statements {
		if err := r.resolve(stmt, res); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) resolveFunction(function *ast.Function, res Resolution, ftype int) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = ftype

	resetCurrentFunction := func() {
		r.currentFunction = enclosingFunction
	}

	defer resetCurrentFunction()

	r.pushScope()
	defer r.popScope(function, res)

	if !function.IsProperty() {
		for _, param := range function.Params {
			if _, err := r.declare(param, nil); err != nil {
				return err
			}
			r.define(param, nil)
		}
	}
	return r.resolveStatements(function.Body, res)
}

func (r *Resolver) resolveLocal(expr ast.Expr, name token.Token, res Resolution) (int, int) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		scope := r.scopes[i]
		index := scopeLookup(name.Lexeme, scope)
		if index >= 0 {
			scope[index].isUsed = true
			return index, len(r.scopes) - i - 1
		}
	}
	return -1, -1
}

func (r *Resolver) pushScope() {
	r.scopes = append(r.scopes, make(rScope, 0))
}

func (r *Resolver) popScope(stmt ast.Stmt, res Resolution) {
	top := r.scopes[len(r.scopes)-1]
	r.scopes = r.scopes[:len(r.scopes)-1]
	for _, info := range top {
		// info.node is nil for function parameters
		if !info.isUsed && info.stmt != nil {
			res.Unused[info.stmt] = true
		}
	}

	if block, ok := stmt.(*ast.Block); ok {
		block.EnvSize = len(top)
	} else if function, ok := stmt.(*ast.Function); ok {
		function.EnvSize = len(top)
	}
}

// node is nil for function params
func (r *Resolver) declare(name token.Token, node ast.Node) (int, error) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		index := scopeLookup(name.Lexeme, scope)
		if index >= 0 {
			return 0, semanticerror.Make(
				fmt.Sprintf("Variable '%s' already declared in this scope.", name.Lexeme))
		}
		scope = append(scope, vInfo{name: name.Lexeme, status: vDeclared, isUsed: false, stmt: node})
		r.scopes[len(r.scopes)-1] = scope
		return len(scope) - 1, nil
	}
	return -1, nil
}

func (r *Resolver) define(name token.Token, node ast.Node) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		index := scopeLookup(name.Lexeme, scope)
		scope[index].status = vDefined
	}
}
