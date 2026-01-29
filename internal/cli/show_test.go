package cli

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestShowCommand_TaskByID_ExplicitType(t *testing.T) {
	expectedTask := &domain.Task{
		ID:          "task123",
		Name:        "Buy groceries",
		Note:        "Milk, eggs, bread",
		ProjectName: "Personal",
		Flagged:     true,
	}

	mockService := &service.MockOmniFocusService{
		Task: expectedTask,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"task123", "--type", "task"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Buy groceries") {
		t.Errorf("Expected output to contain 'Buy groceries', got: %s", output)
	}
}

func TestShowCommand_ProjectByID_ExplicitType(t *testing.T) {
	expectedProject := &domain.Project{
		ID:     "proj123",
		Name:   "Website Redesign",
		Status: "active",
		Note:   "Complete by Q1",
	}

	mockService := &service.MockOmniFocusService{
		Project: expectedProject,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"proj123", "--type", "project"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Website Redesign") {
		t.Errorf("Expected output to contain 'Website Redesign', got: %s", output)
	}
}

func TestShowCommand_TagByID_ExplicitType(t *testing.T) {
	expectedTag := &domain.Tag{
		ID:   "tag123",
		Name: "urgent",
	}

	mockService := &service.MockOmniFocusService{
		Tag: expectedTag,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"tag123", "--type", "tag"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "urgent") {
		t.Errorf("Expected output to contain 'urgent', got: %s", output)
	}
}

func TestShowCommand_AutoDetectTask(t *testing.T) {
	expectedTask := &domain.Task{
		ID:   "task123",
		Name: "Write tests",
	}

	mockService := &service.MockOmniFocusService{
		Task: expectedTask,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"task123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Write tests") {
		t.Errorf("Expected output to contain 'Write tests', got: %s", output)
	}
}

func TestShowCommand_AutoDetectProject_WhenNotFoundAsTask(t *testing.T) {
	expectedProject := &domain.Project{
		ID:     "proj123",
		Name:   "Marketing Campaign",
		Status: "active",
	}

	mockService := &service.MockOmniFocusService{
		Task:    nil,
		TaskErr: fmt.Errorf("task not found: proj123"),
		Project: expectedProject,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"proj123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Marketing Campaign") {
		t.Errorf("Expected output to contain 'Marketing Campaign', got: %s", output)
	}
}

func TestShowCommand_AutoDetectTag_WhenNotFoundAsTaskOrProject(t *testing.T) {
	expectedTag := &domain.Tag{
		ID:   "tag123",
		Name: "work",
	}

	mockService := &service.MockOmniFocusService{
		Task:       nil,
		TaskErr:    fmt.Errorf("task not found: tag123"),
		Project:    nil,
		ProjectErr: fmt.Errorf("project not found: tag123"),
		Tag:        expectedTag,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"tag123"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "work") {
		t.Errorf("Expected output to contain 'work', got: %s", output)
	}
}

