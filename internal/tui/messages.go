package tui

import "github.com/pwojciechowski/lazyfocus/internal/domain"

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
