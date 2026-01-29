package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestPerspectiveCommand_ShowTasks(t *testing.T) {
	// Test showing tasks from a perspective
	mockService := &service.MockOmniFocusService{
		PerspectiveTasks: []domain.Task{
			{ID: "task1", Name: "Review presentation", Flagged: true},
			{ID: "task2", Name: "Send emails", DueDate: timePtr(time.Now())},
		},
	}

	output, exitCode, err := executePerspectiveCommand(mockService, []string{"Forecast"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Review presentation") {
		t.Errorf("Expected output to contain 'Review presentation', got: %s", output)
	}

	if !strings.Contains(output, "Send emails") {
		t.Errorf("Expected output to contain 'Send emails', got: %s", output)
	}
}

func TestPerspectiveCommand_EmptyResults(t *testing.T) {
	// Test perspective with no tasks
	mockService := &service.MockOmniFocusService{
		PerspectiveTasks: []domain.Task{},
	}

	output, exitCode, err := executePerspectiveCommand(mockService, []string{"Today"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "No tasks") {
		t.Errorf("Expected output to indicate no tasks, got: %s", output)
	}
}

func TestPerspectiveCommand_NotFound(t *testing.T) {
	// Test perspective not found error
	mockService := &service.MockOmniFocusService{
		PerspectiveTasksErr: errors.New("perspective not found: NonExistent"),
	}

	_, exitCode, err := executePerspectiveCommand(mockService, []string{"NonExistent"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "perspective not found") {
		t.Errorf("Expected error message about perspective not found, got: %v", err)
	}
}

func TestPerspectiveCommand_ProRequired(t *testing.T) {
	// Test OmniFocus Pro required error
	mockService := &service.MockOmniFocusService{
		PerspectiveTasksErr: errors.New("OmniFocus Pro required for custom perspectives"),
	}

	_, exitCode, err := executePerspectiveCommand(mockService, []string{"CustomView"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "OmniFocus Pro required") {
		t.Errorf("Expected error message about OmniFocus Pro, got: %v", err)
	}
}

func TestPerspectiveCommand_OmniFocusNotRunning(t *testing.T) {
	// Test OmniFocus not running error
	mockService := &service.MockOmniFocusService{
		PerspectiveTasksErr: errors.New("OmniFocus is not running"),
	}

	_, exitCode, err := executePerspectiveCommand(mockService, []string{"Forecast"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "OmniFocus is not running") {
		t.Errorf("Expected error message about OmniFocus not running, got: %v", err)
	}
}

func TestPerspectiveCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	mockService := &service.MockOmniFocusService{
		PerspectiveTasks: []domain.Task{
			{ID: "task1", Name: "Perspective task", Flagged: false},
		},
	}

	output, exitCode, err := executePerspectiveCommand(mockService, []string{"--json", "Forecast"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Check for JSON structure
	if !strings.Contains(output, `"tasks"`) {
		t.Errorf("Expected JSON output to contain 'tasks' field, got: %s", output)
	}

	if !strings.Contains(output, `"Perspective task"`) {
		t.Errorf("Expected JSON output to contain task name, got: %s", output)
	}

	if !strings.Contains(output, `"count"`) {
		t.Errorf("Expected JSON output to contain 'count' field, got: %s", output)
	}
}

func TestPerspectiveCommand_ErrorJSON(t *testing.T) {
	// Test error handling in JSON mode
	mockService := &service.MockOmniFocusService{
		PerspectiveTasksErr: errors.New("perspective not found"),
	}

	output, exitCode, err := executePerspectiveCommand(mockService, []string{"--json", "NonExistent"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	// In JSON mode, error should be in output
	if !strings.Contains(output, `"error"`) {
		t.Errorf("Expected JSON error output to contain 'error' field, got: %s", output)
	}
}

func TestPerspectiveCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	mockService := &service.MockOmniFocusService{
		PerspectiveTasks: []domain.Task{
			{ID: "task1", Name: "Test task", Flagged: false},
		},
	}

	output, exitCode, err := executePerspectiveCommand(mockService, []string{"--quiet", "Forecast"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if output != "" {
		t.Errorf("Expected empty output in quiet mode, got: %s", output)
	}
}

func TestPerspectiveCommand_WithProjectAndTags(t *testing.T) {
	// Test that tasks show project and tags information
	mockService := &service.MockOmniFocusService{
		PerspectiveTasks: []domain.Task{
			{
				ID:          "task1",
				Name:        "Task with metadata",
				ProjectID:   "proj1",
				ProjectName: "Work Project",
				Tags:        []string{"urgent", "meeting"},
			},
		},
	}

	output, exitCode, err := executePerspectiveCommand(mockService, []string{"Forecast"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Task with metadata") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}
}

func TestPerspectiveCommand_NoArguments(t *testing.T) {
	// Test that command requires a perspective name argument
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executePerspectiveCommand(mockService, []string{})

	if err == nil {
		t.Fatal("Expected error for missing argument, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for missing argument, got: %d", exitCode)
	}
}

func TestPerspectiveCommand_TooManyArguments(t *testing.T) {
	// Test that command only accepts one argument
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executePerspectiveCommand(mockService, []string{"Forecast", "Extra"})

	if err == nil {
		t.Fatal("Expected error for too many arguments, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for too many arguments, got: %d", exitCode)
	}
}

// Helper function to execute perspective command and capture output
func executePerspectiveCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Add perspective command
	rootCmd.AddCommand(NewPerspectiveCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "perspective" as first arg
	fullArgs := append([]string{"perspective"}, args...)
	rootCmd.SetArgs(fullArgs)

	// Use ExecuteContext with service in context
	ctx := ContextWithService(context.Background(), mockService)
	err := rootCmd.ExecuteContext(ctx)

	output := buf.String()
	exitCode := 0
	if err != nil {
		exitCode = 1 // Simplified - in real implementation we'd parse specific error types
	}

	return output, exitCode, err
}

// Helper to create time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
