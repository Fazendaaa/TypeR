package compiler

import (
	"fmt"
	"testing"

	"../ast"
	"../code"
	"../lexer"
	"../object"
	"../parser"
)

// compilerTestCase :
type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

// parse :
func parse(input string) *ast.Program {
	l := lexer.InitializeLexer(input)
	p := parser.InitializeParser(l)

	return p.ParseProgram()
}

// concatInstructions :
func concatInstructions(set []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, instruction := range set {
		out = append(out, instruction...)
	}

	return out
}

// testInstructions :
func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)

	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length, want=%q, got%q", len(actual), len(concatted))
	}

	for index, instruction := range concatted {
		if actual[index] != instruction {
			return fmt.Errorf("wrong instruction at %d, want=%q, got=%q", index, concatted, actual)
		}
	}

	return nil
}

// testIntegerObject :
func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)

	if !ok {
		return fmt.Errorf("object is not Integer, got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got=%d, want=%d", result.Value, expected)
	}

	return nil
}

// testStringObject :
func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)

	if !ok {
		return fmt.Errorf("object is not String, got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got=%q, want=%q", result.Value, expected)
	}

	return nil
}

// testConstants :
func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {
	if len(expected) != len(actual) {
		return fmt.Errorf("wrong number of constants, got=%d, want=%d", len(actual), len(expected))
	}

	for index, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[index])

			if nil != err {
				return fmt.Errorf("constant %d - testInteger object failed: %s", index, err)
			}
		case string:
			err := testStringObject(constant, actual[index])

			if nil != err {
				return fmt.Errorf("constant %d - testStringObject failed: %s", index, err)
			}
		}
	}

	return nil
}

// runCompilerTests :
func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		compiler := InitializeCompiler()
		err := compiler.Compile(program)

		if nil != err {
			t.Fatalf("Compiler error: %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)

		if nil != err {
			t.Fatalf("testInstructions error: %s", err)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)

		if nil != err {
			t.Fatalf("testConstants error: %s", err)
		}
	}
}

// TestIntegerArithmetic :
func TestIntegerArithmetic(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "1 + 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1; 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 - 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSubtract),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 * 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpMultiply),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 / 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpDivide),
				code.Make(code.OpPop),
			},
		},
		{
			input: "-1",
			expectedConstants: []interface{}{
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpMinus),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestBooleanExpressions :
func TestBooleanExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "TRUE",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "FALSE",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpFalse),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 > 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 < 2",
			expectedConstants: []interface{}{
				2,
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpGreaterThan),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 == 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input: "1 != 2",
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "TRUE == FALSE",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "TRUE != FALSE",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpFalse),
				code.Make(code.OpNotEqual),
				code.Make(code.OpPop),
			},
		},
		{
			input:             "!TRUE",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpTrue),
				code.Make(code.OpBang),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestConditionals :
func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			if (TRUE) { 10 }; 3333;
			`,
			expectedConstants: []interface{}{
				10,
				3333,
			},
			expectedInstructions: []code.Instructions{
				//  0000
				code.Make(code.OpTrue),
				//  0001
				code.Make(code.OpJumpNotTruthy, 10),
				//  0004
				code.Make(code.OpConstant, 0),
				//  0007
				code.Make(code.OpJump, 11),
				//  00010
				code.Make(code.OpNull),
				//  0011
				code.Make(code.OpPop),
				//  0012
				code.Make(code.OpConstant, 1),
				//  0015
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			if (TRUE) { 10 } else { 20 } 3333;
			`,
			expectedConstants: []interface{}{
				10,
				20,
				3333,
			},
			expectedInstructions: []code.Instructions{
				//  0000
				code.Make(code.OpTrue),
				//  0001
				code.Make(code.OpJumpNotTruthy, 10),
				//  0004
				code.Make(code.OpConstant, 0),
				//  0007
				code.Make(code.OpJump, 13),
				//  0010
				code.Make(code.OpConstant, 1),
				//  0013
				code.Make(code.OpPop),
				//  0014
				code.Make(code.OpConstant, 2),
				//  0017
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestGlobalLetStatements :
func TestGlobalLetStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let one <- 1;
			let two <- 2;
			`,
			expectedConstants: []interface{}{
				1,
				2,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpSetGlobal, 1),
			},
		},
		{
			input: `
			let one <- 1;
			one;
			`,
			expectedConstants: []interface{}{
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let one <- 1;
			let two <- one;
			two;
			`,
			expectedConstants: []interface{}{
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpSetGlobal, 1),
				code.Make(code.OpGetGlobal, 1),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestStringExpressions :
func TestStringExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `"foo"`,
			expectedConstants: []interface{}{
				"foo",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `"foo" + "bar"`,
			expectedConstants: []interface{}{
				"foo",
				"bar",
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestArrayLiterals :
func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpArray, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: "[1, 2, 3]",
			expectedConstants: []interface{}{
				1,
				2,
				3,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
		{
			input: "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{
				1,
				2,
				3,
				4,
				5,
				6,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpSubtract),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpMultiply),
				code.Make(code.OpArray, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
