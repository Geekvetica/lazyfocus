package projects

import (
	"fmt"
	"strings"
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

func (m *MockService) GetProjects(_ string) ([]domain.Project, error) {
	return m.projects, nil
}

func (m *MockService) GetTasksByProject(_ string) ([]domain.Task, error) {
	return m.tasks, nil
}

// Implement other interface methods as stubs...
func (m *MockService) GetInboxTasks() ([]domain.Task, error)                    { return nil, nil }
func (m *MockService) GetAllTasks(_ service.TaskFilters) ([]domain.Task, error) { return nil, nil }
func (m *MockService) GetTasksByTag(_ string) ([]domain.Task, error)            { return nil, nil }
func (m *MockService) GetFlaggedTasks() ([]domain.Task, error)                  { return nil, nil }
func (m *MockService) GetTaskByID(_ string) (*domain.Task, error)               { return nil, nil }
func (m *MockService) CreateTask(_ domain.TaskInput) (*domain.Task, error)      { return nil, nil }
func (m *MockService) ModifyTask(_ string, _ domain.TaskModification) (*domain.Task, error) {
	return nil, nil
}
func (m *MockService) CompleteTask(_ string) (*domain.OperationResult, error) { return nil, nil }
func (m *MockService) DeleteTask(_ string) (*domain.OperationResult, error)   { return nil, nil }
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

// ========================================
// 1. Init and Loading
// ========================================

func TestInit_ReturnsLoadProjectsCommand(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	cmd := m.Init()

	if cmd == nil {
		t.Error("Init should return a command to load projects")
	}
}

func TestLoadProjects_LoadsAllProjects(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	expectedProjects := []domain.Project{
		{ID: "p1", Name: "Work", TaskCount: 5},
		{ID: "p2", Name: "Personal", TaskCount: 3},
		{ID: "p3", Name: "Hobby", TaskCount: 2},
	}
	svc := &MockService{projects: expectedProjects}

	m := New(styles, keys, svc)
	cmd := m.Init()

	// Execute command
	msg := cmd()

	// Verify message
	projectsMsg, ok := msg.(tui.ProjectsLoadedMsg)
	if !ok {
		t.Fatal("Expected ProjectsLoadedMsg")
	}

	if len(projectsMsg.Projects) != 3 {
		t.Errorf("got %d projects, want 3", len(projectsMsg.Projects))
	}
}

// ========================================
// 2. Update Message Dispatch
// ========================================

func TestUpdate_ProjectsLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	projects := []domain.Project{
		{ID: "p1", Name: "Project 1", TaskCount: 5},
		{ID: "p2", Name: "Project 2", TaskCount: 3},
	}
	svc := &MockService{}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: projects})

	if !m.loaded {
		t.Error("should set loaded=true")
	}
	if m.err != nil {
		t.Errorf("should clear error, got: %v", m.err)
	}
	if len(m.projectList.Projects()) != 2 {
		t.Errorf("projectList should have 2 projects, got %d", len(m.projectList.Projects()))
	}
}

func TestUpdate_TasksForProjectLoadedMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	tasks := []domain.Task{
		{ID: "t1", Name: "Task 1"},
		{ID: "t2", Name: "Task 2"},
	}
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: tasks})

	if m.mode != ModeProjectTasks {
		t.Error("should stay in task mode")
	}
	// Tasks are set internally, verify mode is correct
	if m.taskList.SelectedTask() == nil {
		t.Error("taskList should have tasks available")
	}
}

func TestUpdate_WindowSizeMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	if m.width != 100 {
		t.Errorf("width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("height = %d, want 50", m.height)
	}
}

func TestUpdate_WindowSizeMsg_InTaskMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Switch to task mode
	m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 60})

	if m.width != 120 {
		t.Errorf("width = %d, want 120", m.width)
	}
	if m.height != 60 {
		t.Errorf("height = %d, want 60", m.height)
	}
}

func TestUpdate_ErrorMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	testErr := fmt.Errorf("test error")
	m, _ = m.Update(tui.ErrorMsg{Err: testErr})

	if m.err == nil {
		t.Error("should set error")
	}
	if m.err.Error() != "test error" {
		t.Errorf("error = %q, want %q", m.err.Error(), "test error")
	}
}

// ========================================
// 3. Mode Transitions
// ========================================

