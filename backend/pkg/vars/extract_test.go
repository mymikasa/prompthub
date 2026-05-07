package vars

import (
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"empty", "", nil},
		{"no vars", "hello world", nil},
		{"single var", "hello {{name}}", []string{"name"}},
		{"multiple vars", "{{greeting}} {{name}}, welcome to {{place}}", []string{"greeting", "name", "place"}},
		{"duplicate vars", "{{name}} {{name}} {{age}}", []string{"age", "name"}},
		{"underscore", "{{user_name}}", []string{"user_name"}},
		{"number", "{{var1}}", []string{"var1"}},
		{"starts with number ignored", "{{1invalid}}", nil},
		{"nested braces ignored", "{{{name}}}", []string{"name"}},
		{"adjacent", "{{a}}{{b}}", []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Extract(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Extract(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("Extract(%q)[%d] = %q, want %q", tt.input, i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestIsValidName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"simple", "name", true},
		{"underscore", "user_name", true},
		{"number", "var1", true},
		{"starts underscore", "_private", true},
		{"empty", "", false},
		{"too long", string(make([]byte, 101)), false},
		{"starts number", "1invalid", false},
		{"has space", "my name", false},
		{"has hyphen", "my-name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidName(tt.input); got != tt.valid {
				t.Errorf("IsValidName(%q) = %v, want %v", tt.input, got, tt.valid)
			}
		})
	}
}
