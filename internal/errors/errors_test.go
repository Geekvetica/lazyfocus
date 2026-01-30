package errors

import (
	"errors"
	"testing"
)

func TestOmniFocusError(t *testing.T) {
	tests := []struct {
		name               string
		message            string
		suggestion         string
		expectedExitCode   int
		expectedError      string
		expectedSuggestion string
	}{
		{
			name:               "OmniFocus not running",
			message:            "OmniFocus is not running",
			suggestion:         "Please launch OmniFocus and try again",
			expectedExitCode:   ExitOmniFocusError,
			expectedError:      "OmniFocus is not running",
			expectedSuggestion: "Please launch OmniFocus and try again",
		},
		{
			name:               "OmniFocus communication error",
			message:            "failed to communicate with OmniFocus",
			suggestion:         "Ensure OmniFocus is running and try again",
			expectedExitCode:   ExitOmniFocusError,
			expectedError:      "failed to communicate with OmniFocus",
			expectedSuggestion: "Ensure OmniFocus is running and try again",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewOmniFocusError(tt.message, tt.suggestion)

			// Test Error() method
			if err.Error() != tt.expectedError {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expectedError)
			}

			// Test ExitCode() method
			if err.ExitCode() != tt.expectedExitCode {
				t.Errorf("ExitCode() = %v, want %v", err.ExitCode(), tt.expectedExitCode)
			}

			// Test Suggestion() method
			if err.Suggestion() != tt.expectedSuggestion {
				t.Errorf("Suggestion() = %v, want %v", err.Suggestion(), tt.expectedSuggestion)
			}

			// Test that it implements LazyFocusError interface
			var _ LazyFocusError = err
			var _ error = err
		})
	}
}

func TestItemNotFoundError(t *testing.T) {
	tests := []struct {
		name               string
		itemType           string
		itemID             string
		expectedExitCode   int
		expectedError      string
		expectedSuggestion string
	}{
		{
			name:               "Task not found",
			itemType:           "task",
			itemID:             "abc123",
			expectedExitCode:   ExitItemNotFound,
			expectedError:      "task not found: abc123",
			expectedSuggestion: "Verify task ID using 'lazyfocus tasks'",
		},
		{
			name:               "Project not found",
			itemType:           "project",
			itemID:             "proj456",
			expectedExitCode:   ExitItemNotFound,
			expectedError:      "project not found: proj456",
			expectedSuggestion: "Check project name with 'lazyfocus projects'",
		},
		{
			name:               "Tag not found",
			itemType:           "tag",
			itemID:             "tag789",
			expectedExitCode:   ExitItemNotFound,
			expectedError:      "tag not found: tag789",
			expectedSuggestion: "Check tag name with 'lazyfocus tags'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewItemNotFoundError(tt.itemType, tt.itemID)

			// Test Error() method
			if err.Error() != tt.expectedError {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expectedError)
			}

			// Test ExitCode() method
			if err.ExitCode() != tt.expectedExitCode {
				t.Errorf("ExitCode() = %v, want %v", err.ExitCode(), tt.expectedExitCode)
			}

			// Test Suggestion() method
			if err.Suggestion() != tt.expectedSuggestion {
				t.Errorf("Suggestion() = %v, want %v", err.Suggestion(), tt.expectedSuggestion)
			}

			// Test that it implements LazyFocusError interface
			var _ LazyFocusError = err
			var _ error = err
		})
	}
}

func TestValidationError(t *testing.T) {
	tests := []struct {
		name               string
		message            string
		suggestion         string
		expectedExitCode   int
		expectedError      string
		expectedSuggestion string
	}{
		{
			name:               "Empty task name",
			message:            "task name is required",
			suggestion:         "provide a non-empty task name",
			expectedExitCode:   ExitValidationError,
			expectedError:      "task name is required",
			expectedSuggestion: "provide a non-empty task name",
		},
		{
			name:               "Task name too long",
			message:            "task name too long",
			suggestion:         "task name must be less than 500 characters",
			expectedExitCode:   ExitValidationError,
			expectedError:      "task name too long",
			expectedSuggestion: "task name must be less than 500 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.message, tt.suggestion)

			// Test Error() method
			if err.Error() != tt.expectedError {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expectedError)
			}

			// Test ExitCode() method
			if err.ExitCode() != tt.expectedExitCode {
				t.Errorf("ExitCode() = %v, want %v", err.ExitCode(), tt.expectedExitCode)
			}

			// Test Suggestion() method
			if err.Suggestion() != tt.expectedSuggestion {
				t.Errorf("Suggestion() = %v, want %v", err.Suggestion(), tt.expectedSuggestion)
			}

			// Test that it implements LazyFocusError interface
			var _ LazyFocusError = err
			var _ error = err
		})
	}
}

