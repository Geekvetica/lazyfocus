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
