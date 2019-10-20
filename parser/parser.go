package parser

import (
	"github.com/dirkdev98/golox/ast"
	"github.com/dirkdev98/golox/parseerror"
	"github.com/dirkdev98/golox/token"
)

/*
program    -> declaration* EOF ;
declaration -> classDecl
            | varDecl
            | funDecl
			| stmt
varDecl    -> "var" IDENTIFIER ( "=" expression )? ";" ;
funDecl    -> "fun" function ;
classDecl  -> "class" IDENTIFIER  ( "<" IDENTIFIER )? "{" (function|property)* "}" ;
function   -> "class"? IDENTIFIER "(" parameters? ")" block ;
property   -> IDENTIFIER block ;
stmt       -> exprStmt
            | ifStmt
			| printStmt
			| returnStmt
			| whileStmt
			| forStmt
			| breakStmt
			| continueStmt
			| block
breakStmt  -> "break" ";" ;
continueStmt -> "continue" ";" ;
returnStmt -> "return" expression? ";" ;
ifStmt     -> "if" "(" expression ")" statement ( "else " statement )? ;
whileStmt  -> "while" "(" expression ")" statement ;
forStmt    -> "for" "(" ( varDecl | exprStmt | ";" ) expression? ";" expression? ")" statement ;
block      -> "{" declaration* "}"
exprStmt   -> expression ";" ;
printStmt  -> "print" expression ";" ;
expression -> comma ;
comma      -> assignment ( "," assignment ) * ;
assignment -> (call "." )? IDENTIFIER "=" assignment
			| logic_or ;
logic_or   -> logic_and ( "or" logic_and )* ;
logic_and  -> ternary ( "and" ternary ) * ;
ternary    -> equality "?"  expression ":" expression ;
equality   -> comparison ( ( "!=" | "==") comparison )* ;
comparison -> addition ( ( ">" | ">=" | "<" | "<=") addition )*;
addition   -> multiplication ( ( "+" | "-" ) multiplication )*;
multiplication -> unary ( ( "/" | "*" ) unary )*;
unary      -> ( "!" | "-" ) unary;
			| power ;
power      -> call ( "**" unary ) *
call       -> primary ( "(" arguments? ")" | "." IDENTIFIER )* ;
arguments  -> expression ( "," expression )* ;
primary    -> NUMBER | STRING | "false" | "true" | "nil" | "this" | "super"
			| "(" expression ")"
			| IDENTIFIER ;
*/

// Parser will transform an array of tokens to an AST.
// Use parser.New to create a new Parser. Do not create a Parser directly
type Parser struct {
	tokens  []token.Token
	current int
	inloop  bool // used when checking stray break/continue statements
}

// New creates a new parser
func New(tokens []token.Token) Parser {
	return Parser{tokens, 0, false}
}

