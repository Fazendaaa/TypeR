package ast

import (
	"bytes"
	"fmt"
	"strings"

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

// Header :
type Header interface {
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

// ConstStatement :
type ConstStatement struct {
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

// IntegerLiteral :
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

// PrefixExpression :
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

// InfixExpression :
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

// Boolean :
type Boolean struct {
	Token token.Token
	Value bool
}

// BlockStatement :
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

// ConditionalExpression :
type ConditionalExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// IdentifierTypes :
type IdentifierTypes struct {
	Token  token.Token
	Header Header
}

// FunctionLiteral :
type FunctionLiteral struct {
	Token token.Token
	// IdentifierTypes contains the parameters and return values types
	Contract   []*IdentifierTypes
	Parameters []*Identifier
	Body       *BlockStatement
	Name       string
}

// CallExpression :
type CallExpression struct {
	Token      token.Token
	Function   Expression
	Parameters []Expression
}

// StringLiteral :
type StringLiteral struct {
	Token token.Token
	Value string
}

// ArrayLiteral :
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

// IndexExpression :
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

// PointFreeExpression :
type PointFreeExpression struct {
	Token        token.Token
	ToCompose    []*Identifier
	SeedFunction Expression
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
func (cs *ConstStatement) String() string {
	var out bytes.Buffer

	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Name.String())
	out.WriteString(" <- ")

	if nil != cs.Value {
		out.WriteString(cs.Value.String())
	}

	return out.String()
}

// statementNode :
func (cs *ConstStatement) statementNode() {}

// TokenLiteral :
func (cs *ConstStatement) TokenLiteral() string {
	return cs.Token.Literal
}

// String :
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if nil != rs.ReturnValue {
		out.WriteString(rs.ReturnValue.String())
	}

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

// expressionNode :
func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral :
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

// String :
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

// expressionNode :
func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral :
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

// String :
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

// expressionNode :
func (ie *InfixExpression) expressionNode() {}

// TokenLiteral :
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String :
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

// expressionNode :
func (b *Boolean) expressionNode() {}

// TokenLiteral :
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

// String :
func (b *Boolean) String() string {
	return b.Token.Literal
}

// expressionNode :
func (bs *BlockStatement) expressionNode() {}

// TokenLiteral :
func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

// String :
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}

	return out.String()
}

// expressionNode :
func (ce *ConditionalExpression) expressionNode() {}

// TokenLiteral :
func (ce *ConditionalExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String :
func (ce *ConditionalExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ce.Condition.String())
	out.WriteString(" ")
	out.WriteString(ce.Consequence.String())

	if nil != ce.Alternative {
		out.WriteString("else ")
		out.WriteString(ce.Alternative.String())
	}

	return out.String()
}

// expressionNode :
func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral :
func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

// String :
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, p := range fl.Parameters {
		parameters = append(parameters, p.String())
	}

	out.WriteString(fl.TokenLiteral())

	if "" != fl.Name {
		out.WriteString(fmt.Sprintf("<%s>", fl.Name))
	}

	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// expressionNode :
func (ce *CallExpression) expressionNode() {}

// TokenLiteral :
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

// String :
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	parameters := []string{}

	for _, parameter := range ce.Parameters {
		parameters = append(parameters, parameter.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(parameters, ", "))
	out.WriteString(")")

	return out.String()
}

// expressionNode :
func (sl *StringLiteral) expressionNode() {}

// TokenLiteral :
func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.Literal
}

// String :
func (sl *StringLiteral) String() string {
	return sl.Token.Literal
}

// expressionNode :
func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral :
func (al *ArrayLiteral) TokenLiteral() string {
	return al.Token.Literal
}

// String :
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, element := range al.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// expressionNode :
func (ie *IndexExpression) expressionNode() {}

// TokenLiteral :
func (ie *IndexExpression) TokenLiteral() string {
	return ie.Token.Literal
}

// String :
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// expressionNode :
func (pf *PointFreeExpression) expressionNode() {}

// TokenLiteral :
func (pf *PointFreeExpression) TokenLiteral() string {
	return pf.Token.Literal
}

// String :
func (pf *PointFreeExpression) String() string {
	var out bytes.Buffer

	functions := []string{}

	for _, function := range pf.ToCompose {
		functions = append(functions, function.String())
	}

	out.WriteString(strings.Join(functions, " . "))
	out.WriteString(" . ")
	out.WriteString(pf.SeedFunction.String())

	return out.String()
}
