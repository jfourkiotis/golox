package parser

import (
	"golox/ast"
	"golox/scanner"
	"testing"
)

func testIntegerLiteral(expression ast.Expr, expected float64, t *testing.T) {
	literal, ok := expression.(*ast.Literal)
	if !ok {
		t.Fatalf("result is not ast.Literal. Got=%T", expression)
	}

	val, ok := literal.Value.(float64)
	if !ok {
		t.Fatalf("Literal.Value type not float64, got=%T", val)
	}

	if val != expected {
		t.Errorf("literal value not %v. got=%v", expected, val)
	}
}

func TestParseNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"2.4", 2.4},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		testIntegerLiteral(expression, test.expected, t)
	}
}

func TestParseStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"\"hello\"", "hello"},
		{"         \"world\"     ", "world"},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		literal, ok := expression.(*ast.Literal)
		if !ok {
			t.Fatalf("result is not ast.Literal. Got=%T", expression)
		}

		val, ok := literal.Value.(string)
		if !ok {
			t.Fatalf("Literal.Value type not float64, got=%T", val)
		}

		if val != test.expected {
			t.Errorf("literal value not %v. got=%v", test.expected, val)
		}
	}
}

func TestParseBooleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"false", false},
		{"         true  ", true},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		literal, ok := expression.(*ast.Literal)
		if !ok {
			t.Fatalf("result is not ast.Literal. Got=%T", expression)
		}

		val, ok := literal.Value.(bool)
		if !ok {
			t.Fatalf("Literal.Value type not float64, got=%T", val)
		}

		if val != test.expected {
			t.Errorf("literal value not %v. got=%v", test.expected, val)
		}
	}
}

func TestParseNil(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"nil"},
		{"         nil  "},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		literal, ok := expression.(*ast.Literal)
		if !ok {
			t.Fatalf("result is not ast.Literal. Got=%T", expression)
		}

		if literal.Value != nil {
			t.Errorf("literal value not %v. got=%v", nil, literal.Value)
		}
	}
}

func TestParseTernaryOperator(t *testing.T) {
	input := "7 ? 10 : 4"

	scanner := scanner.New(input)
	tokens := scanner.ScanTokens()
	parser := New(tokens)
	expression, _ := parser.expression()

	ternary, ok := expression.(*ast.Ternary)
	if !ok {
		t.Fatalf("Expected ast.Ternary operator. Got=%T", expression)
	}

	if ternary.QMark.Lexeme != "?" {
		t.Errorf("Expected ternary questionmark operator. Got=%v", ternary.QMark.Lexeme)
	}
	if ternary.Colon.Lexeme != ":" {
		t.Errorf("Expected ternary colon operator. Got=%v", ternary.Colon.Lexeme)
	}

	testIntegerLiteral(ternary.Condition, 7, t)
	testIntegerLiteral(ternary.Then, 10, t)
	testIntegerLiteral(ternary.Else, 4, t)
}

func TestParseBinaryOperators(t *testing.T) {
	tests := []struct {
		input    string
		left     float64
		operator string
		right    float64
	}{
		{"1+2", 1, "+", 2},
		{"1-2", 1, "-", 2},
		{"1*2", 1, "*", 2},
		{"1/2", 1, "/", 2},
		{"1**2", 1, "**", 2},
		{"1 != 5", 1, "!=", 5},
		{"1 == 5", 1, "==", 5},
		{"1 >= 5", 1, ">=", 5},
		{"1 <= 5", 1, "<=", 5},
		{"1 >  5", 1, ">", 5},
		{"1 <  5", 1, "<", 5},
		{"1 , 2", 1, ",", 2},
	}

	for i, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		binary, ok := expression.(*ast.Binary)
		if !ok {
			t.Fatalf("test[%v], result is not ast.Binary. Got=%T", i, expression)
		}

		if binary.Operator.Lexeme != test.operator {
			t.Errorf("binary operator value not %v. got=%v", test.operator, binary.Operator.Lexeme)
		}

		testIntegerLiteral(binary.Left, test.left, t)
		testIntegerLiteral(binary.Right, test.right, t)
	}
}

func TestParseUnaryOperators(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		right    interface{}
	}{
		{"!false", "!", false},
		{"!true", "!", true},
		{"-4.2", "-", 4.2},
	}

	for _, test := range tests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		unary, ok := expression.(*ast.Unary)
		if !ok {
			t.Fatalf("result is not ast.Binary. Got=%T", expression)
		}

		if unary.Operator.Lexeme != test.operator {
			t.Errorf("binary operator value not %v. got=%v", test.operator, unary.Operator.Lexeme)
		}

		right, ok := unary.Right.(*ast.Literal)
		if !ok {
			t.Fatalf("lhs not ast.Literal. Got=%T", right)
		}

		rval, ok := right.Value.(float64)
		if !ok {
			rval, ok := right.Value.(bool)
			if !ok {
				t.Fatalf("lhs type not float64 or bool. Got=%T", rval)
			} else if rval != test.right.(bool) {
				t.Errorf("lhs value not %v. got=%v", test.right, rval)
			}
		} else if rval != test.right.(float64) {
			t.Errorf("lhs value not %v. got=%v", test.right, rval)
		}
	}
}

func TestParseGroupedExpressions(t *testing.T) {
	numtests := []struct {
		input    string
		expected float64
	}{
		{"(5)", 5},
		{"(    5)", 5},
		{"     (      5     ) ", 5},
	}

	for _, test := range numtests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		g, ok := expression.(*ast.Grouping)
		if !ok {
			t.Fatalf("result is not ast.Grouping. Got=%T", expression)
		}

		testIntegerLiteral(g.Expression, test.expected, t)
	}
}

