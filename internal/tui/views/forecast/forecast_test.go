package forecast

import (
	"errors"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/filter"
)

// MockService for testing
type MockService struct {
	tasks []domain.Task
}

func (m *MockService) GetAllTasks(_ service.TaskFilters) ([]domain.Task, error) {
	return m.tasks, nil
}

// Stub other methods
func (m *MockService) GetInboxTasks() ([]domain.Task, error)               { return nil, nil }
func (m *MockService) GetTasksByProject(_ string) ([]domain.Task, error)   { return nil, nil }
func (m *MockService) GetTasksByTag(_ string) ([]domain.Task, error)       { return nil, nil }
func (m *MockService) GetFlaggedTasks() ([]domain.Task, error)             { return nil, nil }
func (m *MockService) GetTaskByID(_ string) (*domain.Task, error)          { return nil, nil }
func (m *MockService) CreateTask(_ domain.TaskInput) (*domain.Task, error) { return nil, nil }
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

func TestRenderGroupHeader_IconsCorrect(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	tests := []struct {
		name           string
		group          DueGroup
		collapsed      bool
		expectedIcon   string
		expectedNotice string
	}{
		{
			name:           "expanded group shows down arrow",
			group:          GroupToday,
			collapsed:      false,
			expectedIcon:   "▼",
			expectedNotice: "down arrow indicates expanded (can collapse)",
		},
		{
			name:           "collapsed group shows right arrow",
			group:          GroupToday,
			collapsed:      true,
			expectedIcon:   "▶",
			expectedNotice: "right arrow indicates collapsed (can expand)",
		},
		{
			name:           "overdue expanded shows down arrow",
			group:          GroupOverdue,
			collapsed:      false,
			expectedIcon:   "▼",
			expectedNotice: "down arrow indicates expanded (can collapse)",
		},
		{
			name:           "overdue collapsed shows right arrow",
			group:          GroupOverdue,
			collapsed:      true,
			expectedIcon:   "▶",
			expectedNotice: "right arrow indicates collapsed (can expand)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set collapse state
			m.collapsed = map[DueGroup]bool{tt.group: tt.collapsed}

			// Render the header
			header := m.renderGroupHeader(tt.group, false)

			// Check if the expected icon is present in the rendered output
			if !contains(header, tt.expectedIcon) {
				t.Errorf("expected icon %q in header, got: %s\nNote: %s",
					tt.expectedIcon, header, tt.expectedNotice)
			}
		})
	}
}

// contains checks if a string contains a substring
// (handles styled strings by checking the raw content)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}

// TestInit verifies Init returns a command to load tasks
func TestInit_ReturnsLoadTasksCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{tasks: []domain.Task{{ID: "1", Name: "Task"}}}
	m := New(styles, keys, svc)

	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init() should return a command")
	}

	// Execute the command and verify it loads tasks
	msg := cmd()
	if _, ok := msg.(tui.TasksLoadedMsg); !ok {
		t.Errorf("expected TasksLoadedMsg, got %T", msg)
	}
}

// TestUpdate_TasksLoadedMsg verifies tasks are categorized correctly
func TestUpdate_TasksLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)

	tasks := []domain.Task{
		{ID: "1", Name: "Today task", DueDate: &today},
		{ID: "2", Name: "Tomorrow task", DueDate: &tomorrow},
	}

	m, cmd := m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	if cmd != nil {
		t.Error("TasksLoadedMsg should not return a command")
	}
	if !m.loaded {
		t.Error("model should be marked as loaded")
	}
	if m.err != nil {
		t.Errorf("should not have error, got %v", m.err)
	}
	if len(m.items) == 0 {
		t.Error("items should not be empty")
	}
	// Cursor should be on first task (after header)
	if len(m.items) > 1 && m.items[0].IsHeader && m.cursor != 1 {
		t.Errorf("cursor should be at position 1 (first task), got %d", m.cursor)
	}
}

// TestUpdate_WindowSizeMsg verifies dimensions are updated
func TestUpdate_WindowSizeMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	m, cmd := m.Update(msg)

	if cmd != nil {
		t.Error("WindowSizeMsg should not return a command")
	}
	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}

// TestUpdate_ErrorMsg verifies error state is set
func TestUpdate_ErrorMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	testErr := errors.New("test error")
	errMsg := tui.ErrorMsg{Err: testErr}
	m, cmd := m.Update(errMsg)

	if cmd != nil {
		t.Error("ErrorMsg should not return a command")
	}
	if m.err == nil {
		t.Error("error should be set")
	}
	if m.err != testErr {
		t.Errorf("error = %v, want %v", m.err, testErr)
	}
}

