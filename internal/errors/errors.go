// Package errors provides structured error types for LazyFocus.
package errors

// Exit codes for LazyFocus CLI
const (
	ExitSuccess         = 0
	ExitGeneralError    = 1
	ExitOmniFocusError  = 2
	ExitItemNotFound    = 3
	ExitValidationError = 4
	ExitPermissionError = 5
)

// LazyFocusError is the base interface for all LazyFocus errors.
// It extends the standard error interface with additional methods
// for providing exit codes and user-facing suggestions.
type LazyFocusError interface {
	error
	ExitCode() int
	Suggestion() string
}
