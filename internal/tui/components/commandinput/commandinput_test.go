package commandinput

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	if m.IsVisible() {
		t.Error("new command input should not be visible")
	}
}

func TestShow(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	m = m.Show()

	if !m.IsVisible() {
		t.Error("command input should be visible after Show()")
	}
}

func TestHide(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show()

	m = m.Hide()

	if m.IsVisible() {
		t.Error("command input should not be visible after Hide()")
	}
}

func TestUpdate_Escape_Cancels(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if m.IsVisible() {
		t.Error("should be hidden after Escape")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(CommandCancelledMsg); !ok {
		t.Errorf("expected CommandCancelledMsg, got %T", msg)
	}
}

func TestUpdate_Enter_ExecutesCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	m.input.SetValue("quit")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.IsVisible() {
		t.Error("should be hidden after Enter")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	execMsg, ok := msg.(CommandExecutedMsg)
	if !ok {
		t.Fatalf("expected CommandExecutedMsg, got %T", msg)
	}
	if execMsg.Command.Name != "quit" {
		t.Errorf("command name = %q, want %q", execMsg.Command.Name, "quit")
	}
}

func TestUpdate_Enter_InvalidCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	m.input.SetValue("invalid_command")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.IsVisible() {
		t.Error("should be hidden after Enter")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(CommandErrorMsg); !ok {
		t.Errorf("expected CommandErrorMsg, got %T", msg)
	}
}

func TestUpdate_History_Navigation(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	// Add commands to history
	m.history = []string{"quit", "help", "clear"}
	m = m.Show().SetWidth(80)

	// Press up to get last command
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.input.Value() != "clear" {
		t.Errorf("value = %q, want %q", m.input.Value(), "clear")
	}

	// Press up again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.input.Value() != "help" {
		t.Errorf("value = %q, want %q", m.input.Value(), "help")
	}

	// Press down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.input.Value() != "clear" {
		t.Errorf("value = %q, want %q", m.input.Value(), "clear")
	}
}

func TestUpdate_TabCompletion(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	m.input.SetValue("qu")

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	if m.input.Value() != "quit" {
		t.Errorf("value = %q, want %q", m.input.Value(), "quit")
	}
}
