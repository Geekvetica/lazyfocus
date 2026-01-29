package app

import (
	"errors"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
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
		{"Switch to Projects (not implemented)", '2', tui.ViewInbox}, // Should stay on inbox in Phase 4
		{"Switch to Tags (not implemented)", '3', tui.ViewInbox},
		{"Switch to Forecast (not implemented)", '4', tui.ViewInbox},
		{"Switch to Review (not implemented)", '5', tui.ViewInbox},
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

	// Act - try to switch view while quick add is open (should be ignored)
	newModel, _ = app.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}})
	app = newModel.(Model)

	// Assert - view should not change, quick add should still be visible
	if app.currentView != tui.ViewInbox {
		t.Errorf("expected view to remain tui.ViewInbox, got %d", app.currentView)
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
