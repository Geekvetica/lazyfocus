package service

import (
	"errors"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/bridge"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
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
	expectedJSON := `{"task": {"id": "task-789", "name": "Specific Task", "completed": false}}`

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
			return `{"task": null}`, nil
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
	expectedJSON := `{"project": {"id": "proj-123", "name": "My Project", "status": "active"}}`

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
	expectedJSON := `{"project": {
		"id": "proj-456",
		"name": "Project With Tasks",
		"status": "active",
		"tasks": [
			{"id": "task1", "name": "Task 1", "completed": false}
		]
	}}`

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

func TestGetProjectByID_NotFound_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"project": null}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	project, err := service.GetProjectByID("nonexistent")

	if err == nil {
		t.Fatal("GetProjectByID() error = nil, want error for non-existent project")
	}

	if project != nil {
		t.Error("GetProjectByID() returned non-nil project for non-existent ID")
	}
}

func TestGetProjectWithTasks_NotFound_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"project": null}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	project, err := service.GetProjectWithTasks("nonexistent")

	if err == nil {
		t.Fatal("GetProjectWithTasks() error = nil, want error for non-existent project")
	}

	if project != nil {
		t.Error("GetProjectWithTasks() returned non-nil project for non-existent ID")
	}
}

func TestGetTags_Success_ReturnsTags(t *testing.T) {
	expectedJSON := `{"tags": [
		{"id": "tag1", "name": "Work"},
		{"id": "tag2", "name": "Home"}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tags, err := service.GetTags()

	if err != nil {
		t.Fatalf("GetTags() error = %v, want nil", err)
	}

	if len(tags) != 2 {
		t.Errorf("GetTags() returned %d tags, want 2", len(tags))
	}

	if tags[0].ID != "tag1" {
		t.Errorf("GetTags() first tag ID = %s, want tag1", tags[0].ID)
	}

	if tags[0].Name != "Work" {
		t.Errorf("GetTags() first tag name = %s, want Work", tags[0].Name)
	}
}

func TestGetTags_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetTags()

	if err == nil {
		t.Fatal("GetTags() error = nil, want error")
	}
}

func TestGetTags_ParseError_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `invalid json`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetTags()

	if err == nil {
		t.Fatal("GetTags() error = nil, want error for invalid JSON")
	}
}

func TestGetTagByID_Success_ReturnsTag(t *testing.T) {
	tagID := "tag1"
	expectedJSON := `{"tag": {"id": "tag1", "name": "Work"}}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tag, err := service.GetTagByID(tagID)

	if err != nil {
		t.Fatalf("GetTagByID() error = %v, want nil", err)
	}

	if tag == nil {
		t.Fatal("GetTagByID() returned nil tag")
	}

	if tag.ID != tagID {
		t.Errorf("GetTagByID() tag ID = %s, want %s", tag.ID, tagID)
	}

	if tag.Name != "Work" {
		t.Errorf("GetTagByID() tag name = %s, want Work", tag.Name)
	}
}

func TestGetTagByID_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetTagByID("tag1")

	if err == nil {
		t.Fatal("GetTagByID() error = nil, want error")
	}
}

func TestGetTagCounts_Success_ReturnsCounts(t *testing.T) {
	expectedJSON := `{"counts": {"Work": 5, "Home": 3}}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	counts, err := service.GetTagCounts()

	if err != nil {
		t.Fatalf("GetTagCounts() error = %v, want nil", err)
	}

	if len(counts) != 2 {
		t.Errorf("GetTagCounts() returned %d counts, want 2", len(counts))
	}

	if counts["Work"] != 5 {
		t.Errorf("GetTagCounts() Work count = %d, want 5", counts["Work"])
	}

	if counts["Home"] != 3 {
		t.Errorf("GetTagCounts() Home count = %d, want 3", counts["Home"])
	}
}

func TestGetTagCounts_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetTagCounts()

	if err == nil {
		t.Fatal("GetTagCounts() error = nil, want error")
	}
}

func TestGetPerspectiveTasks_Success_ReturnsTasks(t *testing.T) {
	perspectiveName := "Review"
	expectedJSON := `{"tasks": [
		{"id": "task1", "name": "Task 1", "completed": false},
		{"id": "task2", "name": "Task 2", "completed": false}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetPerspectiveTasks(perspectiveName)

	if err != nil {
		t.Fatalf("GetPerspectiveTasks() error = %v, want nil", err)
	}

	if len(tasks) != 2 {
		t.Errorf("GetPerspectiveTasks() returned %d tasks, want 2", len(tasks))
	}

	if tasks[0].ID != "task1" {
		t.Errorf("GetPerspectiveTasks() first task ID = %s, want task1", tasks[0].ID)
	}
}

func TestGetPerspectiveTasks_Empty_ReturnsEmptySlice(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"tasks": []}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	tasks, err := service.GetPerspectiveTasks("EmptyPerspective")

	if err != nil {
		t.Fatalf("GetPerspectiveTasks() error = %v, want nil", err)
	}

	if len(tasks) != 0 {
		t.Errorf("GetPerspectiveTasks() returned %d tasks, want 0", len(tasks))
	}
}

func TestGetPerspectiveTasks_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	_, err := service.GetPerspectiveTasks("Review")

	if err == nil {
		t.Fatal("GetPerspectiveTasks() error = nil, want error")
	}
}

