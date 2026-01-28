package service

import (
	"errors"
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/domain"
)

// Compile-time check that MockOmniFocusService implements OmniFocusService
var _ OmniFocusService = (*MockOmniFocusService)(nil)

func TestMockOmniFocusService_ImplementsInterface(t *testing.T) {
	var service OmniFocusService = &MockOmniFocusService{}
	if service == nil {
		t.Fatal("MockOmniFocusService should implement OmniFocusService interface")
	}
}

func TestMockOmniFocusService_GetInboxTasks_ReturnsConfiguredValues(t *testing.T) {
	expectedTasks := []domain.Task{
		{ID: "task1", Name: "Test Task"},
	}

	mock := &MockOmniFocusService{
		InboxTasks: expectedTasks,
	}

	tasks, err := mock.GetInboxTasks()
	if err != nil {
		t.Fatalf("GetInboxTasks() error = %v, want nil", err)
	}

	if len(tasks) != 1 {
		t.Errorf("GetInboxTasks() returned %d tasks, want 1", len(tasks))
	}
}

func TestMockOmniFocusService_GetInboxTasks_ReturnsConfiguredError(t *testing.T) {
	expectedErr := errors.New("test error")
	mock := &MockOmniFocusService{
		InboxTasksErr: expectedErr,
	}

	_, err := mock.GetInboxTasks()
	if err != expectedErr {
		t.Errorf("GetInboxTasks() error = %v, want %v", err, expectedErr)
	}
}

func TestMockOmniFocusService_GetProjects_ReturnsConfiguredValues(t *testing.T) {
	expectedProjects := []domain.Project{
		{ID: "proj1", Name: "Test Project", Status: "active"},
	}

	mock := &MockOmniFocusService{
		Projects: expectedProjects,
	}

	projects, err := mock.GetProjects("active")
	if err != nil {
		t.Fatalf("GetProjects() error = %v, want nil", err)
	}

	if len(projects) != 1 {
		t.Errorf("GetProjects() returned %d projects, want 1", len(projects))
	}
}

func TestMockOmniFocusService_GetTaskByID_ReturnsConfiguredTask(t *testing.T) {
	expectedTask := &domain.Task{
		ID:   "task123",
		Name: "Specific Task",
	}

	mock := &MockOmniFocusService{
		Task: expectedTask,
	}

	task, err := mock.GetTaskByID("task123")
	if err != nil {
		t.Fatalf("GetTaskByID() error = %v, want nil", err)
	}

	if task.ID != expectedTask.ID {
		t.Errorf("GetTaskByID() returned task with ID %s, want %s", task.ID, expectedTask.ID)
	}
}

func TestMockOmniFocusService_GetTagCounts_ReturnsConfiguredCounts(t *testing.T) {
	expectedCounts := map[string]int{
		"tag1": 5,
		"tag2": 10,
	}

	mock := &MockOmniFocusService{
		TagCounts: expectedCounts,
	}

	counts, err := mock.GetTagCounts()
	if err != nil {
		t.Fatalf("GetTagCounts() error = %v, want nil", err)
	}

	if len(counts) != 2 {
		t.Errorf("GetTagCounts() returned %d counts, want 2", len(counts))
	}

	if counts["tag1"] != 5 {
		t.Errorf("GetTagCounts()[\"tag1\"] = %d, want 5", counts["tag1"])
	}
}

func TestMockOmniFocusService_AllMethods_CanReturnErrors(t *testing.T) {
	testErr := errors.New("test error")

	tests := []struct {
		name     string
		mockFunc func(*MockOmniFocusService) error
	}{
		{
			name: "GetInboxTasks",
			mockFunc: func(m *MockOmniFocusService) error {
				m.InboxTasksErr = testErr
				_, err := m.GetInboxTasks()
				return err
			},
		},
		{
			name: "GetAllTasks",
			mockFunc: func(m *MockOmniFocusService) error {
				m.AllTasksErr = testErr
				_, err := m.GetAllTasks(TaskFilters{})
				return err
			},
		},
		{
			name: "GetTasksByProject",
			mockFunc: func(m *MockOmniFocusService) error {
				m.ProjectTasksErr = testErr
				_, err := m.GetTasksByProject("proj1")
				return err
			},
		},
		{
			name: "GetTasksByTag",
			mockFunc: func(m *MockOmniFocusService) error {
				m.TagTasksErr = testErr
				_, err := m.GetTasksByTag("tag1")
				return err
			},
		},
		{
			name: "GetFlaggedTasks",
			mockFunc: func(m *MockOmniFocusService) error {
				m.FlaggedTasksErr = testErr
				_, err := m.GetFlaggedTasks()
				return err
			},
		},
		{
			name: "GetTaskByID",
			mockFunc: func(m *MockOmniFocusService) error {
				m.TaskErr = testErr
				_, err := m.GetTaskByID("task1")
				return err
			},
		},
		{
			name: "GetProjects",
			mockFunc: func(m *MockOmniFocusService) error {
				m.ProjectsErr = testErr
				_, err := m.GetProjects("active")
				return err
			},
		},
		{
			name: "GetProjectByID",
			mockFunc: func(m *MockOmniFocusService) error {
				m.ProjectErr = testErr
				_, err := m.GetProjectByID("proj1")
				return err
			},
		},
		{
			name: "GetProjectWithTasks",
			mockFunc: func(m *MockOmniFocusService) error {
				m.ProjectWithTasksErr = testErr
				_, err := m.GetProjectWithTasks("proj1")
				return err
			},
		},
		{
			name: "GetTags",
			mockFunc: func(m *MockOmniFocusService) error {
				m.TagsErr = testErr
				_, err := m.GetTags()
				return err
			},
		},
		{
			name: "GetTagByID",
			mockFunc: func(m *MockOmniFocusService) error {
				m.TagErr = testErr
				_, err := m.GetTagByID("tag1")
				return err
			},
		},
		{
			name: "GetTagCounts",
			mockFunc: func(m *MockOmniFocusService) error {
				m.TagCountsErr = testErr
				_, err := m.GetTagCounts()
				return err
			},
		},
		{
			name: "GetPerspectiveTasks",
			mockFunc: func(m *MockOmniFocusService) error {
				m.PerspectiveTasksErr = testErr
				_, err := m.GetPerspectiveTasks("Today")
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockOmniFocusService{}
			err := tt.mockFunc(mock)
			if err != testErr {
				t.Errorf("%s error = %v, want %v", tt.name, err, testErr)
			}
		})
	}
}
