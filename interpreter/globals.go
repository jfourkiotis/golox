package interpreter

import (
	"golox/env"
	"time"
)

type loxCallable func([]interface{}) (interface{}, error)

type loxFunction struct {
	arity    int
	callable loxCallable
}

// GlobalEnv is the global environment
var GlobalEnv = env.NewGlobal()

func init() {
	GlobalEnv.Define("clock", &loxFunction{
		arity: 0,
		callable: func(args []interface{}) (interface{}, error) {
			return time.Now().Second(), nil
		},
	})
}
