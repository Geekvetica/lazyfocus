package tags

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// MockService for testing
type MockService struct {
	tags   []domain.Tag
	counts map[string]int
	tasks  []domain.Task
}

func (m *MockService) GetTags() ([]domain.Tag, error) {
	return m.tags, nil
}

func (m *MockService) GetTagCounts() (map[string]int, error) {
	return m.counts, nil
}

func (m *MockService) GetTasksByTag(tagID string) ([]domain.Task, error) {
	return m.tasks, nil
}

// Stub other interface methods
func (m *MockService) GetInboxTasks() ([]domain.Task, error)                     { return nil, nil }
func (m *MockService) GetAllTasks(f service.TaskFilters) ([]domain.Task, error)  { return nil, nil }
func (m *MockService) GetTasksByProject(projectID string) ([]domain.Task, error) { return nil, nil }
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
func (m *MockService) GetTagByID(id string) (*domain.Tag, error)               { return nil, nil }
func (m *MockService) GetPerspectiveTasks(name string) ([]domain.Task, error)  { return nil, nil }
func (m *MockService) ResolveProjectName(name string) (string, error)          { return "", nil }

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.Mode() != ModeTagList {
		t.Error("should start in tag list mode")
	}
	if m.loaded {
		t.Error("should not be loaded initially")
	}
}

func TestTagsAndCountsLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}, {ID: "t2", Name: "Tag 2"}},
		counts: map[string]int{"t1": 5, "t2": 3},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(TagsAndCountsLoadedMsg{Tags: svc.tags, Counts: svc.counts})

	if !m.loaded {
		t.Error("should be loaded after TagsAndCountsLoadedMsg")
	}
}

func TestEnterKey_DrillsDown(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
		tasks:  []domain.Task{{ID: "task1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(TagsAndCountsLoadedMsg{Tags: svc.tags, Counts: svc.counts})

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.Mode() != ModeTagTasks {
		t.Error("should switch to task mode after Enter")
	}
	if m.currentTag == nil {
		t.Error("currentTag should be set")
	}
	if cmd == nil {
		t.Error("should return command to load tasks")
	}
}

func TestBackKey_ReturnsToList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(TagsAndCountsLoadedMsg{Tags: svc.tags, Counts: svc.counts})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if m.Mode() != ModeTagList {
		t.Error("should return to tag list mode")
	}
}

func TestHierarchicalTags(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	childTag := domain.Tag{ID: "t2", Name: "Child Tag"}
	parentTag := domain.Tag{ID: "t1", Name: "Parent Tag", Children: []domain.Tag{childTag}}

	svc := &MockService{
		tags:   []domain.Tag{parentTag},
		counts: map[string]int{"t1": 5, "t2": 3},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(TagsAndCountsLoadedMsg{Tags: svc.tags, Counts: svc.counts})

	// Should have 2 tags in flattened list
	tags := m.tagList.Tags()
	if len(tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(tags))
	}

	// Second tag should have depth 1
	if tags[1].Depth != 1 {
		t.Errorf("child tag depth = %d, want 1", tags[1].Depth)
	}
}

func TestSelectedTask_InTaskMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 1},
		tasks:  []domain.Task{{ID: "task1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(TagsAndCountsLoadedMsg{Tags: svc.tags, Counts: svc.counts})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	task := m.SelectedTask()
	if task == nil {
		t.Fatal("should return selected task in task mode")
	}
	if task.ID != "task1" {
		t.Errorf("task ID = %q, want %q", task.ID, "task1")
	}
}
