package main

import (
	"errors"
	"fmt"
)

type Parser struct {
	context *LoxContext
	tokens  []*Token
	current int
}

func MakeParser(context *LoxContext, tokens []*Token) *Parser {
	return &Parser{
		context: context,
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) match(tokenTypes ...TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().tokenType == tokenType
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().tokenType == EOF
}

func (p *Parser) peek() *Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return p.tokens[p.current-1]
}

// equality ;
func (p *Parser) expression() Expr {
	return p.assignment()
}

// IDENTIFIER "=" assignment | logic_or ;
func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(EQUAL) {
		equals := p.previous()
		value := p.assignment()

		switch val := expr.(type) {
		case *VariableExpr:
			name := val.name
			return MakeAssignExpr(name, value)
		}

		p.error(equals, "Invalid assignment target.")
	}

	return expr
}

// logic_and ( "or" logic_and )* ;
func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(OR) {
		operator := p.previous()
		right := p.and()
		expr = MakeLogicalExpr(expr, operator, right)
	}

	return expr
}

// equality ( "and" equality )* ;
func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(AND) {
		operator := p.previous()
		right := p.equality()
		expr = MakeLogicalExpr(expr, operator, right)
	}

	return expr
}

// comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = MakeBinaryExpr(expr, operator, right)
	}
	return expr
}

// term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = MakeBinaryExpr(expr, operator, right)
	}
	return expr
}

// factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = MakeBinaryExpr(expr, operator, right)
	}
	return expr
}

// unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()
		expr = MakeBinaryExpr(expr, operator, right)
	}
	return expr
}

// ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()
		return MakeUnaryExpr(operator, right)
	}

	if p.match(IDENTIFIER) {
		return MakeVariableExpr(p.previous())
	}

	return p.primary()
}

// NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() Expr {
	if p.match(TRUE) {
		return MakeLiteralExpr(true)
	}
	if p.match(FALSE) {
		return MakeLiteralExpr(false)
	}
	if p.match(NIL) {
		return MakeLiteralExpr(nil)
	}

	if p.match(NUMBER, STRING) {
		return MakeLiteralExpr(p.previous().literal)
	}

	if p.match(LEFT_PAREN) {
		expr := p.expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return MakeGroupingExpr(expr)
	}

	panic(p.error(p.peek(), "Expected expression."))
}

func (p *Parser) consume(tokenType TokenType, message string) *Token {
	if p.check(tokenType) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

func (p *Parser) error(token *Token, message string) error {
	p.context.tokenError(token, message)
	return ParseError{token: token, message: message}
}

type ParseError struct {
	token   *Token
	message string
}

func (e ParseError) Error() string {
	if e.token == nil {
		return "Parse error: " + e.message
	}
	if e.token.tokenType == EOF {
		return "Parse error at end: " + e.message
	} else {
		return fmt.Sprintf("Parse error at '%s': %s", e.token.lexme, e.message)
	}
}

func (p *Parser) parse() (result []Stmt, err error) {
	if tryCatch(func() {
		result = p.program()
	}) != nil {
		err = errors.New("parse error")
	}
	return
}

func (p *Parser) program() []Stmt {
	statements := make([]Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() (result Stmt) {
	if tryCatch(func() {
		if p.match(VAR) {
			result = p.varDeclaration()
		} else {
			result = p.statement()
		}
	}) == catch {
		p.synchronize()
		result = nil
	}
	return
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}

	p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	return MakeVarStmt(name, initializer)
}

func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	if p.match(LEFT_BRACE) {
		return MakeBlockStmt(p.block())
	}

	if p.match(IF) {
		return p.ifStatement()
	}

	return p.expressionStatement()
}

// "if" "(" expression ")" statement ( "else" statement )? ;
func (p *Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(ELSE) {
		elseBranch = p.statement()
	}

	return MakeIfStmt(condition, thenBranch, elseBranch)
}

func (p *Parser) block() []Stmt {
	statements := make([]Stmt, 0)

	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return MakePrintStmt(value)
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return MakeExpressionStmt(expr)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().tokenType == SEMICOLON {
			return
		}

		switch p.peek().tokenType {
		case CLASS:
		case FUN:
		case VAR:
		case FOR:
		case IF:
		case WHILE:
		case PRINT:
		case RETURN:
			return
		}

		p.advance()
	}
}
