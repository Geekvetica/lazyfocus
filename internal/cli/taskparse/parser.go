package taskparse

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/dateparse"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

var (
	// Patterns for extracting task components
	tagPattern        = regexp.MustCompile(`#([a-zA-Z0-9_-]+)`)
	projectPattern    = regexp.MustCompile(`@"([^"]+)"|@([a-zA-Z0-9_-]+)`)
	duePattern        = regexp.MustCompile(`due:"([^"]+)"|due:([a-zA-Z0-9_-]+)`)
	deferPattern      = regexp.MustCompile(`defer:"([^"]+)"|defer:([a-zA-Z0-9_-]+)`)
	flagPattern       = regexp.MustCompile(`!`)
	whitespacePattern = regexp.MustCompile(`\s+`)
)

// Parse parses a task input string and extracts structured data.
// Returns TaskInput with extracted fields and remaining text as Name.
func Parse(input string) (domain.TaskInput, error) {
	return ParseWithReference(input, time.Now())
}

// ParseWithReference parses relative to a reference time (for testing).
func ParseWithReference(input string, ref time.Time) (domain.TaskInput, error) {
	if strings.TrimSpace(input) == "" {
		return domain.TaskInput{}, fmt.Errorf("empty task input")
	}

	result := domain.TaskInput{
		TagNames: []string{},
	}

	// Extract tags
	tagMatches := tagPattern.FindAllStringSubmatch(input, -1)
	for _, match := range tagMatches {
		result.TagNames = append(result.TagNames, match[1])
	}

	// Extract project (only first match)
	if projectMatch := projectPattern.FindStringSubmatch(input); projectMatch != nil {
		result.ProjectName = extractValue(projectMatch)
	}

	// Extract due date
	if dueMatch := duePattern.FindStringSubmatch(input); dueMatch != nil {
		dateStr := extractValue(dueMatch)
		dueDate, err := dateparse.ParseWithReference(dateStr, ref)
		if err != nil {
			return domain.TaskInput{}, fmt.Errorf("invalid due date: %w", err)
		}
		result.DueDate = &dueDate
	}

	// Extract defer date
	if deferMatch := deferPattern.FindStringSubmatch(input); deferMatch != nil {
		dateStr := extractValue(deferMatch)
		deferDate, err := dateparse.ParseWithReference(dateStr, ref)
		if err != nil {
			return domain.TaskInput{}, fmt.Errorf("invalid defer date: %w", err)
		}
		result.DeferDate = &deferDate
	}

	// Extract flagged status
	if flagPattern.MatchString(input) {
		flagged := true
		result.Flagged = &flagged
	}

	// Remove all modifiers from input to get task name
	name := input
	name = tagPattern.ReplaceAllString(name, "")
	name = projectPattern.ReplaceAllString(name, "")
	name = duePattern.ReplaceAllString(name, "")
	name = deferPattern.ReplaceAllString(name, "")
	name = flagPattern.ReplaceAllString(name, "")

	// Collapse multiple spaces into one and trim
	name = whitespacePattern.ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)

	if name == "" {
		return domain.TaskInput{}, fmt.Errorf("task name is required")
	}

	result.Name = name

	return result, nil
}

// extractValue extracts a quoted or unquoted value from regex match.
// Assumes match[1] contains quoted value and match[2] contains unquoted value.
func extractValue(match []string) string {
	if match[1] != "" {
		return match[1] // Quoted
	}
	return match[2] // Unquoted
}
