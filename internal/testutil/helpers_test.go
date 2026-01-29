package testutil

import "testing"

func TestStringPtr(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "empty string", input: ""},
		{name: "simple string", input: "hello"},
		{name: "string with spaces", input: "hello world"},
		{name: "unicode string", input: "hello 世界"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := StringPtr(tt.input)

			if ptr == nil {
				t.Fatal("StringPtr() returned nil")
			}

			if *ptr != tt.input {
				t.Errorf("StringPtr() = %q, want %q", *ptr, tt.input)
			}

			// Verify it's a new pointer (not pointing to original)
			original := tt.input
			*ptr = "modified"
			if original != tt.input {
				t.Error("StringPtr() should return pointer to new memory")
			}
		})
	}
}

func TestBoolPtr(t *testing.T) {
	tests := []struct {
		name  string
		input bool
	}{
		{name: "true", input: true},
		{name: "false", input: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ptr := BoolPtr(tt.input)

			if ptr == nil {
				t.Fatal("BoolPtr() returned nil")
			}

			if *ptr != tt.input {
				t.Errorf("BoolPtr() = %v, want %v", *ptr, tt.input)
			}

			// Verify it's a new pointer (not pointing to original)
			original := tt.input
			*ptr = !*ptr
			if original != tt.input {
				t.Error("BoolPtr() should return pointer to new memory")
			}
		})
	}
}
