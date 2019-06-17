package lexer

import (
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

// isLetter : maybe PLUS '?' and '!' as valid also in a near future -- R doesn't allow it
func isLetter(char byte) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_'
}

// isDigit :
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

// peekChar :
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}

	return l.input[l.readPosition]
}

// readChar :
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

// readIt :
func readIt(l *Lexer, isIt func(char byte) bool) string {
	position := l.position

	for isIt(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

// readIdentifier :
func (l *Lexer) readIdentifier() string {
	return readIt(l, isLetter)
}

// readNumber :
func (l *Lexer) readNumber() string {
	return readIt(l, isDigit)
}

// newToken :
func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(char),
	}
}

// newPeekedToken :
func newPeekedToken(l *Lexer, t token.TokenType) token.Token {
	char := l.char
	l.readChar()

	literal := string(char) + string(l.char)

	return token.Token{
		Type:    t,
		Literal: literal,
	}
}

// skipWhitespace :
func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

// NextToken :
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.char {
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
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
	case '>':
		tok = newToken(token.GREATER_THAN, l.char)
	case '<':
		if l.peekChar() == '-' {
			tok = newPeekedToken(l, token.ASSIGN)
		} else if l.peekChar() == '=' {
			tok = newPeekedToken(l, token.LESS_THAN_EQUAL)
		} else {
			tok = newToken(token.LESS_THAN, l.char)
		}
	case '=':
		if l.peekChar() == '=' {
			tok = newPeekedToken(l, token.DOUBLE_EQUAL)
		} else {
			// In lists
			tok = newToken(token.ILLEGAL, l.char)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = newPeekedToken(l, token.DIFFERENT)
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)

			return tok
		} else if isDigit(l.char) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()

			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar()

	return tok
}

// InitializeLexer :
func InitializeLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}
