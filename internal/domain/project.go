package domain

// Project represents a project in OmniFocus
type Project struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"` // "active", "on-hold", "completed", "dropped"
	Note   string `json:"note,omitempty"`
	Tasks  []Task `json:"tasks,omitempty"` // optional, for detailed view
}
