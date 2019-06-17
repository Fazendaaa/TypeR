package lexer

import (
	"testing"

	"../token"
)

// Lexer :
type Lexer struct {
	input string
	// current position in input (points to current char)
	position int
	// current reading position in input (after current char)
	readPosition int
	// current char under examination
	char byte
}

func initLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

// readChar :
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
}

func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
	}
}

// NextToken :
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.char {
	case '<':
		// needs to look futher to better decide whether or not is an assignment or simply a less than or less than equals
		tok = newToken(token.ASSIGN, l.char)
	case '+':
		tok = newToken(token.ADD, l.char)
	case '-':
		tok = newToken(token.SUBTRACT, l.char)
	case '*':
		tok = newToken(token.MULTIPLY, l.char)
	case '/':
		tok = newToken(token.DIVIDE, l.char)
	case '(':
		tok = newToken(token.LEFT_PARENTHESIS, l.char)
	case ')':
		tok = newToken(token.RIGHT_PARENTHESIS, l.char)
	case '{':
		tok = newToken(token.LEFT_BRACE, l.char)
	case '}':
		tok = newToken(token.RIGHT_BRACE, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar()

	return tok
}

// TestNextToken :
func TestNextToken(t *testing.T) {
	input := `<-*+-/(){},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "<-"},
		{token.MULTIPLY, "*"},
		{token.ADD, "+"},
		{token.SUBTRACT, "-"},
		{token.DIVIDE, "/"},
		{token.LEFT_PARENTHESIS, "("},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RIGHT_BRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := initLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong\n\texpected=%q, got=%q", i, tt.expectedLiteral, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong\n\texpected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
