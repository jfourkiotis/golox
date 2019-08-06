package interpreter

import (
	"fmt"
	"golox/runtimeerror"
	"golox/token"
)

// Class ...
type Class struct {
	Callable
	Name    string
	Methods map[string]*UserFunction
}

func (c *Class) String() string {
	return fmt.Sprintf("<class %s>", c.Name)
}

// Call is the operation that executes a class constructor
func (c *Class) Call(arguments []interface{}) (interface{}, error) {
	instance := &ClassInstance{Class: c, fields: make(map[string]interface{})}
	if initializer, prs := c.Methods["init"]; prs {
		_, err := initializer.Bind(instance).Call(arguments)
		if err != nil {
			return nil, err
		}
	}

	return instance, nil
}

// Arity returns the number of allowed parameters in the class constructor
// which is always 0
func (c *Class) Arity() int {
	if initializer, prs := c.Methods["init"]; prs {
		return initializer.Arity()
	}
	return 0
}

// ClassInstance is a user defined class instance
type ClassInstance struct {
	Class  *Class
	fields map[string]interface{}
}

func (c *ClassInstance) String() string {
	return fmt.Sprintf("<class-instance %s>", c.Class.Name)
}

// Get accesses the property
func (c *ClassInstance) Get(name token.Token) (interface{}, error) {
	if v, prs := c.fields[name.Lexeme]; prs {
		return v, nil
	}

	if m, prs := c.Class.Methods[name.Lexeme]; prs {
		return m.Bind(c), nil
	}
	return nil, runtimeerror.Make(name, fmt.Sprintf("Undefined property '%s'", name.Lexeme))
}

// Set accesses the property
func (c *ClassInstance) Set(name token.Token, value interface{}) (interface{}, error) {
	c.fields[name.Lexeme] = value
	return nil, nil
}