// CreateTask Tests

func TestCreateTask_Success_ReturnsCreatedTask(t *testing.T) {
	expectedJSON := `{"task": {"id": "new-task", "name": "New Task", "flagged": false, "completed": false}}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	input := domain.TaskInput{Name: "New Task"}
	task, err := service.CreateTask(input)

	if err != nil {
		t.Fatalf("CreateTask() error = %v, want nil", err)
	}

	if task == nil {
		t.Fatal("CreateTask() returned nil task")
	}

	if task.ID != "new-task" {
		t.Errorf("CreateTask() task ID = %s, want new-task", task.ID)
	}

	if task.Name != "New Task" {
		t.Errorf("CreateTask() task name = %s, want New Task", task.Name)
	}

	if task.Flagged {
		t.Error("CreateTask() task flagged = true, want false")
	}

	if task.Completed {
		t.Error("CreateTask() task completed = true, want false")
	}
}

func TestCreateTask_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	input := domain.TaskInput{Name: "New Task"}
	task, err := service.CreateTask(input)

	if err == nil {
		t.Fatal("CreateTask() error = nil, want error")
	}

	if task != nil {
		t.Error("CreateTask() returned non-nil task on error")
	}
}

func TestCreateTask_InvalidInput_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"task": {"id": "task1", "name": "Task"}}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	input := domain.TaskInput{Name: ""} // Empty name should fail validation
	task, err := service.CreateTask(input)

	if err == nil {
		t.Fatal("CreateTask() error = nil, want validation error")
	}

	if task != nil {
		t.Error("CreateTask() returned non-nil task on validation error")
	}
}

func TestCreateTask_TaskNotCreated_ReturnsError(t *testing.T) {
	expectedJSON := `{"task": null}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	input := domain.TaskInput{Name: "New Task"}
	task, err := service.CreateTask(input)

	if err == nil {
		t.Fatal("CreateTask() error = nil, want error when task is null")
	}

	if task != nil {
		t.Error("CreateTask() returned non-nil task when response has null task")
	}
}

// ModifyTask Tests

func TestModifyTask_Success_ReturnsModifiedTask(t *testing.T) {
	expectedJSON := `{"task": {"id": "task123", "name": "Modified Task", "completed": false}}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	name := "Modified Task"
	mod := domain.TaskModification{Name: &name}
	task, err := service.ModifyTask("task123", mod)

	if err != nil {
		t.Fatalf("ModifyTask() error = %v, want nil", err)
	}

	if task == nil {
		t.Fatal("ModifyTask() returned nil task")
	}

	if task.ID != "task123" {
		t.Errorf("ModifyTask() task ID = %s, want task123", task.ID)
	}

	if task.Name != "Modified Task" {
		t.Errorf("ModifyTask() task name = %s, want Modified Task", task.Name)
	}
}

func TestModifyTask_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	name := "Modified Task"
	mod := domain.TaskModification{Name: &name}
	task, err := service.ModifyTask("task123", mod)

	if err == nil {
		t.Fatal("ModifyTask() error = nil, want error")
	}

	if task != nil {
		t.Error("ModifyTask() returned non-nil task on error")
	}
}

func TestModifyTask_EmptyModification_ReturnsError(t *testing.T) {
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return `{"task": {"id": "task123", "name": "Task"}}`, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	mod := domain.TaskModification{} // Empty modification
	task, err := service.ModifyTask("task123", mod)

	if err == nil {
		t.Fatal("ModifyTask() error = nil, want error for empty modification")
	}

	if task != nil {
		t.Error("ModifyTask() returned non-nil task for empty modification")
	}
}

func TestModifyTask_TaskNotFound_ReturnsError(t *testing.T) {
	expectedJSON := `{"task": null}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	name := "Modified Task"
	mod := domain.TaskModification{Name: &name}
	task, err := service.ModifyTask("nonexistent", mod)

	if err == nil {
		t.Fatal("ModifyTask() error = nil, want error when task not found")
	}

	if task != nil {
		t.Error("ModifyTask() returned non-nil task when response has null task")
	}
}

