package filter

import "testing"

func TestState_IsActive(t *testing.T) {
	tests := []struct {
		name   string
		state  State
		expect bool
	}{
		{"empty state", State{}, false},
		{"with search text", State{SearchText: "test"}, true},
		{"with project", State{ProjectID: "proj1"}, true},
		{"with tag", State{TagID: "tag1"}, true},
		{"with due filter", State{DueFilter: DueToday}, true},
		{"with flagged only", State{FlaggedOnly: true}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.state.IsActive(); got != tt.expect {
				t.Errorf("IsActive() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestState_Clear(t *testing.T) {
	state := State{
		SearchText:  "test",
		ProjectID:   "proj1",
		TagID:       "tag1",
		DueFilter:   DueToday,
		FlaggedOnly: true,
	}

	cleared := state.Clear()

	if cleared.IsActive() {
		t.Error("cleared state should not be active")
	}
}

func TestState_BuilderMethods(t *testing.T) {
	state := State{}.
		WithSearchText("search").
		WithProject("proj1").
		WithTag("tag1").
		WithDueFilter(DueWeek).
		WithFlaggedOnly(true)

	if state.SearchText != "search" {
		t.Errorf("SearchText = %q, want %q", state.SearchText, "search")
	}
	if state.ProjectID != "proj1" {
		t.Errorf("ProjectID = %q, want %q", state.ProjectID, "proj1")
	}
	if state.TagID != "tag1" {
		t.Errorf("TagID = %q, want %q", state.TagID, "tag1")
	}
	if state.DueFilter != DueWeek {
		t.Errorf("DueFilter = %v, want %v", state.DueFilter, DueWeek)
	}
	if !state.FlaggedOnly {
		t.Error("FlaggedOnly = false, want true")
	}
}
