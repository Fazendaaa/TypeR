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

// testStringObject :
func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)

	if !ok {
		return fmt.Errorf("object is not String, got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value, got =%q, want+%q", result.Value, expected)
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

	case string:
		err := testStringObject(expected, actual)

		if nil != err {
			t.Errorf("testStringObject failed: %s", err)
		}

	case []int:
		array, ok := actual.(*object.Array)

		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)

			return
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("wrong number of elements, want=%d, got=%d", len(expected), len(array.Elements))

			return
		}

		for index, expectedElement := range expected {
			err := testIntegerObject(int64(expectedElement), array.Elements[index])

			if nil != err {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}

	case *object.Error:
		errorObject, ok := actual.(*object.Error)

		if !ok {
			t.Errorf("object is not Error: %T (%+v)", actual, actual)

			return
		}

		if errorObject.Message != expected.Message {
			t.Errorf("wrong error message, expected=%q, got=%q", expected.Message, errorObject.Message)
		}

	default:
		t.Errorf("object not defined: %T (%+v)", actual, actual)
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

		//
		// Debug propurses
		//
		// for index, constant := range comp.Bytecode().Constants {
		// 	fmt.Printf("CONSTANT %d %p (%T):\n", index, constant, constant)

		// 	switch constantType := constant.(type) {
		// 	case *object.CompiledFunction:
		// 		fmt.Printf(" Instructions:\n%s", constantType.Instructions)
		// 	case *object.Integer:
		// 		fmt.Printf(" Value: %d\n", constantType.Value)
		// 	}

		// 	fmt.Printf("\n")
		// }

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

// TestGlobalLetStatements :
func TestGlobalLetStatements(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			"let one <-  1; one",
			1,
		},
		{
			"let one <-  1; let two <- 2; one + two",
			3,
		},
		{
			"let one <-  1; let two <- one + one; one + two",
			3,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestStringExpressions :
func TestStringExpressions(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			`"foo"`,
			"foo",
		},
		{
			`"foo" + " bar"`,
			"foo bar",
		},
		{
			`"foo" + " bar" + " baz"`,
			"foo bar baz",
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestArrayLiterals :
func TestArrayLiterals(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			"[]",
			[]int{},
		},
		{
			"[1, 2, 3]",
			[]int{
				1,
				2,
				3,
			},
		},
		{
			"[1 + 2, 3 * 4, 5 + 6]",
			[]int{
				3,
				12,
				11,
			},
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestIndexExpressions :
func TestIndexExpressions(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][0 + 2]",
			3,
		},
		{
			"[[1, 1, 1]][0][0]",
			1,
		},
		{
			"[][0]",
			NULL,
		},
		{
			"[1, 2, 3][99]",
			NULL,
		},
		{
			"[1][-1]",
			NULL,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestCallingFunctionsWithoutParameters :
func TestCallingFunctionsWithoutParameters(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let fivePlusTen <- function() { 5 + 10 }
			fivePlusTen()
			`,
			expected: 15,
		},
		{
			input: `
			let first <- function() { 1 }
			let second <- function() { 2 }
			first() + second()
			`,
			expected: 3,
		},
		{
			input: `
			let a <- function() { 1 }
			let b <- function() { a() + 1 }
			let c <- function() { b() + 1 }
			c()
			`,
			expected: 3,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestFunctionsWithReturnStatement :
func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let earlyExit <- function() { return 99; 100; }
			earlyExit()
			`,
			expected: 99,
		},
		{
			input: `
			let earlyExit <- function() { return 99; return 100; }
			earlyExit()
			`,
			expected: 99,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestFunctionsWithoutReturnValue :
func TestFunctionsWithoutReturnValue(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let noReturn <- function() {}
			noReturn()
			`,
			expected: NULL,
		},
		{
			input: `
			let noReturn <- function() {}
			let noReturnAgain <- function() { noReturn() }

			noReturn()
			noReturnAgain()
			`,
			expected: NULL,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestFirstClassFunctions :
func TestFirstClassFunctions(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let returnOne <- function() { 1 }
			let returnOneReturner <- function() { returnOne }

			returnOneReturner()()
			`,
			expected: 1,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestCallingFunctionsWithBindings :
func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let one <- function() { let one <- 1; one }
			one()
			`,
			expected: 1,
		},
		{
			input: `
			let oneAndTwo <- function() { let one <- 1; let two <- 2; one + two }
			oneAndTwo()
			`,
			expected: 3,
		},
		{
			input: `
			let oneAndTwo <- function() { let one <- 1; let two <- 2; one + two }
			let threeAndFour <- function() { let three <- 3; let four <- 4; three + four }
			oneAndTwo() + threeAndFour()
			`,
			expected: 10,
		},
		{
			input: `
			let firstFoo <- function() { let foo <- 50; foo }
			let secondFoo <- function() { let foo <- 100; foo }
			firstFoo() + secondFoo()
			`,
			expected: 150,
		},
		{
			input: `
			let firstFoo <- function() { let foo <- 50; foo }
			let secondFoo <- function() { let foo <- 100; foo }
			firstFoo() + secondFoo()
			`,
			expected: 150,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestFirsClassFunctions :
func TestFirsClassFunctions(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let returnOneReturns <- function() {
				let returnOne <- function() {
					1
				}

				returnOne
			}
			
			returnOneReturns()()
			`,
			expected: 1,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestCallingFunctionsWithParametersAndBindings :
func TestCallingFunctionsWithParametersAndBindings(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let identity <- function(a) { a }
			identity(4)
			`,
			expected: 4,
		},
		{
			input: `
			let sum <- function(a, b) { a + b }
			sum(1, 2)
			`,
			expected: 3,
		},
		{
			input: `
			let sum <- function(a, b) {
				let c <- a + b;
				c;
			}
			sum(1, 2)
			`,
			expected: 3,
		},
		{
			input: `
			let sum <- function(a, b) {
				let c <- a + b;
				c;
			}
			sum(1, 2) + sum(3, 4);
			`,
			expected: 10,
		},
		{
			input: `
			let sum <- function(a, b) {
				let c <- a + b;
				c;
			}

			let outer <- function() {
				sum(1, 2) + sum(3, 4)
			}

			outer();
			`,
			expected: 10,
		},
		{
			input: `
			let globalNumber <- 10;

			let sum <- function(a, b) {
				let c <- a + b;

				c + globalNumber;
			}

			let outer <- function() {
				sum(1, 2) + sum(3, 4) + globalNumber
			}

			outer() + globalNumber;
			`,
			expected: 50,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestCallingFunctionsWithWrongParameters :
func TestCallingFunctionsWithWrongParameters(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input:    `function() { 1 }(1)`,
			expected: `wrong number of parameters: want=0, got=1`,
		},
		{
			input:    `function(a) { a }()`,
			expected: `wrong number of parameters: want=1, got=0`,
		},
		{
			input:    `function(a, b) { a + b }(1)`,
			expected: `wrong number of parameters: want=2, got=1`,
		},
	}

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.InitializeCompiler()
		err := comp.Compile(program)

		if nil != err {
			t.Fatalf("compiler error: %s", err)
		}

		vm := InitializeVirtualMachine(comp.Bytecode())
		err = vm.Run()

		if nil == err {
			t.Fatalf("expected Virtual Machine error but resulted in none.")
		}

		if err.Error() != tt.expected {
			t.Fatalf("wrong Virtual Machine error: want=%q, got=%q", tt.expected, err)
		}
	}
}

// TestBuiltinFunctions :
func TestBuiltinFunctions(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			`len("")`,
			0,
		},
		{
			`len("four")`,
			4,
		},
		{
			`len("hello, world")`,
			12,
		},
		{
			`len(1)`,
			&object.Error{
				Message: "parameters to `len` not supported, got=INTEGER",
			},
		},
		{
			`len("one", "two")`,
			&object.Error{
				Message: "wrong number of parameters, got=2, want=1",
			},
		},
		{
			`len([])`,
			0,
		},
		{
			`len([1, 2, 3])`,
			3,
		},
		{
			`puts("hello", "world!")`,
			NULL,
		},
		{
			`head([1, 2, 3])`,
			1,
		},
		{
			`head([])`,
			NULL,
		},
		{
			`head(1)`,
			&object.Error{
				Message: "parameter to `head` must be ARRAY, got INTEGER",
			},
		},
		{
			`last([1, 2, 3])`,
			3,
		},
		{
			`last([])`,
			NULL,
		},
		{
			`last(1)`,
			&object.Error{
				Message: "parameter to `last` must be ARRAY, got INTEGER",
			},
		},
		{
			`tail([1, 2, 3])`,
			[]int{
				2,
				3,
			},
		},
		{
			`tail([])`,
			NULL,
		},
		{
			`tail(1)`,
			&object.Error{
				Message: "parameter to `tail` must be ARRAY, got INTEGER",
			},
		},
		{
			`push([1, 2], 3)`,
			[]int{
				1,
				2,
				3,
			},
		},
		{
			`push([], 1)`,
			[]int{
				1,
			},
		},
		{
			`push(1, 1)`,
			&object.Error{
				Message: "parameter to `push` must be ARRAY, got INTEGER",
			},
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestClosures :
func TestClosures(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let newClosure <- function(a) {
				function() {
					a;
				};
			};

			let closure <- newClosure(99);

			closure();
			`,
			expected: 99,
		},
		{
			input: `
			let newAdder <- function(a, b) {
				function(c) {
					a + b + c;
				};
			};

			let adder <- newAdder(1, 2);

			adder(3);
			`,
			expected: 6,
		},
		{
			input: `
			let newAdder <- function(a, b) {
				let c <- a + b

				function(d) {
					c + d;
				};
			};

			let adder <- newAdder(1, 2);

			adder(3);
			`,
			expected: 6,
		},
		{
			input: `
			let newAdderOuter <- function(a, b) {
				let c <- a + b;

				function(d) {
					let e <- c + d;

					function(f) {
						e + f;
					};
				};
			};

			let newAdderInner <- newAdderOuter(1, 2);
			let adder <- newAdderInner(3);

			adder(4);
			`,
			expected: 10,
		},
		{
			input: `
			let a <- 1;

			let newAdderOuter <- function(b) {
				function(c) {
					function(d) {
						a+ b + c + d;
					};
				};
			};

			let newAdderInner <- newAdderOuter(2);
			let adder <- newAdderInner(3);

			adder(4);
			`,
			expected: 10,
		},
		{
			input: `
			let newClosure <- function(a, b) {
				let one <- function() {
					a;
				}
				let two <- function() {
					b;
				}

				function() {
					one() + two();
				}
			}

			let closure <- newClosure(1, 9)

			closure()
			`,
			expected: 10,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestRecursiveClosures :
func TestRecursiveClosures(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let countDown <- function(x) {
				if (0 == x) {
					return 0;
				} 

				countDown(x - 1);
			};

			countDown(1);
			`,
			expected: 0,
		},
		{
			input: `
			let countDown <- function(x) {
				if (0 == x) {
					return 0;
				} 

				countDown(x - 1);
			};

			let wrapper <- function() {
				countDown(1);
			}

			wrapper();
			`,
			expected: 0,
		},
		{
			input: `
			let wrapper <- function() {
				let countDown <- function(x) {
					if (0 == x) {
						return 0;
					} 
	
					countDown(x - 1);
				};

				countDown(1);
			}

			wrapper();
			`,
			expected: 0,
		},
	}

	runVirtualMachineTests(t, tests)
}

// TestRecursiveFibonacci :
func TestRecursiveFibonacci(t *testing.T) {
	tests := []virtualMachineTestCase{
		{
			input: `
			let fibonacci <- function(x) {
				if (0 == x) {
					return 0
				} else {
					if (1 == x) {
						return 1
					} else {
						fibonacci(x - 1) + fibonacci(x -2)
					}
				}
			}

			fibonacci(15)
			`,
			expected: 610,
		},
	}

	runVirtualMachineTests(t, tests)
}
