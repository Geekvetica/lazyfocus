package searchinput

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	if m.IsVisible() {
		t.Error("new search input should not be visible")
	}
}

func TestShow(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	m = m.Show()

	if !m.IsVisible() {
		t.Error("search input should be visible after Show()")
	}
	if m.Value() != "" {
		t.Error("search text should be empty after Show()")
	}
}

func TestHide(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show()

	m = m.Hide()

	if m.IsVisible() {
		t.Error("search input should not be visible after Hide()")
	}
}

func TestUpdate_Typing(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Type a character
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})

	if m.Value() != "t" {
		t.Errorf("value = %q, want %q", m.Value(), "t")
	}
	if cmd == nil {
		t.Error("expected SearchChangedMsg command")
	}
}

func TestUpdate_Escape_Clears(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Type something
	m.input.SetValue("test")

	// Press Escape
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if m.IsVisible() {
		t.Error("should be hidden after Escape")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(SearchClearedMsg); !ok {
		t.Errorf("expected SearchClearedMsg, got %T", msg)
	}
}

func TestUpdate_Enter_Confirms(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Type something
	m.input.SetValue("search text")

	// Press Enter
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.IsVisible() {
		t.Error("should be hidden after Enter")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	confirmed, ok := msg.(SearchConfirmedMsg)
	if !ok {
		t.Fatalf("expected SearchConfirmedMsg, got %T", msg)
	}
	if confirmed.Text != "search text" {
		t.Errorf("text = %q, want %q", confirmed.Text, "search text")
	}
}

func TestView_NotVisible(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	view := m.View()

	if view != "" {
		t.Error("view should be empty when not visible")
	}
}
