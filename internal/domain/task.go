package domain

import "time"

// Task represents a task in OmniFocus
type Task struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Note          string     `json:"note,omitempty"`
	ProjectID     string     `json:"projectId,omitempty"`
	ProjectName   string     `json:"projectName,omitempty"`
	Tags          []string   `json:"tags,omitempty"`
	DueDate       *time.Time `json:"dueDate,omitempty"`
	DeferDate     *time.Time `json:"deferDate,omitempty"`
	Flagged       bool       `json:"flagged"`
	Completed     bool       `json:"completed"`
	CompletedDate *time.Time `json:"completedDate,omitempty"`
}
