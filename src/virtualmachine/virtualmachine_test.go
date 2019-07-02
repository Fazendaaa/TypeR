package virtualmachine

import (
	"fmt"
	"testing"

	"../ast"
	"../compiler"
	"../lexer"
	"../object"
	"../parser"
)

// virtualMachineTestCase :
type virtualMachineTestCase struct {
	input    string
	expected interface{}
}

// parse :
func parse(input string) *ast.Program {
	l := lexer.InitializeLexer(input)
	p := parser.InitializeParser(l)

	return p.ParseProgram()
}

// testIntegerObject :
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)

	if !ok {
		return fmt.Errorf("object is not Integer, got=%T (%+v)", actual, actual)
	}

	if expected != result.Value {
		return fmt.Errorf("object has wrong value, got=%d, want=%d", result.Value, expected)
	}

	return nil
}

// testBooleanObject :
func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)

	if !ok {
		return fmt.Errorf("object is not Boolean, got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got=%t, want=%t", result.Value, expected)
	}

	return nil
}

// testExpectedObject :
func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)

		if nil != err {
			t.Errorf("testIntegerObject failed: %s", err)
		}

	case bool:
		err := testBooleanObject(bool(expected), actual)

		if nil != err {
			t.Errorf("testBooleanObject failed: %s", err)
		}

	case *object.Null:
		if actual != NULL {
			t.Errorf("object is not NULL: %T (%+v)", actual, actual)
		}
	}
}

// runVirtualMachineTests :
func runVirtualMachineTests(t *testing.T, tests []virtualMachineTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)
		comp := compiler.InitializeCompiler()
		err := comp.Compile(program)

		if nil != err {
			t.Fatalf("compiler error: %s", err)
		}

		virtualMachine := InitializeVirtualMachine(comp.Bytecode())
		err = virtualMachine.Run()

		if nil != err {
			t.Fatalf("Virtual Machine error: %s", err)
		}

		stackElement := virtualMachine.LastPoppedStackElement()

		testExpectedObject(t, tt.expected, stackElement)
	}
}

// TestIntegerArithmetic :
func TestIntegerArithmetic(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			"1",
			1,
		},
		{
			"2",
			2,
		},
		{
			"1 + 2",
			3,
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
			"-50 + 100 + -50",
			0,
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10",
			50,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestBooleanExpressions :
func TestBooleanExpressions(t *testing.T) {
	tests := []virtualMachineTestCase{
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
		{
			"!(if (FALSE) { 5 })",
			true,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestConditionals :
func TestConditionals(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			"if (TRUE) { 10 }",
			10,
		},
		{
			"if (TRUE) { 10 } else { 20 }",
			10,
		},
		{
			"if (FALSE) { 10 } else { 20 }",
			20,
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
			"if (1 < 2) { 10 } else { 20 }",
			10,
		},
		{
			"if (1 > 2) { 10 } else { 20 }",
			20,
		},
		{
			"if (1 > 2) { 10 }",
			NULL,
		},
		{
			"if (FALSE) { 10 }",
			NULL,
		},
		{
			"if ((if (FALSE) { 10 })) { 10 } else { 20 }",
			20,
		},
	}

	runVirtualMachineTests(t, tests)
}
