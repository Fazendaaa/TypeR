package repl

import (
	"bufio"
	"fmt"
	"io"

	"../lexer"
	"../parser"
)

// PROMPT :
const PROMPT = ">> "

// printParseErrors :
func printParseErrors(out io.Writer, errors []string) {
	for _, message := range errors {
		io.WriteString(out, "\t"+message+"\n")
	}
}

// Start :
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