// TestHandleKeyPress_NavigationKeys verifies j/k/up/down navigation
func TestHandleKeyPress_NavigationKeys(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", DueDate: &now},
		{ID: "2", Name: "Task 2", DueDate: &now},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	initialCursor := m.cursor

	// Test down key
	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	m, cmd := m.Update(downMsg)
	if cmd != nil {
		t.Error("navigation should not return a command")
	}
	if m.cursor <= initialCursor {
		t.Errorf("cursor should move down from %d, got %d", initialCursor, m.cursor)
	}

	// Test up key
	cursorAfterDown := m.cursor
	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	m, cmd = m.Update(upMsg)
	if cmd != nil {
		t.Error("navigation should not return a command")
	}
	if m.cursor >= cursorAfterDown {
		t.Errorf("cursor should move up from %d, got %d", cursorAfterDown, m.cursor)
	}
}

// TestHandleKeyPress_EnterKey_ToggleCollapse verifies Enter on header toggles collapsed state
func TestHandleKeyPress_EnterKey_ToggleCollapse(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{{ID: "1", Name: "Task", DueDate: &now}}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Move cursor to header (position 0)
	m.cursor = 0

	if !m.items[0].IsHeader {
		t.Fatal("first item should be a header")
	}

	group := m.items[0].Group
	initialCollapsed := m.collapsed[group]

	// Press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, cmd := m.Update(enterMsg)
	if cmd != nil {
		t.Error("enter on header should not return a command")
	}

	// Verify collapse state toggled
	if m.collapsed[group] == initialCollapsed {
		t.Error("collapse state should have toggled")
	}
}

// TestHandleKeyPress_EnterKey_OnTask verifies Enter on task does not toggle
func TestHandleKeyPress_EnterKey_OnTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{{ID: "1", Name: "Task", DueDate: &now}}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Move cursor to task (position 1, after header)
	m.cursor = 1

	if m.items[1].IsHeader {
		t.Fatal("second item should be a task")
	}

	// Store collapse state
	collapsedBefore := make(map[DueGroup]bool)
	for k, v := range m.collapsed {
		collapsedBefore[k] = v
	}

	// Press Enter
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	m, cmd := m.Update(enterMsg)
	if cmd != nil {
		t.Error("enter on task should not return a command")
	}

	// Verify collapse state unchanged
	for group := range collapsedBefore {
		if m.collapsed[group] != collapsedBefore[group] {
			t.Errorf("collapse state for %v should not change when Enter on task", group)
		}
	}
}

// TestHandleKeyPress_EmptyList verifies navigation on empty list doesn't panic
func TestHandleKeyPress_EmptyList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	m.items = []GroupedTask{} // Empty list

	downMsg := tea.KeyMsg{Type: tea.KeyDown}
	_, cmd := m.Update(downMsg)
	if cmd != nil {
		t.Error("navigation on empty list should not return a command")
	}

	upMsg := tea.KeyMsg{Type: tea.KeyUp}
	_, cmd = m.Update(upMsg)
	if cmd != nil {
		t.Error("navigation on empty list should not return a command")
	}
}

// TestCategorizeTask verifies task categorization boundary conditions
func TestCategorizeTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	weekEnd := today.AddDate(0, 0, 7)

	tests := []struct {
		name     string
		due      *time.Time
		expected DueGroup
	}{
		{
			name:     "overdue task - yesterday",
			due:      timePtr(today.AddDate(0, 0, -1)),
			expected: GroupOverdue,
		},
		{
			name:     "today task - noon today",
			due:      timePtr(today.Add(12 * time.Hour)),
			expected: GroupToday,
		},
		{
			name:     "today task - edge case midnight",
			due:      timePtr(today),
			expected: GroupToday,
		},
		{
			name:     "tomorrow task",
			due:      timePtr(tomorrow.Add(12 * time.Hour)),
			expected: GroupTomorrow,
		},
		{
			name:     "this week task - 3 days from now",
			due:      timePtr(today.AddDate(0, 0, 3)),
			expected: GroupThisWeek,
		},
		{
			name:     "this week task - day 6",
			due:      timePtr(today.AddDate(0, 0, 6)),
			expected: GroupThisWeek,
		},
		{
			name:     "later task - next week",
			due:      timePtr(today.AddDate(0, 0, 10)),
			expected: GroupLater,
		},
		{
			name:     "later task - next month",
			due:      timePtr(today.AddDate(0, 1, 0)),
			expected: GroupLater,
		},
		{
			name:     "no due date - nil",
			due:      nil,
			expected: GroupNoDue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := domain.Task{ID: "test", Name: tt.name, DueDate: tt.due}
			group := m.categorizeTask(task, today, tomorrow, weekEnd)
			if group != tt.expected {
				t.Errorf("categorizeTask(%s) = %v, want %v", tt.name, groupName(group), groupName(tt.expected))
			}
		})
	}
}

