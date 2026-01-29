package quickadd

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// TestInitialState verifies the quick add component starts in hidden state
func TestInitialState(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)

	if model.IsVisible() {
		t.Error("Expected quick add to be hidden initially")
	}

	if model.textInput.Value() != "" {
		t.Error("Expected input to be empty initially")
	}

	if model.err != nil {
		t.Error("Expected no error initially")
	}
}

// TestShowHide verifies Show/Hide functionality
func TestShowHide(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)

	// Show the component
	model = model.Show()
	if !model.IsVisible() {
		t.Error("Expected quick add to be visible after Show()")
	}

	// Hide the component
	model = model.Hide()
	if model.IsVisible() {
		t.Error("Expected quick add to be hidden after Hide()")
	}

	// Verify input is cleared on hide
	model = model.Show()
	model.textInput.SetValue("test input")
	model = model.Hide()
	if model.textInput.Value() != "" {
		t.Error("Expected input to be cleared after Hide()")
	}
}

// TestTextInput verifies text input functionality
func TestTextInput(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)
	model = model.Show()

	// Simulate typing
	testInput := "Buy milk"
	for _, ch := range testInput {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{ch}}
		model, _ = model.Update(msg)
	}

	if model.textInput.Value() != testInput {
		t.Errorf("Expected input '%s', got '%s'", testInput, model.textInput.Value())
	}
}

// TestEscapeCancel verifies Escape key cancels without submitting
func TestEscapeCancel(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)
	model = model.Show()
	model.textInput.SetValue("test task")

	// Press Escape
	escapeMsg := tea.KeyMsg{Type: tea.KeyEsc}
	model, _ = model.Update(escapeMsg)

	// Should be hidden and input cleared
	if model.IsVisible() {
		t.Error("Expected quick add to be hidden after Escape")
	}

	if model.textInput.Value() != "" {
		t.Error("Expected input to be cleared after Escape")
	}
}

// TestEnterSubmitsTask verifies Enter key submits and parses correctly
func TestEnterSubmitsTask(t *testing.T) {
	styles := tui.DefaultStyles()

	testTime := time.Date(2024, 1, 15, 17, 0, 0, 0, time.UTC)
	expectedTask := &domain.Task{
		ID:      "task-123",
		Name:    "Buy milk",
		Tags:    []string{"groceries"},
		DueDate: &testTime,
		Flagged: false,
	}

	mockSvc := &service.MockOmniFocusService{
		CreatedTask: expectedTask,
	}

	model := New(styles, mockSvc)
	model = model.Show()
	model.textInput.SetValue("Buy milk #groceries due:2024-01-15")

	// Press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := model.Update(enterMsg)

	// Should be hidden after successful submit
	if model.IsVisible() {
		t.Error("Expected quick add to be hidden after successful submit")
	}

	// Should have cleared input
	if model.textInput.Value() != "" {
		t.Error("Expected input to be cleared after submit")
	}

	// Should return TaskCreatedMsg
	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	// Execute the command to get the message
	msg := cmd()
	taskCreatedMsg, ok := msg.(tui.TaskCreatedMsg)
	if !ok {
		t.Fatalf("Expected TaskCreatedMsg, got %T", msg)
	}

	if taskCreatedMsg.Task.ID != expectedTask.ID {
		t.Errorf("Expected task ID '%s', got '%s'", expectedTask.ID, taskCreatedMsg.Task.ID)
	}

	if taskCreatedMsg.Task.Name != expectedTask.Name {
		t.Errorf("Expected task name '%s', got '%s'", expectedTask.Name, taskCreatedMsg.Task.Name)
	}
}

// TestEnterWithError verifies error handling on task creation failure
func TestEnterWithError(t *testing.T) {
	styles := tui.DefaultStyles()

	mockSvc := &service.MockOmniFocusService{
		CreateTaskErr: errors.New("project not found"),
	}

	model := New(styles, mockSvc)
	model = model.Show()
	model.textInput.SetValue("Task @InvalidProject")

	// Press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := model.Update(enterMsg)

	// Should remain visible on error
	if !model.IsVisible() {
		t.Error("Expected quick add to remain visible on error")
	}

	// Should have error set
	if model.err == nil {
		t.Error("Expected error to be set")
	}

	// Should return ErrorMsg
	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	errorMsg, ok := msg.(tui.ErrorMsg)
	if !ok {
		t.Fatalf("Expected ErrorMsg, got %T", msg)
	}

	if errorMsg.Err == nil {
		t.Error("Expected error in ErrorMsg")
	}
}

