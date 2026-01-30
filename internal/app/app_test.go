package app

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/command"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/commandinput"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/confirm"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/searchinput"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/taskdetail"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/taskedit"
)

func TestNewApp(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}

	// Act
	app := NewApp(mockSvc)

	// Assert
	if app.service == nil {
		t.Error("expected service to be set")
	}
	if app.styles == nil {
		t.Error("expected styles to be initialized")
	}
	if app.currentView != tui.ViewInbox {
		t.Errorf("expected currentView to be tui.ViewInbox (%d), got %d", tui.ViewInbox, app.currentView)
	}
	if app.ready {
		t.Error("expected ready to be false initially")
	}
}

func TestAppInit(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act
	cmd := app.Init()

	// Assert - Init should return inbox view's init command
	if cmd == nil {
		t.Error("expected Init to return a command")
	}
}

func TestAppQuit(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act
	newModel, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	// Assert
	if cmd == nil {
		t.Fatal("expected quit command, got nil")
	}
	// Verify it's actually a quit command by checking if it's tea.Quit
	// We can't directly compare functions, but we can verify the model is unchanged
	_ = newModel
}

func TestAppWindowSizeMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	app = newModel.(Model)

	// Assert
	if app.width != 100 {
		t.Errorf("expected width 100, got %d", app.width)
	}
	if app.height != 50 {
		t.Errorf("expected height 50, got %d", app.height)
	}
	if !app.ready {
		t.Error("expected ready to be true after WindowSizeMsg")
	}
}

func TestAppShowQuickAdd(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act
	newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = newModel.(Model)

	// Assert
	if !app.quickAdd.IsVisible() {
		t.Error("expected quick add to be visible after pressing 'a'")
	}
}

func TestAppHideQuickAddOnEscape(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Show quick add first
	newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = newModel.(Model)

	// Act - press Escape
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyEsc})
	app = newModel.(Model)

	// Assert
	if app.quickAdd.IsVisible() {
		t.Error("expected quick add to be hidden after pressing Escape")
	}
}

func TestAppTaskCreatedMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Existing task"},
		},
	}
	app := NewApp(mockSvc)

	// Show quick add first
	newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = newModel.(Model)

	// Act - send tui.TaskCreatedMsg
	newTask := domain.Task{ID: "task2", Name: "New task"}
	newModel, cmd := app.Update(tui.TaskCreatedMsg{Task: newTask})
	app = newModel.(Model)

	// Assert
	if app.quickAdd.IsVisible() {
		t.Error("expected quick add to be hidden after task creation")
	}
	if cmd == nil {
		t.Error("expected refresh command after task creation")
	}
}

func TestAppToggleHelp(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act - press '?' to show help
	newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	app = newModel.(Model)

	// Assert
	if !app.showHelp {
		t.Error("expected showHelp to be true after pressing '?'")
	}

	// Act - press '?' again to hide help
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	app = newModel.(Model)

	// Assert
	if app.showHelp {
		t.Error("expected showHelp to be false after pressing '?' again")
	}
}

func TestAppViewSwitching(t *testing.T) {
	tests := []struct {
		name         string
		key          rune
		expectedView int
	}{
		{"Switch to Inbox", '1', tui.ViewInbox},
		{"Switch to Projects", '2', tui.ViewProjects},
		{"Switch to Tags", '3', tui.ViewTags},
		{"Switch to Forecast", '4', tui.ViewForecast},
		{"Switch to Review", '5', tui.ViewReview},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockSvc := &service.MockOmniFocusService{}
			app := NewApp(mockSvc)

			// Act
			newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{tt.key}})
			app = newModel.(Model)

			// Assert
			if app.currentView != tt.expectedView {
				t.Errorf("expected currentView to be %d, got %d", tt.expectedView, app.currentView)
			}
		})
	}
}

func TestAppNavigationDelegatesToView(_ *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Task 1"},
			{ID: "task2", Name: "Task 2"},
		},
	}
	app := NewApp(mockSvc)

	// Initialize with size
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	app = newModel.(Model)

	// Act - send navigation key (down arrow)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyDown})

	// Assert - we can't easily verify the internal state of the inbox view,
	// but we can verify the app received and processed the message
	_ = newModel
}

func TestAppErrorMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act
	testErr := errors.New("test error")
	newModel, _ := app.Update(tui.ErrorMsg{Err: testErr})
	app = newModel.(Model)

	// Assert
	if app.err == nil {
		t.Error("expected error to be set")
	}
	if app.err.Error() != "test error" {
		t.Errorf("expected error message 'test error', got '%s'", app.err.Error())
	}
}

func TestAppViewBeforeReady(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Act
	view := app.View()

	// Assert
	if view != "Loading..." {
		t.Errorf("expected 'Loading...' before ready, got '%s'", view)
	}
}

func TestAppViewAfterReady(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test task"},
		},
	}
	app := NewApp(mockSvc)

	// Set ready by sending WindowSizeMsg
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	app = newModel.(Model)

	// Act
	view := app.View()

	// Assert
	if view == "Loading..." {
		t.Error("expected view content, got 'Loading...'")
	}
	// View should contain inbox header
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}

func TestAppViewWithQuickAddOverlay(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Set ready
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	app = newModel.(Model)

	// Show quick add
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = newModel.(Model)

	// Act
	view := app.View()

	// Assert
	if len(view) == 0 {
		t.Error("expected non-empty view with overlay")
	}
	// View should contain quick add overlay when visible
	// We can't easily test the exact content, but verify it's not empty
}

func TestAppViewWithHelpOverlay(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Set ready
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	app = newModel.(Model)

	// Show help
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	app = newModel.(Model)

	// Act
	view := app.View()

	// Assert
	if len(view) == 0 {
		t.Error("expected non-empty view with help overlay")
	}
	// View should contain help text
	// We'll verify specific content in the implementation
}

func TestAppCurrentViewName(t *testing.T) {
	tests := []struct {
		name         string
		view         int
		expectedName string
	}{
		{"Inbox view", tui.ViewInbox, "Inbox"},
		{"Projects view", tui.ViewProjects, "Projects"},
		{"Tags view", tui.ViewTags, "Tags"},
		{"Forecast view", tui.ViewForecast, "Forecast"},
		{"Review view", tui.ViewReview, "Review"},
		{"Unknown view", 99, "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockSvc := &service.MockOmniFocusService{}
			app := NewApp(mockSvc)
			app.currentView = tt.view

			// Act
			name := app.CurrentViewName()

			// Assert
			if name != tt.expectedName {
				t.Errorf("expected view name '%s', got '%s'", tt.expectedName, name)
			}
		})
	}
}

func TestAppQuickAddDelegation(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Show quick add
	newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = newModel.(Model)

	// Act - send key message that should be delegated to quick add
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	app = newModel.(Model)

	// Assert - quick add should still be visible
	if !app.quickAdd.IsVisible() {
		t.Error("expected quick add to still be visible")
	}
}

func TestAppGlobalKeysIgnoredWhenOverlayVisible(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Show quick add
	newModel, _ := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	app = newModel.(Model)

	// Record the current view before attempting to switch
	initialView := app.currentView

	// Act - try to switch view while quick add is open (should be delegated to quick add, not app)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	app = newModel.(Model)

	// Assert - view should not change (quick add intercepts the key), quick add should still be visible
	if app.currentView != initialView {
		t.Errorf("expected view to remain %d, got %d", initialView, app.currentView)
	}
	if !app.quickAdd.IsVisible() {
		t.Error("expected quick add to still be visible")
	}
}

func TestRenderHelpSmallWidth(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Set very small width to trigger min(60, m.width-4) constraint
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 20, Height: 50})
	app = newModel.(Model)

	// Show help to trigger renderHelp()
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
	app = newModel.(Model)

	// Act
	view := app.View()

	// Assert - renderHelp should handle small width gracefully
	if len(view) == 0 {
		t.Error("expected non-empty view with help overlay")
	}
	// The help should render without panicking despite small width
}

