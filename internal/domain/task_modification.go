package domain

import "time"

// TaskModification represents changes to apply to an existing task
// Nil pointer fields are not modified; non-nil fields are set to the value
type TaskModification struct {
	Name       *string    // New name (nil = don't change)
	Note       *string    // New note (nil = don't change)
	ProjectID  *string    // New project ID (nil = don't change, empty string = remove from project)
	AddTags    []string   // Tags to add
	RemoveTags []string   // Tags to remove
	DueDate    *time.Time // New due date (nil = don't change)
	DeferDate  *time.Time // New defer date (nil = don't change)
	Flagged    *bool      // New flagged status (nil = don't change)
	ClearDue   bool       // If true, clear the due date
	ClearDefer bool       // If true, clear the defer date
}

// IsEmpty returns true if no modifications are specified
func (m TaskModification) IsEmpty() bool {
	return m.Name == nil &&
		m.Note == nil &&
		m.ProjectID == nil &&
		len(m.AddTags) == 0 &&
		len(m.RemoveTags) == 0 &&
		m.DueDate == nil &&
		m.DeferDate == nil &&
		m.Flagged == nil &&
		!m.ClearDue &&
		!m.ClearDefer
}

// HasTagChanges returns true if tags are being added or removed
func (m TaskModification) HasTagChanges() bool {
	return len(m.AddTags) > 0 || len(m.RemoveTags) > 0
}
