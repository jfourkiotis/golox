package interpreter

import (
	"golox/env"
	"time"
)

// GlobalEnv is the global environment
var GlobalEnv = env.NewGlobal()

func init() {
	GlobalEnv.Define("clock", &NativeFunction{
		arity: 0,
		nativeCall: func(args []interface{}) (interface{}, error) {
			return time.Now().Second(), nil
		},
	})
}
