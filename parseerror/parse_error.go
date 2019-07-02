package parseerror

import (
	"fmt"
	"os"
)

// HadError is true if a scanner/parser error was encountered
var HadError = false

// Error reports in stderr an error encountered during parsing
func Error(line int, message string) {
	report(line, "", message)
	HadError = true
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", line, where, message)
}
