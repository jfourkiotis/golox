package parser

import (
	"fmt"
	"golox/ast"
	"golox/token"
)

/*
expression -> equality ;
equality   -> comparison ( ( "!=" | "==") comparison )* ;
comparison -> addition ( ( ">" | ">=" | "<" | "<=") addition )*;
addition   -> multiplication ( ( "+" | "-" ) multiplication )*;
multiplication -> unary ( ( "/" | "*" ) unary )*;
unary      -> ( "!" | "-" ) unary
			| power ;
power      -> primary ( "**" unary ) *
primary    -> NUMBER | STRING | "false" | "true" | "nil"
            | "(" expression ")" ;
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
func (p *Parser) Parse() ast.Expr {
	expr := p.expression()
	return expr
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANGEQUAL, token.EQUALEQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.addition()

	for p.match(token.GREATER, token.GREATEREQUAL, token.LESS, token.LESSEQUAL) {
		operator := p.previous()
		right := p.addition()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) addition() ast.Expr {
	expr := p.multiplication()

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right := p.multiplication()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) multiplication() ast.Expr {
	expr := p.unary()

	for p.match(token.STAR, token.SLASH) {
		operator := p.previous()
		right := p.unary()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &ast.Unary{Operator: operator, Right: right}
	}

	return p.power()
}

func (p *Parser) power() ast.Expr {
	expr := p.primary()

	for p.match(token.POWER) {
		operator := p.previous()
		right := p.unary()
		return &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return &ast.Literal{Value: false}
	} else if p.match(token.TRUE) {
		return &ast.Literal{Value: true}
	} else if p.match(token.NIL) {
		return &ast.Literal{Value: nil}
	} else if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{Value: p.previous().Literal}
	} else if p.match(token.LEFTPAREN) {
		expr := p.expression()
		p.consume(token.RIGHTPAREN, "Expected ')' after expression.")
		return &ast.Grouping{Expression: expr}
	}
	panic("Expected expression.")
}

func (p *Parser) consume(tp token.Type, message string) token.Token {
	if p.check(tp) {
		return p.advance()
	}
	panic(fmt.Sprintf("%v: %v", p.peek(), message))
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