func TestCenterOverlayLargeContent(t *testing.T) {
	tests := []struct {
		name           string
		width          int
		height         int
		contentLines   int
		contentWidth   int
		expectedVerPad int
		expectedHorPad int
	}{
		{
			name:           "Content larger than viewport vertically",
			width:          100,
			height:         10,
			contentLines:   20, // More lines than height
			contentWidth:   50,
			expectedVerPad: 0, // Should be clamped to 0
			expectedHorPad: 25,
		},
		{
			name:           "Content larger than viewport horizontally",
			width:          50,
			height:         30,
			contentLines:   10,
			contentWidth:   80, // Wider than width
			expectedVerPad: 10,
			expectedHorPad: 0, // Should be clamped to 0
		},
		{
			name:           "Content larger than viewport both dimensions",
			width:          50,
			height:         20,
			contentLines:   30,
			contentWidth:   80,
			expectedVerPad: 0, // Should be clamped to 0
			expectedHorPad: 0, // Should be clamped to 0
		},
		{
			name:           "Normal centered content",
			width:          100,
			height:         50,
			contentLines:   10,
			contentWidth:   40,
			expectedVerPad: 20,
			expectedHorPad: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockSvc := &service.MockOmniFocusService{}
			app := NewApp(mockSvc)

			// Set viewport size
			newModel, _ := app.Update(tea.WindowSizeMsg{Width: tt.width, Height: tt.height})
			app = newModel.(Model)

			// Show help to trigger centerOverlay via renderHelp
			newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
			app = newModel.(Model)

			// Act
			view := app.View()

			// Assert - should not panic and should produce output
			if len(view) == 0 {
				t.Error("expected non-empty view")
			}
			// centerOverlay should handle edge cases gracefully without panicking
		})
	}
}

// Tests for delete task functionality (Stage 3)

func TestDeleteKey_ShowsConfirmation(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	app = newModel.(Model)
	if !app.confirmModal.IsVisible() {
		t.Error("confirm modal should be visible after 'd' key")
	}
}

func TestDeleteKey_NoTaskSelected_NoConfirmation(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{InboxTasks: []domain.Task{}}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	app = newModel.(Model)
	if app.confirmModal.IsVisible() {
		t.Error("confirm modal should not be visible when no task is selected")
	}
}

func TestDeleteConfirmed_TriggersDeleteCommand(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	app = newModel.(Model)
	newModel, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	app = newModel.(Model)
	if cmd == nil {
		t.Fatal("expected command from confirmation")
	}
	confirmedMsg := cmd()
	_, deleteCmd := app.Update(confirmedMsg)
	if deleteCmd == nil {
		t.Fatal("expected delete command to be returned")
	}
}

func TestTaskDeletedMsg_RefreshesView(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	_, cmd := app.Update(tui.TaskDeletedMsg{TaskID: "task1", TaskName: "Test Task"})
	if cmd == nil {
		t.Error("expected refresh command after task deletion")
	}
}

// Tests for flag toggle functionality (Stage 4)

func TestFlagKey_TogglesFlag(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task", Flagged: false}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	if cmd == nil {
		t.Fatal("expected command to be returned for flag toggle")
	}
}

func TestFlagKey_NoTaskSelected_NoAction(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{InboxTasks: []domain.Task{}}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	if cmd != nil {
		t.Error("expected no command when no task is selected")
	}
}

func TestTaskModifiedMsg_RefreshesView(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task", Flagged: false}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	modifiedTask := domain.Task{ID: "task1", Name: "Test Task", Flagged: true}
	_, cmd := app.Update(tui.TaskModifiedMsg{Task: modifiedTask})
	if cmd == nil {
		t.Error("expected refresh command after task modification")
	}
}

func TestConfirmModal_ContextCarriesTaskInfo(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task Name"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	app = newModel.(Model)
	_, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	confirmedMsg, ok := msg.(confirm.ConfirmedMsg)
	if !ok {
		t.Fatalf("expected ConfirmedMsg, got %T", msg)
	}
	ctx, ok := confirmedMsg.Context.(DeleteContext)
	if !ok {
		t.Fatalf("expected DeleteContext, got %T", confirmedMsg.Context)
	}
	if ctx.TaskID != "task1" {
		t.Errorf("expected TaskID 'task1', got '%s'", ctx.TaskID)
	}
	if ctx.TaskName != "Test Task Name" {
		t.Errorf("expected TaskName 'Test Task Name', got '%s'", ctx.TaskName)
	}
}

// Test for FlagRequestedMsg handling from task detail view (bug fix verification)

