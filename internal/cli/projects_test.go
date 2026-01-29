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

func TestProjectsCommand_DefaultActiveStatus(t *testing.T) {
	// Test that default behavior shows active projects
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Work Project", Status: "active"},
			{ID: "proj2", Name: "Home Renovation", Status: "active"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Work Project") {
		t.Errorf("Expected output to contain 'Work Project', got: %s", output)
	}

	if !strings.Contains(output, "Home Renovation") {
		t.Errorf("Expected output to contain 'Home Renovation', got: %s", output)
	}
}

func TestProjectsCommand_StatusActive(t *testing.T) {
	// Test explicit --status active flag
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Active Project", Status: "active"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--status", "active"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Active Project") {
		t.Errorf("Expected output to contain 'Active Project', got: %s", output)
	}
}

func TestProjectsCommand_StatusOnHold(t *testing.T) {
	// Test --status on-hold filter
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "On Hold Project", Status: "on-hold"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--status", "on-hold"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "On Hold Project") {
		t.Errorf("Expected output to contain 'On Hold Project', got: %s", output)
	}
}

func TestProjectsCommand_StatusCompleted(t *testing.T) {
	// Test --status completed filter
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Completed Project", Status: "completed"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--status", "completed"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Completed Project") {
		t.Errorf("Expected output to contain 'Completed Project', got: %s", output)
	}
}

func TestProjectsCommand_StatusDropped(t *testing.T) {
	// Test --status dropped filter
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Dropped Project", Status: "dropped"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--status", "dropped"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Dropped Project") {
		t.Errorf("Expected output to contain 'Dropped Project', got: %s", output)
	}
}

func TestProjectsCommand_StatusAll(t *testing.T) {
	// Test --status all shows all projects
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Active Project", Status: "active"},
			{ID: "proj2", Name: "On Hold Project", Status: "on-hold"},
			{ID: "proj3", Name: "Completed Project", Status: "completed"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--status", "all"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Active Project") {
		t.Errorf("Expected output to contain 'Active Project', got: %s", output)
	}

	if !strings.Contains(output, "On Hold Project") {
		t.Errorf("Expected output to contain 'On Hold Project', got: %s", output)
	}

	if !strings.Contains(output, "Completed Project") {
		t.Errorf("Expected output to contain 'Completed Project', got: %s", output)
	}
}

func TestProjectsCommand_WithTasks(t *testing.T) {
	// Test --with-tasks flag includes task details
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{
				ID:     "proj1",
				Name:   "Project with Tasks",
				Status: "active",
				Tasks: []domain.Task{
					{ID: "task1", Name: "Task One", Completed: false},
					{ID: "task2", Name: "Task Two", Completed: true},
				},
			},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--with-tasks"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Project with Tasks") {
		t.Errorf("Expected output to contain 'Project with Tasks', got: %s", output)
	}

	if !strings.Contains(output, "Task One") {
		t.Errorf("Expected output to contain 'Task One' when --with-tasks is set, got: %s", output)
	}

	if !strings.Contains(output, "Task Two") {
		t.Errorf("Expected output to contain 'Task Two' when --with-tasks is set, got: %s", output)
	}
}

func TestProjectsCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Test Project", Status: "active"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--json"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Check for JSON structure
	if !strings.Contains(output, `"projects"`) {
		t.Errorf("Expected JSON output to contain 'projects' field, got: %s", output)
	}

	if !strings.Contains(output, `"Test Project"`) {
		t.Errorf("Expected JSON output to contain project name, got: %s", output)
	}

	if !strings.Contains(output, `"count"`) {
		t.Errorf("Expected JSON output to contain 'count' field, got: %s", output)
	}
}

func TestProjectsCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		ProjectsErr: errors.New("OmniFocus is not running"),
	}

	_, exitCode, err := executeProjectsCommand(mockService, []string{})

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

func TestProjectsCommand_ErrorJSON(t *testing.T) {
	// Test error handling in JSON mode
	mockService := &service.MockOmniFocusService{
		ProjectsErr: errors.New("OmniFocus is not running"),
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--json"})

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

func TestProjectsCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{ID: "proj1", Name: "Test Project", Status: "active"},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--quiet"})

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

func TestProjectsCommand_EmptyResults(t *testing.T) {
	// Test empty project list
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "No projects") {
		t.Errorf("Expected output to indicate no projects, got: %s", output)
	}
}

func TestProjectsCommand_WithTasksAndStatus(t *testing.T) {
	// Test combining --with-tasks and --status flags
	mockService := &service.MockOmniFocusService{
		Projects: []domain.Project{
			{
				ID:     "proj1",
				Name:   "Completed Project",
				Status: "completed",
				Tasks: []domain.Task{
					{ID: "task1", Name: "Completed Task", Completed: true},
				},
			},
		},
	}

	output, exitCode, err := executeProjectsCommand(mockService, []string{"--status", "completed", "--with-tasks"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Completed Project") {
		t.Errorf("Expected output to contain 'Completed Project', got: %s", output)
	}

	if !strings.Contains(output, "Completed Task") {
		t.Errorf("Expected output to contain 'Completed Task', got: %s", output)
	}
}

// Helper function to execute projects command and capture output
func executeProjectsCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Add projects command
	rootCmd.AddCommand(NewProjectsCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "projects" as first arg
	fullArgs := append([]string{"projects"}, args...)
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
