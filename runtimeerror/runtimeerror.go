package runtimeerror

import (
	"fmt"
	"golox/token"
	"os"
)

// Print reports a runtime error
func Print(message string) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	HadError = true
}

// MakeRuntimeError creates a new runtime error
func MakeRuntimeError(token token.Token, message string) error {
	return fmt.Errorf("%s\n[line %d]", message, token.Line)
}

// HadError is true if an evaluation error was encountered
var HadError = false
