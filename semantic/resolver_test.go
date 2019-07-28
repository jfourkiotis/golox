package semantic

import (
	"golox/parser"
	"golox/scanner"
	"testing"
)

func TestReturnResolution(t *testing.T) {
	input := `
	return 5;
	`
	s := scanner.New(input)
	tokens := s.ScanTokens()
	p := parser.New(tokens)
	statements := p.Parse()

	_, err := Resolve(statements)
	if err == nil {
		t.Errorf("top-level return not detected.")
	} else if err.Error() != "Cannot return from top-level code." {
		t.Errorf("resolution failed")
	}
}
