package tasklist

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	m := New(styles, keys)

	if m.cursor != 0 {
		t.Errorf("expected cursor to be 0, got %d", m.cursor)
	}

	if len(m.tasks) != 0 {
		t.Errorf("expected tasks to be empty, got %d tasks", len(m.tasks))
	}

	if m.loading {
		t.Error("expected loading to be false initially")
	}

	if !m.empty {
		t.Error("expected empty to be true initially")
	}

	if m.styles != styles {
		t.Error("expected styles to be set")
	}

	// Keys are stored by value, just verify they're not zero
	if m.keys.Up.Keys() == nil {
		t.Error("expected keys to be set")
	}
}

func TestSetTasks(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())

	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}

	m = m.SetTasks(tasks)

	if len(m.tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(m.tasks))
	}

	if m.empty {
		t.Error("expected empty to be false after setting tasks")
	}

	if m.loading {
		t.Error("expected loading to be false after setting tasks")
	}
}

func TestSetTasksEmpty(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())

	// Set some tasks first
	m = m.SetTasks([]domain.Task{{ID: "1", Name: "Task 1"}})

	// Now set empty list
	m = m.SetTasks([]domain.Task{})

	if len(m.tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(m.tasks))
	}

	if !m.empty {
		t.Error("expected empty to be true after setting empty tasks")
	}

	if m.cursor != 0 {
		t.Error("expected cursor to reset to 0 when tasks are empty")
	}
}

func TestSetLoading(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())

	m = m.SetLoading(true)

	if !m.loading {
		t.Error("expected loading to be true")
	}

	m = m.SetLoading(false)

	if m.loading {
		t.Error("expected loading to be false")
	}
}

func TestSelectedTask(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())

	// No tasks - should return nil
	if task := m.SelectedTask(); task != nil {
		t.Error("expected nil when no tasks")
	}

	// Set tasks
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}
	m = m.SetTasks(tasks)

	// First task selected
	task := m.SelectedTask()
	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ID != "1" {
		t.Errorf("expected task ID '1', got '%s'", task.ID)
	}

	// Move cursor
	m.cursor = 1
	task = m.SelectedTask()
	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ID != "2" {
		t.Errorf("expected task ID '2', got '%s'", task.ID)
	}
}

func TestSelectedIndex(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	})

	if m.SelectedIndex() != 0 {
		t.Errorf("expected index 0, got %d", m.SelectedIndex())
	}

	m.cursor = 1
	if m.SelectedIndex() != 1 {
		t.Errorf("expected index 1, got %d", m.SelectedIndex())
	}
}

func TestNavigationDown(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
		{ID: "3", Name: "Task 3"},
	})

	// Start at 0
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}

	// Move down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", m.cursor)
	}

	// Move down again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 2 {
		t.Errorf("expected cursor at 2, got %d", m.cursor)
	}

	// Move down - should wrap to 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 0 {
		t.Errorf("expected cursor to wrap to 0, got %d", m.cursor)
	}
}

func TestNavigationUp(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
		{ID: "3", Name: "Task 3"},
	})

	// Start at 0, move up - should wrap to end
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 2 {
		t.Errorf("expected cursor to wrap to 2, got %d", m.cursor)
	}

	// Move up again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", m.cursor)
	}

	// Move up again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", m.cursor)
	}
}

func TestNavigationVimKeys(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	})

	// Test 'j' for down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if m.cursor != 1 {
		t.Errorf("expected cursor at 1 after 'j', got %d", m.cursor)
	}

	// Test 'k' for up
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if m.cursor != 0 {
		t.Errorf("expected cursor at 0 after 'k', got %d", m.cursor)
	}
}

func TestNavigationNoTasks(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())

	// Try to navigate with no tasks
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0 with no tasks, got %d", m.cursor)
	}

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 0 {
		t.Errorf("expected cursor to stay at 0 with no tasks, got %d", m.cursor)
	}
}

func TestViewLoading(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetLoading(true)
	m.width = 80
	m.height = 24

	view := m.View()

	if !strings.Contains(view, "Loading...") {
		t.Error("expected view to contain 'Loading...'")
	}
}

func TestViewEmpty(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{})
	m.width = 80
	m.height = 24

	view := m.View()

	if !strings.Contains(view, "No tasks") {
		t.Error("expected view to contain 'No tasks'")
	}
}

func TestViewTasksDisplay(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80
	m.height = 24

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)

	tasks := []domain.Task{
		{ID: "1", Name: "Buy groceries", DueDate: &today},
		{ID: "2", Name: "Call dentist", Flagged: true},
		{ID: "3", Name: "Review PR", DueDate: &tomorrow, Completed: false},
		{ID: "4", Name: "Completed task", Completed: true},
	}
	m = m.SetTasks(tasks)

	view := m.View()

	// Check for checkbox icons
	if !strings.Contains(view, "☐") {
		t.Error("expected view to contain checkbox icon")
	}

	if !strings.Contains(view, "☑") {
		t.Error("expected view to contain checked box icon")
	}

	// Check for task names
	if !strings.Contains(view, "Buy groceries") {
		t.Error("expected view to contain 'Buy groceries'")
	}

	if !strings.Contains(view, "Call dentist") {
		t.Error("expected view to contain 'Call dentist'")
	}

	// Check for flag icon
	if !strings.Contains(view, "🚩") {
		t.Error("expected view to contain flag icon")
	}

	// Check for due date formatting
	if !strings.Contains(view, "Today") || !strings.Contains(view, "📅") {
		t.Error("expected view to contain due date with calendar icon")
	}
}

func TestViewSelection(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80
	m.height = 24

	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}
	m = m.SetTasks(tasks)

	view := m.View()

	// The first task should have selected styling
	// We can't test the exact ANSI codes, but we can verify the view is generated
	if view == "" {
		t.Error("expected non-empty view")
	}
}

