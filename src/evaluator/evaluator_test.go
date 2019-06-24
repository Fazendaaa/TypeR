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

// testNullObject :
func testNullObject(t *testing.T, obj object.Object) bool {
	if NULL != obj {
		t.Errorf("object is not NULL, got=%T (%+b)", obj, obj)

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
		{
			"5 + 5 + 5 + 5 - 10",
			10,
		},
		{
			"2 * 2 * 2 * 2 * 2",
			32,
		},
		{
			"-50 + 100 - 50",
			0,
		},
		{
			"5 * 2 + 10",
			20,
		},
		{
			"5 + 2 * 10",
			25,
		},
		{
			"20 + 2 * -10",
			0,
		},
		{
			"50 / 2 * 2 + 10",
			60,
		},
		{
			"2 * (5 + 10)",
			30,
		},
		{
			"3 * 3 * 3 + 10",
			37,
		},
		{
			"3 * (3 * 3) + 10",
			37,
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10",
			50,
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
		{
			"1 < 2",
			true,
		},
		{
			"1 > 2",
			false,
		},
		{
			"1 < 1",
			false,
		},
		{
			"1 > 1",
			false,
		},
		{
			"1 == 1",
			true,
		},
		{
			"1 != 1",
			false,
		},
		{
			"1 == 2",
			false,
		},
		{
			"1 != 2",
			true,
		},
		{
			"TRUE == TRUE",
			true,
		},
		{
			"FALSE == FALSE",
			true,
		},
		{
			"TRUE == FALSE",
			false,
		},
		{
			"TRUE != FALSE",
			true,
		},
		{
			"FALSE != TRUE",
			true,
		},
		{
			"(1 < 2) == TRUE",
			true,
		},
		{
			"(1 < 2) == FALSE",
			false,
		},
		{
			"(1 > 2) == TRUE",
			false,
		},
		{
			"(1 > 2) == FALSE",
			true,
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

// TestConditionalsExpressions :
func TestConditionalsExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"if (TRUE) { 10 }",
			10,
		},
		{
			"if (FALSE) { 10 }",
			nil,
		},
		{
			"if (1) { 10 }",
			10,
		},
		{
			"if (1 < 2) { 10 }",
			10,
		},
		{
			"if (1 > 2) { 10 }",
			nil,
		},
		{
			"if (1 > 2) { 10 } else { 20 }",
			20,
		},
		{
			"if (1 < 2) { 10 } else { 20 }",
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

// TestReturnStatements :
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"return 10",
			10,
		},
		{
			"return 10;",
			10,
		},
		{
			"return 10; 9",
			10,
		},
		{
			"return 10; 9;",
			10,
		},
		{
			"return 2 * 5;",
			10,
		},
		{
			"9; return 2 * 5;",
			10,
		},
		{
			"9; return 2 * 5; 9;",
			10,
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return 10;
				}

				return 1;
			}`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		testIntegerObject(t, evaluated, tt.expected)
	}
}

// TestErrorHandling :
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 * TRUE",
			"type mismatch: INTEGER * BOOLEAN",
		},
		{
			"5 + TRUE; 5",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-TRUE",
			"unknown operator: -BOOLEAN",
		},
		{
			"TRUE + FALSE",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; TRUE + FALSE; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { return TRUE + FALSE; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
				if (10 > 1) {
					return TRUE + FALSE;
				}
			
				return 1;
			}`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errorObject, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error Object returned, got=%T(%+v", evaluated, evaluated)

			continue
		}

		if errorObject.Message != tt.expectedMessage {
			t.Errorf("wrong error message, expected='%s', got='%s'", tt.expectedMessage, errorObject.Message)
		}
	}
}
