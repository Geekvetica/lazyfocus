// Package validation provides input validation for LazyFocus.
package validation

import (
	"strings"
	"unicode"

	"github.com/pwojciechowski/lazyfocus/internal/errors"
)

// Field validation limits
const (
	MaxTaskNameLength    = 500
	MaxNoteLength        = 10000
	MaxProjectNameLength = 200
	MaxTagNameLength     = 100
)

// ValidateTaskName validates a task name.
// Returns a ValidationError if the name is empty, too long, or contains invalid characters.
func ValidateTaskName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.NewValidationError("task name is required", "provide a non-empty task name")
	}
	if len(name) > MaxTaskNameLength {
		return errors.NewValidationError(
			"task name too long",
			"task name must be less than 500 characters",
		)
	}
	if containsControlChars(name) {
		return errors.NewValidationError(
			"task name contains invalid characters",
			"remove control characters from task name",
		)
	}
	return nil
}

// ValidateNote validates a note.
// Returns a ValidationError if the note is too long.
func ValidateNote(note string) error {
	if len(note) > MaxNoteLength {
		return errors.NewValidationError(
			"note too long",
			"note must be less than 10000 characters",
		)
	}
	return nil
}

// ValidateProjectName validates a project name.
// Project names are optional, so empty strings are allowed.
// Returns a ValidationError if the name is too long.
func ValidateProjectName(name string) error {
	if name == "" {
		return nil // Project is optional
	}
	if len(name) > MaxProjectNameLength {
		return errors.NewValidationError(
			"project name too long",
			"project name must be less than 200 characters",
		)
	}
	return nil
}

// ValidateTagName validates a tag name.
// Returns a ValidationError if the name is empty or too long.
func ValidateTagName(name string) error {
	if name == "" {
		return errors.NewValidationError("tag name is required", "provide a non-empty tag name")
	}
	if len(name) > MaxTagNameLength {
		return errors.NewValidationError(
			"tag name too long",
			"tag name must be less than 100 characters",
		)
	}
	return nil
}

// containsControlChars checks if a string contains control characters
// (excluding newlines and tabs which are allowed).
func containsControlChars(s string) bool {
	for _, r := range s {
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return true
		}
	}
	return false
}