// Parse is the driver function that begins parsing
func (p *Parser) Parse() []ast.Stmt {
	statements := make([]ast.Stmt, 0)
	for !p.isAtEnd() {
		// FIXME: p.declaration may return nil
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() ast.Stmt {
	var stmt ast.Stmt
	var err error

	checkError := func() {
		if err != nil {
			p.synchronize()
			parseerror.LogError(err)
			stmt = nil
		}
	}
	defer checkError()

	if p.match(token.CLASS) {
		stmt, err = p.classDeclaration()
	} else if p.match(token.VAR) {
		stmt, err = p.varDeclaration()
	} else if p.match(token.FUN) {
		stmt, err = p.funDeclaration("function")
	} else {
		stmt, err = p.statement()
	}
	return stmt
}

func (p *Parser) classDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expected class name.")
	if err != nil {
		return nil, err
	}

	var superclass *ast.Variable
	if p.match(token.LESS) {
		_, err = p.consume(token.IDENTIFIER, "Expected superclass name.")
		if err != nil {
			return nil, err
		}
		superclass = &ast.Variable{Name: p.previous()}
	}
	_, err = p.consume(token.LEFTBRACE, "Expected '{' before class body.")
	if err != nil {
		return nil, err
	}

	methods := make([]*ast.Function, 0)
	classmethods := make([]*ast.Function, 0)
	for !p.check(token.RIGHTBRACE) && !p.isAtEnd() {
		fun, err2 := p.funDeclaration("method")
		if err2 != nil {
			return nil, err2
		}
		if !fun.IsClassMethod {
			methods = append(methods, fun)
		} else {
			classmethods = append(classmethods, fun)
		}
	}

	_, err = p.consume(token.RIGHTBRACE, "Expected '}' after class body.")
	if err != nil {
		return nil, err
	}

	return &ast.Class{Name: name, Methods: methods, ClassMethods: classmethods, SuperClass: superclass}, nil
}

func (p *Parser) methodArguments(kind string) ([]token.Token, error) {
	_, err := p.consume(token.LEFTPAREN, "Expected '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	parameters := make([]token.Token, 0)
	if !p.check(token.RIGHTPAREN) {
		for {
			if len(parameters) >= 8 {
				return nil, parseerror.MakeError(p.peek(), "Cannot have more than 8 parameters.")
			}

			param, err2 := p.consume(token.IDENTIFIER, "Expected parameter name.")
			if err2 != nil {
				return nil, err2
			}

			parameters = append(parameters, param)

			if !p.match(token.COMMA) {
				break
			}
		}
	}
	_, err = p.consume(token.RIGHTPAREN, "Expected ')' after parameters.")
	return parameters, err
}

func (p *Parser) funDeclaration(kind string) (*ast.Function, error) {
	oldInLoop := p.inloop
	defer p.resetLoop(oldInLoop)
	p.inloop = false

	isClassMethod := false
	if p.match(token.CLASS) {
		isClassMethod = true
	}

	name, err := p.consume(token.IDENTIFIER, "Expected "+kind+" name.")
	if err != nil {
		return nil, err
	}

	var parameters []token.Token
	if p.check(token.LEFTPAREN) {
		parameters, err = p.methodArguments(kind)
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.LEFTBRACE, "Expected '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.Function{Name: name, Params: parameters, Body: body, EnvIndex: -1, IsClassMethod: isClassMethod}, nil
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expected variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expected ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &ast.Var{Name: name, Initializer: initializer, EnvIndex: -1}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.IF) {
		return p.ifStatement()
	} else if p.match(token.WHILE) {
		return p.whileStatement()
	} else if p.match(token.FOR) {
		return p.forStatement()
	} else if p.match(token.PRINT) {
		return p.printStatement()
	} else if p.match(token.RETURN) {
		return p.returnStatement()
	} else if p.match(token.BREAK) {
		return p.breakStatement()
	} else if p.match(token.CONTINUE) {
		return p.continueStatement()
	} else if p.match(token.LEFTBRACE) {
		statements, err := p.block()
		if err == nil {
			return &ast.Block{Statements: statements}, nil
		}
		return nil, err
	}
	return p.expressionStatement()
}

func (p *Parser) breakStatement() (ast.Stmt, error) {
	if !p.inloop {
		return nil, parseerror.MakeError(p.previous(), "Stray break detected.")
	}
	tok := p.previous()
	_, err := p.consume(token.SEMICOLON, "Expected ';' after break")
	if err != nil {
		return nil, err
	}
	return &ast.Break{Token: tok}, nil
}

func (p *Parser) continueStatement() (ast.Stmt, error) {
	if !p.inloop {
		return nil, parseerror.MakeError(p.previous(), "Stray continue detected.")
	}
	tok := p.previous()
	_, err := p.consume(token.SEMICOLON, "Expected ';' after continue")
	if err != nil {
		return nil, err
	}
	return &ast.Continue{Token: tok}, nil
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	oldInLoop := p.inloop
	defer p.resetLoop(oldInLoop)
	p.inloop = true
	_, err := p.consume(token.LEFTPAREN, "Expected '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	// first clause (initializer)
	var initializer ast.Stmt
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}
	// condition
	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}
	// increment
	var increment ast.Expr
	if !p.check(token.RIGHTPAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.RIGHTPAREN, "Expected ')' after for clauses.")
	if err != nil {
		return nil, err
	}
	// for-loop body
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &ast.For{Initializer: initializer, Condition: condition, Increment: increment, Statement: body}, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	oldInLoop := p.inloop
	defer p.resetLoop(oldInLoop)
	p.inloop = true
	_, err := p.consume(token.LEFTPAREN, "Expected '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHTPAREN, "Expected ')' after condition.")
	if err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return &ast.While{Condition: condition, Statement: body}, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	if _, err := p.consume(token.LEFTPAREN, "Expected '(' after 'if'."); err != nil {
		return nil, err
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.RIGHTPAREN, "Expected ')' after 'if' condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	if p.match(token.ELSE) {
		elseBranch, err := p.statement()
		if err != nil {
			return nil, err
		}
		return &ast.If{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
	}
	return &ast.If{Condition: condition, ThenBranch: thenBranch}, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := make([]ast.Stmt, 0)
	for !p.check(token.RIGHTBRACE) && !p.isAtEnd() {
		stmt := p.declaration()
		if stmt == nil {
			return nil, nil // FIXME: should I propagate the declaration error
		}
		statements = append(statements, stmt)
	}
	p.consume(token.RIGHTBRACE, "Expected '}' after block.")
	return statements, nil
}

func (p *Parser) returnStatement() (ast.Stmt, error) {
	keyword := p.previous()

	var value ast.Expr
	var err error
	if !p.check(token.SEMICOLON) {
		value, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expected ';' after return value.")
	if err != nil {
		return nil, err
	}
	return &ast.Return{Keyword: keyword, Value: value}, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expected ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ast.Print{Expression: expr}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expected ';' after value.")
	if err != nil {
		return nil, err
	}
	return &ast.Expression{Expression: expr}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.comma()
}

func (p *Parser) comma() (ast.Expr, error) {
	expr, err := p.assignment()
	if err != nil {
		return nil, err
	}

	for p.match(",") {
		operator := p.previous()
		right, err := p.assignment()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*ast.Variable); ok {
			return &ast.Assign{Name: variable.Name, Value: value, EnvIndex: -1, EnvDepth: -1}, nil
		} else if get, ok := expr.(*ast.Get); ok {
			return &ast.Set{Object: get.Expression, Name: get.Name, Value: value}, nil
		}
		return nil, parseerror.MakeError(equals, "Invalid assignment target.")
	}
	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}
	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.ternary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) ternary() (ast.Expr, error) {
	cond, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match("?") {
		qmark := p.previous()
		thenClause, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err2 := p.consume(token.COLON, "Expected ':' in ternary operator."); err2 != nil {
			return nil, err2
		}
		colon := p.previous()
		elseClause, err := p.expression()
		if err != nil {
			return nil, err
		}
		return &ast.Ternary{Condition: cond, QMark: qmark, Then: thenClause, Colon: colon, Else: elseClause}, nil
	}
	return cond, nil
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANGEQUAL, token.EQUALEQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.addition()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATEREQUAL, token.LESS, token.LESSEQUAL) {
		operator := p.previous()
		right, err := p.addition()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) addition() (ast.Expr, error) {
	expr, err := p.multiplication()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.multiplication()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) multiplication() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Unary{Operator: operator, Right: right}, nil
	}

	return p.power()
}

