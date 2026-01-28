// Package domain defines the core data models for LazyFocus including
// Task, Project, and Tag types that mirror OmniFocus entities.
package domain

// Project represents a project in OmniFocus
type Project struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "active", "on-hold", "completed", "dropped"
	Note   string `json:"note,omitempty"`
	Tasks  []Task `json:"tasks,omitempty"` // optional, for detailed view
}