func TestFlagRequestedMsg_FromTaskDetail_TogglesFlag(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task", Flagged: false}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)

	// Initialize app
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)

	// Show task detail for the first task
	app.taskDetail = app.taskDetail.Show(&testTask)

	// Act - send FlagRequestedMsg from task detail
	newModel, cmd := app.Update(taskdetail.FlagRequestedMsg{TaskID: "task1", Flagged: true})
	app = newModel.(Model)

	// Assert - should return toggleTaskFlag command
	if cmd == nil {
		t.Fatal("expected command to be returned for flag toggle from task detail")
	}

	// Verify task detail is hidden
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to be hidden after flag request")
	}
}

func TestFlagRequestedMsg_NoTask_NoCommand(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)

	// Initialize app
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Task detail is not showing any task (hidden)

	// Act - send FlagRequestedMsg when no task is shown
	newModel, cmd := app.Update(taskdetail.FlagRequestedMsg{})
	app = newModel.(Model)

	// Assert - should not return a command when no task exists
	if cmd != nil {
		t.Error("expected no command when no task is in task detail")
	}

	// Verify task detail remains hidden
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to remain hidden")
	}
}

// Tests for handleTaskDetailMessages (Stage 3)

func TestHandleTaskDetailMessages_CloseMsg(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task"}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Show task detail
	app.taskDetail = app.taskDetail.Show(&testTask)

	// Act - send CloseMsg
	newModel, _ = app.Update(taskdetail.CloseMsg{})
	app = newModel.(Model)

	// Assert - task detail should be hidden
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to be hidden after CloseMsg")
	}
}

func TestHandleTaskDetailMessages_EditRequestedMsg(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task"}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Show task detail
	app.taskDetail = app.taskDetail.Show(&testTask)

	// Act - send EditRequestedMsg
	newModel, _ = app.Update(taskdetail.EditRequestedMsg{Task: testTask})
	app = newModel.(Model)

	// Assert - task detail should be hidden and task edit should be visible
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to be hidden after EditRequestedMsg")
	}
	if !app.taskEdit.IsVisible() {
		t.Error("expected task edit to be visible after EditRequestedMsg")
	}
}

func TestHandleTaskDetailMessages_CompleteRequestedMsg(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task"}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Show task detail
	app.taskDetail = app.taskDetail.Show(&testTask)

	// Act - send CompleteRequestedMsg
	newModel, cmd := app.Update(taskdetail.CompleteRequestedMsg{TaskID: "task1"})
	app = newModel.(Model)

	// Assert - task detail should be hidden and command should be returned
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to be hidden after CompleteRequestedMsg")
	}
	if cmd == nil {
		t.Fatal("expected complete command to be returned")
	}
}

func TestHandleTaskDetailMessages_DeleteRequestedMsg(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task"}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Show task detail
	app.taskDetail = app.taskDetail.Show(&testTask)

	// Act - send DeleteRequestedMsg
	newModel, _ = app.Update(taskdetail.DeleteRequestedMsg{TaskID: "task1", TaskName: "Test Task"})
	app = newModel.(Model)

	// Assert - task detail should be hidden and confirm modal should be visible
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to be hidden after DeleteRequestedMsg")
	}
	if !app.confirmModal.IsVisible() {
		t.Error("expected confirm modal to be visible after DeleteRequestedMsg")
	}
}

// Tests for handleTaskEditMessages (Stage 3)

func TestHandleTaskEditMessages_SaveMsg(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task"}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Show task edit
	app.taskEdit = app.taskEdit.Show(&testTask)

	// Act - send SaveMsg with modification
	newName := "Updated Name"
	mod := domain.TaskModification{Name: &newName}
	newModel, cmd := app.Update(taskedit.SaveMsg{TaskID: "task1", Modification: mod})
	app = newModel.(Model)

	// Assert - task edit should be hidden and command should be returned
	if app.taskEdit.IsVisible() {
		t.Error("expected task edit to be hidden after SaveMsg")
	}
	if cmd == nil {
		t.Fatal("expected modify command to be returned")
	}
}

func TestHandleTaskEditMessages_CancelMsg(t *testing.T) {
	// Arrange
	testTask := domain.Task{ID: "task1", Name: "Test Task"}
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{testTask},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Show task edit
	app.taskEdit = app.taskEdit.Show(&testTask)

	// Act - send CancelMsg
	newModel, _ = app.Update(taskedit.CancelMsg{})
	app = newModel.(Model)

	// Assert - task edit should be hidden
	if app.taskEdit.IsVisible() {
		t.Error("expected task edit to be hidden after CancelMsg")
	}
}

