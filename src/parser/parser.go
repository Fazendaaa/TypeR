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
	INDEX           // myArray[?]
)

var precedences = map[token.TokenType]int{
	token.DOUBLE_EQUAL:     EQUALS,
	token.DIFFERENT:        EQUALS,
	token.LESS_THAN:        LESSGREATER,
	token.GREATER_THAN:     LESSGREATER,
	token.PLUS:             SUM,
	token.MINUS:            SUM,
	token.SLASH:            PRODUCT,
	token.ASTERISK:         PRODUCT,
	token.LEFT_PARENTHESIS: CALL,
	token.LEFT_BRACKET:     INDEX,
}

// Parser :
type Parser struct {
	l      *lexer.Lexer
	errors []string

	previousToken token.Token
	currentToken  token.Token
	peekToken     token.Token

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
	p.previousToken = p.currentToken
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// backToken :
func (p *Parser) backToken() {
	p.peekToken = p.currentToken
	p.currentToken = p.previousToken
	p.previousToken = p.l.PreviousToken()
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

// expectCurrent :
func (p *Parser) expectCurrent(t token.TokenType) bool {
	if p.currentTokenIs(t) {
		p.nextToken()

		return true
	}

	p.currentErrors(t)

	return false
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

	if functionLiteral, ok := statement.Value.(*ast.FunctionLiteral); ok {
		functionLiteral.Name = statement.Name.Value
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

// parseConstStatement :
func (p *Parser) parseConstStatement() *ast.ConstStatement {
	constant := token.Token{
		Type:    token.CONST,
		Literal: "CONST",
	}
	statement := &ast.ConstStatement{
		Token: constant,
	}
	statement.Name = &ast.Identifier{
		Token: constant,
		Value: p.currentToken.Literal,
	}

	if !p.peekTokenIs(token.ASSIGN) {
		return nil
	}

	p.nextToken()
	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if functionLiteral, ok := statement.Value.(*ast.FunctionLiteral); ok {
		functionLiteral.Name = statement.Name.Value
	}

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

// parsePointFreeLiteral :
func (p *Parser) parsePointFreeLiteral() ast.Expression {
	pointFree := &ast.PointFreeExpression{
		Token: token.Token{
			Type:    token.POINT,
			Literal: ".",
		},
	}
	function := &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}

	toCompose := []*ast.Identifier{}
	toCompose = append(toCompose, function)

	for p.currentTokenIs(token.POINT) && !p.currentTokenIs(token.EOF) {
		p.nextToken()

		function = &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		}
		toCompose = append(toCompose, function)

		p.nextToken()
	}

	pointFree.ToCompose = toCompose

	if p.currentTokenIs(token.LEFT_PARENTHESIS) {
		pointFree.Parameters = p.parseExpressionList(token.RIGHT_PARENTHESIS)
	}

	// println(pointFree.String())

	return pointFree
}

// parseIdentifier :
func (p *Parser) parseIdentifier() ast.Expression {
	if p.peekTokenIs(token.POINT) {
		return p.parsePointFreeLiteral()
	}

	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

// parseStatement :
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.IDENTIFIER:
		constant := p.parseConstStatement()

		if nil != constant {
			return constant
		}

		return p.parseExpressionStatement()
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
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()

	p.nextToken()

	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseBoolean :
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}

// parseAnonymousFunctionLiteral :
func (p *Parser) parseAnonymousFunctionLiteral() ast.Expression {
	function := token.Token{
		Type:    token.FUNCTION,
		Literal: "function",
	}
	literal := &ast.FunctionLiteral{
		Token: function,
	}

	p.backToken()
	p.previousToken = function

	literal.Parameters = p.parseFunctionParameters()
	literal.Body = p.parseBlockStatement()

	return literal
}

// parseGroupedExpression :
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	if p.peekTokenIs(token.COMMA) ||
		p.currentTokenIs(token.RIGHT_PARENTHESIS) ||
		p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		return p.parseAnonymousFunctionLiteral()
	}

	expression := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return expression
}

// parseMultipleLinesBlockStatement :
func (p *Parser) parseMultipleLinesBlockStatement() *ast.BlockStatement {
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

// parseOneLinersBlockStatement :
func (p *Parser) parseOneLinersBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: token.Token{
			Type:    token.LEFT_BRACE,
			Literal: "{",
		},
	}
	block.Statements = []ast.Statement{}

	statement := p.parseStatement()

	block.Statements = append(block.Statements, statement)

	return block
}

// parseBlockStatement :
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	p.nextToken()

	if p.currentTokenIs(token.LEFT_BRACE) {
		return p.parseMultipleLinesBlockStatement()
	}

	return p.parseOneLinersBlockStatement()
}

// parseConditionalExpression :
func (p *Parser) parseConditionalExpression() ast.Expression {
	expression := &ast.ConditionalExpression{
		Token: p.currentToken,
	}

	if !p.expectPeek(token.LEFT_PARENTHESIS) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
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
	literal.Body = p.parseBlockStatement()

	return literal
}

// parseCallArguments :
func (p *Parser) parseCallArguments() []ast.Expression {
	arguments := []ast.Expression{}

	if p.peekTokenIs(token.RIGHT_PARENTHESIS) {
		p.nextToken()

		return arguments
	}

	p.nextToken()

	arguments = append(arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		arguments = append(arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RIGHT_PARENTHESIS) {
		return nil
	}

	return arguments
}

// parseExpressionList :
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()

		return list
	}

	p.nextToken()

	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseCallExpression :
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{
		Token:    p.currentToken,
		Function: function,
	}
	expression.Parameters = p.parseExpressionList(token.RIGHT_PARENTHESIS)

	return expression
}

// parseStringLiteral :
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

// parseArrayLiteral :
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{
		Token: p.currentToken,
	}

	array.Elements = p.parseExpressionList(token.RIGHT_BRACKET)

	return array
}

// parseIndexExpression :
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expression := &ast.IndexExpression{
		Token: p.currentToken,
		Left:  left,
	}

	p.nextToken()
	expression.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RIGHT_BRACKET) {
		return nil
	}

	return expression
}

// peekErrors :
func (p *Parser) peekErrors(t token.TokenType) {
	message := fmt.Sprintf("Expected next token to be %s, got '%s' instead", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
}

// currentErrors :
func (p *Parser) currentErrors(t token.TokenType) {
	message := fmt.Sprintf("Expected current token to be %s, got '%s' instead", t, p.currentToken.Type)
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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LEFT_BRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.POINT, p.parsePointFreeLiteral)

	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.DOUBLE_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.DIFFERENT, p.parseInfixExpression)
	p.registerInfix(token.LESS_THAN, p.parseInfixExpression)
	p.registerInfix(token.GREATER_THAN, p.parseInfixExpression)
	p.registerInfix(token.LEFT_PARENTHESIS, p.parseCallExpression)
	p.registerInfix(token.LEFT_BRACKET, p.parseIndexExpression)

	return p
}
