package parser

import (
	"fmt"
	"strings"
	"testing"

	"../ast"
	"../lexer"
)

// testLetStatements :
func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if "let" != s.TokenLiteral() {
		t.Errorf("s.TokenLiteral not 'let', got=%q", s.TokenLiteral())

		return false
	}

	letStatement, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement, got=%T", s)

		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s', got=%s", name, letStatement.Name.Value)

		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("lestStatement.Name.TokenLiteral() not '%s', got=%s", name, letStatement.Name.TokenLiteral())

		return false
	}

	return true
}

// testIdentifier :
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	identifier, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("exp not *ast.Identifier, got=%T", exp)

		return false
	}

	if identifier.Value != value {
		t.Errorf("identifier.Value not %s, got=%s", value, identifier.Value)

		return false
	}

	if identifier.TokenLiteral() != value {
		t.Errorf("identifier.TokenLiteral() not '%s', got=%T", value, identifier.TokenLiteral())

		return false
	}

	return true
}

// testIntegerLiteral :
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integer, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)

		return false
	}

	if integer.Value != value {
		t.Errorf("integer.Value not '%d', got=%d", value, integer.Value)

		return false
	}

	if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer.TokenLiteral() not '%d', got=%T", value, integer.TokenLiteral())

		return false
	}

	return true
}

// testBooleanLiteral :
func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	boolean, ok := expression.(*ast.Boolean)

	if !ok {
		t.Errorf("expression not *ast.Boolean, got=%T", expression)

		return false
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value not '%t', got=%t", boolean.Value, value)

		return false
	}

	// This convertion is needed because TypeR uses booleans values in an upper case
	converted := strings.ToLower(boolean.TokenLiteral())

	if converted != fmt.Sprintf("%t", value) {
		t.Errorf("boolean.TokenLiteral() not '%t', got=%s", value, converted)

		return false
	}

	return true
}

// testLiteralExpresion :
func testLiteralExpresion(t *testing.T, expression ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, expression, int64(v))
	case int64:
		return testIntegerLiteral(t, expression, v)
	case string:
		return testIdentifier(t, expression, v)
	case bool:
		return testBooleanLiteral(t, expression, v)
	}

	t.Errorf("type of expression not handled, got=%T", expression)

	return false
}

// testInfixExpression :
func testInfixExpression(t *testing.T, expression ast.Expression, left interface{}, operator string, right interface{}) bool {
	operatorExpression, ok := expression.(*ast.InfixExpression)

	if !ok {
		t.Errorf("expression is not ast.InfixExpression, got=%T(%s)", expression, expression)

		return false
	}

	if !testLiteralExpresion(t, operatorExpression.Left, left) {
		return false
	}

	if operatorExpression.Operator != operator {
		t.Errorf("exp.Operator is not '%s', got=%q", operator, operatorExpression.Operator)

		return false
	}

	if !testLiteralExpresion(t, operatorExpression.Right, right) {
		return false
	}

	return true
}

// checkParserErrors :
func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()

	if 0 == len(errors) {
		return
	}

	t.Errorf("parser has %d errors", len(errors))

	for _, messsage := range errors {
		t.Errorf("parser error: %q", messsage)
	}

	t.FailNow()
}