func TestResizeHandling(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
	})

	// Send window size message
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	if m.width != 100 {
		t.Errorf("expected width 100, got %d", m.width)
	}

	if m.height != 30 {
		t.Errorf("expected height 30, got %d", m.height)
	}
}

func TestCursorBoundsAfterTaskRemoval(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
		{ID: "3", Name: "Task 3"},
	})

	// Move cursor to last task
	m.cursor = 2

	// Set tasks with fewer items
	m = m.SetTasks([]domain.Task{
		{ID: "1", Name: "Task 1"},
	})

	// Cursor should be clamped to valid range
	if m.cursor >= len(m.tasks) {
		t.Errorf("expected cursor to be clamped, got %d for %d tasks", m.cursor, len(m.tasks))
	}

	if m.cursor != 0 {
		t.Errorf("expected cursor to be 0, got %d", m.cursor)
	}
}

func TestInit(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	cmd := m.Init()

	if cmd != nil {
		t.Error("expected Init() to return nil")
	}
}

func TestUpdateUnhandledMessage(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	originalCursor := m.cursor
	originalTasks := m.tasks

	// Send a message type that isn't handled
	type UnhandledMsg struct{}
	updatedModel, cmd := m.Update(UnhandledMsg{})

	if cmd != nil {
		t.Error("expected Update to return nil cmd for unhandled message")
	}

	if updatedModel.cursor != originalCursor {
		t.Error("expected cursor to remain unchanged for unhandled message")
	}

	if len(updatedModel.tasks) != len(originalTasks) {
		t.Error("expected tasks to remain unchanged for unhandled message")
	}
}

func TestRenderLoadingZeroHeight(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetLoading(true)
	m.height = 0

	view := m.View()

	if view != "Loading..." {
		t.Errorf("expected 'Loading...' for zero height, got '%s'", view)
	}
}

func TestRenderEmptyZeroHeight(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m = m.SetTasks([]domain.Task{})
	m.height = 0

	view := m.View()

	if view != "No tasks" {
		t.Errorf("expected 'No tasks' for zero height, got '%s'", view)
	}
}

func TestFormatDateYesterday(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80
	m.height = 24

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	tasks := []domain.Task{
		{ID: "1", Name: "Yesterday task", DueDate: &yesterday},
	}
	m = m.SetTasks(tasks)

	view := m.View()

	if !strings.Contains(view, "Yesterday") {
		t.Error("expected view to contain 'Yesterday'")
	}
}

func TestFormatDateDifferentYear(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80
	m.height = 24

	now := time.Now()
	differentYear := now.AddDate(-1, 0, 0)

	tasks := []domain.Task{
		{ID: "1", Name: "Last year task", DueDate: &differentYear},
	}
	m = m.SetTasks(tasks)

	view := m.View()

	// Should contain year in format "Jan 2, 2006"
	expectedYear := differentYear.Format("2006")
	if !strings.Contains(view, expectedYear) {
		t.Errorf("expected view to contain year '%s'", expectedYear)
	}

	// Should also contain month
	expectedMonth := differentYear.Format("Jan")
	if !strings.Contains(view, expectedMonth) {
		t.Errorf("expected view to contain month '%s'", expectedMonth)
	}
}

func TestFormatDateSameYearDifferentMonth(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80
	m.height = 24

	now := time.Now()
	// Get a date that's definitely in the same year but different month
	// Add 2 months, or if that crosses year boundary, subtract 2 months
	differentMonth := now.AddDate(0, 2, 0)
	if differentMonth.Year() != now.Year() {
		differentMonth = now.AddDate(0, -2, 0)
	}

	tasks := []domain.Task{
		{ID: "1", Name: "Different month task", DueDate: &differentMonth},
	}
	m = m.SetTasks(tasks)

	view := m.View()

	// Should contain month in format "Jan 2" without year
	expectedFormat := differentMonth.Format("Jan 2")
	if !strings.Contains(view, expectedFormat) {
		t.Errorf("expected view to contain '%s'", expectedFormat)
	}

	// Should NOT contain the year
	year := differentMonth.Format("2006")
	// Split view by lines and check each line doesn't have standalone year
	lines := strings.Split(view, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Different month task") {
			// This line should have the date but not the year suffix
			if strings.Contains(line, ", "+year) {
				t.Errorf("expected same year date to not include year, got '%s'", line)
			}
		}
	}
}

func TestFormatTaskLineSelected(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80

	task := domain.Task{ID: "1", Name: "Test task"}
	line := m.formatTaskLine(task, true)

	// Line should be non-empty and contain task name
	if line == "" {
		t.Error("expected non-empty line for selected task")
	}

	if !strings.Contains(line, "Test task") {
		t.Error("expected line to contain task name")
	}
}

func TestFormatTaskLineCompleted(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80

	task := domain.Task{ID: "1", Name: "Completed task", Completed: true}
	line := m.formatTaskLine(task, false)

	// Line should be non-empty and contain checked checkbox
	if line == "" {
		t.Error("expected non-empty line for completed task")
	}

	if !strings.Contains(line, CheckboxChecked) {
		t.Error("expected line to contain checked checkbox")
	}
}

func TestFormatTaskLineNormal(t *testing.T) {
	m := New(tui.DefaultStyles(), tui.DefaultKeyMap())
	m.width = 80

	task := domain.Task{ID: "1", Name: "Normal task"}
	line := m.formatTaskLine(task, false)

	// Line should be non-empty and contain empty checkbox
	if line == "" {
		t.Error("expected non-empty line for normal task")
	}

	if !strings.Contains(line, CheckboxEmpty) {
		t.Error("expected line to contain empty checkbox")
	}
}
