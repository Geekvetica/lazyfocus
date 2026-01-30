package app

import (
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/searchinput"
	"github.com/pwojciechowski/lazyfocus/internal/tui/filter"
)

// TestFilterIntegration_SearchText tests that search text filters are applied to views
func TestFilterIntegration_SearchText(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Buy groceries", Completed: false},
			{ID: "2", Name: "Write report", Completed: false},
			{ID: "3", Name: "Call grocer", Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Load tasks
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)

	// Set search text
	searchMsg := searchinput.SearchChangedMsg{Text: "grocer"}
	model, _ = app.Update(searchMsg)
	app = model.(Model)

	// Verify filter is applied
	if app.filterState.SearchText != "grocer" {
		t.Errorf("Expected filter state to have search text 'grocer', got '%s'", app.filterState.SearchText)
	}

	// Verify inbox view has filtered tasks
	if app.inboxView.TaskCount() != 2 {
		t.Errorf("Expected 2 filtered tasks (groceries and grocer), got %d", app.inboxView.TaskCount())
	}
}

// TestFilterIntegration_ClearFilter tests that clearing filters works
func TestFilterIntegration_ClearFilter(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Task 1", Completed: false},
			{ID: "2", Name: "Task 2", Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Load tasks
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)

	// Set search text
	model, _ = app.Update(searchinput.SearchChangedMsg{Text: "filter"})
	app = model.(Model)

	// Clear filter
	model, _ = app.Update(searchinput.SearchClearedMsg{})
	app = model.(Model)

	// Verify filter is cleared
	if app.filterState.IsActive() {
		t.Error("Expected filter state to be cleared")
	}

	// Verify all tasks are shown
	if app.inboxView.TaskCount() != 2 {
		t.Errorf("Expected 2 tasks after clearing filter, got %d", app.inboxView.TaskCount())
	}
}

// TestFilterIntegration_DueFilter tests that due date filters work
func TestFilterIntegration_DueFilter(t *testing.T) {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	mockSvc := &service.MockOmniFocusService{}
	mockSvc.AllTasks = []domain.Task{
		{ID: "1", Name: "Due today", DueDate: &today, Completed: false},
		{ID: "2", Name: "Due tomorrow", DueDate: &tomorrow, Completed: false},
		{ID: "3", Name: "No due date", DueDate: nil, Completed: false},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Switch to forecast view and load tasks
	app.currentView = tui.ViewForecast
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.AllTasks})
	app = model.(Model)

	// Apply due filter for today
	app.filterState = app.filterState.WithDueFilter(filter.DueToday)
	app = app.applyFilterToCurrentView()

	// Note: We can't access items directly (unexported), but we can verify through selected task
	// At minimum, the selected task should be the "Due today" task
	selected := app.forecastView.SelectedTask()
	if selected == nil || selected.Name != "Due today" {
		t.Error("Expected selected task to be 'Due today' after filtering")
	}
}

// TestFilterIntegration_FlaggedFilter tests that flagged filter works
func TestFilterIntegration_FlaggedFilter(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		FlaggedTasks: []domain.Task{
			{ID: "1", Name: "Flagged task 1", Flagged: true, Completed: false},
			{ID: "2", Name: "Flagged task 2", Flagged: true, Completed: false},
			{ID: "3", Name: "Not flagged", Flagged: false, Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Switch to review view and load tasks
	app.currentView = tui.ViewReview
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.FlaggedTasks})
	app = model.(Model)

	// Apply flagged filter
	app.filterState = app.filterState.WithFlaggedOnly(true)
	app = app.applyFilterToCurrentView()

	// Verify only flagged tasks are shown
	if app.reviewView.TaskCount() != 2 {
		t.Errorf("Expected 2 flagged tasks, got %d", app.reviewView.TaskCount())
	}
}

// TestFilterIntegration_ProjectFilter tests that project filters work
func TestFilterIntegration_ProjectFilter(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Work task", ProjectID: "proj1", Completed: false},
			{ID: "2", Name: "Personal task", ProjectID: "proj2", Completed: false},
			{ID: "3", Name: "Another work task", ProjectID: "proj1", Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Load tasks
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)

	// Apply project filter
	app.filterState = app.filterState.WithProject("proj1")
	app = app.applyFilterToCurrentView()

	// Verify only work tasks are shown
	if app.inboxView.TaskCount() != 2 {
		t.Errorf("Expected 2 tasks in project 'proj1', got %d", app.inboxView.TaskCount())
	}
}

// TestFilterIntegration_TagFilter tests that tag filters work
func TestFilterIntegration_TagFilter(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Task 1", Tags: []string{"urgent"}, Completed: false},
			{ID: "2", Name: "Task 2", Tags: []string{"low"}, Completed: false},
			{ID: "3", Name: "Task 3", Tags: []string{"urgent", "work"}, Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Load tasks
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)

	// Apply tag filter
	app.filterState = app.filterState.WithTag("urgent")
	app = app.applyFilterToCurrentView()

	// Verify only urgent tasks are shown
	if app.inboxView.TaskCount() != 2 {
		t.Errorf("Expected 2 tasks with tag 'urgent', got %d", app.inboxView.TaskCount())
	}
}

// TestFilterIntegration_MultipleFilters tests combining multiple filters
func TestFilterIntegration_MultipleFilters(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Urgent groceries", Tags: []string{"urgent"}, Flagged: true, Completed: false},
			{ID: "2", Name: "Normal groceries", Tags: []string{"low"}, Flagged: false, Completed: false},
			{ID: "3", Name: "Urgent work", Tags: []string{"urgent"}, Flagged: false, Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Load tasks
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)

	// Apply multiple filters: search + flagged
	app.filterState = app.filterState.WithSearchText("groceries").WithFlaggedOnly(true)
	app = app.applyFilterToCurrentView()

	// Verify only "Urgent groceries" is shown
	if app.inboxView.TaskCount() != 1 {
		t.Errorf("Expected 1 task matching search 'groceries' and flagged, got %d", app.inboxView.TaskCount())
	}
}

// TestFilterIntegration_FilterPersistsAcrossRefresh tests that filters persist after refresh
func TestFilterIntegration_FilterPersistsAcrossRefresh(t *testing.T) {
	mockSvc := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "1", Name: "Task 1", Completed: false},
			{ID: "2", Name: "Special task", Completed: false},
		},
	}

	app := NewApp(mockSvc)
	app.width = 80
	app.height = 24
	app.ready = true

	// Load tasks and apply filter
	model, _ := app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)
	app.filterState = app.filterState.WithSearchText("Special")
	app = app.applyFilterToCurrentView()

	// Verify filter is applied
	if app.inboxView.TaskCount() != 1 {
		t.Errorf("Expected 1 task after filtering, got %d", app.inboxView.TaskCount())
	}

	// Simulate refresh (reload tasks)
	model, _ = app.Update(tui.TasksLoadedMsg{Tasks: mockSvc.InboxTasks})
	app = model.(Model)

	// Verify filter is still applied after refresh
	if app.inboxView.TaskCount() != 1 {
		t.Errorf("Expected filter to persist after refresh, got %d tasks", app.inboxView.TaskCount())
	}
}
