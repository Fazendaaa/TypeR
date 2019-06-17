package parser

import (
	"testing"

	"../ast"
	"../lexer"
)

// testLetStatements :
func testLetStatements(t *testing.T, s ast.Statement, name string) bool {
	if "let" != s.TokenLiteral() {
		t.Errorf("s.TokenLiteral not 'let', got=%q", s.TokenLiteral())

		return false
	}

	letStatement, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement, got=%T", s)

		return false
	}
	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not '%s', got=%s", name, letStatement.Name.Value)

		return false
	}
	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("lestStatement.Name.TokenLiteral() not '%s', got=%s", name, letStatement.Name.TokenLiteral())

		return false
	}

	return true
}

// TestLetStatements :
func TestLetStatements(t *testing.T) {
	input := `
let x <- 5;let y <- 10
let foo <- 2345678
`

	l := lexer.InitializeLexer(input)
	p := InitializeParser(l)

	program := p.ParseProgram()

	if nil == program {
		t.Fatalf("ParseProgram() returned nil")
	}
	if 3 != len(program.Statements) {
		t.Fatalf("program.Statements does not contains three statements, go=%d\n", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		statement := program.Statements[i]

		if !testLetStatements(t, statement, tt.expectedIdentifier) {
			return
		}
	}
}
