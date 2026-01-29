package domain

import (
	"testing"
	"time"
)

func TestTaskModification_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		mod  TaskModification
		want bool
	}{
		{
			name: "empty modification",
			mod:  TaskModification{},
			want: true,
		},
		{
			name: "has name",
			mod: TaskModification{
				Name: stringPtr("New name"),
			},
			want: false,
		},
		{
			name: "has note",
			mod: TaskModification{
				Note: stringPtr("New note"),
			},
			want: false,
		},
		{
			name: "has project ID",
			mod: TaskModification{
				ProjectID: stringPtr("proj-123"),
			},
			want: false,
		},
		{
			name: "has empty project ID to remove from project",
			mod: TaskModification{
				ProjectID: stringPtr(""),
			},
			want: false,
		},
		{
			name: "has add tags",
			mod: TaskModification{
				AddTags: []string{"tag1"},
			},
			want: false,
		},
		{
			name: "has remove tags",
			mod: TaskModification{
				RemoveTags: []string{"tag1"},
			},
			want: false,
		},
		{
			name: "has due date",
			mod: TaskModification{
				DueDate: timePtr(time.Now()),
			},
			want: false,
		},
		{
			name: "has defer date",
			mod: TaskModification{
				DeferDate: timePtr(time.Now()),
			},
			want: false,
		},
		{
			name: "has flagged",
			mod: TaskModification{
				Flagged: boolPtr(true),
			},
			want: false,
		},
		{
			name: "has clear due flag",
			mod: TaskModification{
				ClearDue: true,
			},
			want: false,
		},
		{
			name: "has clear defer flag",
			mod: TaskModification{
				ClearDefer: true,
			},
			want: false,
		},
		{
			name: "has both clear flags",
			mod: TaskModification{
				ClearDue:   true,
				ClearDefer: true,
			},
			want: false,
		},
		{
			name: "has multiple modifications",
			mod: TaskModification{
				Name:     stringPtr("New name"),
				Note:     stringPtr("New note"),
				Flagged:  boolPtr(true),
				AddTags:  []string{"tag1", "tag2"},
				ClearDue: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mod.IsEmpty(); got != tt.want {
				t.Errorf("TaskModification.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskModification_HasTagChanges(t *testing.T) {
	tests := []struct {
		name string
		mod  TaskModification
		want bool
	}{
		{
			name: "no tag changes",
			mod:  TaskModification{},
			want: false,
		},
		{
			name: "has add tags",
			mod: TaskModification{
				AddTags: []string{"tag1"},
			},
			want: true,
		},
		{
			name: "has remove tags",
			mod: TaskModification{
				RemoveTags: []string{"tag1"},
			},
			want: true,
		},
		{
			name: "has both add and remove tags",
			mod: TaskModification{
				AddTags:    []string{"tag1", "tag2"},
				RemoveTags: []string{"tag3"},
			},
			want: true,
		},
		{
			name: "empty add tags slice",
			mod: TaskModification{
				AddTags: []string{},
			},
			want: false,
		},
		{
			name: "empty remove tags slice",
			mod: TaskModification{
				RemoveTags: []string{},
			},
			want: false,
		},
		{
			name: "both empty slices",
			mod: TaskModification{
				AddTags:    []string{},
				RemoveTags: []string{},
			},
			want: false,
		},
		{
			name: "has other changes but no tag changes",
			mod: TaskModification{
				Name:     stringPtr("New name"),
				Flagged:  boolPtr(true),
				ClearDue: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mod.HasTagChanges(); got != tt.want {
				t.Errorf("TaskModification.HasTagChanges() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function for test data
func stringPtr(s string) *string {
	return &s
}
