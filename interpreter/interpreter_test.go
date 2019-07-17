package interpreter

import (
	"golox/ast"
	"golox/env"
	"golox/parser"
	"golox/scanner"
	"golox/token"
	"math"
	"strings"
	"testing"
)

func TestEvalLiteral(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"5;", 5.0},
		{"false;", false},
		{"true;", true},
		{"\"hello\";", "hello"},
		{"(5);", 5.0},
		{"(false);", false},
		{"(true);", true},
		{"(\"hello\");", "hello"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %v", len(statements))
		}

		exprStmt, ok := statements[0].(*ast.Expression)
		if !ok {
			t.Fatalf("Expected *ast.ExpressionStmt. Got=%T", statements[0])
		}

		result, _ := Eval(exprStmt.Expression, env.NewGlobal())
		testLiteralEquality(result, test.expected, t)
	}
}

func TestEvalUnary(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"-5;", -5.0},
		{"!false;", true},
		{"true;", true},
		{"false;", false},
		{"!true;", false},
		{"!5;", false},
		{"!nil;", true},
		{"!\"hello\";", false},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %v", len(statements))
		}

		exprStmt, ok := statements[0].(*ast.Expression)
		if !ok {
			t.Fatalf("Expected *ast.ExpressionStmt. Got=%T", statements[0])
		}

		result, _ := Eval(exprStmt.Expression, env.NewGlobal())
		testLiteralEquality(result, test.expected, t)
	}
}

func TestEvalBinary(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"1 + 2;", 3.0},
		{"1 - 2;", -1.0},
		{"1 / 2;", 0.5},
		{"1 * 2;", 2.0},
		{"2 ** 2;", 4.0},
		{"\"hello \" + \"world\";", "hello world"},
		{"1 > 2;", false},
		{"1 >= 2;", false},
		{"1 < 2;", true},
		{"1 <= 2;", true},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"\"hello\" == 1;", false},
		{"\"hello\" == \"hello\";", true},
		{"nil != nil;", false},
		{"nil == 5;", false},
		{"5.2 == 5.2;", true},
		{"1.2 ** 3.4 ** 0.5 ** 0.9;", math.Pow(1.2, math.Pow(3.4, math.Pow(0.5, 0.9)))},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %v", len(statements))
		}

		exprStmt, ok := statements[0].(*ast.Expression)
		if !ok {
			t.Fatalf("Expected *ast.ExpressionStmt. Got=%T", statements[0])
		}

		result, _ := Eval(exprStmt.Expression, env.NewGlobal())
		testLiteralEquality(result, test.expected, t)
	}
}

func TestEvalBinaryPrecedence(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"1 - 2 - 3;", -4.0},
		{"1 + 2 * 3;", 7.0},
		{"2 ** 3 ** 2;", 512.0},
		{"-2 ** 3 ** -2;", -math.Pow(2.0, math.Pow(3.0, -2.0))},
		{"--2;", 2.0},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %v", len(statements))
		}

		exprStmt, ok := statements[0].(*ast.Expression)
		if !ok {
			t.Fatalf("Expected *ast.ExpressionStmt. Got=%T", statements[0])
		}

		result, _ := Eval(exprStmt.Expression, env.NewGlobal())
		testLiteralEquality(result, test.expected, t)
	}
}

func TestEvalTernary(t *testing.T) {
	tests := []struct {
		literal  string
		expected interface{}
	}{
		{"1 ? 2 : 3;", 2.0},
		{"nil ? 2 : 3;", 3.0},
	}

	for _, test := range tests {
		scanner := scanner.New(test.literal)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		if len(statements) != 1 {
			t.Fatalf("Expected 1 statement. Got %v", len(statements))
		}

		exprStmt, ok := statements[0].(*ast.Expression)
		if !ok {
			t.Fatalf("Expected *ast.ExpressionStmt. Got=%T", statements[0])
		}

		result, _ := Eval(exprStmt.Expression, env.NewGlobal())
		testLiteralEquality(result, test.expected, t)
	}
}

func testLiteralEquality(result interface{}, expected interface{}, t *testing.T) {
	switch r := result.(type) {
	case float64:
		testNumberEquality(r, expected, t)
	case bool:
		testBoolEquality(r, expected, t)
	case string:
		testStringEquality(r, expected, t)
	default:
		t.Fatalf("Unexpected result type. Got=%T", result)
	}
}

func testNumberEquality(lhs float64, expected interface{}, t *testing.T) {
	rhs, ok := expected.(float64)
	if !ok {
		t.Fatalf("Expected number. Got=%T", expected)
	}

	if rhs != lhs {
		t.Errorf("Numbers are not equal. Expected %v. Got %v", lhs, rhs)
	}
}

func testBoolEquality(lhs bool, expected interface{}, t *testing.T) {
	switch rhs := expected.(type) {
	case bool:
		if lhs != rhs {
			t.Errorf("Booleans are not equal. Expected %v. Got %v", lhs, rhs)
		}
	default:
		t.Fatalf("Expected bool. Got=%T", expected)
	}
}

