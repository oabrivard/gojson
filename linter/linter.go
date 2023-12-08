// Package linter provides functionality for linting JSON strings.
package linter

import (
	"fmt"
	"strings"

	"github.com/oabrivard/gojson/lexer"
	"github.com/oabrivard/gojson/parser"
)

// JsonLinter struct holds references to a lexer and a parser for JSON linting.
type JsonLinter struct {
	lexer  *lexer.Lexer   // The lexer to tokenize the input
	parser *parser.Parser // The parser to parse the tokenized input
}

// NewJsonLinter creates and initializes a new JsonLinter with the given input string.
func NewJsonLinter(input string) *JsonLinter {
	l := lexer.NewLexer(input)
	p := parser.NewParser(l)
	return &JsonLinter{lexer: l, parser: p}
}

// Lint performs the linting process on the input JSON.
// It parses the input and then formats it into a nicely structured JSON string.
func (jl *JsonLinter) Lint() (string, error) {
	parsedObject := jl.parser.Parse()

	// If parsing errors are present, return an aggregated error message.
	if len(jl.parser.Errors()) > 0 {
		return "", fmt.Errorf("parsing errors: %v", jl.parser.Errors())
	}

	// Use the custom formatJSON function to format the parsed JSON object.
	formattedJson := formatJSON(parsedObject, "")
	return string(formattedJson), nil
}

// formatJSON formats any JSON value into a nicely indented string.
func formatJSON(obj interface{}, indent string) string {
	// Type switch to handle different types of JSON values.
	switch v := obj.(type) {
	case parser.JsonObject:
		return formatObject(v, indent) // Format a JSON object
	case parser.JsonArray:
		return formatArray(v, indent) // Format a JSON array
	case string:
		return fmt.Sprintf("\"%s\"", v) // Format a JSON string
	case nil:
		return "null" // Format a JSON null
	case bool:
		if v {
			return "true"
		}
		return "false"
	default: // For numbers and other types, use default formatting
		return fmt.Sprintf("%v", v)
	}
}

// formatObject formats a JSON object into a string with proper indentation.
func formatObject(obj map[string]interface{}, indent string) string {
	var result strings.Builder
	result.WriteString("{\n")
	i := 0
	for k, v := range obj {
		// Format each key-value pair in the object.
		result.WriteString(indent + "  \"" + k + "\": " + formatJSON(v, indent+"  "))
		if i < len(obj)-1 {
			result.WriteString(",")
		}
		result.WriteString("\n")
		i++
	}
	result.WriteString(indent + "}")
	return result.String()
}

// formatArray formats a JSON array into a string with proper indentation.
func formatArray(array []interface{}, indent string) string {
	var result strings.Builder
	result.WriteString("[\n")
	for i, v := range array {
		// Format each value in the array.
		result.WriteString(indent + "  " + formatJSON(v, indent+"  "))
		if i < len(array)-1 {
			result.WriteString(",")
		}
		result.WriteString("\n")
	}
	result.WriteString(indent + "]")
	return result.String()
}