// timePtr returns a pointer to a time value
func timePtr(t time.Time) *time.Time {
	return &t
}

// TestView_ShowsGroupHeaders verifies group headers are rendered
func TestView_ShowsGroupHeaders(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)

	tasks := []domain.Task{
		{ID: "1", Name: "Today task", DueDate: &today},
		{ID: "2", Name: "Tomorrow task", DueDate: &tomorrow},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	view := m.View()

	// Check for group names
	if !contains(view, "Today") {
		t.Error("view should contain 'Today' header")
	}
	if !contains(view, "Tomorrow") {
		t.Error("view should contain 'Tomorrow' header")
	}
}

// TestView_ShowsCorrectIcons verifies collapse/expand icons
func TestView_ShowsCorrectIcons(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{{ID: "1", Name: "Task", DueDate: &now}}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Initially expanded - should show down arrow (▼)
	view := m.View()
	if !contains(view, "▼") {
		t.Error("expanded view should contain down arrow (▼)")
	}

	// Collapse the group
	group := m.items[0].Group
	m.collapsed[group] = true
	m.items = m.rebuildItems()

	// Should now show right arrow (▶)
	view = m.View()
	if !contains(view, "▶") {
		t.Error("collapsed view should contain right arrow (▶)")
	}
}

// TestView_ErrorState verifies error message is displayed
func TestView_ErrorState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	testErr := errors.New("test error")
	m, _ = m.Update(tui.ErrorMsg{Err: testErr})

	view := m.View()

	if !contains(view, "Error") {
		t.Error("error view should contain 'Error'")
	}
	if !contains(view, "FORECAST") {
		t.Error("error view should still show header")
	}
}

// TestView_LoadingState verifies loading indicator
func TestView_LoadingState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	// Model not yet loaded
	view := m.View()

	if !contains(view, "Loading") {
		t.Error("unloaded view should contain 'Loading'")
	}
}

// TestView_EmptyState verifies "no tasks" message
func TestView_EmptyState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	// Load empty task list
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: []domain.Task{}})

	view := m.View()

	if !contains(view, "No tasks") {
		t.Error("empty view should contain 'No tasks'")
	}
}

// TestRefresh_ReturnsLoadCommand verifies Refresh returns load command
func TestRefresh_ReturnsLoadCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{tasks: []domain.Task{{ID: "1", Name: "Task"}}}
	m := New(styles, keys, svc)

	cmd := m.Refresh()
	if cmd == nil {
		t.Fatal("Refresh() should return a command")
	}

	// Execute the command and verify it loads tasks
	msg := cmd()
	if _, ok := msg.(tui.TasksLoadedMsg); !ok {
		t.Errorf("expected TasksLoadedMsg, got %T", msg)
	}
}

// TestSelectedTask_OutOfBounds verifies bounds checking
func TestSelectedTask_OutOfBounds(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{{ID: "1", Name: "Task", DueDate: &now}}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Set cursor beyond bounds
	m.cursor = 999

	task := m.SelectedTask()
	if task != nil {
		t.Error("SelectedTask() should return nil for out-of-bounds cursor")
	}
}

// TestGroupTasks_SkipsCompletedTasks verifies completed tasks are filtered
func TestGroupTasks_SkipsCompletedTasks(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{
		{ID: "1", Name: "Active task", DueDate: &now, Completed: false},
		{ID: "2", Name: "Completed task", DueDate: &now, Completed: true},
	}

	items := m.groupTasks(tasks)

	// Count non-header items
	taskCount := 0
	for _, item := range items {
		if !item.IsHeader {
			taskCount++
			if item.Task.Completed {
				t.Error("completed task should be filtered out")
			}
		}
	}

	if taskCount != 1 {
		t.Errorf("expected 1 active task, got %d", taskCount)
	}
}

// TestBuildGroupedItems_EmptyGroups verifies empty groups are skipped
func TestBuildGroupedItems_EmptyGroups(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	// Create groups map with only one non-empty group
	groups := map[DueGroup][]domain.Task{
		GroupOverdue:  {},
		GroupToday:    {{ID: "1", Name: "Task"}},
		GroupTomorrow: {},
		GroupThisWeek: {},
		GroupLater:    {},
		GroupNoDue:    {},
	}

	items := m.buildGroupedItems(groups)

	// Should only have 1 header + 1 task = 2 items
	headerCount := 0
	for _, item := range items {
		if item.IsHeader {
			headerCount++
		}
	}

	if headerCount != 1 {
		t.Errorf("expected 1 header (only non-empty group), got %d", headerCount)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items total (1 header + 1 task), got %d", len(items))
	}
}

