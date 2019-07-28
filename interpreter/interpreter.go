package interpreter

import (
	"fmt"
	"golox/ast"
	"golox/env"
	"golox/runtimeerror"
	"golox/semantic"
	"golox/token"
	"io"
	"math"
	"os"
)

const (
	operandMustBeANumber                 = "Operand must be a number"
	operandsMustBeTwoNumbersOrTwoStrings = "Operands must be two numbers or two strings"
)

// Options contains customization points for the interpreter behavior
type Options struct {
	Writer io.Writer
}

var options = &Options{Writer: os.Stdout}

// return
type returnError struct {
	error
	value interface{}
}

// break
type breakError struct {
	error
}

// continue
type continueError struct {
	error
}

// Interpret tries to calculate the result of an expression, or print a message
// if an error occurs
func Interpret(statements []ast.Stmt, env *env.Environment, res semantic.Resolution) {
	OldGlobalEnv := GlobalEnv
	GlobalEnv = env
	for _, stmt := range statements {
		_, err := Eval(stmt, env, res)
		if err != nil {
			runtimeerror.Print(err.Error())
		}
	}
	GlobalEnv = OldGlobalEnv
}

// Eval evaluates the given AST
func Eval(node ast.Node, environment *env.Environment, res semantic.Resolution) (interface{}, error) {
	switch n := node.(type) {
	case *ast.Literal:
		return n.Value, nil
	case *ast.Grouping:
		return Eval(n.Expression, environment, res)
	case *ast.Unary:
		right, err := Eval(n.Right, environment, res)
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
		left, err := Eval(n.Left, environment, res)
		if err != nil {
			return left, err
		}
		right, err := Eval(n.Right, environment, res)
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
		case token.COMMA:
			_, err := Eval(n.Left, environment, res)
			if err != nil {
				return nil, err
			}
			return Eval(n.Right, environment, res)
		}
	case *ast.Ternary:
		cond, err := Eval(n.Condition, environment, res)
		if err != nil {
			return cond, err
		}
		if isTruthy(cond) {
			return Eval(n.Then, environment, res)
		}
		return Eval(n.Else, environment, res)
	case *ast.Print:
		value, err := Eval(n.Expression, environment, res)
		if err != nil {
			return value, err
		}
		fmt.Fprintln(options.Writer, value)
		return nil, nil
	case *ast.Expression:
		r, err := Eval(n.Expression, environment, res)
		if err != nil {
			return r, err
		}
		return nil, nil
	case *ast.Var:
		if n.Initializer != nil {
			value, err := Eval(n.Initializer, environment, res)
			if err != nil {
				return nil, err
			}
			environment.Define(n.Name.Lexeme, value, n.EnvIndex)
		} else {
			environment.DefineUnitialized(n.Name.Lexeme, n.EnvIndex)
		}
		return nil, nil
	case *ast.Variable:
		if n.EnvDepth >= 0 {
			return environment.GetAt(n.EnvDepth, n.Name, n.EnvIndex)
		}
		return GlobalEnv.Get(n.Name, n.EnvIndex)
	case *ast.Assign:
		value, err := Eval(n.Value, environment, res)
		if err != nil {
			return nil, err
		}

		if n.EnvDepth >= 0 {
			if err := environment.AssignAt(n.EnvDepth, n.EnvIndex, n.Name, value); err == nil {
				return value, nil
			}
			return nil, err
		} else if err := GlobalEnv.Assign(n.Name, n.EnvIndex, value); err == nil {
			return value, nil
		}
		return nil, err
	case *ast.Block:
		newEnvironment := env.NewSized(environment, n.EnvSize)
		for _, stmt := range n.Statements {
			_, err := Eval(stmt, newEnvironment, res)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case *ast.If:
		var err error
		if condValue, err := Eval(n.Condition, environment, res); err == nil {
			if isTruthy(condValue) {
				return Eval(n.ThenBranch, environment, res)
			} else if n.ElseBranch != nil {
				return Eval(n.ElseBranch, environment, res)
			}
		}
		return nil, err
	case *ast.For:
		if n.Initializer != nil {
			_, err := Eval(n.Initializer, environment, res)
			if err != nil {
				return nil, err
			}
		}
		for {

			if n.Condition != nil {
				condition, err := Eval(n.Condition, environment, res)

				if err != nil {
					return nil, err
				}

				if !isTruthy(condition) {
					break
				}
			}

			_, err := Eval(n.Statement, environment, res)

			if err != nil {
				if _, ok := err.(breakError); ok {
					break
				} else if _, ok := err.(continueError); ok {
					if n.Increment != nil {
						_, err := Eval(n.Increment, environment, res)
						if err != nil {
							return nil, err
						}
					}
					continue
				}
				return nil, err
			}

			if n.Increment != nil {
				_, err := Eval(n.Increment, environment, res)
				if err != nil {
					return nil, err
				}
			}
		}
		return nil, nil
	case *ast.While:
		for {
			condition, err := Eval(n.Condition, environment, res)

			if err != nil {
				return nil, err
			}
			if !isTruthy(condition) {
				break
			}
			_, err = Eval(n.Statement, environment, res)

			if err != nil {
				if _, ok := err.(breakError); ok {
					break
				} else if _, ok := err.(continueError); ok {
					continue
				}
				return nil, err
			}
		}
		return nil, nil
	case *ast.Logical:
		left, err := Eval(n.Left, environment, res)
		if err != nil {
			return nil, err
		}
		if n.Operator.Type == token.OR {
			if isTruthy(left) {
				return left, nil
			}
		} else if n.Operator.Type == token.AND {
			if !isTruthy(left) {
				return left, nil
			}
		}
		return Eval(n.Right, environment, res)
	case *ast.Call:
		callee, err := Eval(n.Callee, environment, res)
		if err != nil {
			return nil, err
		}

		args := make([]interface{}, 0)
		for _, arg := range n.Arguments {
			a, err := Eval(arg, environment, res)
			if err == nil {
				args = append(args, a)
			} else {
				return nil, err
			}
		}

		function, ok := callee.(Callable)

		if !ok {
			return nil, runtimeerror.MakeRuntimeError(n.Paren, "Can only call functions and classes.")
		}

		if function.Arity() != len(args) {
			return nil, runtimeerror.MakeRuntimeError(n.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(args)))
		}

		return function.Call(args)
	case *ast.Function:
		function := NewUserFunction(n, environment, res, n.EnvSize)
		environment.Define(n.Name.Lexeme, function, n.EnvIndex)
		return nil, nil
	case *ast.Return:
		var value interface{}
		var err error
		if n.Value != nil {
			value, err = Eval(n.Value, environment, res)
			if err != nil {
				return nil, err
			}
		}
		return nil, returnError{value: value}
	case *ast.Break:
		return nil, breakError{}
	case *ast.Continue:
		return nil, continueError{}
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
