package semantic

import (
	"golox/ast"
	"golox/parser"
	"golox/scanner"
	"testing"
)

func TestVariableResolution(t *testing.T) {
	input := `
	var a = "global";
	{
		fun f() {
			print a;
		}
	
		f();
		var a = "block";
		f();
	}	
	`
	s := scanner.New(input)
	tokens := s.ScanTokens()
	p := parser.New(tokens)
	statements := p.Parse()

	locals, _, err := Resolve(statements)
	if err != nil {
		t.Fatalf("Resolving declarations failed. Got=%q", err.Error())
	}

	for e, i := range locals {
		if variable, ok := e.(*ast.Variable); ok && variable.Name.Lexeme == "a" {
			if i != 2 {
				t.Errorf("Variable %s was declared in depth -%d. Got= -%d", variable.Name.Lexeme, 2, i)
			}
		}
	}
}

func TestReturnResolution(t *testing.T) {
	input := `
	return 5;
	`
	s := scanner.New(input)
	tokens := s.ScanTokens()
	p := parser.New(tokens)
	statements := p.Parse()

	_, _, err := Resolve(statements)
	if err == nil {
		t.Errorf("top-level return not detected.")
	} else if err.Error() != "Cannot return from top-level code." {
		t.Errorf("resolution failed")
	}
}

func TestUnusedLocalVariables(t *testing.T) {
	tests := []struct {
		input          string
		expectedUnused int
	}{
		{
			`
			{
				var a = 5;
				var b = 10;
	
				fun f() {
					var c = 20;
				}
	
				f();
	
				print a;
			}
			`, 2},
		{
			`
				var a = 5;
			`,
			0},
		{
			`
			var a = 5;
			fun f() {
				a = 10;
			}
			`, 0},
	}

	for _, test := range tests {
		s := scanner.New(test.input)
		tokens := s.ScanTokens()
		p := parser.New(tokens)
		statements := p.Parse()

		_, unused, _ := Resolve(statements)

		if len(unused) != test.expectedUnused {
			t.Fatalf("Expected %d unused variables. Got=%d", test.expectedUnused, len(unused))
		}
	}
}
