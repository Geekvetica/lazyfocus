package taskparse

import (
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestParse(t *testing.T) {
	// Use a fixed reference time for date parsing consistency
	ref := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.Local) // Monday

	tests := []struct {
		name      string
		input     string
		want      domain.TaskInput
		wantError bool
	}{
		{
			name:  "simple task",
			input: "Buy milk",
			want: domain.TaskInput{
				Name:     "Buy milk",
				TagNames: []string{},
			},
		},
		{
			name:  "task with single tag",
			input: "Buy milk #groceries",
			want: domain.TaskInput{
				Name:     "Buy milk",
				TagNames: []string{"groceries"},
			},
		},
		{
			name:  "task with multiple tags",
			input: "Task #tag1 #tag2",
			want: domain.TaskInput{
				Name:     "Task",
				TagNames: []string{"tag1", "tag2"},
			},
		},
		{
			name:  "task with tag at start",
			input: "#work Review document",
			want: domain.TaskInput{
				Name:     "Review document",
				TagNames: []string{"work"},
			},
		},
		{
			name:  "task with project",
			input: "Task @Work",
			want: domain.TaskInput{
				Name:        "Task",
				ProjectName: "Work",
				TagNames:    []string{},
			},
		},
		{
			name:  "task with quoted project",
			input: `Task @"My Project"`,
			want: domain.TaskInput{
				Name:        "Task",
				ProjectName: "My Project",
				TagNames:    []string{},
			},
		},
		{
			name:  "task with due date",
			input: "Task due:tomorrow",
			want: domain.TaskInput{
				Name:     "Task",
				DueDate:  timePtr(ref.AddDate(0, 0, 1).Format("2006-01-02")),
				TagNames: []string{},
			},
		},
		{
			name:  "task with quoted due date",
			input: `Task due:"next week"`,
			want: domain.TaskInput{
				Name:     "Task",
				DueDate:  timePtr(ref.AddDate(0, 0, 7).Format("2006-01-02")),
				TagNames: []string{},
			},
		},
		{
			name:  "task with defer date",
			input: "Task defer:tomorrow",
			want: domain.TaskInput{
				Name:      "Task",
				DeferDate: timePtr(ref.AddDate(0, 0, 1).Format("2006-01-02")),
				TagNames:  []string{},
			},
		},
		{
			name:  "task flagged",
			input: "Urgent task !",
			want: domain.TaskInput{
				Name:     "Urgent task",
				Flagged:  boolPtr(true),
				TagNames: []string{},
			},
		},
		{
			name:  "task flagged at start",
			input: "! Very urgent task",
			want: domain.TaskInput{
				Name:     "Very urgent task",
				Flagged:  boolPtr(true),
				TagNames: []string{},
			},
		},
		{
			name:  "combined all features",
			input: "Review PR #work @Development due:\"next friday\" !",
			want: domain.TaskInput{
				Name:        "Review PR",
				TagNames:    []string{"work"},
				ProjectName: "Development",
				DueDate:     timePtr("2024-01-19"), // Next Friday from Monday Jan 15
				Flagged:     boolPtr(true),
			},
		},
		{
			name:  "complex with multiple tags and quoted values",
			input: `Task #tag1 #tag2 @"My Project" due:"next week" defer:tomorrow !`,
			want: domain.TaskInput{
				Name:        "Task",
				TagNames:    []string{"tag1", "tag2"},
				ProjectName: "My Project",
				DueDate:     timePtr(ref.AddDate(0, 0, 7).Format("2006-01-02")),
				DeferDate:   timePtr(ref.AddDate(0, 0, 1).Format("2006-01-02")),
				Flagged:     boolPtr(true),
			},
		},
		{
			name:  "tags in middle of text",
			input: "Buy #groceries milk and #household eggs",
			want: domain.TaskInput{
				Name:     "Buy milk and eggs",
				TagNames: []string{"groceries", "household"},
			},
		},
		{
			name:      "empty input",
			input:     "",
			wantError: true,
		},
		{
			name:      "only whitespace",
			input:     "   ",
			wantError: true,
		},
		{
			name:      "only modifiers no name",
			input:     "#tag @project",
			wantError: true,
		},
		{
			name:      "only flag no name",
			input:     "!",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseWithReference(tt.input, ref)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if got.Name != tt.want.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.want.Name)
			}

			if !equalStringSlices(got.TagNames, tt.want.TagNames) {
				t.Errorf("TagNames = %v, want %v", got.TagNames, tt.want.TagNames)
			}

			if got.ProjectName != tt.want.ProjectName {
				t.Errorf("ProjectName = %q, want %q", got.ProjectName, tt.want.ProjectName)
			}

			if !equalBoolPtrs(got.Flagged, tt.want.Flagged) {
				t.Errorf("Flagged = %v, want %v", got.Flagged, tt.want.Flagged)
			}

			// Compare dates by formatting to date strings (ignore time component)
			if tt.want.DueDate != nil {
				if got.DueDate == nil {
					t.Errorf("Expected DueDate to be set")
				} else if formatDate(*got.DueDate) != formatDate(*tt.want.DueDate) {
					t.Errorf("DueDate = %s, want %s", formatDate(*got.DueDate), formatDate(*tt.want.DueDate))
				}
			} else {
				if got.DueDate != nil {
					t.Errorf("Expected DueDate to be nil, got %v", got.DueDate)
				}
			}

			if tt.want.DeferDate != nil {
				if got.DeferDate == nil {
					t.Errorf("Expected DeferDate to be set")
				} else if formatDate(*got.DeferDate) != formatDate(*tt.want.DeferDate) {
					t.Errorf("DeferDate = %s, want %s", formatDate(*got.DeferDate), formatDate(*tt.want.DeferDate))
				}
			} else {
				if got.DeferDate != nil {
					t.Errorf("Expected DeferDate to be nil, got %v", got.DeferDate)
				}
			}
		})
	}
}

// Helper functions for test data

func timePtr(dateStr string) *time.Time {
	t, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		panic(err)
	}
	// Set to 5 PM to match dateparse behavior
	t = time.Date(t.Year(), t.Month(), t.Day(), 17, 0, 0, 0, time.Local)
	return &t
}

func boolPtr(b bool) *bool {
	return &b
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalBoolPtrs(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
