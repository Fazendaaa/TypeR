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
	environment := object.InitializeEnvironment()

	return Eval(program, environment)
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
		{
			"foo",
			"identifier not found: foo",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
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

//  TestConstStatements :
func TestConstStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"a <- 5; a;",
			5,
		},
		{
			"a <- 5 * 5; a",
			25,
		},
		{
			"a <- 5; b <- a; b;",
			5,
		},
		{
			"a <- 5; b <- a; c <- a + b + 5; c",
			15,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

//  TestLetStatements :
func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let a <- 5; a;",
			5,
		},
		{
			"let a <- 5 * 5; a",
			25,
		},
		{
			"let a <- 5; let b <- a; b;",
			5,
		},
		{
			"let a <- 5; let b <- a; let c <- a + b + 5; c",
			15,
		},
		{
			"let a <- 5 * 5; a <- a + 2; a",
			27,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

// TestFunctionObject :
func TestFunctionObject(t *testing.T) {
	input := "function(x) { x + 2; };"
	expectedBody := "(x + 2)"
	evaluated := testEval(input)
	function, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function, got=%T (%+v)", evaluated, evaluated)
	}

	if 1 != len(function.Parameters) {
		t.Fatalf("function has wrong parameters, expected %d, got=%d", 1, len(function.Parameters))
	}

	if "x" != function.Parameters[0].String() {
		t.Fatalf("parameter is not 'x', got=%q", function.Parameters[0])
	}

	if expectedBody != function.Body.String() {
		t.Fatalf("body is not %q, got=%q", expectedBody, function.Body.String())
	}
}

// TestFunctionApplication
func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"let identity <- function(x) { x; }; identity(5);",
			5,
		},
		{
			"let identity <- function(x) { return x; }; identity(5);",
			5,
		},
		{
			"let double <- function(x) { x * 2; }; double(5)",
			10,
		},
		{
			"let add <- function(x, y) { x + y }; add(5, 5)",
			10,
		},
		{
			"let add <- function(x, y) { x + y }; add(5 + 5, add(5, 5))",
			20,
		},
		{
			"function(x) { x }(5)",
			5,
		},
		{
			"add <- function(x, y) { x + y }; add(5, 5 * 5)",
			30,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

// TestFunctionMemoization :
func TestFunctionMemoization(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`fibonacci <-function(x) {
				if (0 == x) {
					0
				} else {
					if (1 == x) {
						1
					} else {
						fibonacci(x - 1) + fibonacci(x - 2)
					}
				}
			}

			fibonacci(20)
			`,
			6765,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

// TestClosures :
func TestClosures(t *testing.T) {
	input := `
	let newAdder <- function(x) {
		function(y) { x + y };
	};

	let addTwo <- newAdder(2)
	addTwo(2)
	`

	testIntegerObject(t, testEval(input), 4)
}

// TestStringLiteral :
func TestStringLiteral(t *testing.T) {
	input := `"Hello World!";`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}

	if "Hello World!" != str.Value {
		t.Errorf("String has wrong value, got=%q, expected was '%q'", str.Value, "Hello World!")
	}
}

// TestStringConcatenation :
func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String, got=%T (%+v)", evaluated, evaluated)
	}

	if "Hello World!" != str.Value {
		t.Errorf("String has wrong value, got=%q, expected was '%q'", str.Value, "Hello World!")
	}
}

// TestBuiltinFunction :
func TestBuiltinFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`len("")`,
			0,
		},
		{
			`len("four")`,
			4,
		},
		{
			`len("hello world!")`,
			12,
		},
		{
			`len(1)`,
			"parameters to `len` not supported, got=INTEGER",
		},
		{
			`len("one", "two")`,
			"wrong number of parameters, got=2, want=1",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)

			if !ok {
				t.Errorf("object is not Error, got=%T (%+v", evaluated, evaluated)

				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message, expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

// TestArrayLiterals :
func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)

	if !ok {
		t.Fatalf("object is not Array, got=%T (%+v)", evaluated, evaluated)
	}

	if 3 != len(result.Elements) {
		t.Fatalf("array has wrong number of elements, got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

// TestArrayIndexExpressions :
func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let index <- 0; [1][index]",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let myArray <- [1, 2, 3]; myArray[2]",
			3,
		},
		{
			"let myArray <- [1, 2, 3]; myArray[0] + myArray[1] + myArray[2]",
			6,
		},
		{
			"let myArray <- [1, 2, 3]; let i <- myArray[0]; i",
			1,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
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
