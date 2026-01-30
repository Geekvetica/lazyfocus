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

func TestBuildModification_TagCaseSensitivity(t *testing.T) {
	styles := tui.DefaultStyles()

	tests := []struct {
		name           string
		existingTags   []string
		newTagsInput   string
		wantAddTags    []string
		wantRemoveTags []string
	}{
		{
			name:           "preserve uppercase in new tags",
			existingTags:   []string{"urgent"},
			newTagsInput:   "urgent, Work",
			wantAddTags:    []string{"Work"},
			wantRemoveTags: nil,
		},
		{
			name:           "preserve mixed case in new tags",
			existingTags:   []string{"low"},
			newTagsInput:   "HighPriority, UrgentTask",
			wantAddTags:    []string{"HighPriority", "UrgentTask"},
			wantRemoveTags: []string{"low"},
		},
		{
			name:           "preserve original case when removing",
			existingTags:   []string{"Work", "Personal"},
			newTagsInput:   "Work",
			wantAddTags:    nil,
			wantRemoveTags: []string{"Personal"},
		},
		{
			name:           "case-insensitive comparison for duplicates",
			existingTags:   []string{"Work"},
			newTagsInput:   "work",
			wantAddTags:    nil,
			wantRemoveTags: nil,
		},
		{
			name:           "case-insensitive comparison with mixed case input",
			existingTags:   []string{"HighPriority"},
			newTagsInput:   "highpriority",
			wantAddTags:    nil,
			wantRemoveTags: nil,
		},
		{
			name:           "preserve case when adding multiple tags",
			existingTags:   []string{},
			newTagsInput:   "Work, Personal, HighPriority",
			wantAddTags:    []string{"Work", "Personal", "HighPriority"},
			wantRemoveTags: nil,
		},
		{
			name:           "preserve original case when removing multiple tags",
			existingTags:   []string{"Work", "Personal", "HighPriority"},
			newTagsInput:   "",
			wantAddTags:    nil,
			wantRemoveTags: []string{"Work", "Personal", "HighPriority"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(styles)
			task := &domain.Task{
				ID:   "task1",
				Name: "Test",
				Tags: tt.existingTags,
			}
			m = m.Show(task)
			m.inputs[FieldTags].SetValue(tt.newTagsInput)

			mod := m.buildModification()

			// Check AddTags (order-independent)
			if len(mod.AddTags) != len(tt.wantAddTags) {
				t.Errorf("AddTags length = %d, want %d (got: %v, want: %v)",
					len(mod.AddTags), len(tt.wantAddTags), mod.AddTags, tt.wantAddTags)
			} else {
				// Create a set of expected tags
				wantSet := make(map[string]bool)
				for _, tag := range tt.wantAddTags {
					wantSet[tag] = true
				}
				// Verify each tag in AddTags is in the expected set
				for _, tag := range mod.AddTags {
					if !wantSet[tag] {
						t.Errorf("AddTags contains unexpected tag %q, want one of %v", tag, tt.wantAddTags)
					}
				}
			}

			// Check RemoveTags (order-independent)
			if len(mod.RemoveTags) != len(tt.wantRemoveTags) {
				t.Errorf("RemoveTags length = %d, want %d (got: %v, want: %v)",
					len(mod.RemoveTags), len(tt.wantRemoveTags), mod.RemoveTags, tt.wantRemoveTags)
			} else {
				// Create a set of expected tags
				wantSet := make(map[string]bool)
				for _, tag := range tt.wantRemoveTags {
					wantSet[tag] = true
				}
				// Verify each tag in RemoveTags is in the expected set
				for _, tag := range mod.RemoveTags {
					if !wantSet[tag] {
						t.Errorf("RemoveTags contains unexpected tag %q, want one of %v", tag, tt.wantRemoveTags)
					}
				}
			}
		})
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

// Date Validation Edge Cases
func TestDateValidation_ValidFormats(t *testing.T) {
	styles := tui.DefaultStyles()

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"relative date", "tomorrow", true},
		{"next occurrence", "next monday", true},
		{"in N units", "in 3 days", true},
		{"ISO format", "2024-01-15", true},
		{"month day", "Jan 15", true},
		{"empty (clears)", "", true},
		{"invalid string", "not a date", false},
		{"partial", "tomo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(styles)
			task := &domain.Task{ID: "task1", Name: "Test"}
			m = m.Show(task).SetSize(80, 24)

			// Set due date to test input
			m.inputs[FieldDueDate].SetValue(tt.input)

			// Validate
			err := m.validate()
			if tt.valid && err != "" {
				t.Errorf("validation failed for valid date %q: %s", tt.input, err)
			}
			if !tt.valid && err == "" {
				t.Errorf("validation passed for invalid date %q", tt.input)
			}
		})
	}
}

