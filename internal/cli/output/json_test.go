package output

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestJSONFormatter_FormatTasks(t *testing.T) {
	formatter := NewJSONFormatter()
	dueDate := time.Date(2026, 1, 28, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		tasks   []domain.Task
		options TaskFormatOptions
	}{
		{
			name:    "empty task list",
			tasks:   []domain.Task{},
			options: TaskFormatOptions{},
		},
		{
			name: "single task with minimal fields",
			tasks: []domain.Task{
				{ID: "task1", Name: "Buy groceries"},
			},
			options: TaskFormatOptions{},
		},
		{
			name: "task with all fields",
			tasks: []domain.Task{
				{
					ID:          "task1",
					Name:        "Review PR",
					Note:        "Check code style",
					ProjectID:   "proj1",
					ProjectName: "Work",
					Tags:        []string{"urgent", "code-review"},
					DueDate:     &dueDate,
					Flagged:     true,
					Completed:   false,
				},
			},
			options: TaskFormatOptions{ShowProject: true, ShowTags: true},
		},
		{
			name: "multiple tasks",
			tasks: []domain.Task{
				{ID: "task1", Name: "Task 1"},
				{ID: "task2", Name: "Task 2", Flagged: true},
			},
			options: TaskFormatOptions{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTasks(tt.tasks, tt.options)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatTasks() returned invalid JSON: %v", err)
			}

			// Verify tasks array exists
			tasksArray, ok := parsed["tasks"].([]interface{})
			if !ok {
				t.Fatal("FormatTasks() missing 'tasks' array")
			}

			// Verify count matches
			if len(tasksArray) != len(tt.tasks) {
				t.Errorf("FormatTasks() task count = %d, want %d", len(tasksArray), len(tt.tasks))
			}
		})
	}
}

func TestJSONFormatter_FormatTask(t *testing.T) {
	formatter := NewJSONFormatter()
	dueDate := time.Date(2026, 1, 28, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name string
		task domain.Task
	}{
		{
			name: "minimal task",
			task: domain.Task{ID: "task1", Name: "Test"},
		},
		{
			name: "complete task",
			task: domain.Task{
				ID:          "task1",
				Name:        "Complete task",
				Note:        "Notes here",
				ProjectID:   "proj1",
				ProjectName: "Project",
				Tags:        []string{"tag1", "tag2"},
				DueDate:     &dueDate,
				Flagged:     true,
				Completed:   true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTask(tt.task)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatTask() returned invalid JSON: %v", err)
			}

			// Verify task field exists
			taskData, ok := parsed["task"].(map[string]interface{})
			if !ok {
				t.Fatal("FormatTask() missing 'task' field")
			}

			// Verify key fields match
			if taskData["id"] != tt.task.ID {
				t.Errorf("FormatTask() ID = %v, want %v", taskData["id"], tt.task.ID)
			}
			if taskData["name"] != tt.task.Name {
				t.Errorf("FormatTask() Name = %v, want %v", taskData["name"], tt.task.Name)
			}
		})
	}
}

func TestJSONFormatter_FormatProjects(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name     string
		projects []domain.Project
		options  ProjectFormatOptions
	}{
		{
			name:     "empty project list",
			projects: []domain.Project{},
			options:  ProjectFormatOptions{},
		},
		{
			name: "single project",
			projects: []domain.Project{
				{ID: "proj1", Name: "Work", Status: "active"},
			},
			options: ProjectFormatOptions{},
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
					},
				},
			},
			options: ProjectFormatOptions{ShowTasks: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatProjects(tt.projects, tt.options)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatProjects() returned invalid JSON: %v", err)
			}

			// Verify projects array exists
			projectsArray, ok := parsed["projects"].([]interface{})
			if !ok {
				t.Fatal("FormatProjects() missing 'projects' array")
			}

			// Verify count matches
			if len(projectsArray) != len(tt.projects) {
				t.Errorf("FormatProjects() project count = %d, want %d", len(projectsArray), len(tt.projects))
			}
		})
	}
}

func TestJSONFormatter_FormatProject(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name    string
		project domain.Project
	}{
		{
			name:    "minimal project",
			project: domain.Project{ID: "proj1", Name: "Test", Status: "active"},
		},
		{
			name: "project with all fields",
			project: domain.Project{
				ID:     "proj1",
				Name:   "Complete Project",
				Status: "active",
				Note:   "Project notes",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatProject(tt.project)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatProject() returned invalid JSON: %v", err)
			}

			// Verify project field exists
			projectData, ok := parsed["project"].(map[string]interface{})
			if !ok {
				t.Fatal("FormatProject() missing 'project' field")
			}

			// Verify key fields match
			if projectData["id"] != tt.project.ID {
				t.Errorf("FormatProject() ID = %v, want %v", projectData["id"], tt.project.ID)
			}
			if projectData["name"] != tt.project.Name {
				t.Errorf("FormatProject() Name = %v, want %v", projectData["name"], tt.project.Name)
			}
		})
	}
}

func TestJSONFormatter_FormatTags(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name    string
		tags    []domain.Tag
		options TagFormatOptions
	}{
		{
			name:    "empty tag list",
			tags:    []domain.Tag{},
			options: TagFormatOptions{},
		},
		{
			name: "single tag",
			tags: []domain.Tag{
				{ID: "tag1", Name: "urgent"},
			},
			options: TagFormatOptions{},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTags(tt.tags, tt.options)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatTags() returned invalid JSON: %v", err)
			}

			// Verify tags array exists
			tagsArray, ok := parsed["tags"].([]interface{})
			if !ok {
				t.Fatal("FormatTags() missing 'tags' array")
			}

			// Verify count matches
			if len(tagsArray) != len(tt.tags) {
				t.Errorf("FormatTags() tag count = %d, want %d", len(tagsArray), len(tt.tags))
			}
		})
	}
}

func TestJSONFormatter_FormatTag(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name string
		tag  domain.Tag
	}{
		{
			name: "simple tag",
			tag:  domain.Tag{ID: "tag1", Name: "urgent"},
		},
		{
			name: "tag with parent",
			tag:  domain.Tag{ID: "tag2", Name: "meetings", ParentID: "tag1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatTag(tt.tag)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatTag() returned invalid JSON: %v", err)
			}

			// Verify tag field exists
			tagData, ok := parsed["tag"].(map[string]interface{})
			if !ok {
				t.Fatal("FormatTag() missing 'tag' field")
			}

			// Verify key fields match
			if tagData["id"] != tt.tag.ID {
				t.Errorf("FormatTag() ID = %v, want %v", tagData["id"], tt.tag.ID)
			}
			if tagData["name"] != tt.tag.Name {
				t.Errorf("FormatTag() Name = %v, want %v", tagData["name"], tt.tag.Name)
			}
		})
	}
}

func TestJSONFormatter_FormatError(t *testing.T) {
	formatter := NewJSONFormatter()

	tests := []struct {
		name string
		err  error
	}{
		{
			name: "simple error",
			err:  errors.New("something went wrong"),
		},
		{
			name: "formatted error",
			err:  errors.New("OmniFocus is not running"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.FormatError(tt.err)

			// Verify it's valid JSON
			var parsed map[string]interface{}
			err := json.Unmarshal([]byte(result), &parsed)
			if err != nil {
				t.Fatalf("FormatError() returned invalid JSON: %v", err)
			}

			// Verify error field exists
			errorMsg, ok := parsed["error"].(string)
			if !ok {
				t.Fatal("FormatError() missing 'error' field")
			}

			// Verify error message is present
			if errorMsg == "" {
				t.Error("FormatError() error message is empty")
			}
		})
	}
}
