package env

import (
	"fmt"
	"golox/runtimeerror"
	"golox/token"
)

// Environment associates variables to values
type Environment struct {
	values map[string]interface{}
}

// New creates a new environment
func New() *Environment {
	return &Environment{values: make(map[string]interface{})}
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
	return nil, fmt.Errorf("Undefined variable '%v'", name.Lexeme)
}

// Assign sets a new value to an old variable
func (e *Environment) Assign(name token.Token, value interface{}) error {
	if _, prs := e.values[name.Lexeme]; prs {
		e.values[name.Lexeme] = value
		return nil
	}
	return runtimeerror.MakeRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
}
