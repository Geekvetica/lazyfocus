// Package dateparse provides natural language date parsing for LazyFocus.
package dateparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	weekdays = map[string]time.Weekday{
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
		"sunday":    time.Sunday,
	}

	months = map[string]time.Month{
		"january":   time.January,
		"jan":       time.January,
		"february":  time.February,
		"feb":       time.February,
		"march":     time.March,
		"mar":       time.March,
		"april":     time.April,
		"apr":       time.April,
		"may":       time.May,
		"june":      time.June,
		"jun":       time.June,
		"july":      time.July,
		"jul":       time.July,
		"august":    time.August,
		"aug":       time.August,
		"september": time.September,
		"sep":       time.September,
		"sept":      time.September,
		"october":   time.October,
		"oct":       time.October,
		"november":  time.November,
		"nov":       time.November,
		"december":  time.December,
		"dec":       time.December,
	}
)

// Parse parses a natural language date string and returns the time.
// For dates without explicit times, returns 5:00 PM local time.
// Returns error if the format is not recognized.
func Parse(input string) (time.Time, error) {
	return ParseWithReference(input, time.Now())
}

// ParseWithReference parses relative to a reference time (useful for testing).
func ParseWithReference(input string, ref time.Time) (time.Time, error) {
	if input == "" {
		return time.Time{}, fmt.Errorf("empty date string")
	}

	// Normalize input to lowercase for case-insensitive parsing
	normalized := strings.ToLower(strings.TrimSpace(input))

	// Try each parser in order
	parsers := []func(string, time.Time) (time.Time, bool){
		parseRelativeDay,
		parseNextWeekday,
		parseInDaysWeeks,
		parseNextWeek,
		parseISO,
		parseMonthDay,
	}

	for _, parser := range parsers {
		if result, ok := parser(normalized, ref); ok {
			return result, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized date format: %s", input)
}

// parseRelativeDay handles "today", "tomorrow", "yesterday"
func parseRelativeDay(input string, ref time.Time) (time.Time, bool) {
	var days int
	switch input {
	case "today":
		days = 0
	case "tomorrow":
		days = 1
	case "yesterday":
		days = -1
	default:
		return time.Time{}, false
	}

	result := ref.AddDate(0, 0, days)
	return setTo5PM(result), true
}

// parseNextWeek handles "next week"
func parseNextWeek(input string, ref time.Time) (time.Time, bool) {
	if input != "next week" {
		return time.Time{}, false
	}

	result := ref.AddDate(0, 0, 7)
	return setTo5PM(result), true
}

// parseNextWeekday handles "next monday", "next tuesday", etc.
func parseNextWeekday(input string, ref time.Time) (time.Time, bool) {
	if !strings.HasPrefix(input, "next ") {
		return time.Time{}, false
	}

	weekdayStr := strings.TrimPrefix(input, "next ")
	targetWeekday, ok := weekdays[weekdayStr]
	if !ok {
		return time.Time{}, false
	}

	// Find the next occurrence of the target weekday
	currentWeekday := ref.Weekday()
	daysUntil := int(targetWeekday - currentWeekday)
	if daysUntil <= 0 {
		daysUntil += 7 // Move to next week
	}

	result := ref.AddDate(0, 0, daysUntil)
	return setTo5PM(result), true
}

// parseInDaysWeeks handles "in N days" and "in N weeks"
func parseInDaysWeeks(input string, ref time.Time) (time.Time, bool) {
	// Pattern: "in N day(s)" or "in N week(s)"
	re := regexp.MustCompile(`^in\s+(\d+)\s+(day|days|week|weeks)$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return time.Time{}, false
	}

	n, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, false
	}

	unit := matches[2]
	days := n
	if strings.HasPrefix(unit, "week") {
		days = n * 7
	}

	result := ref.AddDate(0, 0, days)
	return setTo5PM(result), true
}

// parseISO handles ISO date format "2024-01-15"
func parseISO(input string, ref time.Time) (time.Time, bool) {
	// Pattern: YYYY-MM-DD
	re := regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return time.Time{}, false
	}

	year, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Time{}, false
	}
	month, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Time{}, false
	}
	day, err := strconv.Atoi(matches[3])
	if err != nil {
		return time.Time{}, false
	}

	result := time.Date(year, time.Month(month), day, 17, 0, 0, 0, time.Local)
	return result, true
}

// parseMonthDay handles "Jan 15", "January 15", "Jan 15 2024", "January 15 2024"
func parseMonthDay(input string, ref time.Time) (time.Time, bool) {
	// Pattern: "monthname day" or "monthname day year"
	parts := strings.Fields(input)
	if len(parts) < 2 || len(parts) > 3 {
		return time.Time{}, false
	}

	monthStr := parts[0]
	month, ok := months[monthStr]
	if !ok {
		return time.Time{}, false
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, false
	}

	year := ref.Year()
	if len(parts) == 3 {
		year, err = strconv.Atoi(parts[2])
		if err != nil {
			return time.Time{}, false
		}
	}

	result := time.Date(year, month, day, 17, 0, 0, 0, time.Local)
	return result, true
}

// setTo5PM sets the time to 5:00 PM (17:00) local time
func setTo5PM(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 17, 0, 0, 0, time.Local)
}