// TestRenderTask_FlaggedTask verifies flagged icon
func TestRenderTask_FlaggedTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	task := domain.Task{ID: "1", Name: "Flagged task", Flagged: true}
	rendered := m.renderTask(task, GroupToday, false)

	if !contains(rendered, "🚩") {
		t.Error("flagged task should contain flag icon (🚩)")
	}
}

// TestRenderTask_CompletedTask verifies completed checkbox
func TestRenderTask_CompletedTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	task := domain.Task{ID: "1", Name: "Completed task", Completed: true}
	rendered := m.renderTask(task, GroupToday, false)

	if !contains(rendered, "☑") {
		t.Error("completed task should contain checked box (☑)")
	}
}

// TestRenderTask_IncompleteTask verifies unchecked checkbox
func TestRenderTask_IncompleteTask(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	task := domain.Task{ID: "1", Name: "Incomplete task", Completed: false}
	rendered := m.renderTask(task, GroupToday, false)

	if !contains(rendered, "☐") {
		t.Error("incomplete task should contain unchecked box (☐)")
	}
}

// TestNextSelectableIndex_Wrapping verifies cursor wraps around
func TestNextSelectableIndex_Wrapping(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", DueDate: &now},
		{ID: "2", Name: "Task 2", DueDate: &now},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	// Move down from last position should wrap to 0
	lastIndex := len(m.items) - 1
	nextIndex := m.nextSelectableIndex(lastIndex, 1)
	if nextIndex != 0 {
		t.Errorf("moving down from last position should wrap to 0, got %d", nextIndex)
	}

	// Move up from first position should wrap to last
	nextIndex = m.nextSelectableIndex(0, -1)
	if nextIndex != lastIndex {
		t.Errorf("moving up from first position should wrap to last (%d), got %d", lastIndex, nextIndex)
	}
}

// TestRenderHeader_TaskCount verifies task count in header
func TestRenderHeader_TaskCount(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", DueDate: &now},
		{ID: "2", Name: "Task 2", DueDate: &now},
		{ID: "3", Name: "Task 3", DueDate: &now},
	}
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	header := m.renderHeader()

	if !contains(header, "3 tasks") {
		t.Errorf("header should contain '3 tasks', got: %s", header)
	}
}

// TestGroupName verifies all group names
func TestGroupName(t *testing.T) {
	tests := []struct {
		group    DueGroup
		expected string
	}{
		{GroupOverdue, "Overdue"},
		{GroupToday, "Today"},
		{GroupTomorrow, "Tomorrow"},
		{GroupThisWeek, "This Week"},
		{GroupLater, "Later"},
		{GroupNoDue, "No Due Date"},
		{DueGroup(999), "Unknown"}, // Invalid group
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := groupName(tt.group)
			if got != tt.expected {
				t.Errorf("groupName(%v) = %q, want %q", tt.group, got, tt.expected)
			}
		})
	}
}

// TestSetFilter verifies filter is applied to tasks
func TestSetFilter(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{
		{ID: "1", Name: "Buy milk", DueDate: &now},
		{ID: "2", Name: "Review code", DueDate: &now},
	}
	m.allTasks = tasks
	m.items = m.groupTasks(tasks)

	// Apply a search filter using the filter package API
	filterState := filter.State{}.WithSearchText("milk")

	m = m.SetFilter(filterState)

	// Should only have items matching "milk"
	taskCount := 0
	for _, item := range m.items {
		if !item.IsHeader {
			taskCount++
			if !contains(item.Task.Name, "milk") {
				t.Errorf("filtered task should contain 'milk', got %q", item.Task.Name)
			}
		}
	}

	if taskCount != 1 {
		t.Errorf("expected 1 filtered task, got %d", taskCount)
	}

	// Cursor should be repositioned
	if len(m.items) > 1 && m.cursor != 1 {
		t.Errorf("cursor should be at position 1 after filter, got %d", m.cursor)
	}
}

// TestSetFilter_EmptyFilter verifies clearing filter shows all tasks
func TestSetFilter_EmptyFilter(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}
	m := New(styles, keys, svc)

	now := time.Now()
	tasks := []domain.Task{
		{ID: "1", Name: "Task 1", DueDate: &now},
		{ID: "2", Name: "Task 2", DueDate: &now},
	}
	m.allTasks = tasks

	// Apply empty filter
	filterState := filter.State{}
	m = m.SetFilter(filterState)

	// Should have all tasks
	taskCount := 0
	for _, item := range m.items {
		if !item.IsHeader {
			taskCount++
		}
	}

	if taskCount != 2 {
		t.Errorf("expected 2 tasks with empty filter, got %d", taskCount)
	}
}
