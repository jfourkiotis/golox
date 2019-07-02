package runtimeerror

import (
	"fmt"
	"os"
)

// Print reports a runtime error
func Print(message string) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	HadError = true
}

// HadError is true if an evaluation error was encountered
var HadError = false
