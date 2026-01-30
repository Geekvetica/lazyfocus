package review

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/filter"
)

// MockService for testing
type MockService struct {
	tasks     []domain.Task
	returnErr error
}

func (m *MockService) GetFlaggedTasks() ([]domain.Task, error) {
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	return m.tasks, nil
}

// Stub other methods
func (m *MockService) GetInboxTasks() ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetAllTasks(_ service.TaskFilters) ([]domain.Task, error) {
	return nil, nil
}
func (m *MockService) GetTasksByProject(_ string) ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetTasksByTag(_ string) ([]domain.Task, error)     { return nil, nil }
func (m *MockService) GetTaskByID(_ string) (*domain.Task, error)        { return nil, nil }
func (m *MockService) CreateTask(_ domain.TaskInput) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) ModifyTask(_ string, _ domain.TaskModification) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) CompleteTask(_ string) (*domain.OperationResult, error) { return nil, nil }
func (m *MockService) DeleteTask(_ string) (*domain.OperationResult, error)   { return nil, nil }
func (m *MockService) GetProjects(_ string) ([]domain.Project, error)         { return nil, nil }
func (m *MockService) GetProjectByID(_ string) (*domain.Project, error)       { return nil, nil }
func (m *MockService) GetProjectWithTasks(_ string) (*domain.Project, error)  { return nil, nil }
func (m *MockService) GetTags() ([]domain.Tag, error)                         { return nil, nil }
func (m *MockService) GetTagByID(_ string) (*domain.Tag, error)               { return nil, nil }
func (m *MockService) GetTagCounts() (map[string]int, error)                  { return nil, nil }
func (m *MockService) GetPerspectiveTasks(_ string) ([]domain.Task, error)    { return nil, nil }
func (m *MockService) ResolveProjectName(_ string) (string, error)            { return "", nil }

// Helper to create a test model with default configuration
func newTestReviewModel() Model {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tasks: []domain.Task{
			{ID: "1", Name: "Important Task", Flagged: true},
			{ID: "2", Name: "Urgent Item", Flagged: true},
		},
	}
	return New(styles, keys, svc)
}

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.loaded {
		t.Error("should not be loaded initially")
	}
	if m.taskCount != 0 {
		t.Errorf("task count = %d, want 0", m.taskCount)
	}
	if m.err != nil {
		t.Errorf("err = %v, want nil", m.err)
	}
}

// 1. Init and Task Loading Tests

func TestInit_ReturnsLoadFlaggedTasksCommand(t *testing.T) {
	m := newTestReviewModel()

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() should return a command to load flagged tasks")
	}

	// Execute the command and verify it returns TasksLoadedMsg
	msg := cmd()
	if _, ok := msg.(tui.TasksLoadedMsg); !ok {
		t.Errorf("Init command returned %T, want tui.TasksLoadedMsg", msg)
	}
}

func TestLoadFlaggedTasks_LoadsOnlyFlaggedTasks(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		tasks: []domain.Task{
			{ID: "1", Name: "Flagged 1", Flagged: true},
			{ID: "2", Name: "Flagged 2", Flagged: true},
		},
	}

	m := New(styles, keys, svc)
	cmd := m.loadFlaggedTasks()

	msg := cmd()
	tasksMsg, ok := msg.(tui.TasksLoadedMsg)
	if !ok {
		t.Fatalf("expected TasksLoadedMsg, got %T", msg)
	}

	if len(tasksMsg.Tasks) != 2 {
		t.Errorf("loaded tasks count = %d, want 2", len(tasksMsg.Tasks))
	}
	for i, task := range tasksMsg.Tasks {
		if !task.Flagged {
			t.Errorf("task[%d] flagged = false, want true", i)
		}
	}
}

func TestLoadFlaggedTasks_ReturnsErrorOnServiceFailure(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	expectedErr := errors.New("service error")
	svc := &MockService{
		returnErr: expectedErr,
	}

	m := New(styles, keys, svc)
	cmd := m.loadFlaggedTasks()

	msg := cmd()
	errMsg, ok := msg.(tui.ErrorMsg)
	if !ok {
		t.Fatalf("expected ErrorMsg, got %T", msg)
	}

	if errMsg.Err != expectedErr {
		t.Errorf("error = %v, want %v", errMsg.Err, expectedErr)
	}
}

