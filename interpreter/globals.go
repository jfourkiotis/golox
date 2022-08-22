package interpreter

import (
	"github.com/jfourkiotis/golox/env"
	"time"
)

// GlobalEnv is the global environment
var GlobalEnv = env.NewGlobal()
var globals = GlobalEnv

func init() {
	GlobalEnv.Define("clock", &NativeFunction{
		arity: 0,
		nativeCall: func(args []interface{}) (interface{}, error) {
			return time.Now().Second(), nil
		},
	}, -1)
}

// ResetGlobalEnv resets the GlobalEnv to its original reference
func ResetGlobalEnv() {
	GlobalEnv = globals
}
