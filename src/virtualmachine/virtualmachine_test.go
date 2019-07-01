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

// testExpectedObject :
func testExpectedObject(t *testing.T, expected interface{}, actual object.Object) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)

		if nil != err {
			t.Errorf("testIntegerObject failed: %s", err)
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

		stackElement := virtualMachine.StackTop()

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
	}

	runVirtualMachineTests(t, tests)
}
