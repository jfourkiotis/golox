package env

import (
	"fmt"
	"golox/runtimeerror"
	"golox/token"
)

// Environment associates variables to values
type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

// New creates a new environment
func New(env *Environment) *Environment {
	return &Environment{values: make(map[string]interface{}), enclosing: env}
}

// NewGlobal creates a new global environment
func NewGlobal() *Environment {
	return New(nil)
}

// Define binds a name to a new value
func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

// Get lookups a variable given a token.Token
func (e *Environment) Get(name token.Token) (interface{}, error) {
	v, prs := e.values[name.Lexeme]
	if prs {
		return v, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, fmt.Errorf("Undefined variable '%v'", name.Lexeme)
}

// Assign sets a new value to an old variable
func (e *Environment) Assign(name token.Token, value interface{}) error {
	if _, prs := e.values[name.Lexeme]; prs {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return runtimeerror.MakeRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
