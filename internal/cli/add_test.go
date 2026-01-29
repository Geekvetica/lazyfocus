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

func TestAddCommand_BasicTask(t *testing.T) {
	// Test creating a simple task with just a name
	createdTask := &domain.Task{
		ID:   "task123",
		Name: "Buy milk",
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"Buy milk"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Created task") {
		t.Errorf("Expected output to contain 'Created task', got: %s", output)
	}

	if !strings.Contains(output, "Buy milk") {
		t.Errorf("Expected output to contain 'Buy milk', got: %s", output)
	}
}

func TestAddCommand_WithTags(t *testing.T) {
	// Test task with tags via natural syntax
	createdTask := &domain.Task{
		ID:   "task123",
		Name: "Buy groceries",
		Tags: []string{"errands", "shopping"},
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"Buy groceries #errands #shopping"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Buy groceries") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}
}

func TestAddCommand_WithProject(t *testing.T) {
	// Test task with project via natural syntax
	createdTask := &domain.Task{
		ID:          "task123",
		Name:        "Review code",
		ProjectID:   "proj1",
		ProjectName: "Work",
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask:       createdTask,
		ResolvedProjectID: "proj1",
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"Review code @Work"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Review code") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}

	if !strings.Contains(output, "Work") {
		t.Errorf("Expected output to contain project name, got: %s", output)
	}
}

func TestAddCommand_WithDueDate(t *testing.T) {
	// Test task with due date via natural syntax
	dueDate := time.Now().AddDate(0, 0, 1) // Tomorrow

	createdTask := &domain.Task{
		ID:      "task123",
		Name:    "Call dentist",
		DueDate: &dueDate,
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"Call dentist due:tomorrow"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Call dentist") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}
}

func TestAddCommand_WithFlagged(t *testing.T) {
	// Test task with flagged via natural syntax
	createdTask := &domain.Task{
		ID:      "task123",
		Name:    "Urgent task",
		Flagged: true,
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"Urgent task !"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Urgent task") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}
}

func TestAddCommand_WithFlags(t *testing.T) {
	// Test task created via command-line flags (not natural syntax)
	dueDate := time.Now().AddDate(0, 0, 1)

	createdTask := &domain.Task{
		ID:          "task123",
		Name:        "Meeting prep",
		DueDate:     &dueDate,
		ProjectName: "Work",
		Flagged:     true,
		Note:        "Prepare slides",
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask:       createdTask,
		ResolvedProjectID: "proj1",
	}

	output, exitCode, err := executeAddCommand(mockService, []string{
		"Meeting prep",
		"--due", "tomorrow",
		"--project", "Work",
		"--flagged",
		"--note", "Prepare slides",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Meeting prep") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}
}

func TestAddCommand_FlagsOverrideParsed(t *testing.T) {
	// Test that command-line flags override natural syntax
	dueDate := time.Now().AddDate(0, 0, 2) // 2 days from now

	createdTask := &domain.Task{
		ID:      "task123",
		Name:    "Task with override",
		DueDate: &dueDate,
		Flagged: false, // Flag overrides the ! in natural syntax
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	// Natural syntax says due:tomorrow and !, but flags override
	output, exitCode, err := executeAddCommand(mockService, []string{
		"Task with override due:tomorrow !",
		"--due", "in 2 days",
	})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Task with override") {
		t.Errorf("Expected output to contain task name, got: %s", output)
	}
}

func TestAddCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	createdTask := &domain.Task{
		ID:   "task123",
		Name: "Test task",
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"--json", "Test task"})

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

	if !strings.Contains(output, `"Test task"`) {
		t.Errorf("Expected JSON output to contain task name, got: %s", output)
	}
}

func TestAddCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		CreateTaskErr: errors.New("OmniFocus is not running"),
	}

	_, exitCode, err := executeAddCommand(mockService, []string{"Test task"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "OmniFocus is not running") {
		t.Errorf("Expected error message about OmniFocus, got: %v", err)
	}
}

func TestAddCommand_NoTaskName(t *testing.T) {
	// Test error when no task name provided
	mockService := &service.MockOmniFocusService{}

	_, exitCode, err := executeAddCommand(mockService, []string{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}
}

func TestAddCommand_ProjectResolutionError(t *testing.T) {
	// Test error when project name cannot be resolved
	mockService := &service.MockOmniFocusService{
		ResolveProjectErr: errors.New("project not found: NonExistent"),
	}

	_, exitCode, err := executeAddCommand(mockService, []string{"Task @NonExistent"})

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

func TestAddCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	createdTask := &domain.Task{
		ID:   "task123",
		Name: "Test task",
	}

	mockService := &service.MockOmniFocusService{
		CreatedTask: createdTask,
	}

	output, exitCode, err := executeAddCommand(mockService, []string{"--quiet", "Test task"})

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

func TestAddCommand_InvalidDueDate(t *testing.T) {
	// Test with --due "invalid" - should return error about invalid due date
	mockService := &service.MockOmniFocusService{}
	_, exitCode, err := executeAddCommand(mockService, []string{"Task name", "--due", "invalid"})

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

func TestAddCommand_InvalidDeferDate(t *testing.T) {
	// Test with --defer "invalid" - should return error about invalid defer date
	mockService := &service.MockOmniFocusService{}
	_, exitCode, err := executeAddCommand(mockService, []string{"Task name", "--defer", "invalid"})

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

// Helper function to execute add command and capture output
func executeAddCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Add the add command
	rootCmd.AddCommand(NewAddCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "add" as first arg
	fullArgs := append([]string{"add"}, args...)
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
