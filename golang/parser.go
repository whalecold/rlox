package main

import "fmt"

type Parser struct {
	tokens  []*Token
	current int
}

func (p *Parser) Parse() (expr Expr) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			expr = nil
		}
	}()
	return p.expression()
}

func (p *Parser) match(types ...TokenType) bool {
	for _, typ := range types {
		if p.check(typ) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(typ TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().typ == typ
}

func (p *Parser) isAtEnd() bool {
	return p.peek().typ == EOF
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = &Binary{expr, operator, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return &Unary{operator, right}
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return &Literal{false}
	}
	if p.match(TRUE) {
		return &Literal{true}
	}
	if p.match(NIL) {
		return &Literal{nil}
	}
	if p.match(NUMBER, STRING) {
		return &Literal{p.previous().literal}
	}
	if p.match(LEFT_PAREN) {
		expression := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &Grouping{expression}
	}
	Panic(p.peek().line, "Expect expression.")
	return nil
}

func (p *Parser) consume(typ TokenType, message string) *Token {
	if p.check(typ) {
		return p.advance()
	}
	Panic(p.peek().line, message)
	return nil
}
