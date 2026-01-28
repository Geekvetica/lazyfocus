package output

import (
	"strings"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestHumanFormatter_FormatTasks(t *testing.T) {
	formatter := NewHumanFormatter()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	tests := []struct {
		name    string
		tasks   []domain.Task
		options TaskFormatOptions
		want    []string // strings that should appear in output
	}{
		{
			name:    "empty task list",
			tasks:   []domain.Task{},
			options: TaskFormatOptions{},
			want:    []string{"0 task"},
		},
		{
			name: "incomplete task",
			tasks: []domain.Task{
				{ID: "task1", Name: "Buy groceries"},
			},
			options: TaskFormatOptions{},
			want:    []string{"☐", "Buy groceries", "1 task"},
		},
		{
			name: "completed task",
			tasks: []domain.Task{
				{ID: "task1", Name: "Done task", Completed: true},
			},
			options: TaskFormatOptions{},
			want:    []string{"☑", "Done task"},
		},
		{
			name: "flagged task",
			tasks: []domain.Task{
				{ID: "task1", Name: "Important", Flagged: true},
			},
			options: TaskFormatOptions{},
			want:    []string{"🚩", "Important"},
		},
		{
			name: "task with due date today",
			tasks: []domain.Task{
				{ID: "task1", Name: "Due today", DueDate: &today},
			},
			options: TaskFormatOptions{},
			want:    []string{"📅", "Today", "Due today"},
		},
		{
			name: "task with tags",
			tasks: []domain.Task{
				{ID: "task1", Name: "Tagged", Tags: []string{"urgent", "work"}},
			},
			options: TaskFormatOptions{ShowTags: true},
			want:    []string{"Tagged", "#urgent", "#work"},
		},
		{
			name: "task with project",
			tasks: []domain.Task{
				{ID: "task1", Name: "Project task", ProjectName: "Work"},
			},
			options: TaskFormatOptions{ShowProject: true},
			want:    []string{"Project task", "Work"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTasks(tt.tasks, tt.options)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatTasks() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestHumanFormatter_FormatTask(t *testing.T) {
	formatter := NewHumanFormatter()
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	tests := []struct {
		name string
		task domain.Task
		want []string // strings that should appear in output
	}{
		{
			name: "minimal task",
			task: domain.Task{ID: "task1", Name: "Simple"},
			want: []string{"☐", "Simple"},
		},
		{
			name: "task with note",
			task: domain.Task{ID: "task1", Name: "With note", Note: "Important note"},
			want: []string{"With note", "Important note"},
		},
		{
			name: "flagged task with due date",
			task: domain.Task{
				ID:      "task1",
				Name:    "Urgent",
				Flagged: true,
				DueDate: &today,
			},
			want: []string{"☐", "Urgent", "🚩", "📅", "Today"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTask(tt.task)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatTask() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestHumanFormatter_FormatProjects(t *testing.T) {
	formatter := NewHumanFormatter()

	tests := []struct {
		name     string
		projects []domain.Project
		options  ProjectFormatOptions
		want     []string
	}{
		{
			name:     "empty project list",
			projects: []domain.Project{},
			options:  ProjectFormatOptions{},
			want:     []string{"0 project"},
		},
		{
			name: "single project",
			projects: []domain.Project{
				{ID: "proj1", Name: "Work", Status: "active"},
			},
			options: ProjectFormatOptions{},
			want:    []string{"Work", "active", "1 project"},
		},
		{
			name: "project with note",
			projects: []domain.Project{
				{ID: "proj1", Name: "Work", Status: "active", Note: "Important project"},
			},
			options: ProjectFormatOptions{ShowNotes: true},
			want:    []string{"Work", "Important project"},
		},
		{
			name: "project with tasks",
			projects: []domain.Project{
				{
					ID:     "proj1",
					Name:   "Work",
					Status: "active",
					Tasks: []domain.Task{
						{ID: "task1", Name: "Task 1"},
						{ID: "task2", Name: "Task 2"},
					},
				},
			},
			options: ProjectFormatOptions{ShowTasks: true},
			want:    []string{"Work", "Task 1", "Task 2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatProjects(tt.projects, tt.options)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatProjects() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestHumanFormatter_FormatProject(t *testing.T) {
	formatter := NewHumanFormatter()

	tests := []struct {
		name    string
		project domain.Project
		want    []string
	}{
		{
			name:    "active project",
			project: domain.Project{ID: "proj1", Name: "Work", Status: "active"},
			want:    []string{"Work", "active"},
		},
		{
			name:    "completed project",
			project: domain.Project{ID: "proj1", Name: "Done", Status: "completed"},
			want:    []string{"Done", "completed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatProject(tt.project)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatProject() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestHumanFormatter_FormatTags(t *testing.T) {
	formatter := NewHumanFormatter()

	tests := []struct {
		name    string
		tags    []domain.Tag
		options TagFormatOptions
		want    []string
	}{
		{
			name:    "empty tag list",
			tags:    []domain.Tag{},
			options: TagFormatOptions{},
			want:    []string{"0 tag"},
		},
		{
			name: "single tag",
			tags: []domain.Tag{
				{ID: "tag1", Name: "urgent"},
			},
			options: TagFormatOptions{},
			want:    []string{"urgent", "1 tag"},
		},
		{
			name: "nested tags",
			tags: []domain.Tag{
				{
					ID:   "tag1",
					Name: "work",
					Children: []domain.Tag{
						{ID: "tag2", Name: "meetings", ParentID: "tag1"},
					},
				},
			},
			options: TagFormatOptions{},
			want:    []string{"work", "meetings"},
		},
		{
			name: "flat tags",
			tags: []domain.Tag{
				{
					ID:   "tag1",
					Name: "work",
					Children: []domain.Tag{
						{ID: "tag2", Name: "meetings", ParentID: "tag1"},
					},
				},
			},
			options: TagFormatOptions{Flat: true},
			want:    []string{"work", "meetings"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTags(tt.tags, tt.options)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatTags() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestHumanFormatter_FormatTag(t *testing.T) {
	formatter := NewHumanFormatter()

	tests := []struct {
		name string
		tag  domain.Tag
		want []string
	}{
		{
			name: "simple tag",
			tag:  domain.Tag{ID: "tag1", Name: "urgent"},
			want: []string{"urgent"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTag(tt.tag)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatTag() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestHumanFormatter_FormatError(t *testing.T) {
	formatter := NewHumanFormatter()

	tests := []struct {
		name string
		err  error
		want []string
	}{
		{
			name: "simple error",
			err:  &testError{msg: "something went wrong"},
			want: []string{"Error:", "something went wrong"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatError(tt.err)

			for _, want := range tt.want {
				if !strings.Contains(result, want) {
					t.Errorf("FormatError() output missing %q\nGot: %s", want, result)
				}
			}
		})
	}
}

func TestFormatDate(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	yesterday := today.AddDate(0, 0, -1)
	thisYear := time.Date(now.Year(), 3, 15, 12, 0, 0, 0, now.Location())
	lastYear := time.Date(now.Year()-1, 12, 25, 12, 0, 0, 0, now.Location())

	tests := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "today",
			date: today,
			want: "Today",
		},
		{
			name: "tomorrow",
			date: tomorrow,
			want: "Tomorrow",
		},
		{
			name: "yesterday",
			date: yesterday,
			want: "Yesterday",
		},
		{
			name: "this year",
			date: thisYear,
			want: "Mar 15",
		},
		{
			name: "last year",
			date: lastYear,
			want: "Dec 25, " + string(rune('0'+lastYear.Year()/1000)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDate(tt.date)
			if !strings.Contains(got, tt.want) {
				t.Errorf("formatDate() = %v, want to contain %v", got, tt.want)
			}
		})
	}
}

// testError is a simple error implementation for testing
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
