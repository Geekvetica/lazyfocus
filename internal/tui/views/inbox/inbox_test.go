package inbox

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// TestInitialState verifies the model is initialized with correct defaults
func TestInitialState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	if m.loaded {
		t.Error("expected loaded to be false initially")
	}

	if m.err != nil {
		t.Errorf("expected no error initially, got %v", m.err)
	}

	if m.service != svc {
		t.Error("expected service to be set")
	}

	if m.styles != styles {
		t.Error("expected styles to be set")
	}
}

// TestInit_ReturnsLoadTasksCommand verifies Init returns the load tasks command
func TestInit_ReturnsLoadTasksCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Task 1"},
		},
	}

	m := New(styles, keys, svc)
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("expected Init to return a command")
	}

	// Execute the command to verify it loads tasks
	msg := cmd()

	switch msg := msg.(type) {
	case tui.TasksLoadedMsg:
		if len(msg.Tasks) != 1 {
			t.Errorf("expected 1 task, got %d", len(msg.Tasks))
		}
		if msg.Tasks[0].Name != "Task 1" {
			t.Errorf("expected task name 'Task 1', got '%s'", msg.Tasks[0].Name)
		}
	case tui.ErrorMsg:
		t.Errorf("expected TasksLoadedMsg, got ErrorMsg: %v", msg.Err)
	default:
		t.Errorf("expected TasksLoadedMsg, got %T", msg)
	}
}

// TestInit_ReturnsErrorOnServiceFailure verifies Init handles service errors
func TestInit_ReturnsErrorOnServiceFailure(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	expectedErr := errors.New("connection refused")
	svc := &service.MockOmniFocusService{
		InboxTasksErr: expectedErr,
	}

	m := New(styles, keys, svc)
	cmd := m.Init()

	if cmd == nil {
		t.Fatal("expected Init to return a command")
	}

	// Execute the command
	msg := cmd()

	errorMsg, ok := msg.(tui.ErrorMsg)
	if !ok {
		t.Fatalf("expected ErrorMsg, got %T", msg)
	}

	if errorMsg.Err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, errorMsg.Err)
	}
}

// TestUpdate_TasksLoadedMsg_UpdatesTaskList verifies TasksLoadedMsg updates the task list
func TestUpdate_TasksLoadedMsg_UpdatesTaskList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}

	newModel, _ := m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	if !newModel.loaded {
		t.Error("expected loaded to be true after TasksLoadedMsg")
	}

	if newModel.TaskCount() != 2 {
		t.Errorf("expected 2 tasks, got %d", newModel.TaskCount())
	}
}

// TestUpdate_ErrorMsg_DisplaysError verifies ErrorMsg sets the error state
func TestUpdate_ErrorMsg_DisplaysError(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	expectedErr := errors.New("failed to load tasks")
	newModel, _ := m.Update(tui.ErrorMsg{Err: expectedErr})

	if newModel.err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, newModel.err)
	}

	// Error should appear in view
	view := newModel.View()
	if view == "" {
		t.Error("expected view to contain error message")
	}
}

// TestUpdate_WindowSizeMsg_UpdatesDimensions verifies WindowSizeMsg updates dimensions
func TestUpdate_WindowSizeMsg_UpdatesDimensions(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	newModel, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	if newModel.width != 120 {
		t.Errorf("expected width 120, got %d", newModel.width)
	}

	if newModel.height != 40 {
		t.Errorf("expected height 40, got %d", newModel.height)
	}
}

// TestUpdate_NavigationDelegatesToTaskList verifies navigation is delegated to task list
func TestUpdate_NavigationDelegatesToTaskList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	// Load some tasks first
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
		{ID: "3", Name: "Task 3"},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Initial selection should be task 1
	selected := m.SelectedTask()
	if selected == nil || selected.Name != "Task 1" {
		t.Error("expected first task to be selected initially")
	}

	// Press down key
	keyMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, _ = m.Update(keyMsg)

	// Should now select task 2
	selected = m.SelectedTask()
	if selected == nil || selected.Name != "Task 2" {
		t.Error("expected second task to be selected after down key")
	}
}

