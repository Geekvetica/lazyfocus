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

func TestDeleteCommand_SingleTaskWithForce(t *testing.T) {
	// Test deleting a single task with --force (no confirmation)
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task deleted",
	}

	mockService := &service.MockOmniFocusService{
		DeleteResult: result,
	}

	output, exitCode, err := executeDeleteCommand(mockService, []string{"--force", "task123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Deleted") {
		t.Errorf("Expected output to contain 'Deleted', got: %s", output)
	}

	if !strings.Contains(output, "task123") {
		t.Errorf("Expected output to contain task ID, got: %s", output)
	}
}

func TestDeleteCommand_MultipleTasksWithForce(t *testing.T) {
	// Test deleting multiple tasks with --force
	result := &domain.OperationResult{
		Success: true,
		ID:      "",
		Message: "Tasks deleted",
	}

	mockService := &service.MockOmniFocusService{
		DeleteResult: result,
	}

	output, exitCode, err := executeDeleteCommand(mockService, []string{"--force", "task1", "task2", "task3"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Should see output for each task
	deletedCount := strings.Count(output, "Deleted")
	if deletedCount < 3 {
		t.Errorf("Expected output for 3 tasks, got: %s", output)
	}
}

func TestDeleteCommand_JSONOutputWithForce(t *testing.T) {
	// Test JSON output format with --force
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task deleted",
	}

	mockService := &service.MockOmniFocusService{
		DeleteResult: result,
	}

	output, exitCode, err := executeDeleteCommand(mockService, []string{"--json", "--force", "task123"})

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

func TestDeleteCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		DeleteTaskErr: errors.New("task not found"),
	}

	_, exitCode, err := executeDeleteCommand(mockService, []string{"--force", "invalid-id"})

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

func TestDeleteCommand_PartialFailure(t *testing.T) {
	// Test that errors are reported but don't stop processing
	mockService := &service.MockOmniFocusService{
		DeleteTaskErr: errors.New("OmniFocus connection failed"),
	}

	output, exitCode, err := executeDeleteCommand(mockService, []string{"--force", "task1", "task2"})

	// Should get an error since all tasks failed
	if err == nil {
		t.Fatal("Expected error when all tasks fail, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	// Verify error message appears
	if !strings.Contains(output, "failed to delete") || !strings.Contains(err.Error(), "OmniFocus connection failed") {
		t.Errorf("Expected error details, got output: %s, err: %v", output, err)
	}
}

func TestDeleteCommand_NoTaskID(t *testing.T) {
	// Test error when no task ID provided
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executeDeleteCommand(mockService, []string{"--force"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}
}

func TestDeleteCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task deleted",
	}

	mockService := &service.MockOmniFocusService{
		DeleteResult: result,
	}

	output, exitCode, err := executeDeleteCommand(mockService, []string{"--quiet", "--force", "task123"})

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

func TestDeleteCommand_JSONModeSkipsConfirmation(t *testing.T) {
	// Test that JSON mode skips confirmation even without --force
	result := &domain.OperationResult{
		Success: true,
		ID:      "task123",
		Message: "Task deleted",
	}

	mockService := &service.MockOmniFocusService{
		DeleteResult: result,
	}

	// No --force flag, but --json should skip confirmation
	output, exitCode, err := executeDeleteCommand(mockService, []string{"--json", "task123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Should have JSON output (not hang waiting for confirmation)
	if !strings.Contains(output, `"success"`) {
		t.Errorf("Expected JSON output, got: %s", output)
	}
}

// Helper function to execute delete command and capture output
func executeDeleteCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Add the delete command
	rootCmd.AddCommand(NewDeleteCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "delete" as first arg
	fullArgs := append([]string{"delete"}, args...)
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
