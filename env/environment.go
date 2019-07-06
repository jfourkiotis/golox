package env

import (
	"fmt"
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