// TestEmptyInputValidation verifies empty input is rejected
func TestEmptyInputValidation(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)
	model = model.Show()

	// Press Enter with empty input
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, cmd := model.Update(enterMsg)

	// Should remain visible
	if !model.IsVisible() {
		t.Error("Expected quick add to remain visible with empty input")
	}

	// Should have error set
	if model.err == nil {
		t.Error("Expected error for empty input")
	}

	// Should return ErrorMsg
	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	if _, ok := msg.(tui.ErrorMsg); !ok {
		t.Fatalf("Expected ErrorMsg, got %T", msg)
	}
}

// TestSetSize verifies size setting functionality
func TestSetSize(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)

	width := 100
	height := 50
	model = model.SetSize(width, height)

	if model.width != width {
		t.Errorf("Expected width %d, got %d", width, model.width)
	}

	if model.height != height {
		t.Errorf("Expected height %d, got %d", height, model.height)
	}
}

// TestNaturalSyntaxParsing verifies various natural syntax patterns
func TestNaturalSyntaxParsing(t *testing.T) {
	styles := tui.DefaultStyles()

	testCases := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "simple task",
			input:       "Buy groceries",
			expectError: false,
		},
		{
			name:        "task with tag",
			input:       "Call mom #personal",
			expectError: false,
		},
		{
			name:        "task with project",
			input:       "Review PR @Work",
			expectError: false,
		},
		{
			name:        "task with due date",
			input:       "Submit report due:tomorrow",
			expectError: false,
		},
		{
			name:        "task with flagged",
			input:       "Urgent meeting !",
			expectError: false,
		},
		{
			name:        "full syntax",
			input:       "Team sync @Work #meeting due:tomorrow !",
			expectError: false,
		},
		{
			name:        "only modifiers (should fail)",
			input:       "#tag @project due:tomorrow !",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &service.MockOmniFocusService{
				CreatedTask: &domain.Task{
					ID:   "test-id",
					Name: "test task",
				},
			}

			if tc.expectError {
				mockSvc.CreateTaskErr = errors.New("task name is required")
			}

			model := New(styles, mockSvc)
			model = model.Show()
			model.textInput.SetValue(tc.input)

			enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
			model, _ = model.Update(enterMsg)

			if tc.expectError {
				if model.err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if model.err != nil {
					t.Errorf("Expected no error but got: %v", model.err)
				}
			}
		})
	}
}

// TestViewRendering verifies the View output contains expected elements
func TestViewRendering(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{}

	model := New(styles, mockSvc)

	// Hidden state should return empty string
	view := model.View()
	if view != "" {
		t.Error("Expected empty view when hidden")
	}

	// Visible state should show content
	model = model.Show()
	model = model.SetSize(80, 40)
	view = model.View()

	if view == "" {
		t.Error("Expected non-empty view when visible")
	}

	// Should contain placeholder text hint
	if !contains(view, "Add task") {
		t.Error("Expected view to contain 'Add task' hint")
	}

	// Should contain help text
	if !contains(view, "Enter") {
		t.Error("Expected view to contain Enter key hint")
	}

	if !contains(view, "Escape") {
		t.Error("Expected view to contain Escape key hint")
	}
}

// TestErrorDisplay verifies error messages are shown in the view
func TestErrorDisplay(t *testing.T) {
	styles := tui.DefaultStyles()
	mockSvc := &service.MockOmniFocusService{
		CreateTaskErr: errors.New("project not found"),
	}

	model := New(styles, mockSvc)
	model = model.Show()
	model = model.SetSize(80, 40)
	model.textInput.SetValue("Task @InvalidProject")

	// Trigger error
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	model, _ = model.Update(enterMsg)

	view := model.View()

	// Should contain error message
	if !contains(view, "project not found") {
		t.Error("Expected view to contain error message")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