// CompleteTask Tests

func TestCompleteTask_Success_ReturnsResult(t *testing.T) {
	expectedJSON := `{"success": true, "id": "task123", "message": "Task completed"}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	result, err := service.CompleteTask("task123")

	if err != nil {
		t.Fatalf("CompleteTask() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("CompleteTask() returned nil result")
	}

	if !result.Success {
		t.Error("CompleteTask() result.Success = false, want true")
	}

	if result.ID != "task123" {
		t.Errorf("CompleteTask() result.ID = %s, want task123", result.ID)
	}

	if result.Message != "Task completed" {
		t.Errorf("CompleteTask() result.Message = %s, want Task completed", result.Message)
	}
}

func TestCompleteTask_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	result, err := service.CompleteTask("task123")

	if err == nil {
		t.Fatal("CompleteTask() error = nil, want error")
	}

	if result != nil {
		t.Error("CompleteTask() returned non-nil result on error")
	}
}

// DeleteTask Tests

func TestDeleteTask_Success_ReturnsResult(t *testing.T) {
	expectedJSON := `{"success": true, "id": "task123", "message": "Task deleted"}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return expectedJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	result, err := service.DeleteTask("task123")

	if err != nil {
		t.Fatalf("DeleteTask() error = %v, want nil", err)
	}

	if result == nil {
		t.Fatal("DeleteTask() returned nil result")
	}

	if !result.Success {
		t.Error("DeleteTask() result.Success = false, want true")
	}

	if result.ID != "task123" {
		t.Errorf("DeleteTask() result.ID = %s, want task123", result.ID)
	}

	if result.Message != "Task deleted" {
		t.Errorf("DeleteTask() result.Message = %s, want Task deleted", result.Message)
	}
}

func TestDeleteTask_ExecutorError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	result, err := service.DeleteTask("task123")

	if err == nil {
		t.Fatal("DeleteTask() error = nil, want error")
	}

	if result != nil {
		t.Error("DeleteTask() returned non-nil result on error")
	}
}

// ResolveProjectName Tests

func TestResolveProjectName_Success_ReturnsProjectID(t *testing.T) {
	projectsJSON := `{"projects": [
		{"id": "proj1", "name": "Work", "status": "active"},
		{"id": "proj2", "name": "Home", "status": "active"}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return projectsJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	projectID, err := service.ResolveProjectName("Work")

	if err != nil {
		t.Fatalf("ResolveProjectName() error = %v, want nil", err)
	}

	if projectID != "proj1" {
		t.Errorf("ResolveProjectName() projectID = %s, want proj1", projectID)
	}
}

func TestResolveProjectName_CaseInsensitive_ReturnsProjectID(t *testing.T) {
	projectsJSON := `{"projects": [
		{"id": "proj1", "name": "Work", "status": "active"},
		{"id": "proj2", "name": "Home", "status": "active"}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return projectsJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	projectID, err := service.ResolveProjectName("work")

	if err != nil {
		t.Fatalf("ResolveProjectName() error = %v, want nil", err)
	}

	if projectID != "proj1" {
		t.Errorf("ResolveProjectName() projectID = %s, want proj1", projectID)
	}
}

func TestResolveProjectName_NotFound_ReturnsError(t *testing.T) {
	projectsJSON := `{"projects": [
		{"id": "proj1", "name": "Work", "status": "active"},
		{"id": "proj2", "name": "Home", "status": "active"}
	]}`

	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return projectsJSON, nil
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	projectID, err := service.ResolveProjectName("NonExistent")

	if err == nil {
		t.Fatal("ResolveProjectName() error = nil, want error for non-existent project")
	}

	if projectID != "" {
		t.Errorf("ResolveProjectName() projectID = %s, want empty string on error", projectID)
	}
}

func TestResolveProjectName_GetProjectsError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("execution failed")
	executor := &mockExecutor{
		executeFunc: func(script string) (string, error) {
			return "", expectedErr
		},
	}

	service := NewOmniFocusService(executor, 30*time.Second)
	projectID, err := service.ResolveProjectName("Work")

	if err == nil {
		t.Fatal("ResolveProjectName() error = nil, want error")
	}

	if projectID != "" {
		t.Errorf("ResolveProjectName() projectID = %s, want empty string on error", projectID)
	}
}