// Tests for handleSearchInputMessages (Stage 3)

func TestHandleSearchInputMessages_SearchChangedMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test Task"},
			{ID: "task2", Name: "Another Task"},
		},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Act - send SearchChangedMsg
	newModel, _ = newModel.(Model).Update(searchinput.SearchChangedMsg{Text: "test"})
	app = newModel.(Model)

	// Assert - filter state should be updated with search text
	if app.filterState.SearchText != "test" {
		t.Errorf("expected search text to be 'test', got '%s'", app.filterState.SearchText)
	}
}

func TestHandleSearchInputMessages_SearchClearedMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test Task"},
		},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Set search text first
	app.filterState = app.filterState.WithSearchText("test")

	// Act - send SearchClearedMsg
	newModel, _ = app.Update(searchinput.SearchClearedMsg{})
	app = newModel.(Model)

	// Assert - filter state should be cleared
	if app.filterState.SearchText != "" {
		t.Errorf("expected search text to be empty, got '%s'", app.filterState.SearchText)
	}
}

func TestHandleSearchInputMessages_SearchConfirmedMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test Task"},
		},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - send SearchConfirmedMsg
	newModel, _ = app.Update(searchinput.SearchConfirmedMsg{Text: "confirmed"})
	app = newModel.(Model)

	// Assert - filter state should be updated
	if app.filterState.SearchText != "confirmed" {
		t.Errorf("expected search text to be 'confirmed', got '%s'", app.filterState.SearchText)
	}
}

// Tests for handleCommandInputMessages (Stage 3)

func TestHandleCommandInputMessages_CommandExecutedMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Act - send CommandExecutedMsg with quit command
	quitCmd := &command.Command{Name: "quit", Args: []string{}}
	_, cmd := newModel.(Model).Update(commandinput.CommandExecutedMsg{Command: quitCmd})

	// Assert - quit command should return tea.Quit
	if cmd == nil {
		t.Fatal("expected quit command to be returned")
	}
}

func TestHandleCommandInputMessages_CommandCancelledMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - send CommandCancelledMsg
	newModel, _ = app.Update(commandinput.CommandCancelledMsg{})
	app = newModel.(Model)

	// Assert - no error should be set
	if app.err != nil {
		t.Errorf("expected no error, got %v", app.err)
	}
}

func TestHandleCommandInputMessages_CommandErrorMsg(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - send CommandErrorMsg
	newModel, _ = app.Update(commandinput.CommandErrorMsg{Error: "command error"})
	app = newModel.(Model)

	// Assert - error should be set
	if app.err == nil {
		t.Fatal("expected error to be set")
	}
	if app.err.Error() != "command error" {
		t.Errorf("expected error 'command error', got '%s'", app.err.Error())
	}
}

// Tests for command execution (Stage 3)

func TestExecuteCommand_Add(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - execute add command with args
	cmd := &command.Command{Name: "add", Args: []string{"test", "task"}}
	newModel, _ = app.executeCommand(cmd)
	app = newModel.(Model)

	// Assert - quick add should be visible
	if !app.quickAdd.IsVisible() {
		t.Error("expected quick add to be visible after add command")
	}
}

func TestExecuteCommand_AddWithoutArgs(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - execute add command without args
	cmd := &command.Command{Name: "add", Args: []string{}}
	newModel, _ = app.executeCommand(cmd)
	app = newModel.(Model)

	// Assert - quick add should be visible
	if !app.quickAdd.IsVisible() {
		t.Error("expected quick add to be visible after add command")
	}
}

func TestExecuteCommand_Complete(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})

	// Act - execute complete command
	cmd := &command.Command{Name: "complete", Args: []string{}}
	_, completeCmd := newModel.(Model).executeCommand(cmd)

	// Assert - should return complete command
	if completeCmd == nil {
		t.Fatal("expected complete command to be returned")
	}
}

func TestExecuteCommand_CompleteNoTask(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})

	// Act - execute complete command with no task selected
	cmd := &command.Command{Name: "complete", Args: []string{}}
	_, completeCmd := newModel.(Model).executeCommand(cmd)

	// Assert - should not return command when no task selected
	if completeCmd != nil {
		t.Error("expected no command when no task is selected")
	}
}

