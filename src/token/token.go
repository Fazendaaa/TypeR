package token

// TokenType : this will work as a PoC only, needs to change it to an int or a byte later on
type TokenType string

// Token : stores the information token related; later on add a line and column to it to make an easier debug later on
type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF     = "EOF"
	ILLEGAL = "ILLEGAL"

	INT    = "INT"
	DOUBLE = "DOUBLE"
	STRING = "STRING"

	IDENTIFIER = "IDENTIFIER"

	EXPORT   = "EXPORT"
	LET      = "LET"
	CONST    = "CONST"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"

	POINT              = "."
	BANG               = "!"
	EQUAL              = "="
	PLUS               = "+"
	MINUS              = "-"
	SLASH              = "/"
	ASTERISK           = "*"
	LESS_THAN          = "<"
	GREATER_THAN       = ">"
	LESS_THAN_EQUAL    = "<="
	GREATER_THAN_EQUAl = ">="

	COMMA             = ","
	SEMICOLON         = ";"
	LEFT_PARENTHESIS  = "("
	RIGHT_PARENTHESIS = ")"
	LEFT_BRACE        = "{"
	RIGHT_BRACE       = "}"
	LEFT_BRACKET      = "["
	RIGHT_BRACKET     = "]"

	ASSIGN       = "<-"
	DOUBLE_EQUAL = "=="
	DIFFERENT    = "!="
)

var keywords = map[string]TokenType{
	"if":       IF,
	"let":      LET,
	"TRUE":     TRUE,
	"else":     ELSE,
	"FALSE":    FALSE,
	"return":   RETURN,
	"<-":       ASSIGN,
	"function": FUNCTION,
	"<":        LESS_THAN,
	"<=":       LESS_THAN_EQUAL,
}

// LookupIdentifier :
func LookupIdentifier(identification string) TokenType {
	if tok, ok := keywords[identification]; ok {
		return tok
	}

	return IDENTIFIER
}
