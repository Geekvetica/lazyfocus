package commandinput

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

func TestUpdate_AllCommandTypes(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantName    string
		wantArgsLen int
	}{
		{"quit command", "quit", "quit", 0},
		{"quit alias q", "q", "quit", 0},
		{"quit alias exit", "exit", "quit", 0},
		{"refresh command", "refresh", "refresh", 0},
		{"refresh alias w", "w", "refresh", 0},
		{"refresh alias sync", "sync", "refresh", 0},
		{"add command", "add Buy milk", "add", 2}, // "Buy" and "milk" are separate args
		{"add alias a", "a Task name", "add", 2},
		{"complete command", "complete", "complete", 0},
		{"complete alias done", "done", "complete", 0},
		{"complete alias c", "c", "complete", 0},
		{"delete command", "delete", "delete", 0},
		{"delete alias del", "del", "delete", 0},
		{"delete alias rm", "rm", "delete", 0},
		{"project filter", "project Work", "project", 1},
		{"project alias p", "p Work", "project", 1},
		{"tag filter", "tag urgent", "tag", 1},
		{"tag alias t", "t urgent", "tag", 1},
		{"due filter", "due today", "due", 1},
		{"due filter tomorrow", "due tomorrow", "due", 1},
		{"flagged filter", "flagged", "flagged", 0},
		{"clear filters", "clear", "clear", 0},
		{"clear alias reset", "reset", "clear", 0},
		{"help command", "help", "help", 0},
		{"help alias ?", "?", "help", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			styles := tui.DefaultStyles()
			m := New(styles).Show().SetWidth(80)
			m.input.SetValue(tt.input)

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
			if execMsg.Command.Name != tt.wantName {
				t.Errorf("command name = %q, want %q", execMsg.Command.Name, tt.wantName)
			}
			if len(execMsg.Command.Args) != tt.wantArgsLen {
				t.Errorf("args len = %d, want %d", len(execMsg.Command.Args), tt.wantArgsLen)
			}
		})
	}
}

func TestUpdate_EmptyCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Empty input
	m.input.SetValue("")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.IsVisible() {
		t.Error("should be hidden after Enter")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(CommandCancelledMsg); !ok {
		t.Errorf("expected CommandCancelledMsg, got %T", msg)
	}
}

func TestTabCompletion_PartialMatch(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Type partial command "re" which matches "refresh" and "reset"
	m.input.SetValue("re")

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Should complete to one of the matching commands
	value := m.input.Value()
	if value != "refresh" && value != "reset" {
		t.Errorf("value = %q, expected refresh or reset", value)
	}

	// Completions should be stored for multiple matches
	if len(m.completions) < 2 {
		t.Errorf("completions len = %d, expected at least 2", len(m.completions))
	}

	// The compIdx should have been incremented
	if m.compIdx == 0 {
		t.Error("expected compIdx to be incremented on first tab")
	}
}

func TestTabCompletion_NoMatch(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Type something that doesn't match
	m.input.SetValue("xyz")
	originalValue := m.input.Value()

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Should not change input
	if m.input.Value() != originalValue {
		t.Errorf("value changed from %q to %q, should remain unchanged", originalValue, m.input.Value())
	}
}

func TestTabCompletion_EmptyInput(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Press Tab on empty input
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Should handle gracefully (either show completions or do nothing)
	// Check that model is still functional
	if !m.IsVisible() {
		t.Error("should remain visible")
	}
}

func TestHistory_UpArrowAtBeginning(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	// Add commands to history
	m.history = []string{"quit", "help", "clear"}
	m = m.Show().SetWidth(80)

	// Press up from beginning (historyIdx = -1)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Should show last command
	if m.input.Value() != "clear" {
		t.Errorf("value = %q, want %q", m.input.Value(), "clear")
	}
	if m.historyIdx != 2 {
		t.Errorf("historyIdx = %d, want 2", m.historyIdx)
	}
}

func TestHistory_UpArrowAtTop(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	// Add commands to history
	m.history = []string{"quit", "help", "clear"}
	m = m.Show().SetWidth(80)

	// Navigate to first item
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp}) // Goes to "clear" (last)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp}) // Goes to "help"
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp}) // Goes to "quit" (first)

	// Now at index 0 (first item)
	if m.historyIdx != 0 {
		t.Fatalf("historyIdx = %d, want 0 for setup", m.historyIdx)
	}

	// Press up from first item
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Should stay at first item (can't go further back)
	if m.historyIdx != 0 {
		t.Errorf("historyIdx = %d, want 0", m.historyIdx)
	}
	if m.input.Value() != "quit" {
		t.Errorf("value = %q, want %q", m.input.Value(), "quit")
	}
}