// 2. Update Message Dispatch Tests

func TestUpdate_TasksLoadedMsg_StoresTasks(t *testing.T) {
	m := newTestReviewModel()
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", Flagged: true},
		{ID: "2", Name: "Task 2", Flagged: true},
	}

	newM, _ := m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	if !newM.loaded {
		t.Error("should be loaded after TasksLoadedMsg")
	}
	if newM.taskCount != 2 {
		t.Errorf("task count = %d, want 2", newM.taskCount)
	}
	if newM.err != nil {
		t.Errorf("err should be cleared, got %v", newM.err)
	}
	if len(newM.allTasks) != 2 {
		t.Errorf("allTasks length = %d, want 2", len(newM.allTasks))
	}
}

func TestUpdate_TasksLoadedMsg_EmptyList(t *testing.T) {
	m := newTestReviewModel()

	newM, _ := m.Update(tui.TasksLoadedMsg{Tasks: []domain.Task{}})

	if !newM.loaded {
		t.Error("should be loaded even with empty tasks")
	}
	if newM.taskCount != 0 {
		t.Errorf("task count = %d, want 0", newM.taskCount)
	}
}

func TestUpdate_WindowSizeMsg_UpdatesDimensions(t *testing.T) {
	m := newTestReviewModel()

	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	if newM.width != 100 {
		t.Errorf("width = %d, want 100", newM.width)
	}
	if newM.height != 50 {
		t.Errorf("height = %d, want 50", newM.height)
	}
}

func TestUpdate_WindowSizeMsg_HandlesSmallHeight(t *testing.T) {
	m := newTestReviewModel()

	// Header height is 3, so this should result in 0 available height
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 2})

	if newM.height != 2 {
		t.Errorf("height = %d, want 2", newM.height)
	}
	// Should not panic with negative height
}

func TestUpdate_ErrorMsg_SetsErrorState(t *testing.T) {
	m := newTestReviewModel()
	expectedErr := errors.New("test error")

	newM, _ := m.Update(tui.ErrorMsg{Err: expectedErr})

	if newM.err == nil {
		t.Fatal("error should be set")
	}
	if newM.err.Error() != expectedErr.Error() {
		t.Errorf("error = %v, want %v", newM.err, expectedErr)
	}
}

func TestUpdate_OtherMessages_PassedToTaskList(t *testing.T) {
	m := newTestReviewModel()
	// Load tasks first
	m, _ = m.Update(tui.TasksLoadedMsg{
		Tasks: []domain.Task{
			{ID: "1", Name: "Task 1", Flagged: true},
		},
	})

	// Pass a key message - should be handled by task list
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}
	newM, _ := m.Update(keyMsg)

	// Model should still be valid (not panicking is success)
	if newM.loaded != m.loaded {
		t.Error("loaded state should not change")
	}
}

// 3. View Rendering Tests

func TestView_ShowsFlaggedTasks(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{
		Tasks: []domain.Task{
			{ID: "1", Name: "Important Task", Flagged: true},
			{ID: "2", Name: "Urgent Item", Flagged: true},
		},
	})

	view := m.View()

	if !strings.Contains(view, "REVIEW") {
		t.Error("view should contain REVIEW header")
	}
	// Task list should be rendered (even if empty in test)
	if view == "" {
		t.Error("view should not be empty")
	}
}

func TestView_HeaderRendering_ShowsCount(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{
		Tasks: []domain.Task{
			{ID: "1", Name: "Task 1", Flagged: true},
			{ID: "2", Name: "Task 2", Flagged: true},
			{ID: "3", Name: "Task 3", Flagged: true},
		},
	})

	view := m.View()

	// Header should show task count
	if !strings.Contains(view, "3)") {
		t.Error("header should show task count of 3")
	}
}

func TestView_HeaderRendering_ShowsSubtext(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: []domain.Task{}})

	view := m.View()

	// Should contain help text - check for key parts that won't be affected by styling
	if !strings.Contains(view, "Review") {
		t.Errorf("header should contain 'Review', got: %q", view)
	}
	if !strings.Contains(view, "flagged") {
		t.Errorf("header should mention flagged tasks, got: %q", view)
	}
	// The styling may affect how the text appears, so just check that the view is not empty
	// and contains basic text
	if len(view) < 10 {
		t.Errorf("view seems too short, got: %q", view)
	}
}