func TestDateParseError(t *testing.T) {
	tests := []struct {
		name               string
		dateStr            string
		reason             string
		expectedExitCode   int
		expectedError      string
		expectedSuggestion string
	}{
		{
			name:               "Invalid date format",
			dateStr:            "xyz",
			reason:             "unrecognized date format",
			expectedExitCode:   ExitValidationError,
			expectedError:      "invalid due date: unrecognized date format: xyz",
			expectedSuggestion: "Use relative (tomorrow), next (next monday), in (in 3 days), or ISO format",
		},
		{
			name:               "Invalid month",
			dateStr:            "2024-13-01",
			reason:             "month out of range",
			expectedExitCode:   ExitValidationError,
			expectedError:      "invalid due date: month out of range: 2024-13-01",
			expectedSuggestion: "Use relative (tomorrow), next (next monday), in (in 3 days), or ISO format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDateParseError(tt.dateStr, tt.reason)

			// Test Error() method
			if err.Error() != tt.expectedError {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expectedError)
			}

			// Test ExitCode() method
			if err.ExitCode() != tt.expectedExitCode {
				t.Errorf("ExitCode() = %v, want %v", err.ExitCode(), tt.expectedExitCode)
			}

			// Test Suggestion() method
			if err.Suggestion() != tt.expectedSuggestion {
				t.Errorf("Suggestion() = %v, want %v", err.Suggestion(), tt.expectedSuggestion)
			}

			// Test that it implements LazyFocusError interface
			var _ LazyFocusError = err
			var _ error = err
		})
	}
}

func TestPermissionError(t *testing.T) {
	tests := []struct {
		name               string
		message            string
		suggestion         string
		expectedExitCode   int
		expectedError      string
		expectedSuggestion string
	}{
		{
			name:               "Automation permission denied",
			message:            "automation permission denied",
			suggestion:         "Allow Terminal/iTerm access in System Preferences > Security > Automation",
			expectedExitCode:   ExitPermissionError,
			expectedError:      "automation permission denied",
			expectedSuggestion: "Allow Terminal/iTerm access in System Preferences > Security > Automation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewPermissionError(tt.message, tt.suggestion)

			// Test Error() method
			if err.Error() != tt.expectedError {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.expectedError)
			}

			// Test ExitCode() method
			if err.ExitCode() != tt.expectedExitCode {
				t.Errorf("ExitCode() = %v, want %v", err.ExitCode(), tt.expectedExitCode)
			}

			// Test Suggestion() method
			if err.Suggestion() != tt.expectedSuggestion {
				t.Errorf("Suggestion() = %v, want %v", err.Suggestion(), tt.expectedSuggestion)
			}

			// Test that it implements LazyFocusError interface
			var _ LazyFocusError = err
			var _ error = err
		})
	}
}

func TestErrorUnwrapping(t *testing.T) {
	// Test that our errors can be unwrapped and type-checked
	originalErr := errors.New("original error")

	omniFocusErr := NewOmniFocusError("OmniFocus error", "Launch OmniFocus")
	validationErr := NewValidationError("validation failed", "fix input")
	itemNotFoundErr := NewItemNotFoundError("task", "abc123")
	dateParseErr := NewDateParseError("bad-date", "invalid format")
	permissionErr := NewPermissionError("permission denied", "check settings")

	// Test that each error is of its specific type
	var lfErr LazyFocusError

	if !errors.As(omniFocusErr, &lfErr) {
		t.Error("OmniFocusError should implement LazyFocusError")
	}

	if !errors.As(validationErr, &lfErr) {
		t.Error("ValidationError should implement LazyFocusError")
	}

	if !errors.As(itemNotFoundErr, &lfErr) {
		t.Error("ItemNotFoundError should implement LazyFocusError")
	}

	if !errors.As(dateParseErr, &lfErr) {
		t.Error("DateParseError should implement LazyFocusError")
	}

	if !errors.As(permissionErr, &lfErr) {
		t.Error("PermissionError should implement LazyFocusError")
	}

	// Test that non-LazyFocus errors don't match
	if errors.As(originalErr, &lfErr) {
		t.Error("Regular errors should not match LazyFocusError")
	}
}
