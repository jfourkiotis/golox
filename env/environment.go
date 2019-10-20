package env

import (
	"fmt"

	"github.com/dirkdev98/golox/runtimeerror"
	"github.com/dirkdev98/golox/token"
)

type uninitialized struct{}

var needsInitialization = &uninitialized{}

// Environment associates variables to values
type Environment struct {
	values        map[string]interface{}
	enclosing     *Environment
	indexedValues []interface{}
}

// New creates a new environment
func New(env *Environment) *Environment {
	return NewSized(env, 0)
}

// NewSized creates a new environment
func NewSized(env *Environment, size int) *Environment {
	return &Environment{values: make(map[string]interface{}), enclosing: env, indexedValues: make([]interface{}, size)}
}

// NewGlobal creates a new global environment
func NewGlobal() *Environment {
	return New(nil)
}

// Define binds a name to a new value
func (e *Environment) Define(name string, value interface{}, index int) {
	if index == -1 {
		e.values[name] = value
	} else {
		e.indexedValues[index] = value
	}
}

// DefineUnitialized creates a new variable. That variable must be initialized before
// used
func (e *Environment) DefineUnitialized(name string, index int) {
	if index == -1 {
		e.values[name] = needsInitialization
	} else {
		e.indexedValues[index] = needsInitialization
	}
}

// Get lookups a variable given a token.Token
func (e *Environment) Get(name token.Token, index int) (interface{}, error) {
	if index == -1 {
		v, prs := e.values[name.Lexeme]
		if prs {
			if v == needsInitialization {
				return nil, runtimeerror.Make(name, fmt.Sprintf("Uninitialized variable access: '%s'", name.Lexeme))
			}
			return v, nil
		}
		if e.enclosing != nil {
			return e.enclosing.Get(name, index)
		}
		return nil, runtimeerror.Make(name, fmt.Sprintf("Undefined variable '%v'", name.Lexeme))
	}
	return e.indexedValues[index], nil
}

// GetAt lookups a variable a certain distance up the chain of environments
func (e *Environment) GetAt(distance int, name token.Token, index int) (interface{}, error) {
	return e.Ancestor(distance).Get(name, index)
}

// Ancestor reaches an environment up the environment chain
func (e *Environment) Ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}
	return env
}

// Assign sets a new value to an old variable
func (e *Environment) Assign(name token.Token, index int, value interface{}) error {
	if index == -1 {
		if _, prs := e.values[name.Lexeme]; prs {
			e.values[name.Lexeme] = value
			return nil
		}
		if e.enclosing != nil {
			return e.enclosing.Assign(name, index, value)
		}
		return runtimeerror.Make(name, fmt.Sprintf("Undefined variable '%s'.", name.Lexeme))
	}
	e.indexedValues[index] = value
	return nil
}

// AssignAt sets a new value to an old variable
func (e *Environment) AssignAt(distance int, index int, name token.Token, value interface{}) error {
	return e.Ancestor(distance).Assign(name, index, value)
}
