package repl

import (
	"bufio"
	"fmt"
	"io"

	"../compiler"
	"../lexer"
	"../object"
	"../parser"
	"../virtualmachine"
)

// PROMPT :
const PROMPT = "λ "

// printParseErrors :
func printParseErrors(out io.Writer, errors []string) {
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}

// Start :
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	constants := []object.Object{}
	globals := make([]object.Object, virtualmachine.GlobalSize)
	symbolTable := compiler.InitializeSymbolTable()

	for index, value := range object.Builtins {
		symbolTable.DefineBuiltin(index, value.Name)
	}

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.InitializeLexer(line)
		p := parser.InitializeParser(l)
		program := p.ParseProgram()

		if 0 != len(p.Errors()) {
			printParseErrors(out, p.Errors())

			continue
		}

		comp := compiler.InitializeWithState(symbolTable, constants)
		err := comp.Compile(program)

		if nil != err {
			fmt.Fprintf(out, "Woops: Compilation failed:\n %s\n", err)

			continue
		}

		code := comp.Bytecode()
		constants = code.Constants

		machine := virtualmachine.InitializeWithGlobalStore(code, globals)
		err = machine.Run()

		if nil != err {
			fmt.Fprintf(out, "Woops: executing bytecode fails:\n %s\n", err)

			continue
		}

		stackTop := machine.LastPoppedStackElement()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}
