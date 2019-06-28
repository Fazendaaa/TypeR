package lexer

import (
	"testing"

	"../token"
)

// TestNextToken :
func TestNextToken(t *testing.T) {
	input := `five <- 5;ten <- 10

add <- function(x, y) x + y

result <- add(five, ten)
!-/*5;
5 < 10 > 5;

if (5 <= 10) {
	return TRUE
} else {
	return FALSE;
}

10 == 10;
10 != 9

let foo <- 10
foo <- 5
"foo"
"bar"
"foo bar"

[1, 2]
`

	test := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "<-"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "<-"},
		{token.INT, "10"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "<-"},
		{token.FUNCTION, "function"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "<-"},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LESS_THAN, "<"},
		{token.INT, "10"},
		{token.GREATER_THAN, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LEFT_PARENTHESIS, "("},
		{token.INT, "5"},
		{token.LESS_THAN_EQUAL, "<="},
		{token.INT, "10"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "TRUE"},
		{token.RIGHT_BRACE, "}"},
		{token.ELSE, "else"},
		{token.LEFT_BRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "FALSE"},
		{token.SEMICOLON, ";"},
		{token.RIGHT_BRACE, "}"},
		{token.INT, "10"},
		{token.DOUBLE_EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.DIFFERENT, "!="},
		{token.INT, "9"},
		{token.LET, "let"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, "<-"},
		{token.INT, "10"},
		{token.IDENTIFIER, "foo"},
		{token.ASSIGN, "<-"},
		{token.INT, "5"},
		{token.STRING, "foo"},
		{token.STRING, "bar"},
		{token.STRING, "foo bar"},
		{token.LEFT_BRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RIGHT_BRACKET, "]"},
		{token.EOF, ""},
	}

	l := InitializeLexer(input)

	for i, tt := range test {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong\n\texpected=%q, got=%q", i, tt.expectedLiteral, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong\n\texpected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
