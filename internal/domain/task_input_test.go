package domain

import (
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/testutil"
)

func TestTaskInput_Validate(t *testing.T) {
	tests := []struct {
		name    string
		input   TaskInput
		wantErr bool
	}{
		{
			name: "valid task with name only",
			input: TaskInput{
				Name: "Buy groceries",
			},
			wantErr: false,
		},
		{
			name: "valid task with all fields",
			input: TaskInput{
				Name:        "Buy groceries",
				Note:        "Get milk and bread",
				ProjectID:   "proj-123",
				ProjectName: "Errands",
				TagNames:    []string{"shopping", "errands"},
				DueDate:     timePtr(time.Now()),
				DeferDate:   timePtr(time.Now().Add(-24 * time.Hour)),
				Flagged:     testutil.BoolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "empty name returns error",
			input: TaskInput{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "whitespace-only name returns error",
			input: TaskInput{
				Name: "   ",
			},
			wantErr: true,
		},
		{
			name: "name with whitespace is valid",
			input: TaskInput{
				Name: "  Buy groceries  ",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.input.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskInput.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskInput_HasProject(t *testing.T) {
	tests := []struct {
		name  string
		input TaskInput
		want  bool
	}{
		{
			name: "has project ID",
			input: TaskInput{
				Name:      "Task",
				ProjectID: "proj-123",
			},
			want: true,
		},
		{
			name: "has project name",
			input: TaskInput{
				Name:        "Task",
				ProjectName: "Errands",
			},
			want: true,
		},
		{
			name: "has both project ID and name",
			input: TaskInput{
				Name:        "Task",
				ProjectID:   "proj-123",
				ProjectName: "Errands",
			},
			want: true,
		},
		{
			name: "no project",
			input: TaskInput{
				Name: "Task",
			},
			want: false,
		},
		{
			name: "empty project ID",
			input: TaskInput{
				Name:      "Task",
				ProjectID: "",
			},
			want: false,
		},
		{
			name: "empty project name",
			input: TaskInput{
				Name:        "Task",
				ProjectName: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.HasProject(); got != tt.want {
				t.Errorf("TaskInput.HasProject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskInput_HasTags(t *testing.T) {
	tests := []struct {
		name  string
		input TaskInput
		want  bool
	}{
		{
			name: "has single tag",
			input: TaskInput{
				Name:     "Task",
				TagNames: []string{"errands"},
			},
			want: true,
		},
		{
			name: "has multiple tags",
			input: TaskInput{
				Name:     "Task",
				TagNames: []string{"errands", "shopping"},
			},
			want: true,
		},
		{
			name: "no tags",
			input: TaskInput{
				Name: "Task",
			},
			want: false,
		},
		{
			name: "empty tag slice",
			input: TaskInput{
				Name:     "Task",
				TagNames: []string{},
			},
			want: false,
		},
		{
			name: "nil tag slice",
			input: TaskInput{
				Name:     "Task",
				TagNames: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.HasTags(); got != tt.want {
				t.Errorf("TaskInput.HasTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions for test data
func timePtr(t time.Time) *time.Time {
	return &t
}
