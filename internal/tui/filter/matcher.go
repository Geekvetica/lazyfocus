// Package filter provides task filtering functionality for the TUI.
package filter

import (
	"strings"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// Matcher filters tasks based on filter state
type Matcher struct {
	state State
}

// NewMatcher creates a new Matcher with the given state
func NewMatcher(state State) *Matcher {
	return &Matcher{state: state}
}

// FilterTasks returns tasks that match the current filter state
func (m *Matcher) FilterTasks(tasks []domain.Task) []domain.Task {
	if !m.state.IsActive() {
		return tasks
	}

	result := make([]domain.Task, 0, len(tasks))
	for _, task := range tasks {
		if m.matches(task) {
			result = append(result, task)
		}
	}
	return result
}

// matches checks if a single task matches the filter state
func (m *Matcher) matches(task domain.Task) bool {
	// Search text filter (case-insensitive)
	if m.state.SearchText != "" {
		searchLower := strings.ToLower(m.state.SearchText)
		nameLower := strings.ToLower(task.Name)
		noteLower := strings.ToLower(task.Note)

		if !strings.Contains(nameLower, searchLower) && !strings.Contains(noteLower, searchLower) {
			return false
		}
	}

	// Project filter
	if m.state.ProjectID != "" && task.ProjectID != m.state.ProjectID {
		return false
	}

	// Tag filter - adapted for []string tags
	if m.state.TagID != "" {
		hasTag := false
		for _, tag := range task.Tags {
			if tag == m.state.TagID {
				hasTag = true
				break
			}
		}
		if !hasTag {
			return false
		}
	}

	// Flagged filter
	if m.state.FlaggedOnly && !task.Flagged {
		return false
	}

	// Due date filter
	if m.state.DueFilter != DueNone {
		if !m.matchesDueFilter(task) {
			return false
		}
	}

	return true
}

// matchesDueFilter checks if task due date matches the due filter
func (m *Matcher) matchesDueFilter(task domain.Task) bool {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	weekEnd := today.AddDate(0, 0, 7)

	switch m.state.DueFilter {
	case DueOverdue:
		return task.DueDate != nil && task.DueDate.Before(today)
	case DueToday:
		if task.DueDate == nil {
			return false
		}
		return !task.DueDate.Before(today) && task.DueDate.Before(tomorrow)
	case DueTomorrow:
		if task.DueDate == nil {
			return false
		}
		dayAfterTomorrow := tomorrow.AddDate(0, 0, 1)
		return !task.DueDate.Before(tomorrow) && task.DueDate.Before(dayAfterTomorrow)
	case DueWeek:
		if task.DueDate == nil {
			return false
		}
		return !task.DueDate.Before(today) && task.DueDate.Before(weekEnd)
	default:
		return true
	}
}
