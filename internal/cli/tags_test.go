package cli

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

func TestTagsCommand_Basic(t *testing.T) {
	// Test that basic tags command shows tags with hierarchy
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work", Children: []domain.Tag{
				{ID: "tag1-1", Name: "Meetings", ParentID: "tag1"},
				{ID: "tag1-2", Name: "Code Review", ParentID: "tag1"},
			}},
			{ID: "tag2", Name: "Personal"},
		},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Work") {
		t.Errorf("Expected output to contain 'Work', got: %s", output)
	}

	if !strings.Contains(output, "Personal") {
		t.Errorf("Expected output to contain 'Personal', got: %s", output)
	}

	if !strings.Contains(output, "Meetings") {
		t.Errorf("Expected output to contain 'Meetings', got: %s", output)
	}
}

func TestTagsCommand_Flat(t *testing.T) {
	// Test --flat flag shows tags without hierarchy indentation
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work", Children: []domain.Tag{
				{ID: "tag1-1", Name: "Meetings", ParentID: "tag1"},
			}},
			{ID: "tag2", Name: "Personal"},
		},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{"--flat"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Work") {
		t.Errorf("Expected output to contain 'Work', got: %s", output)
	}

	if !strings.Contains(output, "Meetings") {
		t.Errorf("Expected output to contain 'Meetings', got: %s", output)
	}

	if !strings.Contains(output, "Personal") {
		t.Errorf("Expected output to contain 'Personal', got: %s", output)
	}
}

func TestTagsCommand_WithCounts(t *testing.T) {
	// Test --with-counts flag shows task counts
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work"},
			{ID: "tag2", Name: "Personal"},
		},
		TagCounts: map[string]int{
			"tag1": 5,
			"tag2": 3,
		},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{"--with-counts"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Work") {
		t.Errorf("Expected output to contain 'Work', got: %s", output)
	}

	if !strings.Contains(output, "Personal") {
		t.Errorf("Expected output to contain 'Personal', got: %s", output)
	}

	// When with-counts is specified, counts should be fetched
	// The formatting might show counts - we'll verify the service was called
	// by checking the output doesn't error
}

func TestTagsCommand_FlatWithCounts(t *testing.T) {
	// Test combining --flat and --with-counts
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work", Children: []domain.Tag{
				{ID: "tag1-1", Name: "Meetings", ParentID: "tag1"},
			}},
			{ID: "tag2", Name: "Personal"},
		},
		TagCounts: map[string]int{
			"tag1":   5,
			"tag1-1": 2,
			"tag2":   3,
		},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{"--flat", "--with-counts"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "Work") {
		t.Errorf("Expected output to contain 'Work', got: %s", output)
	}

	if !strings.Contains(output, "Meetings") {
		t.Errorf("Expected output to contain 'Meetings', got: %s", output)
	}
}

func TestTagsCommand_JSONOutput(t *testing.T) {
	// Test JSON output format
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work"},
			{ID: "tag2", Name: "Personal"},
		},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{"--json"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	// Check for JSON structure
	if !strings.Contains(output, `"tags"`) {
		t.Errorf("Expected JSON output to contain 'tags' field, got: %s", output)
	}

	if !strings.Contains(output, `"Work"`) {
		t.Errorf("Expected JSON output to contain tag name, got: %s", output)
	}

	if !strings.Contains(output, `"count"`) {
		t.Errorf("Expected JSON output to contain 'count' field, got: %s", output)
	}
}

func TestTagsCommand_Error(t *testing.T) {
	// Test error handling
	mockService := &service.MockOmniFocusService{
		TagsErr: errors.New("OmniFocus is not running"),
	}

	_, exitCode, err := executeTagsCommand(mockService, []string{})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "OmniFocus is not running") {
		t.Errorf("Expected error message about OmniFocus, got: %v", err)
	}
}

func TestTagsCommand_ErrorJSON(t *testing.T) {
	// Test error handling in JSON mode
	mockService := &service.MockOmniFocusService{
		TagsErr: errors.New("OmniFocus is not running"),
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{"--json"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	// In JSON mode, error should be in output
	if !strings.Contains(output, `"error"`) {
		t.Errorf("Expected JSON error output to contain 'error' field, got: %s", output)
	}
}

func TestTagsCommand_QuietMode(t *testing.T) {
	// Test quiet mode suppresses output
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work"},
		},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{"--quiet"})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if output != "" {
		t.Errorf("Expected empty output in quiet mode, got: %s", output)
	}
}

func TestTagsCommand_EmptyResults(t *testing.T) {
	// Test empty tags list
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{},
	}

	output, exitCode, err := executeTagsCommand(mockService, []string{})

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got: %d", exitCode)
	}

	if !strings.Contains(output, "No tags") {
		t.Errorf("Expected output to indicate no tags, got: %s", output)
	}
}

func TestTagsCommand_WithCountsError(t *testing.T) {
	// Test error when fetching tag counts fails
	mockService := &service.MockOmniFocusService{
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Work"},
		},
		TagCountsErr: errors.New("failed to fetch tag counts"),
	}

	_, exitCode, err := executeTagsCommand(mockService, []string{"--with-counts"})

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code, got: %d", exitCode)
	}

	if !strings.Contains(err.Error(), "failed to fetch tag counts") {
		t.Errorf("Expected error message about tag counts, got: %v", err)
	}
}

// Helper function to execute tags command and capture output
func executeTagsCommand(mockService service.OmniFocusService, args []string) (string, int, error) {
	// Create a new root command for each test to avoid flag pollution
	rootCmd := newTestRootCommand()

	// Override the service for testing
	Service = mockService
	defer func() { Service = nil }()

	// Add tags command
	rootCmd.AddCommand(NewTagsCommand())

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Prepare args - need to add "tags" as first arg
	fullArgs := append([]string{"tags"}, args...)
	rootCmd.SetArgs(fullArgs)

	// Execute
	err := rootCmd.Execute()

	output := buf.String()
	exitCode := 0
	if err != nil {
		exitCode = 1 // Simplified - in real implementation we'd parse specific error types
	}

	return output, exitCode, err
}
