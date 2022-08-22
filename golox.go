package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jfourkiotis/golox/ast"
	"github.com/jfourkiotis/golox/env"
	"github.com/jfourkiotis/golox/interpreter"
	"github.com/jfourkiotis/golox/parseerror"
	"github.com/jfourkiotis/golox/parser"
	"github.com/jfourkiotis/golox/runtimeerror"
	"github.com/jfourkiotis/golox/scanner"
	"github.com/jfourkiotis/golox/semantic"
	"github.com/jfourkiotis/golox/semanticerror"
	"io/ioutil"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runFile(file string) {
	dat, err := ioutil.ReadFile(file)
	check(err)
	run(string(dat), interpreter.GlobalEnv)
	if parseerror.HadError {
		os.Exit(65)
	} else if runtimeerror.HadError || semanticerror.HadError {
		os.Exit(70)
	}
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)
	env := interpreter.GlobalEnv
	for {
		fmt.Print("> ")
		dat, err := reader.ReadBytes('\n') // there is also ReadString
		check(err)
		run(string(dat), env)
		parseerror.HadError = false
		runtimeerror.HadError = false
		semanticerror.HadError = false
	}
}

func run(src string, env *env.Environment) {
	scanner := scanner.New(src)
	tokens := scanner.ScanTokens()
	parser := parser.New(tokens)
	statements := parser.Parse()
	if parseerror.HadError {
		return
	}
	resolution, err := semantic.Resolve(statements)
	if err != nil || semanticerror.HadError {
		semanticerror.Print(err.Error())
		return
	} else if len(resolution.Unused) != 0 {
		for stmt := range resolution.Unused {
			switch n := stmt.(type) {
			case *ast.Var:
				fmt.Fprintf(os.Stdout, "Unused variable %q [Line: %d]\n", n.Name.Lexeme, n.Name.Line)
			case *ast.Function:
				fmt.Fprintf(os.Stdout, "Unused function %q [Line: %d]\n", n.Name.Lexeme, n.Name.Line)
			case *ast.Class:
				fmt.Fprintf(os.Stdout, "Unused class %q [Line: %d]\n", n.Name.Lexeme, n.Name.Line)
			default:
				panic(fmt.Sprintf("Unexpected ast.Node type %T\n", stmt))
			}
		}
		err = semanticerror.Make(fmt.Sprintf("%d unused local variables/functions found", len(resolution.Unused)))
		return
	}
	interpreter.Interpret(statements, env, resolution)
}

func main() {
	flag.String("file", "", "the script file to execute")
	flag.Parse()

	args := flag.Args()
	if len(args) > 1 {
		fmt.Println("Usage: ./golox [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}
