package bridge

import (
	"testing"
)

func TestParseTasks_ValidJSON(t *testing.T) {
	jsonStr := `{
		"tasks": [
			{
				"id": "abc123",
				"name": "Buy groceries",
				"note": "Remember milk",
				"tags": ["errands"],
				"dueDate": "2025-01-28T17:00:00.000Z",
				"deferDate": null,
				"flagged": true,
				"completed": false,
				"completedDate": null
			},
			{
				"id": "def456",
				"name": "Review PR #142",
				"note": "",
				"projectId": "proj123",
				"projectName": "Work Project",
				"tags": ["work", "code-review"],
				"dueDate": null,
				"deferDate": "2025-01-27T09:00:00.000Z",
				"flagged": false,
				"completed": false,
				"completedDate": null
			}
		]
	}`

	tasks, err := ParseTasks(jsonStr)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(tasks))
	}

	// Verify first task
	task1 := tasks[0]
	if task1.ID != "abc123" {
		t.Errorf("expected ID 'abc123', got '%s'", task1.ID)
	}
	if task1.Name != "Buy groceries" {
		t.Errorf("expected name 'Buy groceries', got '%s'", task1.Name)
	}
	if task1.Note != "Remember milk" {
		t.Errorf("expected note 'Remember milk', got '%s'", task1.Note)
	}
	if !task1.Flagged {
		t.Error("expected task1 to be flagged")
	}
	if task1.Completed {
		t.Error("expected task1 to not be completed")
	}
	if len(task1.Tags) != 1 || task1.Tags[0] != "errands" {
		t.Errorf("expected tags ['errands'], got %v", task1.Tags)
	}
	if task1.DueDate == nil {
		t.Error("expected dueDate to be set")
	}
	if task1.DeferDate != nil {
		t.Error("expected deferDate to be nil")
	}

	// Verify second task
	task2 := tasks[1]
	if task2.ID != "def456" {
		t.Errorf("expected ID 'def456', got '%s'", task2.ID)
	}
	if task2.ProjectID != "proj123" {
		t.Errorf("expected projectId 'proj123', got '%s'", task2.ProjectID)
	}
	if task2.ProjectName != "Work Project" {
		t.Errorf("expected projectName 'Work Project', got '%s'", task2.ProjectName)
	}
	if len(task2.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(task2.Tags))
	}
	if task2.DueDate != nil {
		t.Error("expected dueDate to be nil")
	}
	if task2.DeferDate == nil {
		t.Error("expected deferDate to be set")
	}
}

func TestParseTasks_EmptyArray(t *testing.T) {
	jsonStr := `{"tasks": []}`

	tasks, err := ParseTasks(jsonStr)

	if err != nil {
		t.Fatalf("expected no error for empty array, got %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("expected empty slice, got %d tasks", len(tasks))
	}

	if tasks == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestParseTasks_MalformedJSON(t *testing.T) {
	testCases := []struct {
		name    string
		jsonStr string
	}{
		{
			name:    "Invalid JSON syntax",
			jsonStr: `{"tasks": [}`,
		},
		{
			name:    "Not JSON at all",
			jsonStr: `this is not json`,
		},
		{
			name:    "Incomplete JSON",
			jsonStr: `{"tasks":`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseTasks(tc.jsonStr)

			if err == nil {
				t.Error("expected error for malformed JSON, got nil")
			}
		})
	}
}

func TestParseTasks_ErrorField(t *testing.T) {
	jsonStr := `{"error": "Some error occurred"}`

	_, err := ParseTasks(jsonStr)

	if err == nil {
		t.Fatal("expected error when JSON contains error field")
	}

	expectedMsg := "Some error occurred"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestParseTasks_OmniFocusNotRunning(t *testing.T) {
	jsonStr := `{"error": "OmniFocus is not running"}`

	_, err := ParseTasks(jsonStr)

	if err == nil {
		t.Fatal("expected error when OmniFocus is not running")
	}

	if err != ErrOmniFocusNotRunning {
		t.Errorf("expected ErrOmniFocusNotRunning, got %v", err)
	}
}

func TestParseTasks_CompletedTask(t *testing.T) {
	jsonStr := `{
		"tasks": [
			{
				"id": "completed123",
				"name": "Completed task",
				"note": "",
				"tags": [],
				"dueDate": null,
				"deferDate": null,
				"flagged": false,
				"completed": true,
				"completedDate": "2025-01-27T14:30:00.000Z"
			}
		]
	}`

	tasks, err := ParseTasks(jsonStr)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(tasks))
	}

	task := tasks[0]
	if !task.Completed {
		t.Error("expected task to be completed")
	}
	if task.CompletedDate == nil {
		t.Error("expected completedDate to be set")
	}
}

