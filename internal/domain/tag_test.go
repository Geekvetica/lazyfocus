package domain

import (
	"encoding/json"
	"testing"
)

func TestTag_JSONSerialization(t *testing.T) {
	child1 := Tag{
		ID:   "tag-child-1",
		Name: "Child 1",
	}
	child2 := Tag{
		ID:   "tag-child-2",
		Name: "Child 2",
	}

	tag := Tag{
		ID:       "tag-parent",
		Name:     "Parent Tag",
		ParentID: "tag-grandparent",
		Children: []Tag{child1, child2},
	}

	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal tag: %v", err)
	}

	var unmarshaled Tag
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal tag: %v", err)
	}

	if unmarshaled.ID != tag.ID {
		t.Errorf("ID mismatch: got %s, want %s", unmarshaled.ID, tag.ID)
	}
	if unmarshaled.Name != tag.Name {
		t.Errorf("Name mismatch: got %s, want %s", unmarshaled.Name, tag.Name)
	}
	if unmarshaled.ParentID != tag.ParentID {
		t.Errorf("ParentID mismatch: got %s, want %s", unmarshaled.ParentID, tag.ParentID)
	}
	if len(unmarshaled.Children) != len(tag.Children) {
		t.Errorf("Children length mismatch: got %d, want %d", len(unmarshaled.Children), len(tag.Children))
	}
	if len(unmarshaled.Children) >= 2 {
		if unmarshaled.Children[0].ID != child1.ID {
			t.Errorf("First child ID mismatch: got %s, want %s", unmarshaled.Children[0].ID, child1.ID)
		}
		if unmarshaled.Children[1].ID != child2.ID {
			t.Errorf("Second child ID mismatch: got %s, want %s", unmarshaled.Children[1].ID, child2.ID)
		}
	}
}

func TestTag_OmitEmptyFields(t *testing.T) {
	tag := Tag{
		ID:   "tag-123",
		Name: "Simple Tag",
	}

	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal tag: %v", err)
	}

	jsonString := string(jsonData)

	// Check that empty fields are omitted
	if contains(jsonString, "parentId") {
		t.Error("Expected 'parentId' field to be omitted when empty")
	}
	if contains(jsonString, "children") {
		t.Error("Expected 'children' field to be omitted when empty")
	}

	// Check that required fields are present
	if !contains(jsonString, "id") {
		t.Error("Expected 'id' field to be present")
	}
	if !contains(jsonString, "name") {
		t.Error("Expected 'name' field to be present")
	}
}

func TestTag_WithNestedChildren(t *testing.T) {
	jsonInput := `{
		"id": "tag-parent",
		"name": "Parent",
		"children": [
			{
				"id": "tag-child-1",
				"name": "Child 1",
				"parentId": "tag-parent"
			},
			{
				"id": "tag-child-2",
				"name": "Child 2",
				"parentId": "tag-parent",
				"children": [
					{
						"id": "tag-grandchild",
						"name": "Grandchild",
						"parentId": "tag-child-2"
					}
				]
			}
		]
	}`

	var tag Tag
	err := json.Unmarshal([]byte(jsonInput), &tag)
	if err != nil {
		t.Fatalf("Failed to unmarshal tag: %v", err)
	}

	if tag.ID != "tag-parent" {
		t.Errorf("Tag ID mismatch: got %s, want tag-parent", tag.ID)
	}
	if len(tag.Children) != 2 {
		t.Fatalf("Expected 2 children, got %d", len(tag.Children))
	}

	// Check first child
	if tag.Children[0].ID != "tag-child-1" {
		t.Errorf("First child ID mismatch: got %s, want tag-child-1", tag.Children[0].ID)
	}
	if tag.Children[0].ParentID != "tag-parent" {
		t.Errorf("First child ParentID mismatch: got %s, want tag-parent", tag.Children[0].ParentID)
	}

	// Check second child and its nested child
	if tag.Children[1].ID != "tag-child-2" {
		t.Errorf("Second child ID mismatch: got %s, want tag-child-2", tag.Children[1].ID)
	}
	if len(tag.Children[1].Children) != 1 {
		t.Fatalf("Expected second child to have 1 grandchild, got %d", len(tag.Children[1].Children))
	}
	if tag.Children[1].Children[0].ID != "tag-grandchild" {
		t.Errorf("Grandchild ID mismatch: got %s, want tag-grandchild", tag.Children[1].Children[0].ID)
	}
	if tag.Children[1].Children[0].ParentID != "tag-child-2" {
		t.Errorf("Grandchild ParentID mismatch: got %s, want tag-child-2", tag.Children[1].Children[0].ParentID)
	}
}

func TestTag_EmptyChildrenArray(t *testing.T) {
	jsonInput := `{
		"id": "tag-123",
		"name": "Tag with empty children",
		"children": []
	}`

	var tag Tag
	err := json.Unmarshal([]byte(jsonInput), &tag)
	if err != nil {
		t.Fatalf("Failed to unmarshal tag: %v", err)
	}

	if tag.ID != "tag-123" {
		t.Errorf("Tag ID mismatch: got %s, want tag-123", tag.ID)
	}
	if len(tag.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(tag.Children))
	}
}

func TestTag_RootTagWithoutParent(t *testing.T) {
	tag := Tag{
		ID:   "tag-root",
		Name: "Root Tag",
	}

	jsonData, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("Failed to marshal tag: %v", err)
	}

	var unmarshaled Tag
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal tag: %v", err)
	}

	if unmarshaled.ParentID != "" {
		t.Errorf("Expected ParentID to be empty, got %s", unmarshaled.ParentID)
	}
	if len(unmarshaled.Children) != 0 {
		t.Errorf("Expected 0 children, got %d", len(unmarshaled.Children))
	}
}
