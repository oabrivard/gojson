package linter

import (
	"testing"
)

func TestLintSimpleObject(t *testing.T) {
	input := `{"name": "John", "age": 30, "isStudent": false}`

	jl := NewJsonLinter(input)
	linted, err := jl.Lint()

	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := "{\n  \"name\": \"John\",\n  \"age\": 30,\n  \"isStudent\": false\n}"

	if linted != expected {
		t.Errorf("linted object is not as expected. Got %+v, want %+v", linted, expected)
	}
}

func TestLintComplexObject(t *testing.T) {
	input := `{
		"key": "value",
		"key-n": 101,
		"key-o": {
			"inner key": "inner value"
		},
		"key-l": ["list value"]
	}`

	jl := NewJsonLinter(input)
	linted, err := jl.Lint()

	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := "{\n  \"key\": \"value\",\n  \"key-n\": 101,\n  \"key-o\": {\n    \"inner key\": \"inner value\"\n  },\n  \"key-l\": [\n    \"list value\"\n  ]\n}"

	if linted != expected {
		t.Errorf("linted object is not as expected. Got %+v, want %+v", linted, expected)
	}
}

func TestLintInvalidJson(t *testing.T) {
	input := `{
		"key": "value",
		"key-n": 101,
		"key-o": {
			"inner key": "inner value"
		},
		"key-l": ['list value']
	}`

	jl := NewJsonLinter(input)
	_, err := jl.Lint()

	if err == nil {
		t.Errorf("Expected error(s) during linting")
	}
}
