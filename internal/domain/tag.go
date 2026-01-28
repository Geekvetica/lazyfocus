package domain

// Tag represents a tag in OmniFocus
type Tag struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentId,omitempty"`
	Children []Tag  `json:"children,omitempty"`
}