func TestModeTransition_ProjectList_To_TaskList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1", TaskCount: 5}},
		tasks:    []domain.Task{{ID: "t1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	if m.mode != ModeProjectList {
		t.Error("should start in project list mode")
	}

	// Press Enter to drill down
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.mode != ModeProjectTasks {
		t.Error("should switch to task mode")
	}
	if m.currentProject == nil {
		t.Error("currentProject should be set")
	}
	if m.currentProject.ID != "p1" {
		t.Errorf("currentProject.ID = %q, want %q", m.currentProject.ID, "p1")
	}
	if cmd == nil {
		t.Error("should return command to load tasks")
	}
}

func TestModeTransition_TaskList_To_ProjectList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
		tasks:    []domain.Task{{ID: "t1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down

	if m.mode != ModeProjectTasks {
		t.Error("should be in task mode")
	}

	// Press h to go back
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if m.mode != ModeProjectList {
		t.Error("should return to project list mode")
	}
	if m.currentProject != nil {
		t.Error("currentProject should be cleared")
	}
}

// ========================================
// 4. Key Handling
// ========================================

func TestHandleKeyPress_ProjectListMode_Navigation(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{
			{ID: "p1", Name: "Project 1"},
			{ID: "p2", Name: "Project 2"},
			{ID: "p3", Name: "Project 3"},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	// Test navigation keys are delegated
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if cmd != nil {
		t.Error("navigation should be delegated to projectList")
	}

	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if cmd != nil {
		t.Error("navigation should be delegated to projectList")
	}
}

func TestHandleKeyPress_ProjectListMode_Enter(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.mode != ModeProjectTasks {
		t.Error("Enter should drill down to tasks")
	}
	if cmd == nil {
		t.Error("should return command to load tasks")
	}
}

func TestHandleKeyPress_ProjectListMode_Enter_NoProjectSelected(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{}, // Empty list
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	if m.mode != ModeProjectList {
		t.Error("should stay in project list mode when no project selected")
	}
	if cmd != nil {
		t.Error("should not return command when no project selected")
	}
}

func TestHandleKeyPress_TaskListMode_Navigation(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
		tasks: []domain.Task{
			{ID: "t1", Name: "Task 1"},
			{ID: "t2", Name: "Task 2"},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	// Test navigation keys are delegated
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if cmd != nil {
		t.Error("navigation should be delegated to taskList")
	}

	_, cmd = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if cmd != nil {
		t.Error("navigation should be delegated to taskList")
	}
}

func TestHandleKeyPress_TaskListMode_BackKey(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down

	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if m.mode != ModeProjectList {
		t.Error("h should return to project list")
	}
}

func TestHandleKeyPress_BackKey_AlreadyInProjectList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	if m.mode != ModeProjectList {
		t.Error("should be in project list mode")
	}

	// Press h when already in project list
	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})

	if m.mode != ModeProjectList {
		t.Error("should stay in project list mode")
	}
	if cmd != nil {
		t.Error("should not return command")
	}
}

func TestHandleKeyPress_EscapeKey_InProjectList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	m, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	if m.mode != ModeProjectList {
		t.Error("should stay in project list mode")
	}
	if cmd != nil {
		t.Error("should not return command when already in project list")
	}
}

// ========================================
// 5. Header Rendering
// ========================================

func TestRenderHeader_ProjectListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{
			{ID: "p1", Name: "Project 1"},
			{ID: "p2", Name: "Project 2"},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	header := m.renderHeader()

	if !strings.Contains(header, "PROJECTS") {
		t.Error("header should contain 'PROJECTS'")
	}
	if !strings.Contains(header, "(2)") {
		t.Error("header should contain project count")
	}
	if strings.Contains(header, "back") {
		t.Error("header should not contain back hint in project list mode")
	}
}

func TestRenderHeader_TaskListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "My Project"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down

	header := m.renderHeader()

	if !strings.Contains(header, "My Project") {
		t.Error("header should contain project name")
	}
	if !strings.Contains(header, "back") {
		t.Error("header should contain back hint in task mode")
	}
}

func TestRenderHeader_TaskListMode_NoCurrentProject(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeProjectTasks
	m.currentProject = nil

	header := m.renderHeader()

	if !strings.Contains(header, "PROJECT TASKS") {
		t.Error("header should contain fallback text when no current project")
	}
}

// ========================================
// 6. View Rendering
// ========================================

func TestView_ProjectListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{
			{ID: "p1", Name: "Work", TaskCount: 5},
			{ID: "p2", Name: "Personal", TaskCount: 3},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	view := m.View()

	if !strings.Contains(view, "PROJECTS") {
		t.Error("view should contain header")
	}
	// The exact content depends on projectlist component
}

