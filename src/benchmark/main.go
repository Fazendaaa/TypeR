package main

import (
	"flag"
	"fmt"
	"time"

	"../compiler"
	"../evaluator"
	"../lexer"
	"../object"
	"../parser"
	"../virtualmachine"
)

var engine = flag.String("engine", "virtualmachine", "user 'virtualmachine' or 'evaluator'")

var input = `
fibonacci <-function(x) {
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

fibonacci(40)
`

func main() {
	flag.Parse()

	var duration time.Duration
	var result object.Object

	l := lexer.InitializeLexer(input)
	p := parser.InitializeParser(l)
	program := p.ParseProgram()

	if "virtualmachine" == *engine {
		comp := compiler.InitializeCompiler()
		err := comp.Compile(program)

		if nil != err {
			fmt.Printf("compiler error: %s", err)

			return
		}

		machine := virtualmachine.InitializeVirtualMachine(comp.Bytecode())

		start := time.Now()

		err = machine.Run()

		if nil != err {
			fmt.Printf("Virtual Machine error: %s", err)

			return
		}

		duration = time.Since(start)

		result = machine.LastPoppedStackElement()
	} else {
		env := object.InitializeEnvironment()
		start := time.Now()
		result = evaluator.Eval(program, env)
		duration = time.Since(start)
	}

	fmt.Printf("engine=%s, result=%s, duration=%s\n", *engine, result.Inspect(), duration)
}
