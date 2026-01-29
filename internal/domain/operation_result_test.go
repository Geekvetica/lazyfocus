package domain

import "testing"

func TestNewSuccessResult(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		message string
	}{
		{
			name:    "creates success result with ID and message",
			id:      "task-123",
			message: "Task created successfully",
		},
		{
			name:    "creates success result with empty ID",
			id:      "",
			message: "Operation completed",
		},
		{
			name:    "creates success result with empty message",
			id:      "task-456",
			message: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSuccessResult(tt.id, tt.message)

			if !result.Success {
				t.Errorf("NewSuccessResult() Success = %v, want %v", result.Success, true)
			}

			if result.ID != tt.id {
				t.Errorf("NewSuccessResult() ID = %v, want %v", result.ID, tt.id)
			}

			if result.Message != tt.message {
				t.Errorf("NewSuccessResult() Message = %v, want %v", result.Message, tt.message)
			}
		})
	}
}

func TestNewErrorResult(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "creates error result with message",
			message: "Failed to create task",
		},
		{
			name:    "creates error result with empty message",
			message: "",
		},
		{
			name:    "creates error result with detailed message",
			message: "Failed to create task: OmniFocus is not running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewErrorResult(tt.message)

			if result.Success {
				t.Errorf("NewErrorResult() Success = %v, want %v", result.Success, false)
			}

			if result.ID != "" {
				t.Errorf("NewErrorResult() ID = %v, want %v", result.ID, "")
			}

			if result.Message != tt.message {
				t.Errorf("NewErrorResult() Message = %v, want %v", result.Message, tt.message)
			}
		})
	}
}
