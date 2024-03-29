package runtimeerror

import (
	"fmt"
	"github.com/jfourkiotis/golox/token"
	"os"
)

// Print reports a runtime error
func Print(message string) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	HadError = true
}

// Make creates a new runtime error
func Make(token token.Token, message string) error {
	return fmt.Errorf("%s\n[line %d]", message, token.Line)
}

// HadError is true if an evaluation error was encountered
var HadError = false
