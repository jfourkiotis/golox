package interpreter

import (
	"fmt"
	"golox/ast"
	"golox/env"
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
func Interpret(statements []ast.Stmt, env *env.Environment) {
	for _, stmt := range statements {
		_, err := Eval(stmt, env)
		if err != nil {
			runtimeerror.Print(err.Error())
		}
	}
}

// Eval evaluates the given AST
func Eval(node ast.Node, env *env.Environment) (interface{}, error) {
	switch n := node.(type) {
	case *ast.Literal:
		return n.Value, nil
	case *ast.Grouping:
		return Eval(n.Expression, env)
	case *ast.Unary:
		right, err := Eval(n.Right, env)
		if err != nil {
			return right, err
		} else if n.Operator.Type == token.MINUS {
			err := checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return -right.(float64), nil
		} else if n.Operator.Type == token.BANG {
			return !isTruthy(right), nil
		}
	case *ast.Binary:
		left, err := Eval(n.Left, env)
		if err != nil {
			return left, err
		}
		right, err := Eval(n.Right, env)
		if err != nil {
			return right, err
		}
		switch n.Operator.Type {
		case token.MINUS:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) - right.(float64), nil
		case token.PLUS:
			switch lhs := left.(type) {
			case float64:
				switch rhs := right.(type) {
				case float64:
					return lhs + rhs, nil
				}
			case string:
				switch rhs := right.(type) {
				case string:
					return lhs + rhs, nil
				}
			}
			return nil, fmt.Errorf("%s", operandsMustBeTwoNumbersOrTwoStrings)
		case token.SLASH:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) / right.(float64), nil
		case token.STAR:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) * right.(float64), nil
		case token.POWER:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return math.Pow(left.(float64), right.(float64)), nil
		case token.GREATER:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) > right.(float64), nil
		case token.GREATEREQUAL:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) >= right.(float64), nil
		case token.LESS:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) < right.(float64), nil
		case token.LESSEQUAL:
			err := checkNumberOperand(n.Operator, left, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			err = checkNumberOperand(n.Operator, right, operandMustBeANumber)
			if err != nil {
				return nil, err
			}
			return left.(float64) <= right.(float64), nil
		case token.BANGEQUAL:
			return !isEqual(left, right), nil
		case token.EQUALEQUAL:
			return isEqual(left, right), nil
		}
	case *ast.Ternary:
		cond, err := Eval(n.Condition, env)
		if err != nil {
			return cond, err
		}
		if isTruthy(cond) {
			return Eval(n.Then, env)
		}
		return Eval(n.Else, env)
	case *ast.Print:
		value, err := Eval(n.Expression, env)
		if err != nil {
			return value, err
		}
		fmt.Println(value)
		return nil, nil
	case *ast.Expression:
		r, err := Eval(n.Expression, env)
		if err != nil {
			return r, err
		}
		return nil, nil
	case *ast.Var:
		if n.Initializer != nil {
			value, err := Eval(n.Initializer, env)
			if err != nil {
				return nil, err
			}
			env.Define(n.Name.Lexeme, value)
		} else {
			env.Define(n.Name.Lexeme, nil)
		}
		return nil, nil
	case *ast.Variable:
		return env.Get(n.Name)
	case *ast.Assign:
		value, err := Eval(n.Value, env)
		if err != nil {
			return nil, err
		}
		if err = env.Assign(n.Name, value); err == nil {
			return value, nil
		}
		return nil, err
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

func checkNumberOperand(operator token.Token, value interface{}, msg string) error {
	switch value.(type) {
	case int, float64:
		return nil
	}
	return fmt.Errorf("%v\n[line %v]", msg, operator.Line)
}
