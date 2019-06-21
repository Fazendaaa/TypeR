package parser

import (
	"fmt"
	"strconv"

	"../ast"
	"../lexer"
	"../token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

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

	if !p.expectPeek(token.IDENTIFICATION) {
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

	return statement
}

// parseExpression :
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParserFunction[p.currentToken.Type]

	if nil == prefix {
		return nil
	}

	leftExpression := prefix()

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

// Errors :
func (p *Parser) Errors() []string {
	return p.errors
}

// peekErrors :
func (p *Parser) peekErrors(t token.TokenType) {
	message := fmt.Sprintf("Expected next token to be %s, got '%s' instead", t, p.peekToken.Type)
	p.errors = append(p.errors, message)
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
	p.registerPrefix(token.IDENTIFICATION, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	return p
}
