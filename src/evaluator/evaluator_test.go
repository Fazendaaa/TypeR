package evaluator

import (
	"testing"

	"../lexer"
	"../object"
	"../parser"
)

// testEval :
func testEval(input string) object.Object {
	l := lexer.InitializeLexer(input)
	p := parser.InitializeParser(l)
	program := p.ParseProgram()

	return Eval(program)
}

// testIntegerObject :
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Ingeter, got=%T (%+v)", obj, obj)

		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%d, expected was=%d", result.Value, expected)

		return false
	}

	return true
}

// testBooleanObject :
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean, got=%T (%+v)", obj, obj)

		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value, got=%t, expected was=%t", result.Value, expected)

		return false
	}

	return true
}

// TestEvalIntegerExpression :
func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"5",
			5,
		},
		{
			"10",
			10,
		},
		{
			"-5",
			-5,
		},
		{
			"-10",
			-10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testIntegerObject(t, evaluated, tt.expected)
	}
}

// TestEvalBooleanExpression :
func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
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
		evaluated := testEval(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}

// TestBangOperator :
func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"!TRUE",
			false,
		},
		{
			"!FALSE",
			true,
		},
		{
			"!5",
			false,
		},
		{
			"!!TRUE",
			true,
		},
		{
			"!!FALSE",
			false,
		},
		{
			"!!5",
			true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testBooleanObject(t, evaluated, tt.expected)
	}
}
