package lexer

import (
	"testing"

	"github.com/oabrivard/gojson/token"
)

func TestTokenizeSimpleObject(t *testing.T) {
	input := `{"name": "John", "age": 30, "value": -3.5e+5}`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.BEGIN_OBJECT, "{"},
		{token.STRING, "name"},
		{token.NAME_SEPARATOR, ":"},
		{token.STRING, "John"},
		{token.VALUE_SEPARATOR, ","},
		{token.STRING, "age"},
		{token.NAME_SEPARATOR, ":"},
		{token.NUMBER, "30"},
		{token.VALUE_SEPARATOR, ","},
		{token.STRING, "value"},
		{token.NAME_SEPARATOR, ":"},
		{token.NUMBER, "-3.5e+5"},
		{token.END_OBJECT, "}"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Value)
		}
	}
}
