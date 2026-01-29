package domain

import (
	"errors"
	"strings"
	"time"
)

// TaskInput represents the data needed to create a new task
type TaskInput struct {
	Name        string     // Required: task name
	Note        string     // Optional: task note
	ProjectID   string     // Optional: resolved project ID
	ProjectName string     // Optional: original @project name from input
	TagNames    []string   // Optional: tag names to apply
	DueDate     *time.Time // Optional: due date
	DeferDate   *time.Time // Optional: defer/start date
	Flagged     *bool      // Optional: flagged status
}

// Validate returns error if required fields are missing
func (t TaskInput) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("task name is required")
	}
	return nil
}

// HasProject returns true if a project is specified
func (t TaskInput) HasProject() bool {
	return t.ProjectID != "" || t.ProjectName != ""
}

// HasTags returns true if tags are specified
func (t TaskInput) HasTags() bool {
	return len(t.TagNames) > 0
}
