// Package lexer defines the structure and methods for lexical analysis of JSON.
package lexer

import (
	"github.com/oabrivard/gojson/token"
)

// Lexer struct represents a lexical analyzer with its input, current position,
// next reading position, and current character.
type Lexer struct {
	input        string // the string being scanned
	position     int    // current position in the input (points to current char)
	readPosition int    // current reading position in the input (after current char)
	ch           byte   // current char under examination
	line         int    // current line number
	column       int    // current column number
}

// NewLexer creates and initializes a new Lexer with the given input string.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar() // Initialize the first character
	return l
}

// NextToken reads the next token from the input and returns it.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace() // Skip any whitespace before the next token

	// Switch on the current character to determine the token type
	switch l.ch {
	case '{':
		tok = token.NewToken(token.BEGIN_OBJECT, l.ch, l.line, l.column)
	case '}':
		tok = token.NewToken(token.END_OBJECT, l.ch, l.line, l.column)
	case '[':
		tok = token.NewToken(token.BEGIN_ARRAY, l.ch, l.line, l.column)
	case ']':
		tok = token.NewToken(token.END_ARRAY, l.ch, l.line, l.column)
	case ':':
		tok = token.NewToken(token.NAME_SEPARATOR, l.ch, l.line, l.column)
	case ',':
		tok = token.NewToken(token.VALUE_SEPARATOR, l.ch, l.line, l.column)
	case '"':
		tok = token.NewTokenWithValue(token.STRING, l.readString(), l.line, l.column)
	case 0:
		tok = token.NewTokenWithValue(token.EOF, "", l.line, l.column)
	default:
		// Handle numbers and identifiers or mark as illegal
		if isDigit(l.ch) || l.ch == '-' {
			return token.NewTokenWithValue(token.NUMBER, l.readNumber(), l.line, l.column)
		} else if isLetter(l.ch) {
			s := l.readIdentifier()
			t := token.LookupIdent(s)
			return token.NewTokenWithValue(t, s, l.line, l.column)
		} else {
			tok = token.NewToken(token.ILLEGAL, l.ch, l.line, l.column)
		}
	}

	l.readChar() // Move to the next character
	return tok
}

// readChar advances to the next character in the input.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // End of input
	} else {
		l.ch = l.input[l.readPosition]
	}

	// update line and column number used in error management
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}

	l.position = l.readPosition
	l.readPosition++
}

// skipWhitespace skips over any whitespace characters in the input.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber reads a number (integer or floating point) from the input.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) || l.ch == '.' || l.ch == '-' || l.ch == '+' || l.ch == 'e' || l.ch == 'E' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isDigit checks if a character is a digit.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readString reads a string from the input, handling escaped quotes.
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// readIdentifier reads an identifier from the input.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// isLetter checks if a character is a letter or underscore.
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}