func TestDeferDateValidation(t *testing.T) {
	styles := tui.DefaultStyles()

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"relative date", "tomorrow", true},
		{"next occurrence", "next friday", true},
		{"in N units", "in 2 weeks", true},
		{"ISO format", "2024-02-01", true},
		{"empty (clears)", "", true},
		{"invalid string", "xyz", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(styles)
			task := &domain.Task{ID: "task1", Name: "Test"}
			m = m.Show(task).SetSize(80, 24)

			// Set defer date to test input
			m.inputs[FieldDeferDate].SetValue(tt.input)

			// Validate
			err := m.validate()
			if tt.valid && err != "" {
				t.Errorf("validation failed for valid date %q: %s", tt.input, err)
			}
			if !tt.valid && err == "" {
				t.Errorf("validation passed for invalid date %q", tt.input)
			}
		})
	}
}

// Project Field Handling
func TestProjectField_ChangesProject(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", ProjectName: "Old Project"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Navigate to project field
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab}) // Name -> Note
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab}) // Note -> Project

	// Set new project
	m.inputs[FieldProject].SetValue("New Project")

	// Build modification
	mod := m.buildModification()

	if mod.ProjectID == nil {
		t.Fatal("ProjectID should be set")
	}
	if *mod.ProjectID != "New Project" {
		t.Errorf("ProjectID = %q, want %q", *mod.ProjectID, "New Project")
	}
}

func TestProjectField_ClearsProject(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", ProjectName: "Work"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Clear project field
	m.inputs[FieldProject].SetValue("")

	// Build modification
	mod := m.buildModification()

	if mod.ProjectID == nil {
		t.Fatal("ProjectID should be set to empty string")
	}
	if *mod.ProjectID != "" {
		t.Errorf("ProjectID = %q, want empty string", *mod.ProjectID)
	}
}

// Tag Modification with Special Characters
func TestTags_WithSpaces(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", Tags: []string{"old tag"}}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Add tag with spaces
	m.inputs[FieldTags].SetValue("old tag, new tag with spaces")

	mod := m.buildModification()

	if len(mod.AddTags) != 1 {
		t.Fatalf("AddTags length = %d, want 1", len(mod.AddTags))
	}
	if mod.AddTags[0] != "new tag with spaces" {
		t.Errorf("AddTags[0] = %q, want %q", mod.AddTags[0], "new tag with spaces")
	}
}

func TestTags_WithSpecialChars(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", Tags: []string{}}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Add tags with special characters
	m.inputs[FieldTags].SetValue("@work, #priority, follow-up")

	mod := m.buildModification()

	if len(mod.AddTags) != 3 {
		t.Fatalf("AddTags length = %d, want 3", len(mod.AddTags))
	}

	wantTags := map[string]bool{
		"@work":     true,
		"#priority": true,
		"follow-up": true,
	}

	for _, tag := range mod.AddTags {
		if !wantTags[tag] {
			t.Errorf("unexpected tag %q", tag)
		}
	}
}

// Tab Navigation
func TestTabNavigation_ForwardCycle(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Verify initial field is Name
	if m.focusIndex != FieldName {
		t.Errorf("initial focus = %d, want %d", m.focusIndex, FieldName)
	}

	// Tab through all fields
	fields := []int{FieldName, FieldNote, FieldProject, FieldTags, FieldDueDate, FieldDeferDate, FieldFlagged}
	for i, expected := range fields[1:] {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
		if m.focusIndex != expected {
			t.Errorf("after tab %d: focus = %d, want %d", i+1, m.focusIndex, expected)
		}
	}

	// Tab should cycle back to Name
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	if m.focusIndex != FieldName {
		t.Errorf("after full cycle: focus = %d, want %d", m.focusIndex, FieldName)
	}
}

func TestTabNavigation_BackwardCycle(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Shift+Tab should go to last field
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if m.focusIndex != FieldFlagged {
		t.Errorf("after shift+tab: focus = %d, want %d", m.focusIndex, FieldFlagged)
	}

	// Continue backward
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	if m.focusIndex != FieldDeferDate {
		t.Errorf("after 2nd shift+tab: focus = %d, want %d", m.focusIndex, FieldDeferDate)
	}
}

