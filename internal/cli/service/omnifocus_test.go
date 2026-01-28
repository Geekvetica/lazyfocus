package service

import (
	"errors"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/bridge"
)

// mockExecutor implements bridge.Executor for testing
type mockExecutor struct {
	executeFunc func(script string) (string, error)
}

func (m *mockExecutor) Execute(script string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(script)
	}
	return "", nil
}

func (m *mockExecutor) ExecuteWithTimeout(script string, timeout time.Duration) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(script)
	}
	return "", nil
}

func TestNewOmniFocusService_CreatesServiceWithExecutor(t *testing.T) {
	executor := &mockExecutor{}
	timeout := 30 * time.Second

	service := NewOmniFocusService(executor, timeout)

	if service == nil {
		t.Fatal("NewOmniFocusService() returned nil")
	}

	if service.executor != executor {
		t.Error("NewOmniFocusService() did not set executor correctly")
	}

	if service.timeout != timeout {
		t.Errorf("NewOmniFocusService() timeout = %v, want %v", service.timeout, timeout)
	}
}

func TestGetInboxTasks_Success_ReturnsInboxTasks(t *testing.T) {
	expectedJSON := `{"tasks": [
		{"id": "task1", "name": "Task 1", "flagged": false, "completed": false},
		{"id": "task2", "name": "Task 2", "flagged": true, "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetInboxTasks()

	if err != nil {
		t.Fatalf("GetInboxTasks() error = %v, want nil", err)
	}

	if len(tasks) != 2 {
		t.Errorf("GetInboxTasks() returned %d tasks, want 2", len(tasks))
	}

	if tasks[0].ID != "task1" {
		t.Errorf("GetInboxTasks() first task ID = %s, want task1", tasks[0].ID)
	}
}

func TestGetInboxTasks_OmniFocusNotRunning_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", bridge.ErrOmniFocusNotRunning
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetInboxTasks()

	if !errors.Is(err, bridge.ErrOmniFocusNotRunning) {
		t.Errorf("GetInboxTasks() error = %v, want ErrOmniFocusNotRunning", err)
	}
}

func TestGetInboxTasks_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetInboxTasks()

	if err == nil {
		t.Fatal("GetInboxTasks() error = nil, want error")
	}
}

func TestGetInboxTasks_EmptyResponse_ReturnsEmptySlice(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"tasks": []}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetInboxTasks()

	if err != nil {
		t.Fatalf("GetInboxTasks() error = %v, want nil", err)
	}

	if len(tasks) != 0 {
		t.Errorf("GetInboxTasks() returned %d tasks, want 0", len(tasks))
	}
}

func TestGetProjects_Success_ReturnsProjects(t *testing.T) {
	expectedJSON := `{"projects": [
		{"id": "proj1", "name": "Project 1", "status": "active"},
		{"id": "proj2", "name": "Project 2", "status": "on-hold"}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	projects, err := service.GetProjects("active")

	if err != nil {
		t.Fatalf("GetProjects() error = %v, want nil", err)
	}

	if len(projects) != 2 {
		t.Errorf("GetProjects() returned %d projects, want 2", len(projects))
	}
}

func TestGetFlaggedTasks_Success_ReturnsFlaggedTasks(t *testing.T) {
	expectedJSON := `{"tasks": [
		{"id": "task1", "name": "Important Task", "flagged": true, "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetFlaggedTasks()

	if err != nil {
		t.Fatalf("GetFlaggedTasks() error = %v, want nil", err)
	}

	if len(tasks) != 1 {
		t.Errorf("GetFlaggedTasks() returned %d tasks, want 1", len(tasks))
	}

	if !tasks[0].Flagged {
		t.Error("GetFlaggedTasks() first task not flagged")
	}
}

func TestGetTasksByProject_Success_ReturnsProjectTasks(t *testing.T) {
	projectID := "project-123"
	expectedJSON := `{"tasks": [
		{"id": "task1", "name": "Task 1", "projectId": "project-123", "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetTasksByProject(projectID)

	if err != nil {
		t.Fatalf("GetTasksByProject() error = %v, want nil", err)
	}

	if len(tasks) != 1 {
		t.Errorf("GetTasksByProject() returned %d tasks, want 1", len(tasks))
	}

	if tasks[0].ProjectID != projectID {
		t.Errorf("GetTasksByProject() task projectID = %s, want %s", tasks[0].ProjectID, projectID)
	}
}

func TestGetTasksByTag_Success_ReturnsTaggedTasks(t *testing.T) {
	tagID := "tag-456"
	expectedJSON := `{"tasks": [
		{"id": "task1", "name": "Task 1", "tags": ["tag-456"], "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetTasksByTag(tagID)

	if err != nil {
		t.Fatalf("GetTasksByTag() error = %v, want nil", err)
	}

	if len(tasks) != 1 {
		t.Errorf("GetTasksByTag() returned %d tasks, want 1", len(tasks))
	}
}

func TestGetAllTasks_WithFilters_ReturnsFilteredTasks(t *testing.T) {
	dueDate := time.Now()
	filters := TaskFilters{
		Flagged:  true,
		DueStart: &dueDate,
	}

	expectedJSON := `{"tasks": [
		{"id": "task1", "name": "Task 1", "flagged": true, "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetAllTasks(filters)

	if err != nil {
		t.Fatalf("GetAllTasks() error = %v, want nil", err)
	}

	if len(tasks) != 1 {
		t.Errorf("GetAllTasks() returned %d tasks, want 1", len(tasks))
	}
}

func TestGetTaskByID_Success_ReturnsSingleTask(t *testing.T) {
	taskID := "task-789"
	expectedJSON := `{"tasks": [
		{"id": "task-789", "name": "Specific Task", "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	task, err := service.GetTaskByID(taskID)

	if err != nil {
		t.Fatalf("GetTaskByID() error = %v, want nil", err)
	}

	if task == nil {
		t.Fatal("GetTaskByID() returned nil task")
	}

	if task.ID != taskID {
		t.Errorf("GetTaskByID() task ID = %s, want %s", task.ID, taskID)
	}
}

func TestGetTaskByID_NotFound_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"tasks": []}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	task, err := service.GetTaskByID("nonexistent")

	if err == nil {
		t.Fatal("GetTaskByID() error = nil, want error for non-existent task")
	}

	if task != nil {
		t.Error("GetTaskByID() returned non-nil task for non-existent ID")
	}
}

func TestGetProjectByID_Success_ReturnsSingleProject(t *testing.T) {
	projectID := "proj-123"
	expectedJSON := `{"projects": [
		{"id": "proj-123", "name": "My Project", "status": "active"}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	project, err := service.GetProjectByID(projectID)

	if err != nil {
		t.Fatalf("GetProjectByID() error = %v, want nil", err)
	}

	if project == nil {
		t.Fatal("GetProjectByID() returned nil project")
	}

	if project.ID != projectID {
		t.Errorf("GetProjectByID() project ID = %s, want %s", project.ID, projectID)
	}
}

func TestGetProjectWithTasks_Success_ReturnsProjectWithTasks(t *testing.T) {
	projectID := "proj-456"
	expectedJSON := `{"projects": [
		{
			"id": "proj-456",
			"name": "Project With Tasks",
			"status": "active",
			"tasks": [
				{"id": "task1", "name": "Task 1", "completed": false}
			]
		}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	project, err := service.GetProjectWithTasks(projectID)

	if err != nil {
		t.Fatalf("GetProjectWithTasks() error = %v, want nil", err)
	}

	if project == nil {
		t.Fatal("GetProjectWithTasks() returned nil project")
	}

	if len(project.Tasks) != 1 {
		t.Errorf("GetProjectWithTasks() returned %d tasks, want 1", len(project.Tasks))
	}
}
