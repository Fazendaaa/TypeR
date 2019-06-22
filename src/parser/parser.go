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

// currentPrecedence :
func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return LOWEST
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

// parseBoolean :
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}

// parseGroupedExpression :
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return expression
}

// parseBlockStatement :
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currentToken,
	}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currentTokenIs(token.RIGHT_BRACE) && !p.currentTokenIs(token.EOF) {
		statement := p.parseStatement()

		if nil != statement {
			block.Statements = append(block.Statements, statement)
		}

		p.nextToken()
	}

	return block
}

// parseConditionalExpression :
func (p *Parser) parseConditionalExpression() ast.Expression {
	expresion := &ast.ConditionalExpression{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.LEFT_PARENTHESIS) {
		return nil
	}

	p.nextToken()

	expresion.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	// This might not be true in one statement blocks
	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	expresion.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LEFT_BRACE) {
			return nil
		}

		expresion.Alternative = p.parseBlockStatement()
	}

	return expresion
}

// parseFunctionParameters :
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken()

		return identifiers
	}

	p.nextToken()

	identifier := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
	identifiers = append(identifiers, identifier)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		identifier := &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return identifiers
}

// parseFunctionLiteral :
func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.LEFT_PARENTHESIS) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LEFT_BRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
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
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LEFT_PARENTHESIS, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseConditionalExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

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