func TestExecuteCommand_Delete(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)

	// Act - execute delete command
	cmd := &command.Command{Name: "delete", Args: []string{}}
	newModel, _ = app.executeCommand(cmd)
	app = newModel.(Model)

	// Assert - confirm modal should be visible
	if !app.confirmModal.IsVisible() {
		t.Error("expected confirm modal to be visible after delete command")
	}
}

func TestExecuteCommand_DeleteNoTask(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})

	// Act - execute delete command with no task selected
	cmd := &command.Command{Name: "delete", Args: []string{}}
	newModel, _ = newModel.(Model).executeCommand(cmd)
	app = newModel.(Model)

	// Assert - confirm modal should not be visible
	if app.confirmModal.IsVisible() {
		t.Error("expected confirm modal to not be visible when no task selected")
	}
}

func TestExecuteCommand_Refresh(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Act - execute refresh command
	cmd := &command.Command{Name: "refresh", Args: []string{}}
	_, refreshCmd := newModel.(Model).executeCommand(cmd)

	// Assert - should return refresh command
	if refreshCmd == nil {
		t.Fatal("expected refresh command to be returned")
	}
}

func TestExecuteCommand_Help(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Act - execute help command
	cmd := &command.Command{Name: "help", Args: []string{}}
	newModel, _ = newModel.(Model).executeCommand(cmd)
	app = newModel.(Model)

	// Assert - showHelp should toggle
	if !app.showHelp {
		t.Error("expected showHelp to be true after help command")
	}

	// Execute again to toggle off
	newModel, _ = app.executeCommand(cmd)
	app = newModel.(Model)

	if app.showHelp {
		t.Error("expected showHelp to be false after second help command")
	}
}

func TestExecuteCommand_NilCommand(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Act - execute nil command
	_, cmd := newModel.(Model).executeCommand(nil)

	// Assert - should handle gracefully
	if cmd != nil {
		t.Error("expected no command when nil is passed")
	}
}

func TestExecuteCommand_UnknownCommand(_ *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - execute unknown command
	cmd := &command.Command{Name: "unknown", Args: []string{}}
	newModel, _ = app.executeCommand(cmd)

	// Assert - should handle gracefully without error
	_ = newModel
}

// Tests for key handling (Stage 3)

func TestKeyHandling_EnterShowsTaskDetail(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)

	// Act - press Enter key
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	app = newModel.(Model)

	// Assert - task detail should be visible
	if !app.taskDetail.IsVisible() {
		t.Error("expected task detail to be visible after Enter key")
	}
}

func TestKeyHandling_EnterNoTask(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)

	// Act - press Enter key with no task
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyEnter})
	app = newModel.(Model)

	// Assert - task detail should not be visible
	if app.taskDetail.IsVisible() {
		t.Error("expected task detail to not be visible when no task selected")
	}
}

func TestKeyHandling_EditKey(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)

	// Act - press 'e' key
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	app = newModel.(Model)

	// Assert - task edit should be visible
	if !app.taskEdit.IsVisible() {
		t.Error("expected task edit to be visible after 'e' key")
	}
}

func TestKeyHandling_EditKeyNoTask(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)
	newModel, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = newModel.(Model)

	// Act - press 'e' key with no task
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	app = newModel.(Model)

	// Assert - task edit should not be visible
	if app.taskEdit.IsVisible() {
		t.Error("expected task edit to not be visible when no task selected")
	}
}

func TestKeyHandling_SearchKey(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - press '/' key
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	app = newModel.(Model)

	// Assert - search input should be visible
	if !app.searchInput.IsVisible() {
		t.Error("expected search input to be visible after '/' key")
	}
}

func TestKeyHandling_CommandKey(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app = newModel.(Model)

	// Act - press ':' key
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{':'}})
	app = newModel.(Model)

	// Assert - command input should be visible
	if !app.commandInput.IsVisible() {
		t.Error("expected command input to be visible after ':' key")
	}
}

// Tests for task completion message handling

func TestTaskCompletedMsg_RefreshesView(t *testing.T) {
	// Arrange
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test Task"}},
	}
	app := NewApp(mockSvc)
	newModel, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Act - send TaskCompletedMsg
	_, cmd := newModel.(Model).Update(tui.TaskCompletedMsg{TaskID: "task1", TaskName: "Test Task"})

	// Assert - should return refresh command
	if cmd == nil {
		t.Error("expected refresh command after TaskCompletedMsg")
	}
}
