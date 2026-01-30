package filter

// DueFilter defines due date filtering options
type DueFilter int

// DueFilter constants for filtering tasks by due date.
const (
	DueNone DueFilter = iota
	DueToday
	DueTomorrow
	DueWeek
	DueOverdue
)

// State represents the current filter state
type State struct {
	SearchText  string
	ProjectID   string
	TagID       string
	DueFilter   DueFilter
	FlaggedOnly bool
}

// IsActive returns true if any filter is applied
func (s State) IsActive() bool {
	return s.SearchText != "" ||
		s.ProjectID != "" ||
		s.TagID != "" ||
		s.DueFilter != DueNone ||
		s.FlaggedOnly
}

// Clear returns a State with all filters cleared
func (s State) Clear() State {
	return State{}
}

// WithSearchText returns a State with the search text set
func (s State) WithSearchText(text string) State {
	s.SearchText = text
	return s
}

// WithProject returns a State with the project filter set
func (s State) WithProject(projectID string) State {
	s.ProjectID = projectID
	return s
}

// WithTag returns a State with the tag filter set
func (s State) WithTag(tagID string) State {
	s.TagID = tagID
	return s
}

// WithDueFilter returns a State with the due filter set
func (s State) WithDueFilter(filter DueFilter) State {
	s.DueFilter = filter
	return s
}

// WithFlaggedOnly returns a State with the flagged filter set
func (s State) WithFlaggedOnly(flagged bool) State {
	s.FlaggedOnly = flagged
	return s
}