func TestParseAssignment(t *testing.T) {
	numtests := []struct {
		input            string
		expectedVariable string
		expectedValue    float64
	}{
		{"a = 5", "a", 5},
	}

	for _, test := range numtests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		expression, _ := parser.expression()

		assign, ok := expression.(*ast.Assign)
		if !ok {
			t.Fatalf("result is not *ast.Assign. Got=%T", expression)
		}

		if assign.Name.Lexeme != "a" {
			t.Errorf("Expected variable name to be 'a'. Got=%v", assign.Name.Lexeme)
		}

		testIntegerLiteral(assign.Value, test.expectedValue, t)
	}
}

func TestParseExpressionStatement(t *testing.T) {
	numtests := []struct {
		input    string
		expected float64
	}{
		{"5;", 5},
	}

	for _, test := range numtests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		stmtList := parser.Parse()

		if len(stmtList) != 1 {
			t.Fatalf("Expected 1 statement. Got=%v", len(stmtList))
		}

		exprStmt, ok := stmtList[0].(*ast.Expression)
		if !ok {
			t.Fatalf("Expected *ast.Expression. Got=%T", stmtList[0])
		}
		testIntegerLiteral(exprStmt.Expression, 5, t)
	}
}

func TestParseBlockStatement(t *testing.T) {
	numtests := []struct {
		input    string
		expected interface{}
	}{
		{"{{print 10;}}", 10},
	}

	for _, test := range numtests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		stmtList := parser.Parse()

		if len(stmtList) != 1 {
			t.Fatalf("Expected 1 statement. Got=%v", len(stmtList))
		}

		block1, ok := stmtList[0].(*ast.Block)
		if !ok {
			t.Fatalf("Expected *ast.Block. Got=%T", stmtList[0])
		}

		block2, ok := block1.Statements[0].(*ast.Block)
		if !ok {
			t.Fatalf("Expected *ast.Block. Got=%T", block1.Statements[0])
		}

		printStmt, ok := block2.Statements[0].(*ast.Print)
		if !ok {
			t.Fatalf("Expected *ast.Print. Got=%T", block2.Statements[0])
		}

		testIntegerLiteral(printStmt.Expression, 10, t)
	}
}

func TestParsePrintStatement(t *testing.T) {
	numtests := []struct {
		input    string
		expected interface{}
	}{
		{"print 5;", 5},
	}

	for _, test := range numtests {
		scanner := scanner.New(test.input)
		tokens := scanner.ScanTokens()
		parser := New(tokens)
		stmtList := parser.Parse()

		if len(stmtList) != 1 {
			t.Fatalf("Expected 1 statement. Got=%v", len(stmtList))
		}

		printStmt, ok := stmtList[0].(*ast.Print)
		if !ok {
			t.Fatalf("Expected *ast.Print. Got=%T", stmtList[0])
		}
		testIntegerLiteral(printStmt.Expression, 5, t)
	}
}

func TestParseUnaryPowerExpressions(t *testing.T) {
	input1 := "-5**2"

	s1 := scanner.New(input1)
	t1 := s1.ScanTokens()
	p1 := New(t1)
	e1, _ := p1.expression()

	u1, ok := e1.(*ast.Unary)

	if !ok {
		t.Fatalf("Unary expression expected. Got=%T", e1)
	}

	pow1 := u1.Right
	b1, ok := pow1.(*ast.Binary)

	if !ok {
		t.Fatalf("Binary expression expected. Got=%T", pow1)
	}

	if b1.Operator.Lexeme != "**" {
		t.Errorf("Power operator expected. Got=%v", b1.Operator)
	}

	testIntegerLiteral(b1.Left, 5, t)
	testIntegerLiteral(b1.Right, 2, t)

	input2 := "-5**-2"

	s2 := scanner.New(input2)
	t2 := s2.ScanTokens()
	p2 := New(t2)
	e2, _ := p2.expression()

	u2, ok := e2.(*ast.Unary)

	if !ok {
		t.Fatalf("Unary expression expected. Got=%T", e1)
	}

	pow2 := u2.Right
	b2, ok := pow2.(*ast.Binary)

	if !ok {
		t.Fatalf("Binary expression expected. Got=%T", pow2)
	}

	if b2.Operator.Lexeme != "**" {
		t.Errorf("Power operator expected. Got=%v", b2.Operator)
	}

	testIntegerLiteral(b1.Left, 5, t)

	u3, ok := b2.Right.(*ast.Unary)
	if !ok {
		t.Fatalf("Unary expression expected. Got=%T", b2.Right)
	}

	testIntegerLiteral(u3.Right, 2, t)

	input3 := "5**2**5"

	s3 := scanner.New(input3)
	t3 := s3.ScanTokens()
	p3 := New(t3)
	e3, _ := p3.expression()

	b3, ok := e3.(*ast.Binary)
	if !ok {
		t.Fatalf("Unary expression expected. Got=%T", b3)
	}

	testIntegerLiteral(b3.Left, 5, t)

	b4, ok := b3.Right.(*ast.Binary)
	if !ok {
		t.Fatalf("Unary expression expected. Got=%T", b3.Right)
	}

	testIntegerLiteral(b4.Left, 2, t)
	testIntegerLiteral(b4.Right, 5, t)
}
