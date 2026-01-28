// Package service provides the business logic layer that sits between
// the CLI/TUI and the bridge layer. It orchestrates script execution
// and response parsing to provide high-level OmniFocus operations.
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/bridge"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// TaskFilters defines filtering criteria for task queries
type TaskFilters struct {
	Inbox     bool
	ProjectID string
	TagID     string
	Flagged   bool
	DueStart  *time.Time
	DueEnd    *time.Time
	Completed bool
}

// OmniFocusService defines the interface for interacting with OmniFocus
type OmniFocusService interface {
	// Tasks
	GetInboxTasks() ([]domain.Task, error)
	GetAllTasks(filters TaskFilters) ([]domain.Task, error)
	GetTasksByProject(projectID string) ([]domain.Task, error)
	GetTasksByTag(tagID string) ([]domain.Task, error)
	GetFlaggedTasks() ([]domain.Task, error)
	GetTaskByID(id string) (*domain.Task, error)

	// Projects
	GetProjects(status string) ([]domain.Project, error)
	GetProjectByID(id string) (*domain.Project, error)
	GetProjectWithTasks(id string) (*domain.Project, error)

	// Tags
	GetTags() ([]domain.Tag, error)
	GetTagByID(id string) (*domain.Tag, error)
	GetTagCounts() (map[string]int, error)

	// Perspectives
	GetPerspectiveTasks(name string) ([]domain.Task, error)
}

// DefaultOmniFocusService implements OmniFocusService using the bridge layer
type DefaultOmniFocusService struct {
	executor bridge.Executor
	timeout  time.Duration
}

// NewOmniFocusService creates a new OmniFocusService instance
func NewOmniFocusService(executor bridge.Executor, timeout time.Duration) *DefaultOmniFocusService {
	return &DefaultOmniFocusService{
		executor: executor,
		timeout:  timeout,
	}
}

// GetInboxTasks retrieves all tasks from the OmniFocus inbox
func (s *DefaultOmniFocusService) GetInboxTasks() ([]domain.Task, error) {
	script, err := bridge.GetScript("get_inbox_tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to load inbox tasks script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute inbox tasks script: %w", err)
	}

	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inbox tasks: %w", err)
	}

	return tasks, nil
}

// GetAllTasks retrieves all tasks matching the provided filters
func (s *DefaultOmniFocusService) GetAllTasks(filters TaskFilters) ([]domain.Task, error) {
	// For now, use inbox tasks script as base
	// TODO: Implement filtering logic once we have more scripts
	script, err := bridge.GetScript("get_inbox_tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tasks script: %w", err)
	}

	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tasks: %w", err)
	}

	return tasks, nil
}

// GetTasksByProject retrieves all tasks for a specific project
func (s *DefaultOmniFocusService) GetTasksByProject(projectID string) ([]domain.Task, error) {
	// TODO: Implement with project-specific script
	params := map[string]string{
		"ProjectID": projectID,
	}

	script, err := bridge.GetScriptWithParams("get_inbox_tasks", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load project tasks script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute project tasks script: %w", err)
	}

	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project tasks: %w", err)
	}

	return tasks, nil
}

// GetTasksByTag retrieves all tasks with a specific tag
func (s *DefaultOmniFocusService) GetTasksByTag(tagID string) ([]domain.Task, error) {
	// TODO: Implement with tag-specific script
	params := map[string]string{
		"TagID": tagID,
	}

	script, err := bridge.GetScriptWithParams("get_inbox_tasks", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load tag tasks script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tag tasks script: %w", err)
	}

	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tag tasks: %w", err)
	}

	return tasks, nil
}

// GetFlaggedTasks retrieves all flagged tasks
func (s *DefaultOmniFocusService) GetFlaggedTasks() ([]domain.Task, error) {
	// TODO: Implement with flagged tasks script
	script, err := bridge.GetScript("get_inbox_tasks")
	if err != nil {
		return nil, fmt.Errorf("failed to load flagged tasks script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute flagged tasks script: %w", err)
	}

	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse flagged tasks: %w", err)
	}

	return tasks, nil
}

// GetTaskByID retrieves a single task by its ID
func (s *DefaultOmniFocusService) GetTaskByID(id string) (*domain.Task, error) {
	// TODO: Implement with task-by-id script
	params := map[string]string{
		"TaskID": id,
	}

	script, err := bridge.GetScriptWithParams("get_inbox_tasks", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load task script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute task script: %w", err)
	}

	tasks, err := bridge.ParseTasks(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse task: %w", err)
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("task not found: %s", id)
	}

	return &tasks[0], nil
}

// GetProjects retrieves projects filtered by status
func (s *DefaultOmniFocusService) GetProjects(status string) ([]domain.Project, error) {
	script, err := bridge.GetScript("get_projects")
	if err != nil {
		return nil, fmt.Errorf("failed to load projects script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute projects script: %w", err)
	}

	projects, err := bridge.ParseProjects(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projects, nil
}

// GetProjectByID retrieves a single project by its ID
func (s *DefaultOmniFocusService) GetProjectByID(id string) (*domain.Project, error) {
	// TODO: Implement with project-by-id script
	params := map[string]string{
		"ProjectID": id,
	}

	script, err := bridge.GetScriptWithParams("get_projects", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load project script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute project script: %w", err)
	}

	projects, err := bridge.ParseProjects(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project not found: %s", id)
	}

	return &projects[0], nil
}

// GetProjectWithTasks retrieves a project with all its tasks
func (s *DefaultOmniFocusService) GetProjectWithTasks(id string) (*domain.Project, error) {
	// TODO: Implement with project-with-tasks script
	params := map[string]string{
		"ProjectID": id,
	}

	script, err := bridge.GetScriptWithParams("get_projects", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load project script: %w", err)
	}

	output, err := s.executor.ExecuteWithTimeout(script, s.timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to execute project script: %w", err)
	}

	projects, err := bridge.ParseProjects(output)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project not found: %s", id)
	}

	return &projects[0], nil
}

// GetTags retrieves all tags from OmniFocus
func (s *DefaultOmniFocusService) GetTags() ([]domain.Tag, error) {
	// TODO: Implement once we have get_tags script
	return []domain.Tag{}, errors.New("not implemented")
}

// GetTagByID retrieves a single tag by its ID
func (s *DefaultOmniFocusService) GetTagByID(id string) (*domain.Tag, error) {
	// TODO: Implement once we have tag-by-id script
	return nil, errors.New("not implemented")
}

// GetTagCounts retrieves the count of tasks for each tag
func (s *DefaultOmniFocusService) GetTagCounts() (map[string]int, error) {
	// TODO: Implement once we have tag counts script
	return nil, errors.New("not implemented")
}

// GetPerspectiveTasks retrieves tasks from a named perspective
func (s *DefaultOmniFocusService) GetPerspectiveTasks(name string) ([]domain.Task, error) {
	// TODO: Implement once we have perspective script
	return []domain.Task{}, errors.New("not implemented")
}
