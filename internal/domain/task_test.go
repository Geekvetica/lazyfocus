package domain

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTask_JSONSerialization(t *testing.T) {
	dueDate := time.Date(2026, 1, 30, 15, 0, 0, 0, time.UTC)
	deferDate := time.Date(2026, 1, 28, 9, 0, 0, 0, time.UTC)
	completedDate := time.Date(2026, 1, 27, 14, 30, 0, 0, time.UTC)

	task := Task{
		ID:            "task-123",
		Name:          "Buy groceries",
		Note:          "Need milk and eggs",
		ProjectID:     "proj-456",
		ProjectName:   "Personal Errands",
		Tags:          []string{"errands", "shopping"},
		DueDate:       &dueDate,
		DeferDate:     &deferDate,
		Flagged:       true,
		Completed:     false,
		CompletedDate: &completedDate,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	var unmarshaled Task
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if unmarshaled.ID != task.ID {
		t.Errorf("ID mismatch: got %s, want %s", unmarshaled.ID, task.ID)
	}
	if unmarshaled.Name != task.Name {
		t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, task.Name)
	}
	if unmarshaled.Note != task.Note {
		t.Errorf("Note mismatch: got %s, want %s", unmarshaled.Note, task.Note)
	}
	if unmarshaled.ProjectID != task.ProjectID {
		t.Errorf("ProjectID mismatch: got %s, want %s", unmarshaled.ProjectID, task.ProjectID)
	}
	if unmarshaled.ProjectName != task.ProjectName {
		t.Errorf("ProjectName mismatch: got %s, want %s", unmarshaled.ProjectName, task.ProjectName)
	}
	if len(unmarshaled.Tags) != len(task.Tags) {
		t.Errorf("Tags length mismatch: got %d, want %d", len(unmarshaled.Tags), len(task.Tags))
	}
	if unmarshaled.Flagged != task.Flagged {
		t.Errorf("Flagged mismatch: got %v, want %v", unmarshaled.Flagged, task.Flagged)
	}
	if unmarshaled.Completed != task.Completed {
		t.Errorf("Completed mismatch: got %v, want %v", unmarshaled.Completed, task.Completed)
	}
	if !unmarshaled.DueDate.Equal(*task.DueDate) {
		t.Errorf("DueDate mismatch: got %v, want %v", unmarshaled.DueDate, task.DueDate)
	}
	if !unmarshaled.DeferDate.Equal(*task.DeferDate) {
		t.Errorf("DeferDate mismatch: got %v, want %v", unmarshaled.DeferDate, task.DeferDate)
	}
	if !unmarshaled.CompletedDate.Equal(*task.CompletedDate) {
		t.Errorf("CompletedDate mismatch: got %v, want %v", unmarshaled.CompletedDate, task.CompletedDate)
	}
}

func TestTask_OmitEmptyFields(t *testing.T) {
	task := Task{
		ID:        "task-123",
		Name:      "Simple task",
		Flagged:   false,
		Completed: false,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	jsonString := string(jsonData)

	// Check that empty fields are omitted
	if contains(jsonString, "note") {
		t.Error("Expected 'note' field to be omitted when empty")
	}
	if contains(jsonString, "projectId") {
		t.Error("Expected 'projectId' field to be omitted when empty")
	}
	if contains(jsonString, "projectName") {
		t.Error("Expected 'projectName' field to be omitted when empty")
	}
	if contains(jsonString, "tags") {
		t.Error("Expected 'tags' field to be omitted when empty")
	}
	if contains(jsonString, "dueDate") {
		t.Error("Expected 'dueDate' field to be omitted when nil")
	}
	if contains(jsonString, "deferDate") {
		t.Error("Expected 'deferDate' field to be omitted when nil")
	}
	if contains(jsonString, "completedDate") {
		t.Error("Expected 'completedDate' field to be omitted when nil")
	}

	// Check that required fields are present
	if !contains(jsonString, "id") {
		t.Error("Expected 'id' field to be present")
	}
	if !contains(jsonString, "name") {
		t.Error("Expected 'name' field to be present")
	}
	if !contains(jsonString, "flagged") {
		t.Error("Expected 'flagged' field to be present")
	}
	if !contains(jsonString, "completed") {
		t.Error("Expected 'completed' field to be present")
	}
}

func TestTask_ISO8601DateFormat(t *testing.T) {
	jsonInput := `{
		"id": "task-123",
		"name": "Test task",
		"dueDate": "2026-01-30T15:00:00Z",
		"deferDate": "2026-01-28T09:00:00Z",
		"flagged": false,
		"completed": false
	}`

	var task Task
	err := json.Unmarshal([]byte(jsonInput), &task)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	expectedDue := time.Date(2026, 1, 30, 15, 0, 0, 0, time.UTC)
	expectedDefer := time.Date(2026, 1, 28, 9, 0, 0, 0, time.UTC)

	if task.DueDate == nil {
		t.Fatal("Expected DueDate to be non-nil")
	}
	if !task.DueDate.Equal(expectedDue) {
		t.Errorf("DueDate mismatch: got %v, want %v", task.DueDate, expectedDue)
	}

	if task.DeferDate == nil {
		t.Fatal("Expected DeferDate to be non-nil")
	}
	if !task.DeferDate.Equal(expectedDefer) {
		t.Errorf("DeferDate mismatch: got %v, want %v", task.DeferDate, expectedDefer)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