// TestLetStatements :
func TestLetStatements(t *testing.T) {
	input := `
let x <- 5;let y<-10
let foo <- 2345678
`

	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if nil == program {
		t.Fatalf("ParseProgram() returned nil")
	}
	if 3 != len(program.Statements) {
		t.Fatalf("program.Statements does not contains three statements, got=%d\n", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]

		if !testLetStatements(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10
return 1230987123
`
	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if 3 != len(program.Statements) {
		t.Fatalf("program.Statements does not contain three statements, got=%d", len(program.Statements))
	}

	for _, statement := range program.Statements {
		returnStatement, ok := statement.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("statement not *ast.ReturnStatement, got=%T", statement)

			continue
		}
		if "return" != returnStatement.TokenLiteral() {
			t.Errorf("returnStatement.TokenLiteral() not 'return', got=%q", returnStatement.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar"

	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if 1 != len(program.Statements) {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	identifier, ok := statement.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Expression not *ast.Identifier, got=%T", statement.Expression)
	}
	if "foobar" != identifier.Value {
		t.Errorf("identifier.Value not '%s', got=%s", "foobar", identifier.Value)
	}
	if "foobar" != identifier.TokenLiteral() {
		t.Errorf("identifier.TokenLiteral() not '%s', got=%s", "foobar", identifier.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if 1 != len(program.Statements) {
		t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not as.ExpressionStatement, got=%T", program.Statements[0])
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("expresion not *ast.IntegralLiteral, got=%T", statement.Expression)
	}

	if 5 != literal.Value {
		t.Errorf("literal.Value not '%d', got=%d", 5, literal.Value)
	}

	if "5" != literal.TokenLiteral() {
		t.Errorf("literal.TokenLiteral not '%s', got=%s", "5", literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{
			"!5;",
			"!",
			5,
		},
		{
			"-15",
			"-",
			15,
		},
		{
			"!TRUE",
			"!",
			true,
		},
		{
			"!FALSE",
			"!",
			false,
		},
	}

	for _, tt := range prefixTests {
		l := lexer.InitializeLexer(tt.input)
		p := InitializeParser(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if 1 != len(program.Statements) {
			t.Fatalf("program.Statements does not contain '%d' statements, got=%d", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		expresion, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression, got=%T", statement.Expression)
		}

		if expresion.Operator != tt.operator {
			t.Fatalf("expresion.Operator is not '%s', got=%s", tt.operator, expresion.Operator)
		}

		if !testLiteralExpresion(t, expresion.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTest := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{
			"5 + 5",
			5,
			"+",
			5,
		},
		{
			"5 - 5",
			5,
			"-",
			5,
		},
		{
			"5 * 5",
			5,
			"*",
			5,
		},
		{
			"5 / 5",
			5,
			"/",
			5,
		},
		{
			"5 > 5",
			5,
			">",
			5,
		},
		{
			"5 < 5",
			5,
			"<",
			5,
		},
		{
			"5 == 5",
			5,
			"==",
			5,
		},
		{
			"5 != 5",
			5,
			"!=",
			5,
		},
		{
			"TRUE == TRUE",
			true,
			"==",
			true,
		},
		{
			"TRUE != FALSE",
			true,
			"!=",
			false,
		},
		{
			"FALSE == FALSE",
			false,
			"==",
			false,
		},
	}

	for _, tt := range infixTest {
		l := lexer.InitializeLexer(tt.input)
		p := InitializeParser(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if 1 != len(program.Statements) {
			t.Fatalf("program.Statements does not contain %d statements, got=%d", 1, len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, statement.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a * b + c",
			"((a * b) + c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"TRUE",
			"TRUE",
		},
		{
			"FALSE",
			"FALSE",
		},
		{
			"3 > 5 == FALSE",
			"((3 > 5) == FALSE)",
		},
		{
			"3 < 5 == TRUE",
			"((3 < 5) == TRUE)",
		},
	}

	for _, tt := range tests {
		l := lexer.InitializeLexer(tt.input)
		p := InitializeParser(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{
			"TRUE",
			true,
		},
		{
			"FALSE",
			false,
		},
	}

	for _, tt := range tests {
		l := lexer.InitializeLexer(tt.input)
		p := InitializeParser(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if 1 != len(program.Statements) {
			t.Fatalf("program has not enough statements, got=%d", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got =%T", program.Statements[0])
		}

		boolean, ok := statement.Expression.(*ast.Boolean)

		if !ok {
			t.Fatalf("expression not *ast.Boolean, got=%T", statement.Expression)
		}

		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not '%t', got='%t", tt.expectedBoolean, boolean.Value)
		}
	}
}
