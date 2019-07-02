package interpreter

import (
	"golox/ast"
	"golox/token"
	"math"
)

// Eval evaluates the given AST
func Eval(node ast.Expr) interface{} {
	switch n := node.(type) {
	case *ast.Literal:
		return n.Value
	case *ast.Grouping:
		return Eval(n.Expression)
	case *ast.Unary:
		right := Eval(n.Right)
		if n.Operator.Type == token.MINUS {
			return -right.(float64)
		} else if n.Operator.Type == token.BANG {
			return !isTruthy(right)
		}
	case *ast.Binary:
		left := Eval(n.Left)
		right := Eval(n.Right)
		switch n.Operator.Type {
		case token.MINUS:
			return left.(float64) - right.(float64)
		case token.PLUS:
			switch lhs := left.(type) {
			case float64:
				switch rhs := right.(type) {
				case float64:
					return lhs + rhs
				}
			case string:
				switch rhs := right.(type) {
				case string:
					return lhs + rhs
				}
			}
		case token.SLASH:
			return left.(float64) / right.(float64)
		case token.STAR:
			return left.(float64) * right.(float64)
		case token.POWER:
			return math.Pow(left.(float64), right.(float64))
		case token.GREATER:
			return left.(float64) > right.(float64)
		case token.GREATEREQUAL:
			return left.(float64) >= right.(float64)
		case token.LESS:
			return left.(float64) < right.(float64)
		case token.LESSEQUAL:
			return left.(float64) <= right.(float64)
		case token.BANGEQUAL:
			return !isEqual(left, right)
		case token.EQUALEQUAL:
			return isEqual(left, right)
		}
	case *ast.Ternary:
		cond := Eval(n.Condition)
		if isTruthy(cond) {
			return Eval(n.Then)
		}
		return Eval(n.Else)
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