func TestView_TaskListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
		tasks: []domain.Task{
			{ID: "t1", Name: "Task 1"},
			{ID: "t2", Name: "Task 2"},
		},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	view := m.View()

	if !strings.Contains(view, "Project 1") {
		t.Error("view should contain project name in header")
	}
	// The exact task content depends on tasklist component
}

func TestView_ErrorState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ErrorMsg{Err: fmt.Errorf("test error message")})

	view := m.View()

	if !strings.Contains(view, "Error") {
		t.Error("view should contain 'Error'")
	}
	if !strings.Contains(view, "test error message") {
		t.Error("view should contain error message")
	}
}

func TestView_LoadingState(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.loaded {
		t.Error("should not be loaded initially")
	}

	view := m.View()

	if !strings.Contains(view, "PROJECTS") {
		t.Error("view should contain header even when loading")
	}
}

func TestView_EmptyProjectList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{projects: []domain.Project{}}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	view := m.View()

	if !strings.Contains(view, "PROJECTS (0)") {
		t.Error("view should show 0 projects")
	}
}

func TestView_EmptyTaskList(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Empty Project"}},
		tasks:    []domain.Task{}, // No tasks
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down
	m, _ = m.Update(tui.TasksLoadedMsg{Tasks: svc.tasks})

	view := m.View()

	if !strings.Contains(view, "Empty Project") {
		t.Error("view should show project name")
	}
	// The empty state message is handled by tasklist component
}

// ========================================
// 7. Refresh Functionality
// ========================================

func TestRefresh_ProjectListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})

	cmd := m.Refresh()

	if cmd == nil {
		t.Error("Refresh should return command")
	}

	msg := cmd()
	if _, ok := msg.(tui.ProjectsLoadedMsg); !ok {
		t.Error("Refresh should return ProjectsLoadedMsg")
	}
}

func TestRefresh_TaskListMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{
		projects: []domain.Project{{ID: "p1", Name: "Project 1"}},
		tasks:    []domain.Task{{ID: "t1", Name: "Task 1"}},
	}

	m := New(styles, keys, svc)
	m, _ = m.Update(tui.ProjectsLoadedMsg{Projects: svc.projects})
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter}) // Drill down

	if m.mode != ModeProjectTasks {
		t.Error("should be in task mode")
	}

	cmd := m.Refresh()

	if cmd == nil {
		t.Error("Refresh should return command")
	}

	msg := cmd()
	if _, ok := msg.(tui.TasksLoadedMsg); !ok {
		t.Error("Refresh in task mode should return TasksLoadedMsg")
	}
}

func TestRefresh_TaskListMode_NoCurrentProject(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.mode = ModeProjectTasks
	m.currentProject = nil

	cmd := m.Refresh()

	if cmd == nil {
		t.Error("Refresh should return command even with no current project")
	}

	// Should fall back to loading projects
	msg := cmd()
	if _, ok := msg.(tui.ProjectsLoadedMsg); !ok {
		t.Error("Refresh without current project should fall back to ProjectsLoadedMsg")
	}
}

// ========================================
// 8. Additional Edge Cases
// ========================================

func TestMode_ReturnsCurrentMode(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)

	if m.Mode() != ModeProjectList {
		t.Errorf("Mode() = %v, want ModeProjectList", m.Mode())
	}

	m.mode = ModeProjectTasks
	if m.Mode() != ModeProjectTasks {
		t.Errorf("Mode() = %v, want ModeProjectTasks", m.Mode())
	}
}

func TestRenderError_WithWidth(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.width = 80
	m, _ = m.Update(tui.ErrorMsg{Err: fmt.Errorf("error with width")})

	view := m.renderError()

	if !strings.Contains(view, "error with width") {
		t.Error("should contain error message")
	}
}

func TestRenderError_WithoutWidth(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m.width = 0 // No width set
	m, _ = m.Update(tui.ErrorMsg{Err: fmt.Errorf("error without width")})

	view := m.renderError()

	if !strings.Contains(view, "error without width") {
		t.Error("should contain error message")
	}
	// Should use default width of 40
}

func TestWindowSizeMsg_NegativeHeight(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	svc := &MockService{}

	m := New(styles, keys, svc)
	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 1})

	// With header height of 2, available height should be clamped to 0
	if m.height != 1 {
		t.Errorf("height = %d, want 1", m.height)
	}
}