// Flagged Field Toggle
func TestFlaggedToggle_SpaceKey(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", Flagged: false}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Navigate to flagged field
	for i := 0; i < FieldFlagged; i++ {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	if m.focusIndex != FieldFlagged {
		t.Fatalf("focus = %d, want %d", m.focusIndex, FieldFlagged)
	}

	// Press Enter to toggle (space is not handled, Enter toggles on flagged field)
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if !m.flagged {
		t.Error("flagged should be true after toggle")
	}

	// Build modification
	mod := m.buildModification()
	if mod.Flagged == nil {
		t.Fatal("Flagged should be set")
	}
	if !*mod.Flagged {
		t.Error("Flagged modification should be true")
	}
}

func TestFlaggedToggle_EnterKey(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", Flagged: true}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Navigate to flagged field
	for i := 0; i < FieldFlagged; i++ {
		m, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	}

	// Press Enter to toggle from true to false
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.flagged {
		t.Error("flagged should be false after toggle")
	}

	// Should still be visible (toggle, not submit)
	if !m.IsVisible() {
		t.Error("overlay should still be visible after flagged toggle")
	}
}

// Save and Cancel
func TestSave_EmitsModification(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Make a change
	m.inputs[FieldName].SetValue("Updated Test")

	// Press Enter to save (on a non-flagged field)
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if newM.IsVisible() {
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
	if saveMsg.Modification.Name == nil || *saveMsg.Modification.Name != "Updated Test" {
		t.Error("modification should include name change")
	}
}

func TestCancel_NoChangesEmitted(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Make some changes
	m.inputs[FieldName].SetValue("Changed Name")

	// Press Escape to cancel
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if newM.IsVisible() {
		t.Error("overlay should be hidden after cancel")
	}
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	if _, ok := msg.(CancelMsg); !ok {
		t.Errorf("expected CancelMsg, got %T", msg)
	}
}

func TestSave_NoChanges_StillEmitsMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Don't make any changes
	// Press Enter to save
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if newM.IsVisible() {
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
	// Modification might be empty, but SaveMsg should still be emitted
}

// View State
func TestShow_HidesOnVisible(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)

	if m.IsVisible() {
		t.Error("should not be visible initially")
	}

	m = m.Show(task)

	if !m.IsVisible() {
		t.Error("should be visible after Show()")
	}
}

func TestHide_ReturnsHiddenModel(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task)

	if !m.IsVisible() {
		t.Error("should be visible after Show()")
	}

	m = m.Hide()

	if m.IsVisible() {
		t.Error("should not be visible after Hide()")
	}
}

// Additional edge cases
func TestUpdate_NotVisible_ReturnsUnchanged(t *testing.T) {
	styles := tui.DefaultStyles()
	m := New(styles)

	// Update when not visible should do nothing
	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if cmd != nil {
		t.Error("should not return command when not visible")
	}
	if newM.IsVisible() {
		t.Error("should remain not visible")
	}
}

func TestDateClearing_DueDateAndDeferDate(t *testing.T) {
	styles := tui.DefaultStyles()
	dueDate := time.Now().Add(24 * time.Hour)
	deferDate := time.Now().Add(12 * time.Hour)
	task := &domain.Task{
		ID:        "task1",
		Name:      "Test",
		DueDate:   &dueDate,
		DeferDate: &deferDate,
	}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Clear both dates
	m.inputs[FieldDueDate].SetValue("")
	m.inputs[FieldDeferDate].SetValue("")

	mod := m.buildModification()

	if !mod.ClearDue {
		t.Error("ClearDue should be true")
	}
	if !mod.ClearDefer {
		t.Error("ClearDefer should be true")
	}
}

func TestBuildModification_NoteChange(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test", Note: "Original note"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Change note
	m.inputs[FieldNote].SetValue("Updated note")

	mod := m.buildModification()

	if mod.Note == nil {
		t.Fatal("Note should be set")
	}
	if *mod.Note != "Updated note" {
		t.Errorf("Note = %q, want %q", *mod.Note, "Updated note")
	}
}

func TestView_WithError(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Set an error
	m.err = "Test error message"

	view := m.View()

	if !strings.Contains(view, "Test error message") {
		t.Error("view should contain error message")
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

func TestWindowSizeMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	task := &domain.Task{ID: "task1", Name: "Test"}
	m := New(styles)
	m = m.Show(task).SetSize(80, 24)

	// Send WindowSizeMsg
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if newM.width != 120 {
		t.Errorf("width = %d, want 120", newM.width)
	}
	if newM.height != 40 {
		t.Errorf("height = %d, want 40", newM.height)
	}
}