func TestShowCommand_ItemNotFound(t *testing.T) {
	mockService := &service.MockOmniFocusService{
		Task:       nil,
		TaskErr:    fmt.Errorf("task not found: unknown123"),
		Project:    nil,
		ProjectErr: fmt.Errorf("project not found: unknown123"),
		Tag:        nil,
		TagErr:     fmt.Errorf("tag not found: unknown123"),
	}

	_, exitCode, err := executeShowCommand(mockService, []string{"unknown123"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check it's an ItemNotFoundError
	if _, ok := err.(*ItemNotFoundError); !ok {
		t.Errorf("Expected ItemNotFoundError, got: %T", err)
	}

	if !strings.Contains(err.Error(), "item not found: unknown123") {
		t.Errorf("Expected error message to contain 'item not found: unknown123', got: %v", err)
	}

	// Exit code 3 is for item not found
	if exitCode != 3 {
		t.Errorf("Expected exit code 3, got: %d", exitCode)
	}
}

func TestShowCommand_JSONOutput_Task(t *testing.T) {
	dueDate := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	expectedTask := &domain.Task{
		ID:          "task123",
		Name:        "Complete report",
		Note:        "Include charts",
		ProjectName: "Q1 Reports",
		DueDate:     &dueDate,
		Flagged:     true,
	}

	mockService := &service.MockOmniFocusService{
		Task: expectedTask,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"task123", "--json"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Check for JSON structure
	if !strings.Contains(output, `"id"`) {
		t.Errorf("Expected JSON output to contain 'id' field, got: %s", output)
	}

	if !strings.Contains(output, `"Complete report"`) {
		t.Errorf("Expected JSON output to contain task name, got: %s", output)
	}
}

func TestShowCommand_JSONOutput_Project(t *testing.T) {
	expectedProject := &domain.Project{
		ID:     "proj123",
		Name:   "Infrastructure",
		Status: "on-hold",
		Note:   "Waiting for budget approval",
	}

	mockService := &service.MockOmniFocusService{
		Project: expectedProject,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"proj123", "--type", "project", "--json"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, `"Infrastructure"`) {
		t.Errorf("Expected JSON output to contain 'Infrastructure', got: %s", output)
	}

	if !strings.Contains(output, `"on-hold"`) {
		t.Errorf("Expected JSON output to contain 'on-hold', got: %s", output)
	}
}

func TestShowCommand_JSONOutput_Tag(t *testing.T) {
	expectedTag := &domain.Tag{
		ID:   "tag123",
		Name: "personal",
	}

	mockService := &service.MockOmniFocusService{
		Tag: expectedTag,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"tag123", "--type", "tag", "--json"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, `"personal"`) {
		t.Errorf("Expected JSON output to contain 'personal', got: %s", output)
	}
}

func TestShowCommand_InvalidType(t *testing.T) {
	mockService := &service.MockOmniFocusService{}

	_, _, err := executeShowCommand(mockService, []string{"item123", "--type", "invalid"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "unknown type") {
		t.Errorf("Expected error message about unknown type, got: %v", err)
	}
}

func TestShowCommand_NoID(t *testing.T) {
	mockService := &service.MockOmniFocusService{}

	_, _, err := executeShowCommand(mockService, []string{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestShowCommand_ServiceError_Task(t *testing.T) {
	mockService := &service.MockOmniFocusService{
		TaskErr: fmt.Errorf("OmniFocus is not running"),
	}

	_, _, err := executeShowCommand(mockService, []string{"task123", "--type", "task"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "OmniFocus is not running") {
		t.Errorf("Expected error message about OmniFocus, got: %v", err)
	}
}

func TestShowCommand_QuietMode_Task(t *testing.T) {
	expectedTask := &domain.Task{
		ID:   "task123",
		Name: "Test task",
	}

	mockService := &service.MockOmniFocusService{
		Task: expectedTask,
	}

	output, exitCode, err := executeShowCommand(mockService, []string{"task123", "--quiet"})

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

func TestShowCommand_QuietMode_NotFound(t *testing.T) {
	mockService := &service.MockOmniFocusService{
		Task:       nil,
		TaskErr:    fmt.Errorf("task not found: unknown123"),
		Project:    nil,
		ProjectErr: fmt.Errorf("project not found: unknown123"),
		Tag:        nil,
		TagErr:     fmt.Errorf("tag not found: unknown123"),
	}

	_, exitCode, err := executeShowCommand(mockService, []string{"unknown123", "--quiet"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check it's an ItemNotFoundError
	if _, ok := err.(*ItemNotFoundError); !ok {
		t.Errorf("Expected ItemNotFoundError, got: %T", err)
	}

	if exitCode != 3 {
		t.Errorf("Expected exit code 3, got: %d", exitCode)
	}
}

// Helper function to execute show command and capture output
func executeShowCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := NewRootCommand()

	// Add show command
	rootCmd.AddCommand(NewShowCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "show" as first arg
	fullArgs := append([]string{"show"}, args...)
	rootCmd.SetArgs(fullArgs)

	// Use ExecuteContext with service in context
	ctx := ContextWithService(context.Background(), mockService)
	err := rootCmd.ExecuteContext(ctx)

	output := buf.String()
	exitCode := 0
	if err != nil {
		// Check if it's ItemNotFoundError
		if _, ok := err.(*ItemNotFoundError); ok {
			exitCode = 3
		} else {
			exitCode = 1
		}
	}

	return output, exitCode, err
}
