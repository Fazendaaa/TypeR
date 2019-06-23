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
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{
			"let x <- 5;",
			"x",
			5,
		},
		{
			"let y<-10",
			"y",
			10,
		},
		{
			"let foo <- y",
			"foo",
			"y",
		},
	}

	for _, tt := range tests {
		l := lexer.InitializeLexer(tt.input)
		p := InitializeParser(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		if 1 != len(program.Statements) {
			t.Fatalf("program.Statements does not contains %d statements, got=%d\n", 1, len(program.Statements))
		}

		statement := program.Statements[0]

		if !testLetStatements(t, statement, tt.expectedIdentifier) {
			return
		}

		value := statement.(*ast.LetStatement).Value

		if !testLiteralExpresion(t, value, tt.expectedValue) {
			return
		}
	}
}

// TestReturnStatements :
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

// TestIdentifierExpression :
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

// TestIntegerLiteralExpression :
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
		t.Fatalf("expression not *ast.IntegralLiteral, got=%T", statement.Expression)
	}

	if 5 != literal.Value {
		t.Errorf("literal.Value not '%d', got=%d", 5, literal.Value)
	}

	if "5" != literal.TokenLiteral() {
		t.Errorf("literal.TokenLiteral not '%s', got=%s", "5", literal.TokenLiteral())
	}
}

// TestParsingPrefixExpressions :
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

		expression, ok := statement.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("statement is not ast.PrefixExpression, got=%T", statement.Expression)
		}

		if expression.Operator != tt.operator {
			t.Fatalf("expression.Operator is not '%s', got=%s", tt.operator, expression.Operator)
		}

		if !testLiteralExpresion(t, expression.Right, tt.value) {
			return
		}
	}
}

// TestParsingInfixExpressions :
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

// TestOperatorPrecedenceParsing  :
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
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(TRUE == TRUE)",
			"(!(TRUE == TRUE))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

// TestBooleanExpression :
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

// TestConditionalIfOnlyExpressions :
func TestConditionalIfOnlyExpressions(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.InitializeLexer(input)
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

	expression, ok := statement.Expression.(*ast.ConditionalExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.ConditionalExpression, got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if 1 != len(expression.Consequence.Statements) {
		t.Errorf("consequence is not 1 statement, got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not set.ExpressionStatement, got%T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if nil != expression.Alternative {
		t.Errorf("exp.Alternative was not nil, got=%+v", expression.Alternative)
	}
}

// TestConditionalIfElseExpressions :
func TestConditionalIfElseExpressions(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.InitializeLexer(input)
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

	expression, ok := statement.Expression.(*ast.ConditionalExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.ConditionalExpression, got=%T", statement.Expression)
	}

	if !testInfixExpression(t, expression.Condition, "x", "<", "y") {
		return
	}

	if 1 != len(expression.Consequence.Statements) {
		t.Errorf("consequence is not 1 statement, got=%d", len(expression.Consequence.Statements))
	}

	consequence, ok := expression.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got%T", expression.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if nil == expression.Alternative {
		t.Errorf("exp.Alternative was nil")
	}

	if 1 != len(expression.Alternative.Statements) {
		t.Errorf("alternative is not 1 stamentent, got=%d", len(expression.Alternative.Statements))
	}

	alternative, ok := expression.Alternative.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement, got=%T", expression.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

// TestFunctionLiteral :
func TestFunctionLiteral(t *testing.T) {
	input := `function(x, y) { x + y; }`

	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if 1 != len(program.Statements) {
		t.Fatalf("program.Statements does not contain %d statements, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statement[0] is not ast.ExpressionStatement, got=%T", program.Statements[0])
	}

	function, ok := statement.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("statement.Expression is not ast.FunctionLiteral, got=%T", statement.Expression)
	}

	if 2 != len(function.Parameters) {
		t.Fatalf("function literal parameters wrong, want %d, got=%d", 2, len(function.Parameters))
	}

	testLiteralExpresion(t, function.Parameters[0], "x")
	testLiteralExpresion(t, function.Parameters[1], "y")

	if 1 != len(function.Body.Statements) {
		t.Fatalf("function.Body.Statements has not %d statements, got=%d", 1, len(function.Body.Statements))
	}

	bodyStatements, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body statement is not ast.ExpressionStatement, got=%T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatements.Expression, "x", "+", "y")
}

// TestFunctionParametersParsing :
func TestFunctionParametersParsing(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
	}{
		{
			input:              "function() {};",
			expectedParameters: []string{},
		},
		{
			input: "function(x) {};",
			expectedParameters: []string{
				"x",
			},
		},
		{
			input: "function(x, y, z) {};",
			expectedParameters: []string{
				"x",
				"y",
				"z",
			},
		},
	}

	for _, tt := range tests {
		l := lexer.InitializeLexer(tt.input)
		p := InitializeParser(l)
		program := p.ParseProgram()

		checkParserErrors(t, p)

		statement := program.Statements[0].(*ast.ExpressionStatement)
		function := statement.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParameters) {
			t.Errorf("length parameters wrong, want %d, got=%d", len(function.Parameters), len(tt.expectedParameters))
		}

		for i, identifier := range tt.expectedParameters {
			testLiteralExpresion(t, function.Parameters[i], identifier)
		}
	}
}

// TestCallExporessionParsing :
func TestCallExporessionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5)"

	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)
	program := p.ParseProgram()

	checkParserErrors(t, p)

	if 1 != len(program.Statements) {
		t.Fatalf("program.Statements does not contain %d statements, got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("statement is not ExpressionStatement, got=%T", program.Statements[0])
	}

	expression, ok := statement.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("statement.Expression is not ast.CallExpression, got=%T", statement.Expression)
	}

	if !testIdentifier(t, expression.Function, "add") {
		return
	}

	if 3 != len(expression.Arguments) {
		t.Fatalf("wrong length of arguments, got=%d", len(expression.Arguments))
	}

	testLiteralExpresion(t, expression.Arguments[0], 1)
	testInfixExpression(t, expression.Arguments[1], 2, "*", 3)
	testInfixExpression(t, expression.Arguments[2], 4, "+", 5)
}
