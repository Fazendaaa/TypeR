package parser

import (
	"fmt"
	"strconv"

	"../ast"
	"../lexer"
	"../token"
)

const (
	_           int = iota
	LOWEST          // Starting condition
	EQUALS          // ==
	LESSGREATER     // > or <
	SUM             // +
	PRODUCT         // *
	PREFIX          // -X or !X
	CALL            // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.DOUBLE_EQUAL: EQUALS,
	token.DIFFERENT:    EQUALS,
	token.LESS_THAN:    LESSGREATER,
	token.GREATER_THAN: LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
}

// Parser :
type Parser struct {
	l      *lexer.Lexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParserFunction map[token.TokenType]prefixParserFunction
	infixParserFunction  map[token.TokenType]infixParserFunction
}

type (
	prefixParserFunction func() ast.Expression
	infixParserFunction  func(ast.Expression) ast.Expression
)

// registerPrefix :
func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParserFunction) {
	p.prefixParserFunction[tt] = fn
}

// registerInfix :
func (p *Parser) registerInfix(tt token.TokenType, fn infixParserFunction) {
	p.infixParserFunction[tt] = fn
}

// nextToken :
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// currentTokenIs :
func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

// peekTokenIs :
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// peekPrecedence :
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// expectPeek :
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()

		return true
	}

	p.peekErrors(t)

	return false
}

// parseLetStatement :
func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// parseReturnStatement :
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{
		Token: p.currentToken,
	}

	p.nextToken()

	statement.ReturnValue = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// parseIdentifier :
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

// parseStatement :
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseIntegerLiteral :
func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: p.currentToken,
	}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if nil != err {
		message := fmt.Sprintf("could not parse '%q' as integer", p.currentToken.Literal)
		p.errors = append(p.errors, message)

		return nil
	}

	literal.Value = value

	return literal
}

// noPrefixParserFnError :
func (p *Parser) noPrefixParserFnError(t token.TokenType) {
	message := fmt.Sprintf("no prefix parse function for '%s' was found", t)
	p.errors = append(p.errors, message)
}

// parseExpression :
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParserFunction[p.currentToken.Type]

	if nil == prefix {
		p.noPrefixParserFnError(p.currentToken.Type)

		return nil
	}

	leftExpression := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParserFunction[p.peekToken.Type]

		if nil == infix {
			return leftExpression
		}

		p.nextToken()

		leftExpression = infix(leftExpression)
	}

	return leftExpression
}

// parseExpressionStatement :
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{
		Token: p.currentToken,
	}

	statement.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// parsePrefixExpression :
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression :
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expresion := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()

	p.nextToken()

	expresion.Right = p.parseExpression(precedence)

	return expresion
}

// currentPrecedence :
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
}

// peekErrors :
func (p *Parser) peekErrors(t token.TokenType) {
	message := fmt.Sprintf("Expected next token to be %s, got '%s' instead", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
}

// Errors :
func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram :
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()

		if nil != statement {
			program.Statements = append(program.Statements, statement)
		}

		p.nextToken()
	}

	return program
}

// InitializeParser :
func InitializeParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Sets the current and peek tokens
	p.nextToken()
	p.nextToken()

	p.prefixParserFunction = make(map[token.TokenType]prefixParserFunction)
	p.infixParserFunction = make(map[token.TokenType]infixParserFunction)

	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.DOUBLE_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.DIFFERENT, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)

	return p
}
