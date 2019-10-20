package interpreter

import (
	"fmt"

	"github.com/dirkdev98/golox/ast"
	"github.com/dirkdev98/golox/env"
	"github.com/dirkdev98/golox/semantic"
	"github.com/dirkdev98/golox/token"
)

type loxCallable func([]interface{}) (interface{}, error)

// Callable is the generic interface for functions/classes in Lox
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

// String returns the name of the native function
func (n *NativeFunction) String() string {
	return fmt.Sprintf("<native/%p>", n.nativeCall)
}

// UserFunction are functions defined in Lox code
type UserFunction struct {
	Callable
	Definition    *ast.Function
	Closure       *env.Environment
	Resolution    semantic.Resolution
	IsInitializer bool
	envSize       int
}

// NewUserFunction creates a new UserFunction
func NewUserFunction(def *ast.Function, closure *env.Environment, res semantic.Resolution, envSize int) *UserFunction {
	return &UserFunction{Definition: def, Closure: closure, Resolution: res, envSize: envSize, IsInitializer: false}
}

// Call executes a user-defined Lox function
func (u *UserFunction) Call(arguments []interface{}) (interface{}, error) {
	env := env.NewSized(u.Closure, u.envSize)

	if !u.Definition.IsProperty() {
		for i, param := range u.Definition.Params {
			env.Define(param.Lexeme, arguments[i], i)
		}
	}

	for _, stmt := range u.Definition.Body {
		_, err := Eval(stmt, env, u.Resolution)

		if err != nil {
			if r, ok := err.(returnError); ok {
				if u.IsInitializer {
					return u.Closure.GetAt(0, token.Token{Lexeme: "this"}, 0)
				}
				return r.value, nil
			}
			return nil, err
		}
	}

	if u.IsInitializer {
		return u.Closure.GetAt(0, token.Token{Lexeme: "this"}, 0)
	}
	return nil, nil
}

// Arity returns the number of arguments of the user-defined function
func (u *UserFunction) Arity() int {
	return len(u.Definition.Params)
}

// String returns the name of the user-function
func (u *UserFunction) String() string {
	return u.Definition.Name.Lexeme
}

// Bind creates a new instance method
func (u *UserFunction) Bind(instance *ClassInstance) *UserFunction {
	thisEnv := env.NewSized(u.Closure, 1)
	thisEnv.Define("this", instance, 0)
	return &UserFunction{Definition: u.Definition, Closure: thisEnv, Resolution: u.Resolution, envSize: u.envSize, IsInitializer: u.IsInitializer}
}
