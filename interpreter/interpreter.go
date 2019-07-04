package interpreter

import (
	"fmt"
	"golox/ast"
	"golox/runtimeerror"
	"golox/token"
	"math"
)

const (
	operandMustBeANumber                 = "Operand must be a number"
	operandsMustBeTwoNumbersOrTwoStrings = "Operands must be two numbers or two strings"
)

// Interpret tries to calculate the result of an expression, or print a message
// if an error occurs
func Interpret(statements []ast.Stmt) {
	for stmt := range statements {
		r, ok := Eval(stmt)
		if !ok {
			runtimeerror.Print(r.(string))
		}
	}
}

// Eval evaluates the given AST
func Eval(node ast.Node) (interface{}, bool) {
	switch n := node.(type) {
	case *ast.Literal:
		return n.Value, true
	case *ast.Grouping:
		return Eval(n.Expression)
	case *ast.Unary:
		right, ok := Eval(n.Right)
		if !ok {
			return right, ok
		} else if n.Operator.Type == token.MINUS {
			msg, ok := checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return -right.(float64), true
		} else if n.Operator.Type == token.BANG {
			return !isTruthy(right), true
		}
	case *ast.Binary:
		left, ok := Eval(n.Left)
		if !ok {
			return left, ok
		}
		right, ok := Eval(n.Right)
		if !ok {
			return right, ok
		}
		switch n.Operator.Type {
		case token.MINUS:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) - right.(float64), true
		case token.PLUS:
			switch lhs := left.(type) {
			case float64:
				switch rhs := right.(type) {
				case float64:
					return lhs + rhs, true
				}
			case string:
				switch rhs := right.(type) {
				case string:
					return lhs + rhs, true
				}
			}
			return operandsMustBeTwoNumbersOrTwoStrings, false
		case token.SLASH:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) / right.(float64), true
		case token.STAR:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) * right.(float64), true
		case token.POWER:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return math.Pow(left.(float64), right.(float64)), true
		case token.GREATER:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) > right.(float64), true
		case token.GREATEREQUAL:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) >= right.(float64), true
		case token.LESS:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) < right.(float64), true
		case token.LESSEQUAL:
			msg, ok := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			msg, ok = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if !ok {
				return msg, ok
			}
			return left.(float64) <= right.(float64), true
		case token.BANGEQUAL:
			return !isEqual(left, right), true
		case token.EQUALEQUAL:
			return isEqual(left, right), true
		}
	case *ast.Ternary:
		cond, ok := Eval(n.Condition)
		if !ok {
			return cond, ok
		}
		if isTruthy(cond) {
			return Eval(n.Then)
		}
		return Eval(n.Else)
	case *ast.Print:
		value, ok := Eval(n.Expression)
		if !ok {
			return value, ok
		}
		fmt.Println(value)
		return nil, true
	case *ast.Expression:
		r, ok := Eval(n.Expression)
		if !ok {
			return r, ok
		}
		return nil, ok
	}
	panic("Fatal error")
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	} else if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func isEqual(left interface{}, right interface{}) bool {
	// nil is only equal to nil
	if left == nil && right == nil {
		return true
	}
	if left == nil {
		return false
	}
	return left == right
}

func checkNumberOperand(operator token.Token, value interface{}, msg string) (string, bool) {
	switch value.(type) {
	case int, float64:
		return "", true
	}
	return fmt.Sprintf("%v\n[line %v]", msg, operator.Line), false
}
