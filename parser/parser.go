package parser

import (
	"golox/ast"
	"golox/parseerror"
	"golox/token"
)

/*
program    -> declaration* EOF ;
declaration -> varDecl
			| stmt
varDecl    -> "var" IDENTIFIER ( "=" expression )? ";" ;
stmt       -> exprStmt
			| printStmt
exprStmt   -> expression ";" ;
printStmt  -> "print" expression ";" ;
expression -> comma ;
comma      -> ternary ( "," ternary ) * ;
ternary    -> equality "?"  expression ":" expression ;
equality   -> comparison ( ( "!=" | "==") comparison )* ;
comparison -> addition ( ( ">" | ">=" | "<" | "<=") addition )*;
addition   -> multiplication ( ( "+" | "-" ) multiplication )*;
multiplication -> unary ( ( "/" | "*" ) unary )*;
unary      -> ( "!" | "-" ) unary
			| power ;
power      -> primary ( "**" unary ) *
primary    -> NUMBER | STRING | "false" | "true" | "nil"
			| "(" expression ")"
			| IDENTIFIER ;
*/

// Parser will transform an array of tokens to an AST.
// Use parser.New to create a new Parser. Do not create a Parser directly
type Parser struct {
	tokens  []token.Token
	current int
}

// New creates a new parser
func New(tokens []token.Token) Parser {
	return Parser{tokens, 0}
}

// Parse is the driver function that begins parsing
func (p *Parser) Parse() []ast.Stmt {
	statements := make([]ast.Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(token.VAR) {
		stmt, err := p.varDeclaration()
		if err != nil {
			parseerror.LogError(err)
			return nil
		}
		return stmt
	}
	stmt, err := p.statement()
	if err != nil {
		parseerror.LogError(err)
		return nil
	}
	return stmt
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
	return &ast.Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
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
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	for p.match(",") {
		operator := p.previous()
		right, err := p.ternary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
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
		p.match(":") // TODO: error handling
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
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for p.match(token.POWER) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return &ast.Binary{Left: expr, Operator: operator, Right: right}, nil
	}
	return expr, nil
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
		return ast.Variable{Name: p.previous()}, nil
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
