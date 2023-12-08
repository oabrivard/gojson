// Package parser defines the structure and methods for parsing JSON.
package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/oabrivard/gojson/lexer"
	"github.com/oabrivard/gojson/token"
)

// Parser struct represents a parser with a lexer, current and peek tokens,
// and a slice to store parsing errors.
type Parser struct {
	lexer *lexer.Lexer // the lexer from which the parser receives tokens

	curToken  token.Token // current token under examination
	peekToken token.Token // next token in the input

	errors []string // slice to store errors encountered during parsing
}

// NewParser creates and initializes a new Parser with the given lexer.
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	// Initialize curToken and peekToken
	p.nextToken()
	p.nextToken()
	return p
}

// nextToken advances both curToken and peekToken.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// JsonObject and JsonArray are types to represent JSON objects and arrays, respectively.
type JsonObject map[string]interface{}
type JsonArray []interface{}

// Parse starts the parsing process and returns the top-level JSON object.
func (p *Parser) Parse() JsonObject {
	return p.parseObject()
}

// parseObject parses a JSON object from the token stream.
func (p *Parser) parseObject() JsonObject {
	object := make(JsonObject)

	// Ensure the current token is the beginning of an object
	if !p.curTokenIs(token.BEGIN_OBJECT) {
		p.addError(fmt.Sprintf("expected '{' at line %d, column %d, got '%s'", p.curToken.Line, p.curToken.Column, p.curToken.Value))
		return nil
	}

	// Move to the next token
	p.nextToken()

	// Loop until the end of the object is reached
	for !p.curTokenIs(token.END_OBJECT) && !p.curTokenIs(token.EOF) {
		key := p.parseObjectKey()
		if key == "" {
			return nil
		}

		// Ensure a name separator (:) follows the key
		if !p.expectPeek(token.NAME_SEPARATOR) {
			return nil
		}

		// Move to the value token
		p.nextToken()

		// Parse the value
		value, err := p.parseValue()
		if err != nil {
			return nil
		}

		object[key] = value

		// Move past the value
		p.nextToken()

		// Handle comma separation for multiple key-value pairs
		if p.curTokenIs(token.VALUE_SEPARATOR) {
			if p.peekToken.Type == token.END_OBJECT { // No comma just before the end of the object
				p.addError(fmt.Sprintf("No ',' before '}' at line %d, column %d", p.curToken.Line, p.curToken.Column))
				return nil
			}

			p.nextToken()
		}
	}

	// Ensure the end of the object is reached
	if !p.curTokenIs(token.END_OBJECT) {
		p.addError(fmt.Sprintf("expected '}' at line %d, column %d, got '%s'", p.curToken.Line, p.curToken.Column, p.curToken.Value))
		return nil
	}

	return object
}

// parseArray parses a JSON array from the token stream.
func (p *Parser) parseArray() JsonArray {
	array := JsonArray{}

	// Ensure the current token is the beginning of an array
	if !p.curTokenIs(token.BEGIN_ARRAY) {
		p.addError(fmt.Sprintf("expected '[' at line %d, column %d, got '%s'", p.curToken.Line, p.curToken.Column, p.curToken.Value))
		return nil
	}

	// Move to the next token
	p.nextToken()

	// Loop until the end of the array is reached
	for !p.curTokenIs(token.END_ARRAY) {
		// Parse the value
		value, err := p.parseValue()
		if err != nil {
			return nil
		}

		array = append(array, value)

		// Move past the value
		p.nextToken()

		// Handle comma separation for multiple values
		if p.curTokenIs(token.VALUE_SEPARATOR) {
			p.nextToken()
		}
	}

	// Ensure the end of the array is reached
	if !p.curTokenIs(token.END_ARRAY) {
		return nil
	}

	return array
}

// addError appends an error message to the parser's errors slice.
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// parseObjectKey parses and returns the key of an object field.
func (p *Parser) parseObjectKey() string {
	if p.curToken.Type != token.STRING {
		p.addError(fmt.Sprintf("expected string for key at line %d, column %d, got '%s'", p.curToken.Line, p.curToken.Column, p.curToken.Value))
		return ""
	}
	return p.curToken.Value
}

// parseValue parses a JSON value based on the current token type.
func (p *Parser) parseValue() (interface{}, error) {
	switch p.curToken.Type {
	case token.STRING:
		return p.curToken.Value, nil
	case token.NUMBER:
		return p.parseNumber(), nil
	case token.TRUE, token.FALSE:
		return p.parseBoolean(), nil
	case token.NULL:
		return nil, nil
	case token.BEGIN_OBJECT:
		return p.parseObject(), nil
	case token.BEGIN_ARRAY:
		return p.parseArray(), nil
	default:
		p.addError(fmt.Sprintf("unexpected token '%s' at line %d, column %d", p.curToken.Value, p.curToken.Line, p.curToken.Column))
		return nil, errors.New("unexpected token")
	}
}

// parseNumber parses a number token into an appropriate Go numeric type.
func (p *Parser) parseNumber() interface{} {
	numStr := p.curToken.Value

	// Check for float or integer representation
	if strings.Contains(numStr, ".") || strings.ContainsAny(numStr, "eE") {
		// Parse as float
		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			p.addError(fmt.Sprintf("could not parse %q as float at line %d, column %d", numStr, p.curToken.Line, p.curToken.Column))
			return nil
		}
		return val
	}

	// Parse as integer
	val, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		p.addError(fmt.Sprintf("could not parse %q as integer at line %d, column %d", numStr, p.curToken.Line, p.curToken.Column))
		return nil
	}
	return val
}

// parseBoolean returns a boolean value based on the current token.
func (p *Parser) parseBoolean() bool {
	return p.curToken.Type == token.TRUE
}

// expectPeek checks if the next token is of the expected type.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.addError(fmt.Sprintf("expected next token to be %v, got %v instead, at line %d, column %d", t, p.peekToken.Type, p.curToken.Line, p.curToken.Column))
		return false
	}
}

func (p *Parser) Errors() []string {
	return p.errors
}

// curTokenIs checks if the current token is of a specific type.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
