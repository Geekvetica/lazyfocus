package validation

import (
	"errors"
	"strings"
	"testing"

	lferrors "github.com/pwojciechowski/lazyfocus/internal/errors"
)

func TestValidateTaskName(t *testing.T) {
	tests := []struct {
		name        string
		taskName    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid task name",
			taskName:    "Buy groceries",
			expectError: false,
		},
		{
			name:        "Valid task name with special characters",
			taskName:    "Review PR #142 @work",
			expectError: false,
		},
		{
			name:        "Valid task name with newlines",
			taskName:    "Multi-line\ntask name",
			expectError: false,
		},
		{
			name:        "Valid task name with tabs",
			taskName:    "Task with\ttabs",
			expectError: false,
		},
		{
			name:        "Empty task name",
			taskName:    "",
			expectError: true,
			errorMsg:    "task name is required",
		},
		{
			name:        "Whitespace-only task name",
			taskName:    "   \t\n   ",
			expectError: true,
			errorMsg:    "task name is required",
		},
		{
			name:        "Task name at max length",
			taskName:    strings.Repeat("a", MaxTaskNameLength),
			expectError: false,
		},
		{
			name:        "Task name exceeds max length",
			taskName:    strings.Repeat("a", MaxTaskNameLength+1),
			expectError: true,
			errorMsg:    "task name too long",
		},
		{
			name:        "Task name with control characters",
			taskName:    "Task with \x00 null character",
			expectError: true,
			errorMsg:    "task name contains invalid characters",
		},
		{
			name:        "Task name with bell character",
			taskName:    "Task with \a bell",
			expectError: true,
			errorMsg:    "task name contains invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskName(tt.taskName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}

				// Verify it's a LazyFocusError
				var lfErr lferrors.LazyFocusError
				if !errors.As(err, &lfErr) {
					t.Errorf("Expected LazyFocusError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateNote(t *testing.T) {
	tests := []struct {
		name        string
		note        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty note",
			note:        "",
			expectError: false,
		},
		{
			name:        "Valid short note",
			note:        "This is a short note",
			expectError: false,
		},
		{
			name:        "Valid long note",
			note:        strings.Repeat("a", 5000),
			expectError: false,
		},
		{
			name:        "Note at max length",
			note:        strings.Repeat("a", MaxNoteLength),
			expectError: false,
		},
		{
			name:        "Note exceeds max length",
			note:        strings.Repeat("a", MaxNoteLength+1),
			expectError: true,
			errorMsg:    "note too long",
		},
		{
			name:        "Note with special characters",
			note:        "Note with\nnewlines\nand\ttabs",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNote(tt.note)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}

				// Verify it's a LazyFocusError
				var lfErr lferrors.LazyFocusError
				if !errors.As(err, &lfErr) {
					t.Errorf("Expected LazyFocusError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty project name (optional)",
			projectName: "",
			expectError: false,
		},
		{
			name:        "Valid project name",
			projectName: "Work",
			expectError: false,
		},
		{
			name:        "Valid project name with spaces",
			projectName: "Big Project Alpha",
			expectError: false,
		},
		{
			name:        "Project name at max length",
			projectName: strings.Repeat("a", MaxProjectNameLength),
			expectError: false,
		},
		{
			name:        "Project name exceeds max length",
			projectName: strings.Repeat("a", MaxProjectNameLength+1),
			expectError: true,
			errorMsg:    "project name too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.projectName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}

				// Verify it's a LazyFocusError
				var lfErr lferrors.LazyFocusError
				if !errors.As(err, &lfErr) {
					t.Errorf("Expected LazyFocusError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestValidateTagName(t *testing.T) {
	tests := []struct {
		name        string
		tagName     string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Empty tag name",
			tagName:     "",
			expectError: true,
			errorMsg:    "tag name is required",
		},
		{
			name:        "Valid tag name",
			tagName:     "urgent",
			expectError: false,
		},
		{
			name:        "Valid tag name with spaces",
			tagName:     "work from home",
			expectError: false,
		},
		{
			name:        "Tag name at max length",
			tagName:     strings.Repeat("a", MaxTagNameLength),
			expectError: false,
		},
		{
			name:        "Tag name exceeds max length",
			tagName:     strings.Repeat("a", MaxTagNameLength+1),
			expectError: true,
			errorMsg:    "tag name too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTagName(tt.tagName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain %q, got %q", tt.errorMsg, err.Error())
				}

				// Verify it's a LazyFocusError
				var lfErr lferrors.LazyFocusError
				if !errors.As(err, &lfErr) {
					t.Errorf("Expected LazyFocusError, got %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestContainsControlChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "No control characters",
			input:    "Normal text",
			expected: false,
		},
		{
			name:     "Newline is allowed",
			input:    "Text with\nnewline",
			expected: false,
		},
		{
			name:     "Tab is allowed",
			input:    "Text with\ttab",
			expected: false,
		},
		{
			name:     "Null character",
			input:    "Text with \x00 null",
			expected: true,
		},
		{
			name:     "Bell character",
			input:    "Text with \a bell",
			expected: true,
		},
		{
			name:     "Backspace character",
			input:    "Text with \b backspace",
			expected: true,
		},
		{
			name:     "Escape character",
			input:    "Text with \x1b escape",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsControlChars(tt.input)
			if result != tt.expected {
				t.Errorf("containsControlChars(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
