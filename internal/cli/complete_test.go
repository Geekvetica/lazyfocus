package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestCompleteCommand_SingleTask(t *testing.T) {
	// Test completing a single task
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task completed",
	}

	mockService := &service.MockOmniFocusService{
		CompleteResult: result,
	}

	output, exitCode, err := executeCompleteCommand(mockService, []string{"task123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Completed") {
		t.Errorf("Expected output to contain 'Completed', got: %s", output)
	}

	if !strings.Contains(output, "task123") {
		t.Errorf("Expected output to contain task ID, got: %s", output)
	}
}

func TestCompleteCommand_MultipleTasks(t *testing.T) {
	// Test completing multiple tasks
	result := &domain.OperationResult{
		Success: true,
		ID:      "",
		Message: "Tasks completed",
	}

	mockService := &service.MockOmniFocusService{
		CompleteResult: result,
	}

	output, exitCode, err := executeCompleteCommand(mockService, []string{"task1", "task2", "task3"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Should see output for each task
	completedCount := strings.Count(output, "Completed")
	if completedCount < 3 {
		t.Errorf("Expected output for 3 tasks, got: %s", output)
	}
}

func TestCompleteCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task completed",
	}

	mockService := &service.MockOmniFocusService{
		CompleteResult: result,
	}

	output, exitCode, err := executeCompleteCommand(mockService, []string{"--json", "task123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Check for JSON structure
	if !strings.Contains(output, `"success"`) {
		t.Errorf("Expected JSON output to contain 'success' field, got: %s", output)
	}

	if !strings.Contains(output, `"id"`) {
		t.Errorf("Expected JSON output to contain 'id' field, got: %s", output)
	}
}

func TestCompleteCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		CompleteTaskErr: errors.New("task not found"),
	}

	_, exitCode, err := executeCompleteCommand(mockService, []string{"invalid-id"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "task not found") {
		t.Errorf("Expected error about task not found, got: %v", err)
	}
}

func TestCompleteCommand_PartialFailure(t *testing.T) {
	// Test that errors are reported but don't stop processing
	// Set the mock to return an error - this will apply to all calls
	mockService := &service.MockOmniFocusService{
		CompleteTaskErr: errors.New("OmniFocus connection failed"),
	}

	output, exitCode, err := executeCompleteCommand(mockService, []string{"task1", "task2"})

	// Should get an error since all tasks failed
	if err == nil {
		t.Fatal("Expected error when all tasks fail, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	// Verify error message appears in output (non-quiet mode)
	if !strings.Contains(output, "failed to complete") || !strings.Contains(err.Error(), "OmniFocus connection failed") {
		t.Errorf("Expected error details in output/error, got output: %s, err: %v", output, err)
	}
}

func TestCompleteCommand_NoTaskID(t *testing.T) {
	// Test error when no task ID provided
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executeCompleteCommand(mockService, []string{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}
}

func TestCompleteCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task completed",
	}

	mockService := &service.MockOmniFocusService{
		CompleteResult: result,
	}

	output, exitCode, err := executeCompleteCommand(mockService, []string{"--quiet", "task123"})

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

// Helper function to execute complete command and capture output
func executeCompleteCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Add the complete command
	rootCmd.AddCommand(NewCompleteCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "complete" as first arg
	fullArgs := append([]string{"complete"}, args...)
	rootCmd.SetArgs(fullArgs)

	// Use ExecuteContext with service in context
	ctx := ContextWithService(context.Background(), mockService)
	err := rootCmd.ExecuteContext(ctx)

	output := buf.String()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}

	return output, exitCode, err
}