func TestParseProjects_ValidJSON(t *testing.T) {
	jsonStr := `{
		"projects": [
			{
				"id": "xyz789",
				"name": "Home Renovation",
				"status": "active",
				"note": "Kitchen project"
			},
			{
				"id": "abc111",
				"name": "Old Project",
				"status": "completed",
				"note": ""
			}
		]
	}`

	projects, err := ParseProjects(jsonStr)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	// Verify first project
	proj1 := projects[0]
	if proj1.ID != "xyz789" {
		t.Errorf("expected ID 'xyz789', got '%s'", proj1.ID)
	}
	if proj1.Name != "Home Renovation" {
		t.Errorf("expected name 'Home Renovation', got '%s'", proj1.Name)
	}
	if proj1.Status != "active" {
		t.Errorf("expected status 'active', got '%s'", proj1.Status)
	}
	if proj1.Note != "Kitchen project" {
		t.Errorf("expected note 'Kitchen project', got '%s'", proj1.Note)
	}

	// Verify second project
	proj2 := projects[1]
	if proj2.ID != "abc111" {
		t.Errorf("expected ID 'abc111', got '%s'", proj2.ID)
	}
	if proj2.Status != "completed" {
		t.Errorf("expected status 'completed', got '%s'", proj2.Status)
	}
}

func TestParseProjects_EmptyArray(t *testing.T) {
	jsonStr := `{"projects": []}`

	projects, err := ParseProjects(jsonStr)

	if err != nil {
		t.Fatalf("expected no error for empty array, got %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("expected empty slice, got %d projects", len(projects))
	}

	if projects == nil {
		t.Error("expected non-nil slice, got nil")
	}
}

func TestParseProjects_MalformedJSON(t *testing.T) {
	testCases := []struct {
		name    string
		jsonStr string
	}{
		{
			name:    "Invalid JSON syntax",
			jsonStr: `{"projects": [}`,
		},
		{
			name:    "Not JSON at all",
			jsonStr: `this is not json`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ParseProjects(tc.jsonStr)

			if err == nil {
				t.Error("expected error for malformed JSON, got nil")
			}
		})
	}
}

func TestParseProjects_ErrorField(t *testing.T) {
	jsonStr := `{"error": "Database connection failed"}`

	_, err := ParseProjects(jsonStr)

	if err == nil {
		t.Fatal("expected error when JSON contains error field")
	}

	expectedMsg := "Database connection failed"
	if err.Error() != expectedMsg {
		t.Errorf("expected error message '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestParseProjects_OmniFocusNotRunning(t *testing.T) {
	jsonStr := `{"error": "OmniFocus is not running"}`

	_, err := ParseProjects(jsonStr)

	if err == nil {
		t.Fatal("expected error when OmniFocus is not running")
	}

	if err != ErrOmniFocusNotRunning {
		t.Errorf("expected ErrOmniFocusNotRunning, got %v", err)
	}
}

func TestParseProjects_WithTasks(t *testing.T) {
	jsonStr := `{
		"projects": [
			{
				"id": "proj456",
				"name": "Project with tasks",
				"status": "active",
				"note": "Has nested tasks",
				"tasks": [
					{
						"id": "task1",
						"name": "Task 1",
						"note": "",
						"tags": [],
						"dueDate": null,
						"deferDate": null,
						"flagged": false,
						"completed": false,
						"completedDate": null
					}
				]
			}
		]
	}`

	projects, err := ParseProjects(jsonStr)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	proj := projects[0]
	if len(proj.Tasks) != 1 {
		t.Fatalf("expected 1 task in project, got %d", len(proj.Tasks))
	}

	task := proj.Tasks[0]
	if task.ID != "task1" {
		t.Errorf("expected task ID 'task1', got '%s'", task.ID)
	}
	if task.Name != "Task 1" {
		t.Errorf("expected task name 'Task 1', got '%s'", task.Name)
	}
}
