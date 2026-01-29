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

func TestModifyCommand_Name(t *testing.T) {
	// Test modifying task name
	modifiedTask := &domain.Task{
		ID:   "task123",
		Name: "New task name",
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--name", "New task name"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}

	if !strings.Contains(output, "New task name") {
		t.Errorf("Expected output to contain new name, got: %s", output)
	}
}

func TestModifyCommand_Note(t *testing.T) {
	// Test modifying task note
	modifiedTask := &domain.Task{
		ID:   "task123",
		Name: "Task",
		Note: "Updated note",
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--note", "Updated note"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_Project(t *testing.T) {
	// Test modifying task project
	modifiedTask := &domain.Task{
		ID:          "task123",
		Name:        "Task",
		ProjectID:   "proj1",
		ProjectName: "Work",
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask:      modifiedTask,
		ResolvedProjectID: "proj1",
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--project", "Work"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_AddTags(t *testing.T) {
	// Test adding tags
	modifiedTask := &domain.Task{
		ID:   "task123",
		Name: "Task",
		Tags: []string{"urgent", "work"},
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{
		"task123",
		"--add-tag", "urgent",
		"--add-tag", "work",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_RemoveTags(t *testing.T) {
	// Test removing tags
	modifiedTask := &domain.Task{
		ID:   "task123",
		Name: "Task",
		Tags: []string{},
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{
		"task123",
		"--remove-tag", "old",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_DueDate(t *testing.T) {
	// Test setting due date
	dueDate := time.Now().AddDate(0, 0, 1)
	modifiedTask := &domain.Task{
		ID:      "task123",
		Name:    "Task",
		DueDate: &dueDate,
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--due", "tomorrow"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_DeferDate(t *testing.T) {
	// Test setting defer date
	deferDate := time.Now().AddDate(0, 0, 2)
	modifiedTask := &domain.Task{
		ID:        "task123",
		Name:      "Task",
		DeferDate: &deferDate,
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--defer", "in 2 days"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_Flagged(t *testing.T) {
	// Test setting flagged status
	modifiedTask := &domain.Task{
		ID:      "task123",
		Name:    "Task",
		Flagged: true,
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--flagged", "true"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}

	if !strings.Contains(output, "Flagged") {
		t.Errorf("Expected output to show flagged status, got: %s", output)
	}
}

func TestModifyCommand_ClearDue(t *testing.T) {
	// Test clearing due date
	modifiedTask := &domain.Task{
		ID:      "task123",
		Name:    "Task",
		DueDate: nil,
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--clear-due"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_ClearDefer(t *testing.T) {
	// Test clearing defer date
	modifiedTask := &domain.Task{
		ID:        "task123",
		Name:      "Task",
		DeferDate: nil,
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--clear-defer"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_MultipleChanges(t *testing.T) {
	// Test modifying multiple fields at once
	dueDate := time.Now().AddDate(0, 0, 1)
	modifiedTask := &domain.Task{
		ID:      "task123",
		Name:    "Updated name",
		DueDate: &dueDate,
		Flagged: true,
		Note:    "Updated note",
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{
		"task123",
		"--name", "Updated name",
		"--due", "tomorrow",
		"--flagged", "true",
		"--note", "Updated note",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Modified") {
		t.Errorf("Expected output to contain 'Modified', got: %s", output)
	}
}

func TestModifyCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	modifiedTask := &domain.Task{
		ID:   "task123",
		Name: "Updated task",
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{
		"--json",
		"task123",
		"--name", "Updated task",
	})

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

	if !strings.Contains(output, `"task"`) {
		t.Errorf("Expected JSON output to contain 'task' field, got: %s", output)
	}
}

func TestModifyCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		ModifyTaskErr: errors.New("task not found"),
	}

	_, exitCode, err := executeModifyCommand(mockService, []string{"invalid-id", "--name", "New name"})

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

func TestModifyCommand_NoModifications(t *testing.T) {
	// Test error when no modifications specified
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executeModifyCommand(mockService, []string{"task123"})

	if err == nil {
		t.Fatal("Expected error when no modifications specified, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}
}

func TestModifyCommand_NoTaskID(t *testing.T) {
	// Test error when no task ID provided
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executeModifyCommand(mockService, []string{"--name", "New name"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}
}

func TestModifyCommand_TooManyTaskIDs(t *testing.T) {
	// Test error when multiple task IDs provided
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executeModifyCommand(mockService, []string{"task1", "task2", "--name", "New name"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}
}

func TestModifyCommand_ProjectResolutionError(t *testing.T) {
	// Test error when project name cannot be resolved
	mockService := &service.MockOmniFocusService{
		ResolveProjectErr: errors.New("project not found: NonExistent"),
	}

	_, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--project", "NonExistent"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "project not found") {
		t.Errorf("Expected error about project not found, got: %v", err)
	}
}

func TestModifyCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	modifiedTask := &domain.Task{
		ID:   "task123",
		Name: "Updated task",
	}

	mockService := &service.MockOmniFocusService{
		ModifiedTask: modifiedTask,
	}

	output, exitCode, err := executeModifyCommand(mockService, []string{
		"--quiet",
		"task123",
		"--name", "Updated task",
	})

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

func TestModifyCommand_InvalidDueDate(t *testing.T) {
	// Test with --due "invalid"
	mockService := &service.MockOmniFocusService{}
	_, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--due", "invalid"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "invalid due date") {
		t.Errorf("Expected error about invalid due date, got: %v", err)
	}

	if !strings.Contains(err.Error(), "unrecognized date format") {
		t.Errorf("Expected error about unrecognized date format, got: %v", err)
	}
}

func TestModifyCommand_InvalidDeferDate(t *testing.T) {
	// Test with --defer "invalid"
	mockService := &service.MockOmniFocusService{}
	_, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--defer", "invalid"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "invalid defer date") {
		t.Errorf("Expected error about invalid defer date, got: %v", err)
	}

	if !strings.Contains(err.Error(), "unrecognized date format") {
		t.Errorf("Expected error about unrecognized date format, got: %v", err)
	}
}

func TestModifyCommand_InvalidFlaggedValue(t *testing.T) {
	// Test with --flagged "invalid" (not true/false)
	mockService := &service.MockOmniFocusService{}
	_, exitCode, err := executeModifyCommand(mockService, []string{"task123", "--flagged", "invalid"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "invalid flagged value") {
		t.Errorf("Expected error about invalid flagged value, got: %v", err)
	}
}

// Helper function to execute modify command and capture output
func executeModifyCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Add the modify command
	rootCmd.AddCommand(NewModifyCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "modify" as first arg
	fullArgs := append([]string{"modify"}, args...)
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
