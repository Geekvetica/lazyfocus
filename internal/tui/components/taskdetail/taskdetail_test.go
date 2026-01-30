package taskdetail

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

	if m.IsVisible() {
		t.Error("new view should not be visible")
	}
	if m.Task() != nil {
		t.Error("new view should have no task")
	}
}

func TestShow(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	task := &domain.Task{ID: "task1", Name: "Test Task"}
	m = m.Show(task)

	if !m.IsVisible() {
		t.Error("view should be visible after Show()")
	}
	if m.Task() == nil {
		t.Error("task should be set")
	}
	if m.Task().ID != "task1" {
		t.Errorf("task ID = %q, want %q", m.Task().ID, "task1")
	}
}

func TestHide(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	task := &domain.Task{ID: "task1", Name: "Test Task"}
	m = m.Show(task).Hide()

	if m.IsVisible() {
		t.Error("view should not be visible after Hide()")
	}
}

func TestUpdate_Escape_Closes(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	task := &domain.Task{ID: "task1", Name: "Test Task"}
	m := New(styles, keys).Show(task).SetSize(80, 24)

	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	m, cmd := m.Update(keyMsg)

	if m.IsVisible() {
		t.Error("view should be hidden after Escape")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(CloseMsg); !ok {
		t.Errorf("expected CloseMsg, got %T", msg)
	}
}

func TestUpdate_EditKey(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	task := &domain.Task{ID: "task1", Name: "Test Task"}
	m := New(styles, keys).Show(task).SetSize(80, 24)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	_, cmd := m.Update(keyMsg)

	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if req, ok := msg.(EditRequestedMsg); !ok {
		t.Errorf("expected EditRequestedMsg, got %T", msg)
	} else if req.Task.ID != "task1" {
		t.Errorf("task ID = %q, want %q", req.Task.ID, "task1")
	}
}

func TestUpdate_CompleteKey(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	task := &domain.Task{ID: "task1", Name: "Test Task"}
	m := New(styles, keys).Show(task).SetSize(80, 24)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	_, cmd := m.Update(keyMsg)

	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if req, ok := msg.(CompleteRequestedMsg); !ok {
		t.Errorf("expected CompleteRequestedMsg, got %T", msg)
	} else if req.TaskID != "task1" {
		t.Errorf("task ID = %q, want %q", req.TaskID, "task1")
	}
}

func TestUpdate_DeleteKey(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	task := &domain.Task{ID: "task1", Name: "Test Task"}
	m := New(styles, keys).Show(task).SetSize(80, 24)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	_, cmd := m.Update(keyMsg)

	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if req, ok := msg.(DeleteRequestedMsg); !ok {
		t.Errorf("expected DeleteRequestedMsg, got %T", msg)
	} else {
		if req.TaskID != "task1" {
			t.Errorf("task ID = %q, want %q", req.TaskID, "task1")
		}
		if req.TaskName != "Test Task" {
			t.Errorf("task name = %q, want %q", req.TaskName, "Test Task")
		}
	}
}

func TestUpdate_FlagKey(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	task := &domain.Task{ID: "task1", Name: "Test Task", Flagged: false}
	m := New(styles, keys).Show(task).SetSize(80, 24)

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}}
	_, cmd := m.Update(keyMsg)

	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if req, ok := msg.(FlagRequestedMsg); !ok {
		t.Errorf("expected FlagRequestedMsg, got %T", msg)
	} else {
		if req.TaskID != "task1" {
			t.Errorf("task ID = %q, want %q", req.TaskID, "task1")
		}
		if !req.Flagged {
			t.Error("Flagged should be true (toggled from false)")
		}
	}
}

func TestUpdate_NotVisible_IgnoresInput(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys) // Not visible

	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	m, cmd := m.Update(keyMsg)

	if cmd != nil {
		t.Error("should not return command when not visible")
	}
}

func TestView_NotVisible(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	view := m.View()

	if view != "" {
		t.Error("view should be empty when not visible")
	}
}

func TestView_Visible_ShowsTaskInfo(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	dueDate := time.Now().Add(24 * time.Hour)
	task := &domain.Task{
		ID:          "task1",
		Name:        "Test Task",
		ProjectName: "Test Project",
		Tags:        []string{"urgent", "work"},
		DueDate:     &dueDate,
		Note:        "This is a test note",
		Flagged:     true,
	}

	m := New(styles, keys).Show(task).SetSize(80, 24)

	view := m.View()

	if view == "" {
		t.Error("view should not be empty")
	}
	if !strings.Contains(view, "Test Task") {
		t.Error("view should contain task name")
	}
	if !strings.Contains(view, "Test Project") {
		t.Error("view should contain project name")
	}
	if !strings.Contains(view, "urgent") {
		t.Error("view should contain tag")
	}
	if !strings.Contains(view, "This is a test note") {
		t.Error("view should contain note")
	}
	if !strings.Contains(view, "🚩") {
		t.Error("view should contain flag icon")
	}
}

func TestSetSize(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	m = m.SetSize(100, 50)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}
