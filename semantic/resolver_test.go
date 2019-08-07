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

func TestReturnFromInitializer(t *testing.T) {
	input := `
	class Foo {
		init() {
			return 5;
		}
	}
	`
	s := scanner.New(input)
	tokens := s.ScanTokens()
	p := parser.New(tokens)
	statements := p.Parse()

	_, err := Resolve(statements)
	expected := "Cannot return a value from an initializer."
	if err == nil {
		t.Fatalf("Expected error.")
	} else if err.Error() != expected {
		t.Errorf("Expected error %q", expected)
	}
}

func TestResolveThis(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			`
			print this;
			`,
			"Cannot use 'this' outside of a class.",
		},
		{
			`
			fun foo() {
				print this;
			}
			`,
			"Cannot use 'this' outside of a class.",
		},
		{
			`
			fun foo() {
				fun bar(n) {
					return n + this;
				}
				return bar(10);
			}
			`,
			"Cannot use 'this' outside of a class.",
		},
		{
			`
			class Foo {
				class foo() {
					return this + 5;
				}
			}
			`,
			"Cannot use 'this' outside instance initializers or methods.",
		},
	}

	for _, test := range tests {
		s := scanner.New(test.input)
		tokens := s.ScanTokens()
		p := parser.New(tokens)
		statements := p.Parse()

		_, err := Resolve(statements)
		if err == nil {
			t.Fatalf("Expected error.")
		}
		if err.Error() != test.expected {
			t.Errorf("Expected error %q. Got %q", test.expected, err.Error())
		}
	}
}
