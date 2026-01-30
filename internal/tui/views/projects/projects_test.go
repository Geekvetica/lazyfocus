package projects

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// MockService for testing
type MockService struct {
	projects []domain.Project
	tasks    []domain.Task
}

func (m *MockService) GetProjects(status string) ([]domain.Project, error) {
	return m.projects, nil
}

func (m *MockService) GetTasksByProject(projectID string) ([]domain.Task, error) {
	return m.tasks, nil
}

// Implement other interface methods as stubs...
func (m *MockService) GetInboxTasks() ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetAllTasks(f service.TaskFilters) ([]domain.Task, error) {
	return nil, nil
}
func (m *MockService) GetTasksByTag(tagID string) ([]domain.Task, error)   { return nil, nil }
func (m *MockService) GetFlaggedTasks() ([]domain.Task, error)             { return nil, nil }
func (m *MockService) GetTaskByID(id string) (*domain.Task, error)         { return nil, nil }
func (m *MockService) CreateTask(input domain.TaskInput) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) ModifyTask(id string, mod domain.TaskModification) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) CompleteTask(id string) (*domain.OperationResult, error) {
	return nil, nil
}
func (m *MockService) DeleteTask(id string) (*domain.OperationResult, error) { return nil, nil }
func (m *MockService) GetProjectByID(id string) (*domain.Project, error)      { return nil, nil }
func (m *MockService) GetProjectWithTasks(id string) (*domain.Project, error) { return nil, nil }
func (m *MockService) GetTags() ([]domain.Tag, error)                         { return nil, nil }
func (m *MockService) GetTagByID(id string) (*domain.Tag, error)              { return nil, nil }
func (m *MockService) GetTagCounts() (map[string]int, error)                  { return nil, nil }
func (m *MockService) GetPerspectiveTasks(name string) ([]domain.Task, error) { return nil, nil }
func (m *MockService) ResolveProjectName(name string) (string, error)         { return "", nil }

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.Mode() != ModeProjectList {
		t.Error("should start in project list mode")
	}
	if m.loaded {
		t.Error("should not be loaded initially")
	}
}

func TestProjectsLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{
			{ID: "p1", Name: "Project 1", TaskCount: 5},
			{ID: "p2", Name: "Project 2", TaskCount: 3},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	if !m.loaded {
		t.Error("should be loaded after ProjectsLoadedMsg")
	}
}

func TestEnterKey_DrillsDown(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
		tasks:    []domain.Task{{ID: "t1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	// Press Enter to drill down
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.Mode() != ModeProjectTasks {
		t.Error("should switch to task mode after Enter")
	}
	if m.currentProject == nil {
		t.Error("currentProject should be set")
	}
	if cmd == nil {
		t.Error("should return command to load tasks")
	}
}

func TestBackKey_ReturnsToList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down

	// Press h to go back
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if m.Mode() != ModeProjectList {
		t.Error("should return to project list mode")
	}
}

func TestEscapeKey_ReturnsToList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down

	// Press Escape to go back
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if m.Mode() != ModeProjectList {
		t.Error("should return to project list mode")
	}
}

func TestSelectedTask_InTaskMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
		tasks:    []domain.Task{{ID: "t1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	task := m.SelectedTask()
	if task == nil {
		t.Fatal("should return selected task in task mode")
	}
	if task.ID != "t1" {
		t.Errorf("task ID = %q, want %q", task.ID, "t1")
	}
}

func TestSelectedTask_InProjectMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	task := m.SelectedTask()
	if task != nil {
		t.Error("should return nil in project list mode")
	}
}
