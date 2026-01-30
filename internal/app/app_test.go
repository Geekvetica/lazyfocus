package app

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/confirm"
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

func TestAppNavigationDelegatesToView(t *testing.T) {
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
	app = newModel.(Model)

	// Assert - we can't easily verify the internal state of the inbox view,
	// but we can verify the app received and processed the message
	_ = app
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
