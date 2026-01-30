package filter

import (
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestMatcher_FilterTasks_SearchText(t *testing.T) {
	tasks := []domain.Task{
		{ID: "1", Name: "Buy groceries"},
		{ID: "2", Name: "Review PR"},
		{ID: "3", Name: "Meeting notes", Note: "need groceries for dinner"},
	}

	matcher := NewMatcher(State{SearchText: "groceries"})
	result := matcher.FilterTasks(tasks)

	if len(result) != 2 {
		t.Errorf("got %d tasks, want 2", len(result))
		for _, task := range result {
			t.Logf("  Found: ID=%s, Name=%s, Note=%s", task.ID, task.Name, task.Note)
		}
	}
}

func TestMatcher_FilterTasks_Project(t *testing.T) {
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", ProjectID: "proj1"},
		{ID: "2", Name: "Task 2", ProjectID: "proj2"},
		{ID: "3", Name: "Task 3", ProjectID: "proj1"},
	}

	matcher := NewMatcher(State{ProjectID: "proj1"})
	result := matcher.FilterTasks(tasks)

	if len(result) != 2 {
		t.Errorf("got %d tasks, want 2", len(result))
	}
}

func TestMatcher_FilterTasks_Tag(t *testing.T) {
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", Tags: []string{"tag1"}},
		{ID: "2", Name: "Task 2", Tags: []string{"tag2"}},
		{ID: "3", Name: "Task 3", Tags: []string{"tag1", "tag2"}},
	}

	matcher := NewMatcher(State{TagID: "tag1"})
	result := matcher.FilterTasks(tasks)

	if len(result) != 2 {
		t.Errorf("got %d tasks, want 2", len(result))
	}
}

func TestMatcher_FilterTasks_Flagged(t *testing.T) {
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", Flagged: true},
		{ID: "2", Name: "Task 2", Flagged: false},
		{ID: "3", Name: "Task 3", Flagged: true},
	}

	matcher := NewMatcher(State{FlaggedOnly: true})
	result := matcher.FilterTasks(tasks)

	if len(result) != 2 {
		t.Errorf("got %d tasks, want 2", len(result))
	}
}

func TestMatcher_FilterTasks_DueToday(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tasks := []domain.Task{
		{ID: "1", Name: "Yesterday", DueDate: &yesterday},
		{ID: "2", Name: "Today", DueDate: &today},
		{ID: "3", Name: "Tomorrow", DueDate: &tomorrow},
		{ID: "4", Name: "No due"},
	}

	matcher := NewMatcher(State{DueFilter: DueToday})
	result := matcher.FilterTasks(tasks)

	if len(result) != 1 {
		t.Errorf("got %d tasks, want 1", len(result))
	}
	if result[0].Name != "Today" {
		t.Errorf("got task %q, want %q", result[0].Name, "Today")
	}
}

func TestMatcher_FilterTasks_NoFilter(t *testing.T) {
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}

	matcher := NewMatcher(State{})
	result := matcher.FilterTasks(tasks)

	if len(result) != 2 {
		t.Errorf("got %d tasks, want 2", len(result))
	}
}
