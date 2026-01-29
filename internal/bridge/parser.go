package bridge

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// TasksResponse represents the JSON response from get_inbox_tasks.js
type TasksResponse struct {
	Tasks []domain.Task `json:"tasks"`
	Error string        `json:"error,omitempty"`
}

// ProjectsResponse represents the JSON response from get_projects.js
type ProjectsResponse struct {
	Projects []domain.Project `json:"projects"`
	Error    string           `json:"error,omitempty"`
}

// TaskResponse represents a single task response
type TaskResponse struct {
	Task  *domain.Task `json:"task,omitempty"`
	Error string       `json:"error,omitempty"`
}

// ProjectResponse represents a single project response
type ProjectResponse struct {
	Project *domain.Project `json:"project,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// TagResponse represents a single tag response
type TagResponse struct {
	Tag   *domain.Tag `json:"tag,omitempty"`
	Error string      `json:"error,omitempty"`
}

// TagsResponse represents an array of tags response
type TagsResponse struct {
	Tags  []domain.Tag `json:"tags"`
	Error string       `json:"error,omitempty"`
}

// TagCountsResponse represents tag counts response
type TagCountsResponse struct {
	Counts map[string]int `json:"counts"`
	Error  string         `json:"error,omitempty"`
}

// OperationResultResponse represents the response from write operations
type OperationResultResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// checkResponseError checks if a response contains an error field
// Returns ErrOmniFocusNotRunning if the error is "OmniFocus is not running"
// Returns error for any other error message
func checkResponseError(errorMsg string) error {
	if errorMsg == "" {
		return nil
	}

	if errorMsg == "OmniFocus is not running" {
		return ErrOmniFocusNotRunning
	}

	return errors.New(errorMsg)
}

// ParseTasks parses JSON output into a slice of Tasks
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON
func ParseTasks(jsonStr string) ([]domain.Task, error) {
	var response TasksResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	// Return empty slice if no tasks (not nil)
	if response.Tasks == nil {
		return []domain.Task{}, nil
	}

	return response.Tasks, nil
}

// ParseProjects parses JSON output into a slice of Projects
func ParseProjects(jsonStr string) ([]domain.Project, error) {
	var response ProjectsResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse projects JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	// Return empty slice if no projects (not nil)
	if response.Projects == nil {
		return []domain.Project{}, nil
	}

	return response.Projects, nil
}

// ParseTask parses JSON output into a single Task
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON
func ParseTask(jsonStr string) (*domain.Task, error) {
	var response TaskResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	return response.Task, nil
}

// ParseProject parses JSON output into a single Project
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON
func ParseProject(jsonStr string) (*domain.Project, error) {
	var response ProjectResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	return response.Project, nil
}

// ParseTag parses JSON output into a single Tag
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON
func ParseTag(jsonStr string) (*domain.Tag, error) {
	var response TagResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	return response.Tag, nil
}

// ParseTags parses JSON output into a slice of Tags
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON
func ParseTags(jsonStr string) ([]domain.Tag, error) {
	var response TagsResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tags JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	// Return empty slice if no tags (not nil)
	if response.Tags == nil {
		return []domain.Tag{}, nil
	}

	return response.Tags, nil
}

// ParseTagCounts parses JSON output into a map of tag names to counts
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON
func ParseTagCounts(jsonStr string) (map[string]int, error) {
	var response TagCountsResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag counts JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	// Return empty map if no counts (not nil)
	if response.Counts == nil {
		return map[string]int{}, nil
	}

	return response.Counts, nil
}

// ParseOperationResult parses JSON output into an OperationResult
// Returns ErrOmniFocusNotRunning if the JSON contains an error about OmniFocus not running
// Returns parsing error for malformed JSON or operation failure
func ParseOperationResult(jsonStr string) (*domain.OperationResult, error) {
	var response OperationResultResponse

	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse operation result JSON: %w", err)
	}

	// Check if response contains an error
	if err := checkResponseError(response.Error); err != nil {
		return nil, err
	}

	// Create domain OperationResult
	result := &domain.OperationResult{
		Success: response.Success,
		ID:      response.ID,
		Message: response.Message,
	}

	return result, nil
}