func testStringEquality(lhs string, expected interface{}, t *testing.T) {
	switch rhs := expected.(type) {
	case string:
		if lhs != rhs {
			t.Errorf("Strings are not equal. Expected %v. Got %v", lhs, rhs)
		}
	default:
		t.Fatalf("Expected string. Got=%T", expected)
	}
}

func TestEnvironment(t *testing.T) {
	input := `
		var a = 5;
		var b = 10;
		var c = a * b;
	`
	scanner := scanner.New(input)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()

	env := GlobalEnv
	Interpret(statements, env)

	if a, err := env.Get(token.Token{Lexeme: "a"}); err != nil {
		t.Fatalf("Expected variable 'a' in env")
	} else if a.(float64) != 5.0 {
		t.Errorf("Expected a = 5. Got %v", a.(float64))
	}
	if b, err := env.Get(token.Token{Lexeme: "b"}); err != nil {
		t.Fatalf("Expected variable 'b' in env")
	} else if b.(float64) != 10.0 {
		t.Errorf("Expected b = 10. Got %v", b.(float64))
	}
	if c, err := env.Get(token.Token{Lexeme: "c"}); err != nil {
		t.Fatalf("Expected variable 'c' in env")
	} else if c.(float64) != 50.0 {
		t.Errorf("Expected c = 50. Got %v", c.(float64))
	}
}

func TestEvalAssignment(t *testing.T) {
	input := `
		var a = 5;
		var b = 10;
		var c = a * b;
		c = 20;
		b = 200;
		a = 2000;
	`
	scanner := scanner.New(input)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()

	env := env.NewGlobal()
	Interpret(statements, env)

	if a, err := env.Get(token.Token{Lexeme: "a"}); err != nil {
		t.Fatalf("Expected variable 'a' in env")
	} else if a.(float64) != 2000.0 {
		t.Errorf("Expected a = 2000. Got %v", a.(float64))
	}
	if b, err := env.Get(token.Token{Lexeme: "b"}); err != nil {
		t.Fatalf("Expected variable 'b' in env")
	} else if b.(float64) != 200.0 {
		t.Errorf("Expected b = 200. Got %v", b.(float64))
	}
	if c, err := env.Get(token.Token{Lexeme: "c"}); err != nil {
		t.Fatalf("Expected variable 'c' in env")
	} else if c.(float64) != 20.0 {
		t.Errorf("Expected c = 20. Got %v", c.(float64))
	}
}

func TestEvalWhileStatement(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{
			`
				var a = 5;
				var b = 0;
				while (a > 0) {
					a = a - 1;
					b = b + 1;
				} 
				print b;
			`, "5"},
		{
			`
				var a = 5;
				var b = 0;
				for( ; a > 0; a=a-1) {
					b = b + 1;
				} 
				print b;
			`, "5"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		out := &strings.Builder{}
		options.Writer = out
		env := env.NewGlobal()
		for _, stmt := range statements {
			_, err := Eval(stmt, env)
			if err != nil {
				t.Errorf("Runtime error when evaluating if-statement: %s", err.Error())
			}
		}

		outStr := strings.TrimSuffix(out.String(), "\n")
		if outStr != test.expectedOutput {
			t.Errorf("Expected <%s>. Got <%s>", test.expectedOutput, outStr)
		}
	}
}

func TestEvalGlobals(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"clock();"},
	}

	for _, test := range tests {
		s := scanner.New(test.input)
		tokens := s.ScanTokens()
		p := parser.New(tokens)
		statements := p.Parse()

		e, _ := statements[0].(*ast.Expression)
		v, _ := Eval(e.Expression, GlobalEnv)

		if v.(int) < 0 || v.(int) > 59 {
			t.Errorf("Expected a number in [0, 59]")
		}

	}
}

func TestEvalIfStatement(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{
			`
				var a = 5;
				if (a > 5) {
					print "yes";
				} else {
					print "no";
				}
			`, "no"},
		{
			`
				var a = 10;
				if (a > 5) {
					print "yes";
				}
			`, "yes"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		out := &strings.Builder{}
		options.Writer = out
		env := env.NewGlobal()
		for _, stmt := range statements {
			_, err := Eval(stmt, env)
			if err != nil {
				t.Errorf("Runtime error when evaluating if-statement: %s", err.Error())
			}
		}

		outStr := strings.TrimSuffix(out.String(), "\n")
		if outStr != test.expectedOutput {
			t.Errorf("Expected <%s>. Got <%s>", test.expectedOutput, outStr)
		}
	}
}

func TestEvalLogicalOperators(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput string
	}{
		{"print \"hi\" or 2;", "hi"},
		{"print nil or \"yes\";", "yes"},
		{"print 2 and \"la\";", "la"},
		{"print false and 2;", "false"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := parser.New(tokens)
		statements := parser.Parse()

		out := &strings.Builder{}
		options.Writer = out
		env := env.NewGlobal()
		for _, stmt := range statements {
			_, err := Eval(stmt, env)
			if err != nil {
				t.Errorf("Runtime error when evaluating logical operator: %s", err.Error())
			}
		}

		outStr := strings.TrimSuffix(out.String(), "\n")
		if outStr != test.expectedOutput {
			t.Errorf("Expected <%s>. Got <%s>", test.expectedOutput, outStr)
		}
	}
}
