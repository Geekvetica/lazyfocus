package taskedit

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
	m := New(styles)

	if m.IsVisible() {
		t.Error("new overlay should not be visible")
	}
	if len(m.inputs) != NumFields {
		t.Errorf("inputs count = %d, want %d", len(m.inputs), NumFields)
	}
}

func TestShow(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	dueDate := time.Now().Add(24 * time.Hour)
	task := &domain.Task{
		ID:          "task1",
		Name:        "Test Task",
		Note:        "Test note",
		ProjectName: "Test Project",
		Tags:        []string{"urgent", "work"},
		DueDate:     &dueDate,
		Flagged:     true,
	}

	m = m.Show(task)

	if !m.IsVisible() {
		t.Error("overlay should be visible after Show()")
	}
	if m.inputs[FieldName].Value() != "Test Task" {
		t.Errorf("name = %q, want %q", m.inputs[FieldName].Value(), "Test Task")
	}
	if m.inputs[FieldNote].Value() != "Test note" {
		t.Errorf("note = %q, want %q", m.inputs[FieldNote].Value(), "Test note")
	}
	if m.inputs[FieldProject].Value() != "Test Project" {
		t.Errorf("project = %q, want %q", m.inputs[FieldProject].Value(), "Test Project")
	}
	if !strings.Contains(m.inputs[FieldTags].Value(), "urgent") {
		t.Error("tags should contain 'urgent'")
	}
	if !m.flagged {
		t.Error("flagged should be true")
	}
}

func TestHide(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Test"}
	m = m.Show(task).Hide()

	if m.IsVisible() {
		t.Error("overlay should not be visible after Hide()")
	}
}

func TestUpdate_Escape_Cancels(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Test"}
	m = m.Show(task).SetSize(80, 24)

	keyMsg := tea.KeyMsg{Type: tea.KeyEscape}
	m, cmd := m.Update(keyMsg)

	if m.IsVisible() {
		t.Error("overlay should be hidden after Escape")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(CancelMsg); !ok {
		t.Errorf("expected CancelMsg, got %T", msg)
	}
}

func TestUpdate_Tab_Navigation(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Test"}
	m = m.Show(task).SetSize(80, 24)

	if m.focusIndex != 0 {
		t.Errorf("initial focus = %d, want 0", m.focusIndex)
	}

	// Tab forward
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if m.focusIndex != 1 {
		t.Errorf("focus after tab = %d, want 1", m.focusIndex)
	}

	// Shift+Tab backward
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if m.focusIndex != 0 {
		t.Errorf("focus after shift+tab = %d, want 0", m.focusIndex)
	}
}

func TestUpdate_Enter_Saves(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Original", Flagged: false}
	m = m.Show(task).SetSize(80, 24)

	// Change the name
	m.inputs[FieldName].SetValue("Updated Name")

	// Press Enter
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.IsVisible() {
		t.Error("overlay should be hidden after save")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	saveMsg, ok := msg.(SaveMsg)
	if !ok {
		t.Fatalf("expected SaveMsg, got %T", msg)
	}
	if saveMsg.TaskID != "task1" {
		t.Errorf("TaskID = %q, want %q", saveMsg.TaskID, "task1")
	}
	if saveMsg.Modification.Name == nil || *saveMsg.Modification.Name != "Updated Name" {
		t.Error("modification should include name change")
	}
}

func TestValidation_EmptyName(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Test"}
	m = m.Show(task).SetSize(80, 24)

	// Clear name
	m.inputs[FieldName].SetValue("")

	// Press Enter
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should still be visible with error
	if !m.IsVisible() {
		t.Error("overlay should remain visible when validation fails")
	}
	if m.err == "" {
		t.Error("error message should be set")
	}
}

func TestFlaggedToggle(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Test", Flagged: false}
	m = m.Show(task).SetSize(80, 24)

	// Navigate to flagged field
	for i := 0; i < FieldFlagged; i++ {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	if m.focusIndex != FieldFlagged {
		t.Errorf("focus = %d, want %d", m.focusIndex, FieldFlagged)
	}
	if m.flagged {
		t.Error("flagged should be false initially")
	}

	// Press Enter to toggle
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !m.flagged {
		t.Error("flagged should be true after toggle")
	}
	if !m.IsVisible() {
		t.Error("overlay should still be visible (toggle, not submit)")
	}
}

func TestBuildModification_TagChanges(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{
		ID:   "task1",
		Name: "Test",
		Tags: []string{"urgent", "old"},
	}
	m = m.Show(task)

	// Change tags: remove "old", keep "urgent", add "new"
	m.inputs[FieldTags].SetValue("urgent, new")

	mod := m.buildModification()

	if len(mod.AddTags) != 1 || mod.AddTags[0] != "new" {
		t.Errorf("AddTags = %v, want [new]", mod.AddTags)
	}
	if len(mod.RemoveTags) != 1 || mod.RemoveTags[0] != "old" {
		t.Errorf("RemoveTags = %v, want [old]", mod.RemoveTags)
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

func TestView_Visible(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	task := &domain.Task{ID: "task1", Name: "Test"}
	m = m.Show(task).SetSize(80, 24)

	view := m.View()

	if view == "" {
		t.Error("view should not be empty when visible")
	}
	if !strings.Contains(view, "Edit Task") {
		t.Error("view should contain title")
	}
	if !strings.Contains(view, "Name") {
		t.Error("view should contain field labels")
	}
}
