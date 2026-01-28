package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

func TestTasksCommand_Inbox(t *testing.T) {
	// Test that default behavior shows inbox tasks
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Buy groceries", Flagged: false},
			{ID: "task2", Name: "Call dentist", Flagged: true},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Buy groceries") {
		t.Errorf("Expected output to contain 'Buy groceries', got: %s", output)
	}

	if !strings.Contains(output, "Call dentist") {
		t.Errorf("Expected output to contain 'Call dentist', got: %s", output)
	}
}

func TestTasksCommand_InboxFlag(t *testing.T) {
	// Test explicit --inbox flag
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Inbox task", Flagged: false},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--inbox"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Inbox task") {
		t.Errorf("Expected output to contain 'Inbox task', got: %s", output)
	}
}

func TestTasksCommand_All(t *testing.T) {
	// Test --all flag shows all tasks
	mockService := &service.MockOmniFocusService{
		AllTasks: []domain.Task{
			{ID: "task1", Name: "Task 1", Flagged: false},
			{ID: "task2", Name: "Task 2", Flagged: true},
			{ID: "task3", Name: "Task 3", Completed: true},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--all"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Task 1") {
		t.Errorf("Expected output to contain 'Task 1', got: %s", output)
	}

	if !strings.Contains(output, "Task 3") {
		t.Errorf("Expected output to contain completed 'Task 3', got: %s", output)
	}
}

func TestTasksCommand_Flagged(t *testing.T) {
	// Test --flagged filter
	mockService := &service.MockOmniFocusService{
		FlaggedTasks: []domain.Task{
			{ID: "task1", Name: "Important task", Flagged: true},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--flagged"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Important task") {
		t.Errorf("Expected output to contain 'Important task', got: %s", output)
	}
}

func TestTasksCommand_Project(t *testing.T) {
	// Test --project filter
	mockService := &service.MockOmniFocusService{
		ProjectTasks: []domain.Task{
			{ID: "task1", Name: "Project task", ProjectID: "proj1", ProjectName: "Work"},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--project", "proj1"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Project task") {
		t.Errorf("Expected output to contain 'Project task', got: %s", output)
	}
}

func TestTasksCommand_Tag(t *testing.T) {
	// Test --tag filter
	mockService := &service.MockOmniFocusService{
		TagTasks: []domain.Task{
			{ID: "task1", Name: "Tagged task", Tags: []string{"urgent"}},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--tag", "urgent"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Tagged task") {
		t.Errorf("Expected output to contain 'Tagged task', got: %s", output)
	}
}

func TestTasksCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test task", Flagged: false},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--json"})

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

	if !strings.Contains(output, `"Test task"`) {
		t.Errorf("Expected JSON output to contain task name, got: %s", output)
	}

	if !strings.Contains(output, `"count"`) {
		t.Errorf("Expected JSON output to contain 'count' field, got: %s", output)
	}
}

func TestTasksCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		InboxTasksErr: errors.New("OmniFocus is not running"),
	}

	_, exitCode, err := executeTasksCommand(mockService, []string{})

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

func TestTasksCommand_ErrorJSON(t *testing.T) {
	// Test error handling in JSON mode
	mockService := &service.MockOmniFocusService{
		InboxTasksErr: errors.New("OmniFocus is not running"),
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--json"})

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

func TestTasksCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test task", Flagged: false},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--quiet"})

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

func TestTasksCommand_DueDateToday(t *testing.T) {
	// Test --due today
	today := time.Now()
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Due today", DueDate: &today},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--due", "today"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Due today") {
		t.Errorf("Expected output to contain 'Due today', got: %s", output)
	}
}

func TestTasksCommand_DueDateTomorrow(t *testing.T) {
	// Test --due tomorrow
	tomorrow := time.Now().AddDate(0, 0, 1)
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Due tomorrow", DueDate: &tomorrow},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--due", "tomorrow"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Due tomorrow") {
		t.Errorf("Expected output to contain 'Due tomorrow', got: %s", output)
	}
}

func TestTasksCommand_DueDateSpecific(t *testing.T) {
	// Test --due with specific date
	specificDate, _ := time.Parse("2006-01-02", "2024-12-25")
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Christmas task", DueDate: &specificDate},
		},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{"--due", "2024-12-25"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Christmas task") {
		t.Errorf("Expected output to contain 'Christmas task', got: %s", output)
	}
}

func TestTasksCommand_EmptyResults(t *testing.T) {
	// Test empty task list
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{},
	}

	output, exitCode, err := executeTasksCommand(mockService, []string{})

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

// Helper function to execute tasks command and capture output
func executeTasksCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Override the service for testing
	Service = mockService
	defer func() { Service = nil }()

	// Add tasks command
	rootCmd.AddCommand(NewTasksCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "tasks" as first arg
	fullArgs := append([]string{"tasks"}, args...)
	rootCmd.SetArgs(fullArgs)

	// Execute
	err := rootCmd.Execute()

	output := buf.String()
	exitCode := 0
	if err != nil {
		exitCode = 1 // Simplified - in real implementation we'd parse specific error types
	}

	return output, exitCode, err
}

// newTestRootCommand creates a simplified root command for testing
func newTestRootCommand() *cobra.Command {
	return NewRootCommand()
}