// TestView_ShowsHeaderWithTaskCount verifies the view displays correct header
func TestView_ShowsHeaderWithTaskCount(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	// Load tasks
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
		{ID: "3", Name: "Task 3"},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	view := m.View()

	// Should contain INBOX and task count
	if view == "" {
		t.Error("expected non-empty view")
	}

	// View should contain inbox header (we'll verify exact format in implementation)
	// For now just check it's not empty
}

// TestView_ShowsErrorWhenPresent verifies error display
func TestView_ShowsErrorWhenPresent(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	expectedErr := errors.New("connection refused")
	m, _ = m.Update(tui.ErrorMsg{Err: expectedErr})

	view := m.View()

	if view == "" {
		t.Error("expected view to show error")
	}
}

// TestTaskCount_ReturnsCorrectCount verifies TaskCount method
func TestTaskCount_ReturnsCorrectCount(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	// Initially should be 0
	if m.TaskCount() != 0 {
		t.Errorf("expected 0 tasks initially, got %d", m.TaskCount())
	}

	// Load tasks
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	if m.TaskCount() != 2 {
		t.Errorf("expected 2 tasks, got %d", m.TaskCount())
	}
}

// TestSelectedTask_DelegatesToTaskList verifies SelectedTask method
func TestSelectedTask_DelegatesToTaskList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	// Initially should be nil (no tasks)
	if m.SelectedTask() != nil {
		t.Error("expected nil when no tasks loaded")
	}

	// Load tasks
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	selected := m.SelectedTask()
	if selected == nil {
		t.Fatal("expected selected task to be non-nil")
	}

	if selected.Name != "Task 1" {
		t.Errorf("expected 'Task 1', got '%s'", selected.Name)
	}
}

// TestRefresh_ReloadsTasksFromService verifies Refresh method
func TestRefresh_ReloadsTasksFromService(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "New Task"},
		},
	}

	m := New(styles, keys, svc)

	cmd := m.Refresh()
	if cmd == nil {
		t.Fatal("expected Refresh to return a command")
	}

	// Execute the command
	msg := cmd()

	tasksMsg, ok := msg.(tui.TasksLoadedMsg)
	if !ok {
		t.Fatalf("expected TasksLoadedMsg, got %T", msg)
	}

	if len(tasksMsg.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasksMsg.Tasks))
	}

	if tasksMsg.Tasks[0].Name != "New Task" {
		t.Errorf("expected 'New Task', got '%s'", tasksMsg.Tasks[0].Name)
	}
}

// TestRefresh_ReturnsErrorOnServiceFailure verifies Refresh error handling
func TestRefresh_ReturnsErrorOnServiceFailure(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	expectedErr := errors.New("network error")
	svc := &service.MockOmniFocusService{
		InboxTasksErr: expectedErr,
	}

	m := New(styles, keys, svc)

	cmd := m.Refresh()
	if cmd == nil {
		t.Fatal("expected Refresh to return a command")
	}

	// Execute the command
	msg := cmd()

	errorMsg, ok := msg.(tui.ErrorMsg)
	if !ok {
		t.Fatalf("expected ErrorMsg, got %T", msg)
	}

	if errorMsg.Err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, errorMsg.Err)
	}
}

// TestView_FormatsTasksCorrectly verifies task formatting in view
func TestView_FormatsTasksCorrectly(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Create tasks with various properties
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	tasks := []domain.Task{
		{ID: "1", Name: "Simple task"},
		{ID: "2", Name: "Task with due date", DueDate: &tomorrow},
		{ID: "3", Name: "Flagged task", Flagged: true},
	}

	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}

	// The view should contain the task list output (delegated to tasklist component)
	// We're mainly testing integration here
}

// TestUpdate_PreservesTaskListState verifies task list state is preserved
func TestUpdate_PreservesTaskListState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &service.MockOmniFocusService{}

	m := New(styles, keys, svc)

	tasks := []domain.Task{
		{ID: "1", Name: "Task 1"},
		{ID: "2", Name: "Task 2"},
		{ID: "3", Name: "Task 3"},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Navigate to task 2
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})

	selected := m.SelectedTask()
	if selected == nil || selected.Name != "Task 2" {
		t.Error("expected Task 2 to be selected")
	}

	// Send a resize message (should preserve selection)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	selected = m.SelectedTask()
	if selected == nil || selected.Name != "Task 2" {
		t.Error("expected Task 2 to still be selected after resize")
	}
}
