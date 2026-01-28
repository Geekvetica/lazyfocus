package domain

import (
	"encoding/json"
	"testing"
)

func TestProject_JSONSerialization(t *testing.T) {
	task1 := Task{
		ID:        "task-1",
		Name:      "First task",
		Flagged:   false,
		Completed: false,
	}
	task2 := Task{
		ID:        "task-2",
		Name:      "Second task",
		Flagged:   true,
		Completed: false,
	}

	project := Project{
		ID:     "proj-123",
		Name:   "My Project",
		Status: "active",
		Note:   "Project notes here",
		Tasks:  []Task{task1, task2},
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("Failed to marshal project: %v", err)
	}

	var unmarshaled Project
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal project: %v", err)
	}

	if unmarshaled.ID != project.ID {
		t.Errorf("ID mismatch: got %s, want %s", unmarshaled.ID, project.ID)
	}
	if unmarshaled.Name != project.Name {
		t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, project.Name)
	}
	if unmarshaled.Status != project.Status {
		t.Errorf("Status mismatch: got %s, want %s", unmarshaled.Status, project.Status)
	}
	if unmarshaled.Note != project.Note {
		t.Errorf("Note mismatch: got %s, want %s", unmarshaled.Note, project.Note)
	}
	if len(unmarshaled.Tasks) != len(project.Tasks) {
		t.Errorf("Tasks length mismatch: got %d, want %d", len(unmarshaled.Tasks), len(project.Tasks))
	}
	if len(unmarshaled.Tasks) >= 2 {
		if unmarshaled.Tasks[0].ID != task1.ID {
			t.Errorf("First task ID mismatch: got %s, want %s", unmarshaled.Tasks[0].ID, task1.ID)
		}
		if unmarshaled.Tasks[1].ID != task2.ID {
			t.Errorf("Second task ID mismatch: got %s, want %s", unmarshaled.Tasks[1].ID, task2.ID)
		}
	}
}

func TestProject_OmitEmptyFields(t *testing.T) {
	project := Project{
		ID:     "proj-123",
		Name:   "Simple Project",
		Status: "active",
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("Failed to marshal project: %v", err)
	}

	jsonString := string(jsonData)

	// Check that empty fields are omitted
	if contains(jsonString, "note") {
		t.Error("Expected 'note' field to be omitted when empty")
	}
	if contains(jsonString, "tasks") {
		t.Error("Expected 'tasks' field to be omitted when empty")
	}

	// Check that required fields are present
	if !contains(jsonString, "id") {
		t.Error("Expected 'id' field to be present")
	}
	if !contains(jsonString, "name") {
		t.Error("Expected 'name' field to be present")
	}
	if !contains(jsonString, "status") {
		t.Error("Expected 'status' field to be present")
	}
}

func TestProject_StatusValues(t *testing.T) {
	testCases := []struct {
		name   string
		status string
	}{
		{"Active project", "active"},
		{"On-hold project", "on-hold"},
		{"Completed project", "completed"},
		{"Dropped project", "dropped"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			project := Project{
				ID:     "proj-123",
				Name:   tc.name,
				Status: tc.status,
			}

			jsonData, err := json.Marshal(project)
			if err != nil {
				t.Fatalf("Failed to marshal project: %v", err)
			}

			var unmarshaled Project
			err = json.Unmarshal(jsonData, &unmarshaled)
			if err != nil {
				t.Fatalf("Failed to unmarshal project: %v", err)
			}

			if unmarshaled.Status != tc.status {
				t.Errorf("Status mismatch: got %s, want %s", unmarshaled.Status, tc.status)
			}
		})
	}
}

func TestProject_WithNestedTasks(t *testing.T) {
	jsonInput := `{
		"id": "proj-123",
		"name": "Test Project",
		"status": "active",
		"tasks": [
			{
				"id": "task-1",
				"name": "First task",
				"flagged": true,
				"completed": false
			},
			{
				"id": "task-2",
				"name": "Second task",
				"flagged": false,
				"completed": true
			}
		]
	}`

	var project Project
	err := json.Unmarshal([]byte(jsonInput), &project)
	if err != nil {
		t.Fatalf("Failed to unmarshal project: %v", err)
	}

	if len(project.Tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(project.Tasks))
	}

	if project.Tasks[0].ID != "task-1" {
		t.Errorf("First task ID mismatch: got %s, want task-1", project.Tasks[0].ID)
	}
	if project.Tasks[0].Name != "First task" {
		t.Errorf("First task name mismatch: got %s, want First task", project.Tasks[0].Name)
	}
	if !project.Tasks[0].Flagged {
		t.Error("Expected first task to be flagged")
	}
	if project.Tasks[0].Completed {
		t.Error("Expected first task to not be completed")
	}

	if project.Tasks[1].ID != "task-2" {
		t.Errorf("Second task ID mismatch: got %s, want task-2", project.Tasks[1].ID)
	}
	if !project.Tasks[1].Completed {
		t.Error("Expected second task to be completed")
	}
}
