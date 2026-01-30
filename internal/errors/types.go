package errors

import "fmt"

// OmniFocusError represents errors related to OmniFocus communication.
// This includes OmniFocus not running, connection issues, or script execution failures.
type OmniFocusError struct {
	message    string
	suggestion string
}

// NewOmniFocusError creates a new OmniFocusError with the given message and suggestion.
func NewOmniFocusError(message, suggestion string) *OmniFocusError {
	return &OmniFocusError{
		message:    message,
		suggestion: suggestion,
	}
}

func (e *OmniFocusError) Error() string {
	return e.message
}

func (e *OmniFocusError) ExitCode() int {
	return ExitOmniFocusError
}

func (e *OmniFocusError) Suggestion() string {
	return e.suggestion
}

// ItemNotFoundError represents errors when a task, project, or tag is not found.
type ItemNotFoundError struct {
	itemType   string
	itemID     string
	suggestion string
}

// NewItemNotFoundError creates a new ItemNotFoundError with the given item type and ID.
func NewItemNotFoundError(itemType, itemID string) *ItemNotFoundError {
	var suggestion string
	switch itemType {
	case "task":
		suggestion = "Verify task ID using 'lazyfocus tasks'"
	case "project":
		suggestion = "Check project name with 'lazyfocus projects'"
	case "tag":
		suggestion = "Check tag name with 'lazyfocus tags'"
	default:
		suggestion = "Verify the item exists in OmniFocus"
	}

	return &ItemNotFoundError{
		itemType:   itemType,
		itemID:     itemID,
		suggestion: suggestion,
	}
}

func (e *ItemNotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.itemType, e.itemID)
}

func (e *ItemNotFoundError) ExitCode() int {
	return ExitItemNotFound
}

func (e *ItemNotFoundError) Suggestion() string {
	return e.suggestion
}

// ValidationError represents input validation failures.
type ValidationError struct {
	message    string
	suggestion string
}

// NewValidationError creates a new ValidationError with the given message and suggestion.
func NewValidationError(message, suggestion string) *ValidationError {
	return &ValidationError{
		message:    message,
		suggestion: suggestion,
	}
}

func (e *ValidationError) Error() string {
	return e.message
}

func (e *ValidationError) ExitCode() int {
	return ExitValidationError
}

func (e *ValidationError) Suggestion() string {
	return e.suggestion
}

// DateParseError represents date parsing failures.
type DateParseError struct {
	dateStr    string
	reason     string
	suggestion string
}

// NewDateParseError creates a new DateParseError with the given date string and reason.
func NewDateParseError(dateStr, reason string) *DateParseError {
	return &DateParseError{
		dateStr:    dateStr,
		reason:     reason,
		suggestion: "Use relative (tomorrow), next (next monday), in (in 3 days), or ISO format",
	}
}

func (e *DateParseError) Error() string {
	return fmt.Sprintf("invalid due date: %s: %s", e.reason, e.dateStr)
}

func (e *DateParseError) ExitCode() int {
	return ExitValidationError
}

func (e *DateParseError) Suggestion() string {
	return e.suggestion
}

// PermissionError represents automation permission denied errors.
type PermissionError struct {
	message    string
	suggestion string
}

// NewPermissionError creates a new PermissionError with the given message and suggestion.
func NewPermissionError(message, suggestion string) *PermissionError {
	return &PermissionError{
		message:    message,
		suggestion: suggestion,
	}
}

func (e *PermissionError) Error() string {
	return e.message
}

func (e *PermissionError) ExitCode() int {
	return ExitPermissionError
}

func (e *PermissionError) Suggestion() string {
	return e.suggestion
}
