package parser

import (
	"reflect"
	"testing"

	"github.com/oabrivard/gojson/lexer"
)

func TestParseSimpleObject(t *testing.T) {
	input := `{"name": "John", "age": 30, "isStudent": false}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{
		"name":      "John",
		"age":       int64(30), // Assuming numbers are parsed as float64
		"isStudent": false,
	}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep1Valid(t *testing.T) {
	input := `{}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep1Invalid(t *testing.T) {
	input := ``

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 1 || p.errors[0] != "expected '{' at line 1, column 1, got ''" {
		t.Errorf("Not the expected error(s) during parsing, got %v", p.errors)
	}

	if parsed != nil {
		t.Errorf("expected a nil result from parsing an empty input")
	}
}

func TestParseStep2Valid1(t *testing.T) {
	input := `{"key": "value"}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{
		"key": "value",
	}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep2Valid2(t *testing.T) {
	input := `{
		"key": "value",
		"key2": "value"
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{
		"key":  "value",
		"key2": "value",
	}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep2Invalid1(t *testing.T) {
	input := `{"key": "value",}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 1 || p.errors[0] != "No ',' before '}' at line 1, column 16" {
		t.Errorf("Not the expected error(s) during parsing, got %v", p.errors)
	}

	if parsed != nil {
		t.Errorf("expected a nil result from parsing an empty input")
	}
}

func TestParseStep2Invalid2(t *testing.T) {
	input := `{
		"key": "value",
		key2: "value"
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 1 || p.errors[0] != "expected string for key at line 3, column 6, got 'key'" {
		t.Errorf("Not the expected error(s) during parsing, got %v", p.errors)
	}

	if parsed != nil {
		t.Errorf("expected a nil result from parsing an empty input")
	}
}

func TestParseStep3Valid(t *testing.T) {
	input := `{
		"key1": true,
		"key2": false,
		"key3": null,
		"key4": "value",
		"key5": 101
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{
		"key1": true,
		"key2": false,
		"key3": nil,
		"key4": "value",
		"key5": int64(101),
	}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep3Invalid(t *testing.T) {
	input := `{
		"key1": true,
		"key2": False,
		"key3": null,
		"key4": "value",
		"key5": 101
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 1 || p.errors[0] != "unexpected token 'False' at line 3, column 16" {
		t.Errorf("Not the expected error(s) during parsing, got %v", p.errors)
	}

	if parsed != nil {
		t.Errorf("expected a nil result from parsing an empty input")
	}
}

func TestParseStep4Valid1(t *testing.T) {
	input := `{
		"key": "value",
		"key-n": 101,
		"key-o": {},
		"key-l": []
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{
		"key":   "value",
		"key-n": int64(101),
		"key-o": JsonObject{},
		"key-l": JsonArray{},
	}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep4Valid2(t *testing.T) {
	input := `{
		"key": "value",
		"key-n": 101,
		"key-o": {
			"inner key": "inner value"
		},
		"key-l": ["list value"]
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 0 {
		errMsg := ""
		for _, s := range p.errors {
			errMsg += s + "\n"
		}
		t.Fatalf(errMsg)
	}

	expected := JsonObject{
		"key":   "value",
		"key-n": int64(101),
		"key-o": JsonObject{
			"inner key": "inner value",
		},
		"key-l": JsonArray{"list value"},
	}

	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("parsed object is not as expected. Got %+v, want %+v", parsed, expected)
	}
}

func TestParseStep4Invalid(t *testing.T) {
	input := `{
		"key": "value",
		"key-n": 101,
		"key-o": {
			"inner key": "inner value"
		},
		"key-l": ['list value']
	}`

	l := lexer.NewLexer(input)
	p := NewParser(l)
	parsed := p.Parse()

	if len(p.errors) != 2 || p.errors[0] != "unexpected token ''' at line 7, column 13" || p.errors[1] != "expected string for key at line 7, column 18, got 'list'" {
		t.Errorf("Not the expected error(s) during parsing, got %v", p.errors)
	}

	if parsed != nil {
		t.Errorf("expected a nil result from parsing an empty input")
	}
}
