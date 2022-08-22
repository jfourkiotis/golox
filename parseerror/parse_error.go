package parseerror

import (
	"fmt"
	"github.com/jfourkiotis/golox/token"
	"os"
)

// HadError is true if a scanner/parser error was encountered
var HadError = false

// LogMessage reports in stderr an error encountered during parsing
func LogMessage(line int, message string) {
	report(line, "", message)
	HadError = true
}

// LogError reports in stderr an error encountered during parsing
func LogError(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err.Error())
	HadError = true
}

// MakeError renders an parsing error as a string
func MakeError(tok token.Token, message string) error {
	if tok.Type == token.EOF {
		return fmt.Errorf("[line %v] Error at end: %s", tok.Line, message)
	}
	return fmt.Errorf("[line %v] Error at '%s': %s", tok.Line, tok.Lexeme, message)
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s: %s\n", line, where, message)
}
