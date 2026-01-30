package review

import (
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// MockService for testing
type MockService struct {
	tasks []domain.Task
}

func (m *MockService) GetFlaggedTasks() ([]domain.Task, error) {
	return m.tasks, nil
}

// Stub other methods
func (m *MockService) GetInboxTasks() ([]domain.Task, error)                     { return nil, nil }
func (m *MockService) GetAllTasks(f service.TaskFilters) ([]domain.Task, error)  { return nil, nil }
func (m *MockService) GetTasksByProject(projectID string) ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetTasksByTag(tagID string) ([]domain.Task, error)         { return nil, nil }
func (m *MockService) GetTaskByID(id string) (*domain.Task, error)               { return nil, nil }
func (m *MockService) CreateTask(input domain.TaskInput) (*domain.Task, error)   { return nil, nil }
func (m *MockService) ModifyTask(id string, mod domain.TaskModification) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) CompleteTask(id string) (*domain.OperationResult, error) { return nil, nil }
func (m *MockService) DeleteTask(id string) (*domain.OperationResult, error)   { return nil, nil }
func (m *MockService) GetProjects(status string) ([]domain.Project, error)     { return nil, nil }
func (m *MockService) GetProjectByID(id string) (*domain.Project, error)       { return nil, nil }
func (m *MockService) GetProjectWithTasks(id string) (*domain.Project, error)  { return nil, nil }
func (m *MockService) GetTags() ([]domain.Tag, error)                          { return nil, nil }
func (m *MockService) GetTagByID(id string) (*domain.Tag, error)               { return nil, nil }
func (m *MockService) GetTagCounts() (map[string]int, error)                   { return nil, nil }
func (m *MockService) GetPerspectiveTasks(name string) ([]domain.Task, error)  { return nil, nil }
func (m *MockService) ResolveProjectName(name string) (string, error)          { return "", nil }

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.loaded {
		t.Error("should not be loaded initially")
	}
}

func TestTasksLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tasks: []domain.Task{
			{ID: "1", Name: "Flagged 1", Flagged: true},
			{ID: "2", Name: "Flagged 2", Flagged: true},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	if !m.loaded {
		t.Error("should be loaded")
	}
	if m.TaskCount() != 2 {
		t.Errorf("task count = %d, want 2", m.TaskCount())
	}
}

func TestSelectedTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tasks: []domain.Task{{ID: "1", Name: "Task 1", Flagged: true}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	task := m.SelectedTask()
	if task == nil {
		t.Fatal("should return selected task")
	}
	if task.ID != "1" {
		t.Errorf("task ID = %q, want %q", task.ID, "1")
	}
}
