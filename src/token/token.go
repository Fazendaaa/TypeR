package token

// TokenType : this will work as a PoC only, needs to change it to an int or a byte later on
type TokenType string

// Token : stores the information token related; later on add a line and column to it to make an easier debug later on
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	INDENTATION = "INDENTATION"
	UNKNOWN     = "UNKNOWN"
	INT         = "INT"
	DOUBLE      = "DOUBLE"

	ASSIGN   = "<-"
	ADD      = "+"
	SUBTRACT = "-"
	DIVIDE   = "/"
	MULTIPLY = "*"

	COMMA             = ","
	SEMICOLON         = ";"
	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACE        = "{"
	RIGHT_BRACE       = "}"

	FUNCTION = "function"
	LET      = "let"
)
