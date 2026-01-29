package service

import (
	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// MockOmniFocusService is a mock implementation of OmniFocusService for testing
type MockOmniFocusService struct {
	// Tasks - Read Operations
	InboxTasks      []domain.Task
	InboxTasksErr   error
	AllTasks        []domain.Task
	AllTasksErr     error
	ProjectTasks    []domain.Task
	ProjectTasksErr error
	TagTasks        []domain.Task
	TagTasksErr     error
	FlaggedTasks    []domain.Task
	FlaggedTasksErr error
	Task            *domain.Task
	TaskErr         error

	// Tasks - Write Operations
	CreatedTask     *domain.Task
	CreateTaskErr   error
	ModifiedTask    *domain.Task
	ModifyTaskErr   error
	CompleteResult  *domain.OperationResult
	CompleteTaskErr error
	DeleteResult    *domain.OperationResult
	DeleteTaskErr   error

	// Projects
	Projects            []domain.Project
	ProjectsErr         error
	Project             *domain.Project
	ProjectErr          error
	ProjectWithTasks    *domain.Project
	ProjectWithTasksErr error

	// Tags
	Tags         []domain.Tag
	TagsErr      error
	Tag          *domain.Tag
	TagErr       error
	TagCounts    map[string]int
	TagCountsErr error

	// Perspectives
	PerspectiveTasks    []domain.Task
	PerspectiveTasksErr error

	// Helper Methods
	ResolvedProjectID string
	ResolveProjectErr error
}

// GetInboxTasks returns configured inbox tasks or error
func (m *MockOmniFocusService) GetInboxTasks() ([]domain.Task, error) {
	if m.InboxTasksErr != nil {
		return nil, m.InboxTasksErr
	}
	return m.InboxTasks, nil
}

// GetAllTasks returns configured tasks or error
func (m *MockOmniFocusService) GetAllTasks(filters TaskFilters) ([]domain.Task, error) {
	if m.AllTasksErr != nil {
		return nil, m.AllTasksErr
	}
	return m.AllTasks, nil
}

// GetTasksByProject returns configured project tasks or error
func (m *MockOmniFocusService) GetTasksByProject(projectID string) ([]domain.Task, error) {
	if m.ProjectTasksErr != nil {
		return nil, m.ProjectTasksErr
	}
	return m.ProjectTasks, nil
}

// GetTasksByTag returns configured tag tasks or error
func (m *MockOmniFocusService) GetTasksByTag(tagID string) ([]domain.Task, error) {
	if m.TagTasksErr != nil {
		return nil, m.TagTasksErr
	}
	return m.TagTasks, nil
}

// GetFlaggedTasks returns configured flagged tasks or error
func (m *MockOmniFocusService) GetFlaggedTasks() ([]domain.Task, error) {
	if m.FlaggedTasksErr != nil {
		return nil, m.FlaggedTasksErr
	}
	return m.FlaggedTasks, nil
}

// GetTaskByID returns configured task or error
func (m *MockOmniFocusService) GetTaskByID(id string) (*domain.Task, error) {
	if m.TaskErr != nil {
		return nil, m.TaskErr
	}
	return m.Task, nil
}

// GetProjects returns configured projects or error
func (m *MockOmniFocusService) GetProjects(status string) ([]domain.Project, error) {
	if m.ProjectsErr != nil {
		return nil, m.ProjectsErr
	}
	return m.Projects, nil
}

// GetProjectByID returns configured project or error
func (m *MockOmniFocusService) GetProjectByID(id string) (*domain.Project, error) {
	if m.ProjectErr != nil {
		return nil, m.ProjectErr
	}
	return m.Project, nil
}

// GetProjectWithTasks returns configured project with tasks or error
func (m *MockOmniFocusService) GetProjectWithTasks(id string) (*domain.Project, error) {
	if m.ProjectWithTasksErr != nil {
		return nil, m.ProjectWithTasksErr
	}
	return m.ProjectWithTasks, nil
}

// GetTags returns configured tags or error
func (m *MockOmniFocusService) GetTags() ([]domain.Tag, error) {
	if m.TagsErr != nil {
		return nil, m.TagsErr
	}
	return m.Tags, nil
}

// GetTagByID returns configured tag or error
func (m *MockOmniFocusService) GetTagByID(id string) (*domain.Tag, error) {
	if m.TagErr != nil {
		return nil, m.TagErr
	}
	return m.Tag, nil
}

// GetTagCounts returns configured tag counts or error
func (m *MockOmniFocusService) GetTagCounts() (map[string]int, error) {
	if m.TagCountsErr != nil {
		return nil, m.TagCountsErr
	}
	return m.TagCounts, nil
}

// GetPerspectiveTasks returns configured perspective tasks or error
func (m *MockOmniFocusService) GetPerspectiveTasks(name string) ([]domain.Task, error) {
	if m.PerspectiveTasksErr != nil {
		return nil, m.PerspectiveTasksErr
	}
	return m.PerspectiveTasks, nil
}

// CreateTask returns configured created task or error
func (m *MockOmniFocusService) CreateTask(input domain.TaskInput) (*domain.Task, error) {
	if m.CreateTaskErr != nil {
		return nil, m.CreateTaskErr
	}
	return m.CreatedTask, nil
}

// ModifyTask returns configured modified task or error
func (m *MockOmniFocusService) ModifyTask(id string, mod domain.TaskModification) (*domain.Task, error) {
	if m.ModifyTaskErr != nil {
		return nil, m.ModifyTaskErr
	}
	return m.ModifiedTask, nil
}

// CompleteTask returns configured completion result or error
func (m *MockOmniFocusService) CompleteTask(id string) (*domain.OperationResult, error) {
	if m.CompleteTaskErr != nil {
		return nil, m.CompleteTaskErr
	}
	return m.CompleteResult, nil
}

// DeleteTask returns configured deletion result or error
func (m *MockOmniFocusService) DeleteTask(id string) (*domain.OperationResult, error) {
	if m.DeleteTaskErr != nil {
		return nil, m.DeleteTaskErr
	}
	return m.DeleteResult, nil
}

// ResolveProjectName returns configured project ID or error
func (m *MockOmniFocusService) ResolveProjectName(name string) (string, error) {
	if m.ResolveProjectErr != nil {
		return "", m.ResolveProjectErr
	}
	return m.ResolvedProjectID, nil
}
