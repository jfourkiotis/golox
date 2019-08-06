package semanticerror

import (
	"fmt"
	"os"
)

// Print reports a semantic error
func Print(message string) {
	fmt.Fprintf(os.Stderr, "%v\n", message)
	HadError = true
}

// Make creates a new semantic error
func Make(message string) error {
	return fmt.Errorf("%s", message)
}

// HadError is true if an evaluation error was encountered
var HadError = false
