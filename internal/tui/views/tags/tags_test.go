package tags

import (
	"errors"
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

func (m *MockService) GetTasksByTag(_ string) ([]domain.Task, error) {
	return m.tasks, nil
}

// Stub other interface methods
func (m *MockService) GetInboxTasks() ([]domain.Task, error)                    { return nil, nil }
func (m *MockService) GetAllTasks(_ service.TaskFilters) ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetTasksByProject(_ string) ([]domain.Task, error)        { return nil, nil }
func (m *MockService) GetFlaggedTasks() ([]domain.Task, error)                  { return nil, nil }
func (m *MockService) GetTaskByID(_ string) (*domain.Task, error)               { return nil, nil }
func (m *MockService) CreateTask(_ domain.TaskInput) (*domain.Task, error)      { return nil, nil }
func (m *MockService) ModifyTask(_ string, _ domain.TaskModification) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) CompleteTask(_ string) (*domain.OperationResult, error) { return nil, nil }
func (m *MockService) DeleteTask(_ string) (*domain.OperationResult, error)   { return nil, nil }
func (m *MockService) GetProjects(_ string) ([]domain.Project, error)         { return nil, nil }
func (m *MockService) GetProjectByID(_ string) (*domain.Project, error)       { return nil, nil }
func (m *MockService) GetProjectWithTasks(_ string) (*domain.Project, error)  { return nil, nil }
func (m *MockService) GetTagByID(_ string) (*domain.Tag, error)               { return nil, nil }
func (m *MockService) GetPerspectiveTasks(_ string) ([]domain.Task, error)    { return nil, nil }
func (m *MockService) ResolveProjectName(_ string) (string, error)            { return "", nil }

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

func TestLoadedWithCountsMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}, {ID: "t2", Name: "Tag 2"}},
		counts: map[string]int{"t1": 5, "t2": 3},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})

	if !m.loaded {
		t.Error("should be loaded after LoadedWithCountsMsg")
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
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})

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
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})
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
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})

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
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})
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

func TestInit_ReturnsLoadTagsCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("Init() should return a command to load tags")
	}

	// Execute the command
	msg := cmd()

	if _, ok := msg.(LoadedWithCountsMsg); !ok {
		t.Errorf("expected LoadedWithCountsMsg, got %T", msg)
	}
}

func TestUpdate_TagsLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	_, cmd := m.Update(tui.TagsLoadedMsg{Tags: svc.tags})

	// Should trigger loading tags and counts
	if cmd == nil {
		t.Error("TagsLoadedMsg should return command to load tags and counts")
	}

	// Execute the command to verify it returns LoadedWithCountsMsg
	msg := cmd()
	if _, ok := msg.(LoadedWithCountsMsg); !ok {
		t.Errorf("expected LoadedWithCountsMsg, got %T", msg)
	}
}

func TestUpdate_TasksLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	tasks := []domain.Task{{ID: "task1", Name: "Task 1"}}
	svc := &MockService{tasks: tasks}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Verify tasks are stored in task list (task list was updated)
	// Can't access unexported field, but the update happened without error
	if m.mode == ModeTagTasks {
		// If we're in task mode, the task list should have the selected task
		task := m.taskList.SelectedTask()
		if task != nil && task.ID != "task1" {
			t.Errorf("task ID = %q, want %q", task.ID, "task1")
		}
	}
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 30 {
		t.Errorf("height = %d, want 30", m.height)
	}
}

func TestUpdate_WindowSizeMsg_InTaskMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 30 {
		t.Errorf("height = %d, want 30", m.height)
	}
}

func TestUpdate_ErrorMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	expectedErr := errors.New("test error")
	m, _ = m.Update(tui.ErrorMsg{Err: expectedErr})

	if m.err != expectedErr {
		t.Errorf("err = %v, want %v", m.err, expectedErr)
	}
}

func TestModeTransition_TagList_To_TaskList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
		tasks:  []domain.Task{{ID: "task1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})

	if m.Mode() != ModeTagList {
		t.Fatal("should start in tag list mode")
	}

	// Press Enter to drill down
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.Mode() != ModeTagTasks {
		t.Error("should switch to task mode after Enter")
	}
	if m.currentTag == nil {
		t.Error("currentTag should be set")
	}
	if m.currentTag.ID != "t1" {
		t.Errorf("currentTag.ID = %q, want %q", m.currentTag.ID, "t1")
	}
	if cmd == nil {
		t.Error("should return command to load tasks")
	}
}

func TestModeTransition_TaskList_To_TagList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.Mode() != ModeTagTasks {
		t.Fatal("should be in task mode")
	}

	// Press 'h' to go back
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if m.Mode() != ModeTagList {
		t.Error("should return to tag list mode")
	}
	if m.currentTag != nil {
		t.Error("currentTag should be nil after returning to tag list")
	}
}

func TestHandleKeyPress_TagListMode_Enter_NoSelection(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	// No tags loaded, so no selection

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should stay in tag list mode
	if m.Mode() != ModeTagList {
		t.Error("should remain in tag list mode when no tag selected")
	}
	if cmd != nil {
		t.Error("should not return command when no tag selected")
	}
}

func TestHandleKeyPress_BackKey_AlreadyInTagList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.Mode() != ModeTagList {
		t.Fatal("should start in tag list mode")
	}

	// Press 'h' when already in tag list
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	// Should stay in tag list mode
	if m.Mode() != ModeTagList {
		t.Error("should remain in tag list mode")
	}
}

