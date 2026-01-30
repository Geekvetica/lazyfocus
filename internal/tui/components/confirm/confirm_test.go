package confirm

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	if m.IsVisible() {
		t.Error("new modal should not be visible")
	}
}

func TestShow(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	m = m.Show("Test Title", "Test message")

	if !m.IsVisible() {
		t.Error("modal should be visible after Show()")
	}
	if m.title != "Test Title" {
		t.Errorf("title = %q, want %q", m.title, "Test Title")
	}
	if m.message != "Test message" {
		t.Errorf("message = %q, want %q", m.message, "Test message")
	}
}

func TestShowWithContext(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	ctx := struct{ TaskID string }{"task123"}
	m = m.ShowWithContext("Delete", "Are you sure?", ctx)

	if !m.IsVisible() {
		t.Error("modal should be visible")
	}
	if m.context == nil {
		t.Error("context should be set")
	}
}

func TestHide(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show("Title", "Message")

	m = m.Hide()

	if m.IsVisible() {
		t.Error("modal should not be visible after Hide()")
	}
}

func TestUpdate_Confirm_Y(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show("Title", "Message")

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	m, cmd := m.Update(keyMsg)

	if m.IsVisible() {
		t.Error("modal should be hidden after confirm")
	}

	if cmd == nil {
		t.Fatal("cmd should not be nil")
	}

	// Execute the command and check the message type
	msg := cmd()
	if _, ok := msg.(ConfirmedMsg); !ok {
		t.Errorf("expected ConfirmedMsg, got %T", msg)
	}
}

func TestUpdate_Confirm_Enter(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show("Title", "Message")

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, cmd := m.Update(keyMsg)

	if m.IsVisible() {
		t.Error("modal should be hidden after confirm")
	}

	if cmd == nil {
		t.Fatal("cmd should not be nil")
	}

	msg := cmd()
	if _, ok := msg.(ConfirmedMsg); !ok {
		t.Errorf("expected ConfirmedMsg, got %T", msg)
	}
}

func TestUpdate_Cancel_N(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show("Title", "Message")

	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	m, cmd := m.Update(keyMsg)

	if m.IsVisible() {
		t.Error("modal should be hidden after cancel")
	}

	if cmd == nil {
		t.Fatal("cmd should not be nil")
	}

	msg := cmd()
	if _, ok := msg.(CancelledMsg); !ok {
		t.Errorf("expected CancelledMsg, got %T", msg)
	}
}

func TestUpdate_Cancel_Escape(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show("Title", "Message")

	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	m, cmd := m.Update(keyMsg)

	if m.IsVisible() {
		t.Error("modal should be hidden after cancel")
	}

	if cmd != nil {
		msg := cmd()
		if _, ok := msg.(CancelledMsg); !ok {
			t.Errorf("expected CancelledMsg, got %T", msg)
		}
	}
}

func TestUpdate_ContextPassedThrough(t *testing.T) {
	styles := tui.DefaultStyles()
	ctx := "task123"
	m := New(styles).ShowWithContext("Delete", "Sure?", ctx)

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, cmd := m.Update(keyMsg)

	if cmd == nil {
		t.Fatal("cmd should not be nil")
	}

	msg := cmd()
	confirmed, ok := msg.(ConfirmedMsg)
	if !ok {
		t.Fatalf("expected ConfirmedMsg, got %T", msg)
	}
	if confirmed.Context != ctx {
		t.Errorf("context = %v, want %v", confirmed.Context, ctx)
	}
}

func TestUpdate_NotVisibleIgnoresInput(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles) // Not visible

	keyMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, cmd := m.Update(keyMsg)

	if cmd != nil {
		t.Error("should not return command when not visible")
	}
}

func TestView_NotVisible(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	view := m.View()

	if view != "" {
		t.Errorf("view should be empty when not visible, got %q", view)
	}
}

func TestView_Visible(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show("Delete Task", "Are you sure?").SetSize(80, 24)

	view := m.View()

	if view == "" {
		t.Error("view should not be empty when visible")
	}
	if !strings.Contains(view, "Delete Task") {
		t.Error("view should contain title")
	}
	if !strings.Contains(view, "Are you sure?") {
		t.Error("view should contain message")
	}
	if !strings.Contains(view, "Confirm") {
		t.Error("view should contain confirm hint")
	}
	if !strings.Contains(view, "Cancel") {
		t.Error("view should contain cancel hint")
	}
}

func TestSetSize(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	m = m.SetSize(100, 50)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}