func (p *Parser) power() (ast.Expr, error) {
	expr, err := p.call()
	if err != nil {
		return nil, err
	}

	for p.match(token.POWER) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()

	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFTPAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else if p.match(token.DOT) {
			name, err := p.consume(token.IDENTIFIER, "Expected property name after '.'")
			if err != nil {
				return nil, err
			}
			expr = &ast.Get{Expression: expr, Name: name}
		} else {
			break
		}
	}
	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	args := make([]ast.Expr, 0)
	if !p.check(token.RIGHTPAREN) {
		for {
			arg, err := p.assignment() // we don't want the comma operator here
			if err != nil {
				return nil, err
			}
			if len(args) >= 8 {
				return nil, parseerror.MakeError(p.peek(), "Cannot have more than 8 arguments.")
			}
			args = append(args, arg)
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(token.RIGHTPAREN, "Expected ')' after arguments.")
	if err != nil {
		return nil, err
	}
	return &ast.Call{Callee: callee, Paren: paren, Arguments: args}, nil
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}, nil
	} else if p.match(token.TRUE) {
		return &ast.Literal{Value: true}, nil
	} else if p.match(token.NIL) {
		return &ast.Literal{Value: nil}, nil
	} else if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}, nil
	} else if p.match(token.SUPER) {
		keyword := p.previous()
		_, err := p.consume(token.DOT, "Expected '.' after 'super'.")
		if err != nil {
			return nil, err
		}
		method, err := p.consume(token.IDENTIFIER, "Expected superclass method name.")
		if err != nil {
			return nil, err
		}
		return &ast.Super{Keyword: keyword, Method: method}, nil
	} else if p.match(token.THIS) {
		return &ast.This{Keyword: p.previous(), EnvIndex: -1, EnvDepth: -1}, nil
	} else if p.match(token.LEFTPAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHTPAREN, "Expected ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Grouping{Expression: expr}, nil
	} else if p.match(token.IDENTIFIER) {
		return &ast.Variable{Name: p.previous(), EnvIndex: -1, EnvDepth: -1}, nil
	}
	return nil, parseerror.MakeError(p.peek(), "Expected expression")
}

func (p *Parser) consume(tp token.Type, message string) (token.Token, error) {
	if p.check(tp) {
		return p.advance(), nil
	}
	return p.previous(), parseerror.MakeError(p.peek(), message)
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) match(types ...token.Type) bool {
	for _, tp := range types {
		if p.check(tp) {
			p.advance()
			return true
		}
	}
	return false
}

// check checks if the next token is of the given type
func (p *Parser) check(tp token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tp
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// peek returns the next token
func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

// peek returns the current token
func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) resetLoop(val bool) {
	p.inloop = val
}
