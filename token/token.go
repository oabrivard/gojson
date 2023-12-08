package token

type TokenType int

const (
	// Special tokens
	EOF     TokenType = iota // Represents the end of the file/input
	ILLEGAL                  // Represents any character or sequence of characters that doesn't form a valid token in JSON

	// Symbols and structure tokens
	BEGIN_ARRAY     // [
	END_ARRAY       // ]
	BEGIN_OBJECT    // {
	END_OBJECT      // }
	NAME_SEPARATOR  // :
	VALUE_SEPARATOR // ,

	// Whitespace
	WHITESPACE // Represents whitespace (spaces, tabs, line feeds, carriage returns).

	// Literal types
	STRING // Represents a string literal
	NUMBER // Represents a number
	TRUE   // Represents the boolean value "true"
	FALSE  // Represents the boolean value "false"
	NULL   // Represents the "null" value
)

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func NewToken(tokenType TokenType, ch byte, l int, c int) Token {
	return Token{Type: tokenType, Value: string(ch), Line: l, Column: c}
}

func NewTokenWithValue(tokenType TokenType, val string, l int, c int) Token {
	return Token{Type: tokenType, Value: val, Line: l, Column: c}
}

var keywords = map[string]TokenType{
	"true":  TRUE,
	"false": FALSE,
	"null":  NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return ILLEGAL
}
