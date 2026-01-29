// Package output provides formatting functionality for LazyFocus CLI output.
// It supports both human-readable and JSON output formats.
package output

import "github.com/pwojciechowski/lazyfocus/internal/domain"

// Exit codes used by LazyFocus CLI
const (
	ExitSuccess             = 0 // Successful execution
	ExitGeneralError        = 1 // General error
	ExitOmniFocusNotRunning = 2 // OmniFocus is not running
	ExitItemNotFound        = 3 // Requested item not found
)

// Formatter defines the interface for formatting LazyFocus output
type Formatter interface {
	// FormatTasks formats a list of tasks with the given options
	FormatTasks(tasks []domain.Task, options TaskFormatOptions) string

	// FormatProjects formats a list of projects with the given options
	FormatProjects(projects []domain.Project, options ProjectFormatOptions) string

	// FormatTags formats a list of tags with the given options
	FormatTags(tags []domain.Tag, options TagFormatOptions) string

	// FormatTask formats a single task
	FormatTask(task domain.Task) string

	// FormatProject formats a single project
	FormatProject(project domain.Project) string

	// FormatTag formats a single tag
	FormatTag(tag domain.Tag) string

	// FormatError formats an error message
	FormatError(err error) string

	// FormatCreatedTask formats a newly created task
	FormatCreatedTask(task domain.Task) string

	// FormatModifiedTask formats a modified task
	FormatModifiedTask(task domain.Task) string

	// FormatCompletedTask formats a completed task operation result
	FormatCompletedTask(result domain.OperationResult) string

	// FormatDeletedTask formats a deleted task operation result
	FormatDeletedTask(result domain.OperationResult) string
}

// TaskFormatOptions contains options for formatting tasks
type TaskFormatOptions struct {
	ShowCompleted bool // Include completed tasks in output
	ShowProject   bool // Show project name for each task
	ShowTags      bool // Show tags for each task
}

// ProjectFormatOptions contains options for formatting projects
type ProjectFormatOptions struct {
	ShowTasks bool // Include tasks in project output
	ShowNotes bool // Show project notes
}

// TagFormatOptions contains options for formatting tags
type TagFormatOptions struct {
	ShowCounts bool // Show task counts for each tag
	Flat       bool // Show tags in flat list (no hierarchy)
}