func TestView_ErrorState_DisplaysError(t *testing.T) {
	m := newTestReviewModel()
	testErr := errors.New("failed to load tasks")
	m, _ = m.Update(tui.ErrorMsg{Err: testErr})

	view := m.View()

	if !strings.Contains(view, "Error") {
		t.Error("view should contain 'Error' text")
	}
	if !strings.Contains(view, "failed to load tasks") {
		t.Error("view should contain error message")
	}
}

func TestView_ErrorState_ShowsSeparator(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tea.WindowSizeMsg{Width: 50, Height: 20})
	m, _ = m.Update(tui.ErrorMsg{Err: errors.New("test error")})

	view := m.View()

	// Error view should contain separator
	if !strings.Contains(view, "─") {
		t.Error("error view should contain separator line")
	}
}

func TestView_ErrorState_WithZeroWidth(t *testing.T) {
	m := newTestReviewModel()
	// Don't set width (defaults to 0)
	m, _ = m.Update(tui.ErrorMsg{Err: errors.New("test error")})

	view := m.View()

	// Should use default width (40) for separator
	if !strings.Contains(view, "─") {
		t.Error("error view should contain separator even with zero width")
	}
}

// 4. Refresh Functionality Tests

func TestRefresh_ReloadsFlaggedTasks(t *testing.T) {
	m := newTestReviewModel()

	cmd := m.Refresh()
	if cmd == nil {
		t.Fatal("Refresh() should return a command")
	}

	// Execute the command
	msg := cmd()
	if _, ok := msg.(tui.TasksLoadedMsg); !ok {
		t.Errorf("Refresh command returned %T, want tui.TasksLoadedMsg", msg)
	}
}

// 5. Filter Functionality Tests

func TestSetFilter_AppliesSearchFilter(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{
		Tasks: []domain.Task{
			{ID: "1", Name: "Important Task", Flagged: true},
			{ID: "2", Name: "Urgent Item", Flagged: true},
		},
	})

	// Apply search filter
	filterState := filter.State{}.WithSearchText("Important")
	newM := m.SetFilter(filterState)

	// Task count should be updated based on filter
	// Note: actual filtering is done by filter.Matcher, we just verify the filter was applied
	if newM.filter.SearchText != "Important" {
		t.Errorf("filter search text = %q, want %q", newM.filter.SearchText, "Important")
	}
}

func TestSetFilter_ClearFilter(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{
		Tasks: []domain.Task{
			{ID: "1", Name: "Task 1", Flagged: true},
			{ID: "2", Name: "Task 2", Flagged: true},
		},
	})

	// Apply filter first
	filterState := filter.State{}.WithSearchText("test")
	m = m.SetFilter(filterState)

	// Clear filter
	clearedFilter := filter.State{}
	newM := m.SetFilter(clearedFilter)

	if newM.filter.IsActive() {
		t.Error("filter should not be active after clearing")
	}
}

// 6. SelectedTask Tests

func TestSelectedTask_ReturnsCurrentSelection(t *testing.T) {
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

func TestSelectedTask_ReturnsNilWhenEmpty(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: []domain.Task{}})

	task := m.SelectedTask()
	if task != nil {
		t.Error("should return nil when no tasks")
	}
}

// 7. TaskCount Tests

func TestTaskCount_ReturnsCorrectCount(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{
		Tasks: []domain.Task{
			{ID: "1", Name: "Task 1", Flagged: true},
			{ID: "2", Name: "Task 2", Flagged: true},
			{ID: "3", Name: "Task 3", Flagged: true},
		},
	})

	count := m.TaskCount()
	if count != 3 {
		t.Errorf("TaskCount() = %d, want 3", count)
	}
}

func TestTaskCount_ReturnsZeroWhenEmpty(t *testing.T) {
	m := newTestReviewModel()
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: []domain.Task{}})

	count := m.TaskCount()
	if count != 0 {
		t.Errorf("TaskCount() = %d, want 0", count)
	}
}
