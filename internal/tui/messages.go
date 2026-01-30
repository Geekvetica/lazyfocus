package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// View constants define the different views in the TUI
const (
	ViewInbox    = 1
	ViewProjects = 2
	ViewTags     = 3
	ViewForecast = 4
	ViewReview   = 5
)

// Data Loading Messages

// TasksLoadedMsg is sent when tasks are loaded asynchronously
type TasksLoadedMsg struct {
	Tasks []domain.Task
}

// ProjectsLoadedMsg is sent when projects are loaded asynchronously
type ProjectsLoadedMsg struct {
	Projects []domain.Project
}

// TagsLoadedMsg is sent when tags are loaded asynchronously
type TagsLoadedMsg struct {
	Tags []domain.Tag
}

// Task Action Messages

// TaskCreatedMsg is sent when a new task is created
type TaskCreatedMsg struct {
	Task domain.Task
}

// TaskCompletedMsg is sent when a task is marked as completed
type TaskCompletedMsg struct {
	TaskID   string
	TaskName string
}

// TaskDeletedMsg is sent when a task is deleted
type TaskDeletedMsg struct {
	TaskID   string
	TaskName string
}

// TaskModifiedMsg is sent when a task is modified
type TaskModifiedMsg struct {
	Task domain.Task
}

// UI Messages

// ErrorMsg is sent when an error occurs during an operation
type ErrorMsg struct {
	Err error
}

// ViewChangedMsg is sent when the user switches to a different view
type ViewChangedMsg struct {
	View int
}

// ClearErrorMsg is sent to clear any displayed error
type ClearErrorMsg struct{}

// Drill-down Navigation Messages

// ProjectSelectedMsg is sent when a project is selected for drill-down
type ProjectSelectedMsg struct {
	ProjectID   string
	ProjectName string
}

// TagSelectedMsg is sent when a tag is selected for drill-down
type TagSelectedMsg struct {
	TagID   string
	TagName string
}

// DrillBackMsg is sent when navigating back from a drill-down view
type DrillBackMsg struct{}

// Task Action Messages (additional)

// TaskFlagToggledMsg is sent when a task's flag status is toggled
type TaskFlagToggledMsg struct {
	TaskID   string
	TaskName string
	Flagged  bool
}

// ShowTaskDetailMsg is sent to show task detail view
type ShowTaskDetailMsg struct {
	Task domain.Task
}

// ShowEditOverlayMsg is sent to show the edit task overlay
type ShowEditOverlayMsg struct {
	Task domain.Task
}

// ShowDeleteConfirmMsg is sent to show delete confirmation
type ShowDeleteConfirmMsg struct {
	TaskID   string
	TaskName string
}

// Search/Filter Messages

// SearchChangedMsg is sent when search text changes
type SearchChangedMsg struct {
	Text string
}

// SearchClearedMsg is sent when search is cleared
type SearchClearedMsg struct{}

// Command Mode Messages

// CommandExecutedMsg is sent when a command is executed
type CommandExecutedMsg struct {
	Command string
	Args    []string
}

// Error State Messages

// ErrorOccurredMsg indicates an error occurred
type ErrorOccurredMsg struct {
	Error     error
	Retryable bool
	RetryCmd  tea.Cmd
}

// ErrorDismissedMsg indicates the user dismissed an error.
//
// Deprecated: Prefer using errorstate.ErrorDismissedMsg emitted by the
// error state component. This type exists for backward compatibility.
type ErrorDismissedMsg struct{}

// RetryRequestedMsg indicates the user requested a retry.
//
// Deprecated: The errorstate component handles retries via RetryCmd.
// This type exists for backward compatibility.
type RetryRequestedMsg struct{}
