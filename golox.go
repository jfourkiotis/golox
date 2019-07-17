package main

import (
	"bufio"
	"flag"
	"fmt"
	"golox/env"
	"golox/interpreter"
	"golox/parseerror"
	"golox/parser"
	"golox/runtimeerror"
	"golox/scanner"
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
	} else if runtimeerror.HadError {
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
	interpreter.Interpret(statements, env)
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
