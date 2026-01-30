package forecast

import (
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// MockService for testing
type MockService struct {
	tasks []domain.Task
}

func (m *MockService) GetAllTasks(f service.TaskFilters) ([]domain.Task, error) {
	return m.tasks, nil
}

// Stub other methods
func (m *MockService) GetInboxTasks() ([]domain.Task, error)                     { return nil, nil }
func (m *MockService) GetTasksByProject(projectID string) ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetTasksByTag(tagID string) ([]domain.Task, error)         { return nil, nil }
func (m *MockService) GetFlaggedTasks() ([]domain.Task, error)                   { return nil, nil }
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

func TestGroupTasks(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	nextWeek := today.AddDate(0, 0, 3)
	farFuture := today.AddDate(0, 1, 0)

	tasks := []domain.Task{
		{ID: "1", Name: "Overdue task", DueDate: &yesterday},
		{ID: "2", Name: "Today task", DueDate: &today},
		{ID: "3", Name: "Tomorrow task", DueDate: &tomorrow},
		{ID: "4", Name: "This week task", DueDate: &nextWeek},
		{ID: "5", Name: "Later task", DueDate: &farFuture},
		{ID: "6", Name: "No due task"},
	}

	items := m.groupTasks(tasks)

	// Should have 6 headers + 6 tasks = 12 items (or less if some groups are empty)
	headerCount := 0
	taskCount := 0
	for _, item := range items {
		if item.IsHeader {
			headerCount++
		} else {
			taskCount++
		}
	}

	if taskCount != 6 {
		t.Errorf("expected 6 tasks, got %d", taskCount)
	}
	if headerCount < 1 {
		t.Error("expected at least 1 header")
	}
}

func TestSelectedTask_OnHeader(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	today := time.Now()
	svc := &MockService{
		tasks: []domain.Task{{ID: "1", Name: "Task 1", DueDate: &today}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})
	m.cursor = 0 // Header position

	task := m.SelectedTask()
	if task != nil {
		t.Error("should return nil when cursor is on header")
	}
}

func TestSelectedTask_OnTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	today := time.Now()
	svc := &MockService{
		tasks: []domain.Task{{ID: "1", Name: "Task 1", DueDate: &today}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})
	m.cursor = 1 // Task position (after header)

	task := m.SelectedTask()
	if task == nil {
		t.Fatal("should return task")
	}
	if task.ID != "1" {
		t.Errorf("task ID = %q, want %q", task.ID, "1")
	}
}
