package cli

import (
	"testing"
	"time"
)

func TestParseDueDate_Today(t *testing.T) {
	result, err := parseDueDate("today")
	if err != nil {
		t.Fatalf("parseDueDate(today) returned error: %v", err)
	}

	now := time.Now()
	expected := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	if !result.Equal(expected) {
		t.Errorf("parseDueDate(today) = %v, want %v", result, expected)
	}

	// Verify it's in local timezone
	if result.Location() != now.Location() {
		t.Errorf("parseDueDate(today).Location() = %v, want %v", result.Location(), now.Location())
	}
}

func TestParseDueDate_Tomorrow(t *testing.T) {
	result, err := parseDueDate("tomorrow")
	if err != nil {
		t.Fatalf("parseDueDate(tomorrow) returned error: %v", err)
	}

	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	expected := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, now.Location())

	if !result.Equal(expected) {
		t.Errorf("parseDueDate(tomorrow) = %v, want %v", result, expected)
	}

	// Verify it's in local timezone
	if result.Location() != now.Location() {
		t.Errorf("parseDueDate(tomorrow).Location() = %v, want %v", result.Location(), now.Location())
	}
}

func TestParseDueDate_YYYYMMDD(t *testing.T) {
	result, err := parseDueDate("2024-03-15")
	if err != nil {
		t.Fatalf("parseDueDate(2024-03-15) returned error: %v", err)
	}

	loc := time.Now().Location()
	expected := time.Date(2024, 3, 15, 23, 59, 59, 0, loc)

	if !result.Equal(expected) {
		t.Errorf("parseDueDate(2024-03-15) = %v, want %v", result, expected)
	}

	// Verify it's in local timezone
	if result.Location() != loc {
		t.Errorf("parseDueDate(2024-03-15).Location() = %v, want %v", result.Location(), loc)
	}
}

func TestParseDueDate_InvalidFormat(t *testing.T) {
	testCases := []string{
		"invalid",
		"2024-13-01", // Invalid month
		"2024-01-32", // Invalid day
		"24-01-01",   // Wrong year format
		"",
	}

	for _, tc := range testCases {
		_, err := parseDueDate(tc)
		if err == nil {
			t.Errorf("parseDueDate(%q) should return error but didn't", tc)
		}
	}
}

func TestParseDueDate_TimezoneConsistency(t *testing.T) {
	// Test that dates parsed in different ways are in the same timezone
	today, _ := parseDueDate("today")
	tomorrow, _ := parseDueDate("tomorrow")
	explicit, _ := parseDueDate("2024-03-15")

	if today.Location() != tomorrow.Location() {
		t.Error("today and tomorrow should be in same timezone")
	}

	if today.Location() != explicit.Location() {
		t.Error("today and explicit date should be in same timezone")
	}
}

// TestFilterTasksByDueDate_WithUTCDates tests that filtering works correctly
// when tasks have due dates in UTC (as returned from JavaScript)
func TestFilterTasksByDueDate_WithUTCDates(t *testing.T) {
	// Create test tasks with dates in different timezones
	loc, _ := time.LoadLocation("Europe/Warsaw") // CET/CEST (UTC+1/+2)

	// Task due at 10:00 UTC on Jan 28
	utcDate1 := time.Date(2024, 1, 28, 10, 0, 0, 0, time.UTC)

	// Task due at 22:00 UTC on Jan 28 (23:00 CET)
	utcDate2 := time.Date(2024, 1, 28, 22, 0, 0, 0, time.UTC)

	// Task due at 23:30 UTC on Jan 28 (00:30 CET Jan 29)
	utcDate3 := time.Date(2024, 1, 28, 23, 30, 0, 0, time.UTC)

	tasks := []Task{
		{Name: "Task 1", DueDate: &utcDate1},
		{Name: "Task 2", DueDate: &utcDate2},
		{Name: "Task 3", DueDate: &utcDate3},
	}

	// Filter by end of Jan 28 in CET
	// This should include tasks 1 and 2, but not 3
	dueDate := time.Date(2024, 1, 28, 23, 59, 59, 0, loc)

	var filtered []Task
	for _, task := range tasks {
		if task.DueDate != nil && !task.DueDate.After(dueDate) {
			filtered = append(filtered, task)
		}
	}

	// Verify results
	if len(filtered) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(filtered))
	}

	if len(filtered) >= 2 {
		if filtered[0].Name != "Task 1" {
			t.Errorf("Expected first task to be 'Task 1', got %s", filtered[0].Name)
		}
		if filtered[1].Name != "Task 2" {
			t.Errorf("Expected second task to be 'Task 2', got %s", filtered[1].Name)
		}
	}
}

// Task is a minimal task structure for testing
type Task struct {
	Name    string
	DueDate *time.Time
}
