package ast

import "../token"

// Node :
type Node interface {
	TokenLiteral() string
}

// Statement :
type Statement interface {
	Node
	statementNode()
}

// Expression :
type Expression interface {
	Node
	expressionNode()
}

// Program :
type Program struct {
	Statements []Statement
}

// Identifier :
type Identifier struct {
	Token token.Token
	Value string
}

// LetStatement :
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

// TokenLiteral :
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// statementNode :
func (ls *LetStatement) statementNode() {}

// TokenLiteral :
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// TokenLiteral :
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
