package interpreter

import (
	"golox/ast"
	"golox/env"
)

type loxCallable func([]interface{}) (interface{}, error)

// Callable is the generic interface for functions in Lox
type Callable interface {
	Arity() int
	Call([]interface{}) (interface{}, error)
}

// NativeFunction is a builtin Lox function
type NativeFunction struct {
	Callable
	nativeCall loxCallable
	arity      int
}

// Call is the operation that executes a builtin function
func (n *NativeFunction) Call(arguments []interface{}) (interface{}, error) {
	return n.nativeCall(arguments)
}

// Arity returns the number of allowed parameters for the native function
func (n *NativeFunction) Arity() int {
	return n.arity
}

// UserFunction are functions defined in Lox code
type UserFunction struct {
	Callable
	Definition *ast.Function
}

// Call executes a user-defined Lox function
func (u *UserFunction) Call(arguments []interface{}) (interface{}, error) {
	env := env.New(GlobalEnv)
	for i, param := range u.Definition.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	for _, stmt := range u.Definition.Body {
		_, err := Eval(stmt, env)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

// Arity returns the number of arguments of the user-defined function
func (u *UserFunction) Arity() int {
	return len(u.Definition.Params)
}
