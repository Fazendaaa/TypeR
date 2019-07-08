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

		case []code.Instructions:
			fn, ok := actual[index].(*object.CompiledFunction)

			if !ok {
				return fmt.Errorf("constant %d - not a function: %T", index, actual[index])
			}

			err := testInstructions(constant, fn.Instructions)

			if nil != err {
				return fmt.Errorf("constant %d - testInstructions failed: %s", index, err)
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

// TestIndexExpressions :
func TestIndexExpressions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "[1, 2, 3][1 + 1]",
			expectedConstants: []interface{}{
				1,
				2,
				3,
				1,
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
		{
			input: "[1, 2 + 2, 3][1 + 1]",
			expectedConstants: []interface{}{
				1,
				2,
				2,
				3,
				1,
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpAdd),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpArray, 3),
				code.Make(code.OpConstant, 4),
				code.Make(code.OpConstant, 5),
				code.Make(code.OpAdd),
				code.Make(code.OpIndex),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestFunctions :
func TestFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `function() { return 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `function() { 5 + 10 }`,
			expectedConstants: []interface{}{
				5,
				10,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `function() { 1; 2 }`,
			expectedConstants: []interface{}{
				1,
				2,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestCompilerScopes :
func TestCompilerScopes(t *testing.T) {
	compiler := InitializeCompiler()

	if 0 != compiler.scopeIndex {
		t.Errorf("scopeIndex wrong, got=%d, want=%d", compiler.scopeIndex, 0)
	}

	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMultiply)

	compiler.enterScope()

	if 1 != compiler.scopeIndex {
		t.Errorf("scopeIndex wrong, got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSubtract)

	if 1 != len(compiler.scopes[compiler.scopeIndex].instructions) {
		t.Errorf("instructions length wrong, got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction

	if code.OpSubtract != last.Opcode {
		t.Errorf("lastInstruction.Opcode wrong, got=%d, want=%d", last.Opcode, code.OpSubtract)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.leaveScope()

	if 0 != compiler.scopeIndex {
		t.Errorf("scopeIndex wrong, got=%d, want=%d", compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}

	if nil != compiler.symbolTable.Outer {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(code.OpAdd)

	if 2 != len(compiler.scopes[compiler.scopeIndex].instructions) {
		t.Errorf("instructions length wrong, got=%d", len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction

	if code.OpAdd != last.Opcode {
		t.Errorf("lastInstruction.OpCode wrong, got=%d, want=%d", last.Opcode, code.OpAdd)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction

	if code.OpMultiply != previous.Opcode {
		t.Errorf("previousInstruction.OpCode wrong, got=%d, want=%d", previous.Opcode, code.OpMultiply)
	}
}

// TestFunctionsWithoutReturnValue :
func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `function() { }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpReturn),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestFunctionCalls :
func TestFunctionCalls(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `function() { 24 }()`,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let noParam <- function() { 24 }
			noParam()
			`,
			expectedConstants: []interface{}{
				24,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let oneParameter <- function(a) { a }
			oneParameter(24)
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
				24,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let manyArg <- function(a, b, c) { a; b; c }
			manyArg(24, 25, 26)
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetLocal, 2),
					code.Make(code.OpReturnValue),
				},
				24,
				25,
				26,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 1),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpConstant, 3),
				code.Make(code.OpCall, 3),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestLetStatScopes :
func TestLetStatScopes(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let number <- 55;

			function() { number; }
			`,
			expectedConstants: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			function() {
				let number <- 55;
				
				number;
			}
			`,
			expectedConstants: []interface{}{
				55,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			function() {
				let a <- 55;
				let b <- 77;
				
				a + b;
			}
			`,
			expectedConstants: []interface{}{
				55,
				77,
				[]code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 1),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpGetLocal, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestBultins :
func TestBultins(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			len([])
			push([], 1)
			`,
			expectedConstants: []interface{}{
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpGetBuiltin, 1),
				code.Make(code.OpArray, 0),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
				code.Make(code.OpGetBuiltin, 5),
				code.Make(code.OpArray, 0),
				code.Make(code.OpConstant, 0),
				code.Make(code.OpCall, 2),
				code.Make(code.OpPop),
			},
		},
		{
			input: `function () { len([]) }`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetBuiltin, 1),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 0, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestClosures :
func TestClosures(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			function(a) {
				function(b) {
					a + b
				}
			}
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetFreeVariable, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 0, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			function(a) {
				function(b) {
					function(c) {
						a + b + c
					}
				}
			}
			`,
			expectedConstants: []interface{}{
				[]code.Instructions{
					code.Make(code.OpGetFreeVariable, 0),
					code.Make(code.OpGetFreeVariable, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetFreeVariable, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 0, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 1, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 2, 0),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let global <- 55

			function() {
				let a <- 66

				function() {
					let b <- 77

					function() {
						let c <- 88

						global + a + b + c
					}
				}
			}
			`,
			expectedConstants: []interface{}{
				55,
				66,
				77,
				88,
				[]code.Instructions{
					code.Make(code.OpConstant, 3),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpGetFreeVariable, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpGetFreeVariable, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpAdd),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 2),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetFreeVariable, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 4, 2),
					code.Make(code.OpReturnValue),
				},
				[]code.Instructions{
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpClosure, 5, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpConstant, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpClosure, 6, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}

// TestRecursiveFunctions :
func TestRecursiveFunctions(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
			let countDown <- function(x) { countDown(x - 1) }

			countDown(1)
			`,
			expectedConstants: []interface{}{
				1,
				[]code.Instructions{
					code.Make(code.OpCurrentClosure),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSubtract),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
				1,
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 1, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpConstant, 2),
				code.Make(code.OpCall, 1),
				code.Make(code.OpPop),
			},
		},
		{
			input: `
			let wrapper <- function() {
				let countDown <- function(x) { countDown(x - 1) };

				countDown(1);
			}

			wrapper();
			`,
			expectedConstants: []interface{}{
				1,
				[]code.Instructions{
					code.Make(code.OpCurrentClosure),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSubtract),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
				1,
				[]code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetLocal, 0),
					code.Make(code.OpGetLocal, 0),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpCall, 1),
					code.Make(code.OpReturnValue),
				},
			},
			expectedInstructions: []code.Instructions{
				code.Make(code.OpClosure, 3, 0),
				code.Make(code.OpSetGlobal, 0),
				code.Make(code.OpGetGlobal, 0),
				code.Make(code.OpCall, 0),
				code.Make(code.OpPop),
			},
		},
	}

	runCompilerTests(t, tests)
}
