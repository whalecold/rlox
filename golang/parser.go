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

func (p *Parser) and() Expr {
	expr := p.equality()
	for p.match(AND) {
		operator := p.previous()
		// TODO why not and?
		right := p.equality()
		expr = &Logical{expr, operator, right}
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()
	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = &Logical{expr, operator, right}
	}
	return expr
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()
		if ident, ok := expr.(*Variable); ok {
			name := ident.name
			return &Assign{name, value}
		} else if get, ok := expr.(*Get); ok {
			return &Set{get.object, get.name, value}
		}
		Panic(equals.line, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) expression() Expr {
	return p.assignment()
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
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else if p.match(DOT) {
			name := p.consume(IDENTIFIER, "Expect property name after '.'.")
			expr = &Get{expr, name}
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	args := []Expr{}
	if !p.check(RIGHT_PAREN) {
		for {
			if len(args) >= 255 {
				Panic(p.peek().line, "Can't have more than 255 arguments.")
			}
			args = append(args, p.expression())
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
	return &Call{callee, paren, args}
}

func (p *Parser) primary() Expr {
	if p.match(THIS) {
		return &This{p.previous()}
	}
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
	if p.match(IDENTIFIER) {
		return &Variable{p.previous()}
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

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().typ == SEMICOLON {
			return
		}
		switch p.peek().typ {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		default:
			p.advance()
		}
	}
}

func (p *Parser) ParseStmts() []Stmt {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var stmts []Stmt
	for !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	return stmts
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			p.synchronize()
			hadError = false
		}
	}()
	if p.match(CLASS) {
		return p.classDeclaration()
	}
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) classDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect class name.")
	p.consume(LEFT_BRACE, "Expect '{' before class body.")

	var methods []*Function
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		methods = append(methods, p.function("method").(*Function))
	}
	p.consume(RIGHT_BRACE, "Expect '}' after class body.")
	return &Class{name, methods}
}

func (p *Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect kind %s name.", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))

	var parameters []*Token
	if !p.check(RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				Panic(p.peek().line, "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(IDENTIFIER, "Expect parameter name."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return &Function{name, parameters, body}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")
	var initializer Expr
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return &Var{name, initializer}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	if !p.check(SEMICOLON) {
		value = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return &Return{keyword, value}
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(LEFT_BRACE) {
		return &Block{p.block()}
	}
	return p.exprStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.exprStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(RIGHT_PAREN) {
		increment = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &Block{[]Stmt{body, &Expression{increment}}}
	}
	if condition == nil {
		condition = &Literal{true}
	}
	body = &While{condition, body}
	if initializer != nil {
		body = &Block{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	body := p.statement()
	return &While{condition, body}
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return &If{condition, thenBranch, elseBranch}
}

func (p *Parser) block() []Stmt {
	var stmts []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmts = append(stmts, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return stmts
}

func (p *Parser) printStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &Print{expr}
}

func (p *Parser) exprStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &Expression{expr}
}
