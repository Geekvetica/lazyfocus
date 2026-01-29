package tui

import (
	"errors"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestTasksLoadedMsg(t *testing.T) {
	tasks := []domain.Task{
		{ID: "task1", Name: "Test Task 1"},
		{ID: "task2", Name: "Test Task 2"},
	}

	msg := TasksLoadedMsg{Tasks: tasks}

	if len(msg.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(msg.Tasks))
	}
	if msg.Tasks[0].ID != "task1" {
		t.Errorf("expected task ID 'task1', got '%s'", msg.Tasks[0].ID)
	}
}

func TestProjectsLoadedMsg(t *testing.T) {
	projects := []domain.Project{
		{ID: "proj1", Name: "Project 1", Status: "active"},
		{ID: "proj2", Name: "Project 2", Status: "on-hold"},
	}

	msg := ProjectsLoadedMsg{Projects: projects}

	if len(msg.Projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(msg.Projects))
	}
	if msg.Projects[0].Name != "Project 1" {
		t.Errorf("expected project name 'Project 1', got '%s'", msg.Projects[0].Name)
	}
}

func TestTagsLoadedMsg(t *testing.T) {
	tags := []domain.Tag{
		{ID: "tag1", Name: "work"},
		{ID: "tag2", Name: "urgent"},
	}

	msg := TagsLoadedMsg{Tags: tags}

	if len(msg.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(msg.Tags))
	}
	if msg.Tags[1].Name != "urgent" {
		t.Errorf("expected tag name 'urgent', got '%s'", msg.Tags[1].Name)
	}
}

func TestTaskCreatedMsg(t *testing.T) {
	now := time.Now()
	task := domain.Task{
		ID:      "new-task",
		Name:    "New Task",
		Flagged: true,
		DueDate: &now,
	}

	msg := TaskCreatedMsg{Task: task}

	if msg.Task.ID != "new-task" {
		t.Errorf("expected task ID 'new-task', got '%s'", msg.Task.ID)
	}
	if !msg.Task.Flagged {
		t.Error("expected task to be flagged")
	}
}

func TestTaskCompletedMsg(t *testing.T) {
	msg := TaskCompletedMsg{
		TaskID:   "completed-task",
		TaskName: "Completed Task",
	}

	if msg.TaskID != "completed-task" {
		t.Errorf("expected task ID 'completed-task', got '%s'", msg.TaskID)
	}
	if msg.TaskName != "Completed Task" {
		t.Errorf("expected task name 'Completed Task', got '%s'", msg.TaskName)
	}
}

func TestTaskDeletedMsg(t *testing.T) {
	msg := TaskDeletedMsg{
		TaskID:   "deleted-task",
		TaskName: "Deleted Task",
	}

	if msg.TaskID != "deleted-task" {
		t.Errorf("expected task ID 'deleted-task', got '%s'", msg.TaskID)
	}
	if msg.TaskName != "Deleted Task" {
		t.Errorf("expected task name 'Deleted Task', got '%s'", msg.TaskName)
	}
}

func TestTaskModifiedMsg(t *testing.T) {
	task := domain.Task{
		ID:   "modified-task",
		Name: "Modified Task",
		Tags: []string{"urgent", "work"},
	}

	msg := TaskModifiedMsg{Task: task}

	if msg.Task.ID != "modified-task" {
		t.Errorf("expected task ID 'modified-task', got '%s'", msg.Task.ID)
	}
	if len(msg.Task.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(msg.Task.Tags))
	}
}

func TestErrorMsg(t *testing.T) {
	err := errors.New("something went wrong")
	msg := ErrorMsg{Err: err}

	if msg.Err == nil {
		t.Error("expected error to be set")
	}
	if msg.Err.Error() != "something went wrong" {
		t.Errorf("expected error message 'something went wrong', got '%s'", msg.Err.Error())
	}
}

func TestViewChangedMsg(t *testing.T) {
	tests := []struct {
		name     string
		view     int
		expected int
	}{
		{"inbox view", ViewInbox, 1},
		{"projects view", ViewProjects, 2},
		{"tags view", ViewTags, 3},
		{"forecast view", ViewForecast, 4},
		{"review view", ViewReview, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := ViewChangedMsg{View: tt.view}
			if msg.View != tt.expected {
				t.Errorf("expected view %d, got %d", tt.expected, msg.View)
			}
		})
	}
}

func TestClearErrorMsg(t *testing.T) {
	msg := ClearErrorMsg{}
	// Just verify it can be instantiated
	_ = msg
}

func TestViewConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant int
		expected int
	}{
		{"ViewInbox", ViewInbox, 1},
		{"ViewProjects", ViewProjects, 2},
		{"ViewTags", ViewTags, 3},
		{"ViewForecast", ViewForecast, 4},
		{"ViewReview", ViewReview, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("expected %s to be %d, got %d", tt.name, tt.expected, tt.constant)
			}
		})
	}
}

func TestViewConstantsAreUnique(t *testing.T) {
	views := []int{ViewInbox, ViewProjects, ViewTags, ViewForecast, ViewReview}
	seen := make(map[int]bool)

	for _, view := range views {
		if seen[view] {
			t.Errorf("duplicate view constant: %d", view)
		}
		seen[view] = true
	}

	if len(seen) != 5 {
		t.Errorf("expected 5 unique view constants, got %d", len(seen))
	}
}