func TestHistory_DownArrowAtEnd(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	// Add commands to history
	m.history = []string{"quit", "help", "clear"}
	m.historyIdx = 2 // At last item
	m = m.Show().SetWidth(80)

	// Press down from last item
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Should clear input and reset historyIdx
	if m.input.Value() != "" {
		t.Errorf("value = %q, want empty", m.input.Value())
	}
	if m.historyIdx != -1 {
		t.Errorf("historyIdx = %d, want -1", m.historyIdx)
	}
}

func TestHistory_DownArrowWhenNotNavigating(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	// Add commands to history
	m.history = []string{"quit", "help", "clear"}
	m.historyIdx = -1 // Not navigating
	m = m.Show().SetWidth(80)

	// Press down when not in history navigation
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Should do nothing
	if m.historyIdx != -1 {
		t.Errorf("historyIdx = %d, want -1", m.historyIdx)
	}
}

func TestHistory_EmptyHistory(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Press up with empty history
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})

	// Should do nothing
	if m.input.Value() != "" {
		t.Errorf("value = %q, want empty", m.input.Value())
	}
}

func TestHistory_AddsCommandOnExecute(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	m.input.SetValue("quit")
	historyLen := len(m.history)

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// History should grow by 1
	if len(m.history) != historyLen+1 {
		t.Errorf("history len = %d, want %d", len(m.history), historyLen+1)
	}
	if m.history[len(m.history)-1] != "quit" {
		t.Errorf("last history item = %q, want %q", m.history[len(m.history)-1], "quit")
	}
}

func TestFocused_WhenVisible(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show()

	// Model should be visible and input focused
	if !m.IsVisible() {
		t.Error("should be visible")
	}
	if !m.input.Focused() {
		t.Error("input should be focused")
	}
}

func TestNotFocused_WhenHidden(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show()
	m = m.Hide()

	// Model should not be visible and input not focused
	if m.IsVisible() {
		t.Error("should not be visible")
	}
	if m.input.Focused() {
		t.Error("input should not be focused")
	}
}

func TestNotFocused_IgnoresInput(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles) // Not shown/focused

	// Try to update with a key press
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	// Should return nil command (no action)
	if cmd != nil {
		t.Error("expected nil command when not visible")
	}
	if m.input.Value() != "" {
		t.Error("input should remain empty when not visible")
	}
}

func TestView_NotVisible_ReturnsEmpty(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles) // Not shown

	view := m.View()
	if view != "" {
		t.Error("view should be empty when not visible")
	}
}

func TestView_Visible_ShowsInput(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	view := m.View()
	if view == "" {
		t.Error("view should not be empty when visible")
	}
	// Should contain the prompt
	if !strings.Contains(view, ":") {
		t.Error("view should contain prompt")
	}
}

func TestView_ShowsCompletions(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Type partial command and trigger completions
	m.input.SetValue("re")
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Set completions manually to test view rendering
	m.completions = []string{"refresh", "reset"}

	view := m.View()
	// View should show completions hint
	if !strings.Contains(view, "refresh") || !strings.Contains(view, "reset") {
		t.Error("view should show completions when multiple matches exist")
	}
}

func TestSetWidth(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).SetWidth(100)

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.input.Width != 96 { // 100 - 4
		t.Errorf("input width = %d, want 96", m.input.Width)
	}
}

func TestInit(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	cmd := m.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestUpdate_WithColon_ParsesCorrectly(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	// Commands can be entered with or without leading colon
	m.input.SetValue(":quit")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

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

func TestUpdate_WithQuotedArgs(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles).Show().SetWidth(80)

	m.input.SetValue("add \"Buy milk and eggs\"")

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	execMsg, ok := msg.(CommandExecutedMsg)
	if !ok {
		t.Fatalf("expected CommandExecutedMsg, got %T", msg)
	}
	if execMsg.Command.Name != "add" {
		t.Errorf("command name = %q, want %q", execMsg.Command.Name, "add")
	}
	if len(execMsg.Command.Args) != 1 {
		t.Fatalf("args len = %d, want 1", len(execMsg.Command.Args))
	}
	if execMsg.Command.Args[0] != "Buy milk and eggs" {
		t.Errorf("args[0] = %q, want %q", execMsg.Command.Args[0], "Buy milk and eggs")
	}
}
