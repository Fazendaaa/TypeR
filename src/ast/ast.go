package ast

import (
	"bytes"

	"../token"
)

// Node :
type Node interface {
	TokenLiteral() string
	String() string
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

// ReturnStatement :
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// ExpressionStatement :
type ExpressionStatement struct {
	// The first token of the expression
	Token      token.Token
	Expression Expression
}

// statementNode :
func (i *Identifier) statementNode() {}

// statementNode :
func (i *Identifier) expressionNode() {}

// String :
func (i *Identifier) String() string {
	return i.Value
}

// TokenLiteral :
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

// String :
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// TokenLiteral :
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// String :
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" <- ")

	if nil != ls.Value {
		out.WriteString(ls.Value.String())
	}

	// Optional
	// out.WriteString(";")

	return out.String()
}

// statementNode :
func (ls *LetStatement) statementNode() {}

// TokenLiteral :
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

// String :
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if nil != rs.ReturnValue {
		out.WriteString(rs.ReturnValue.String())
	}

	// Optional
	// out.WriteString(";")

	return out.String()
}

// statementNode :
func (rs *ReturnStatement) statementNode() {}

// TokenLiteral :
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

// String :
func (es *ExpressionStatement) String() string {
	if nil != es.Expression {
		return es.Expression.String()
	}

	return ""
}

// statementNode :
func (es *ExpressionStatement) statementNode() {}

// TokenLiteral :
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
