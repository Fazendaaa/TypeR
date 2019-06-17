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
`

	test := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFICATION, "five"},
		{token.ASSIGN, "<-"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFICATION, "ten"},
		{token.ASSIGN, "<-"},
		{token.INT, "10"},
		{token.IDENTIFICATION, "add"},
		{token.ASSIGN, "<-"},
		{token.FUNCTION, "function"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFICATION, "x"},
		{token.COMMA, ","},
		{token.IDENTIFICATION, "y"},
		{token.RIGHT_PARENTHESIS, ")"},
		{token.IDENTIFICATION, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFICATION, "y"},
		{token.IDENTIFICATION, "result"},
		{token.ASSIGN, "<-"},
		{token.IDENTIFICATION, "add"},
		{token.LEFT_PARENTHESIS, "("},
		{token.IDENTIFICATION, "five"},
		{token.COMMA, ","},
		{token.IDENTIFICATION, "ten"},
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
		{token.EOF, ""},
	}

	l := InitLexer(input)

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
