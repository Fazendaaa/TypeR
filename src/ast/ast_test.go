package ast

import (
	"testing"

	"../token"
)

// TestString :
func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFICATION,
						Literal: "myVar",
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:    token.IDENTIFICATION,
						Literal: "anotherVar",
					},
					Value: "anotherVar",
				},
			},
		},
	}

	if "let myVar <- anotherVar" != program.String() {
		t.Errorf("program.String() wrong, got=%q", program.String())
	}
}