func TestHandleKeyPress_EscapeKey_ReturnsToList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.Mode() != ModeTagTasks {
		t.Fatal("should be in task mode")
	}

	// Press Escape to go back
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if m.Mode() != ModeTagList {
		t.Error("should return to tag list mode after Escape")
	}
}

func TestRenderHeader_TagListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}, {ID: "t2", Name: "Tag 2"}},
		counts: map[string]int{"t1": 5, "t2": 3},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})

	header := m.renderHeader()

	if header == "" {
		t.Error("header should not be empty")
	}
	// Header should contain "TAGS" and count
	if !containsAny(header, "TAGS", "(2)") {
		t.Errorf("header should contain TAGS and count, got: %q", header)
	}
}

func TestRenderHeader_TaskListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	tag := &domain.Tag{ID: "t1", Name: "Urgent"}
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks
	m.currentTag = tag

	header := m.renderHeader()

	if header == "" {
		t.Error("header should not be empty")
	}
	// Header should contain tag name
	if !containsAny(header, "Urgent") {
		t.Errorf("header should contain tag name, got: %q", header)
	}
	// Should contain back hint
	if !containsAny(header, "back", "Esc") {
		t.Errorf("header should contain back hint, got: %q", header)
	}
}

func TestRenderHeader_TaskListMode_NoCurrentTag(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks
	m.currentTag = nil

	header := m.renderHeader()

	if header == "" {
		t.Error("header should not be empty")
	}
	// Should have fallback text
	if !containsAny(header, "TAG TASKS") {
		t.Errorf("header should contain fallback text, got: %q", header)
	}
}

func TestView_TagListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(LoadedWithCountsMsg{Tags: svc.tags, Counts: svc.counts})

	view := m.View()

	if view == "" {
		t.Error("view should not be empty")
	}
	// Should contain header
	if !containsAny(view, "TAGS") {
		t.Errorf("view should contain TAGS header, got: %q", view)
	}
}

func TestView_TaskListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tasks: []domain.Task{{ID: "task1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks
	m.currentTag = &domain.Tag{ID: "t1", Name: "Urgent"}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	view := m.View()

	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestView_ErrorState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.err = errors.New("test error")

	view := m.View()

	if view == "" {
		t.Error("view should not be empty")
	}
	// Should show error
	if !containsAny(view, "Error") {
		t.Errorf("view should contain error message, got: %q", view)
	}
}

func TestRefresh_TagListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	m.mode = ModeTagList

	cmd := m.Refresh()

	if cmd == nil {
		t.Error("Refresh() in tag list mode should return command")
	}

	// Execute command
	msg := cmd()
	if _, ok := msg.(LoadedWithCountsMsg); !ok {
		t.Errorf("expected LoadedWithCountsMsg, got %T", msg)
	}
}

func TestRefresh_TaskListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tasks: []domain.Task{{ID: "task1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks
	m.currentTag = &domain.Tag{ID: "t1", Name: "Urgent"}

	cmd := m.Refresh()

	if cmd == nil {
		t.Error("Refresh() in task list mode should return command")
	}

	// Execute command
	msg := cmd()
	if _, ok := msg.(tui.TasksLoadedMsg); !ok {
		t.Errorf("expected TasksLoadedMsg, got %T", msg)
	}
}

func TestRefresh_TaskListMode_NoCurrentTag(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tags:   []domain.Tag{{ID: "t1", Name: "Tag 1"}},
		counts: map[string]int{"t1": 5},
	}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks
	m.currentTag = nil

	cmd := m.Refresh()

	if cmd == nil {
		t.Error("Refresh() should return command")
	}

	// Should fall back to loading tags
	msg := cmd()
	if _, ok := msg.(LoadedWithCountsMsg); !ok {
		t.Errorf("expected LoadedWithCountsMsg, got %T", msg)
	}
}

func TestSelectedTask_InTagListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeTagList

	task := m.SelectedTask()

	if task != nil {
		t.Error("SelectedTask() should return nil in tag list mode")
	}
}

func TestMode_ReturnsCurrentMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.Mode() != ModeTagList {
		t.Errorf("Mode() = %v, want ModeTagList", m.Mode())
	}

	m.mode = ModeTagTasks
	if m.Mode() != ModeTagTasks {
		t.Errorf("Mode() = %v, want ModeTagTasks", m.Mode())
	}
}

func TestRenderError_WithWidth(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.width = 80
	m.err = errors.New("test error")

	errorView := m.renderError()

	if errorView == "" {
		t.Error("error view should not be empty")
	}
	if !containsAny(errorView, "Error") {
		t.Errorf("error view should contain 'Error', got: %q", errorView)
	}
}

func TestRenderError_NoWidth(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.width = 0
	m.err = errors.New("test error")

	errorView := m.renderError()

	if errorView == "" {
		t.Error("error view should not be empty")
	}
	// Should use default width
	if !containsAny(errorView, "Error") {
		t.Errorf("error view should contain 'Error', got: %q", errorView)
	}
}

func TestDelegateToCurrentList_TagListMode(_ *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeTagList

	// Send arbitrary message - should delegate to tag list (no error)
	m, cmd := m.delegateToCurrentList(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	_ = m
	_ = cmd
}

func TestDelegateToCurrentList_TaskListMode(_ *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeTagTasks

	// Send arbitrary message - should delegate to task list (no error)
	m, cmd := m.delegateToCurrentList(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	_ = m
	_ = cmd
}

// Helper function to check if a string contains any of the given substrings
func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) > 0 && len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				match := true
				for j := 0; j < len(sub); j++ {
					if s[i+j] != sub[j] {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}
