package projectlist

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

func TestNew(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	m := New(styles, keys)

	if len(m.projects) != 0 {
		t.Error("should start with empty projects list")
	}
	if !m.empty {
		t.Error("should start in empty state")
	}
	if m.loading {
		t.Error("should not be loading initially")
	}
}

func TestSetProjects(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	projects := []domain.Project{
		{ID: "p1", Name: "Project 1", Status: "active", TaskCount: 5},
		{ID: "p2", Name: "Project 2", Status: "on-hold", TaskCount: 3},
	}

	m = m.SetProjects(projects)

	if len(m.projects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(m.projects))
	}
	if m.empty {
		t.Error("should not be empty after setting projects")
	}
	if m.loading {
		t.Error("should not be loading after setting projects")
	}
}

func TestSetLoading(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	m = m.SetLoading(true)

	if !m.loading {
		t.Error("should be loading")
	}

	m = m.SetLoading(false)

	if m.loading {
		t.Error("should not be loading")
	}
}

func TestSelectedProject(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	// No projects - should return nil
	if m.SelectedProject() != nil {
		t.Error("should return nil when no projects")
	}

	// With projects
	projects := []domain.Project{
		{ID: "p1", Name: "Project 1", Status: "active", TaskCount: 5},
		{ID: "p2", Name: "Project 2", Status: "on-hold", TaskCount: 3},
	}
	m = m.SetProjects(projects)

	selected := m.SelectedProject()
	if selected == nil {
		t.Fatal("should return selected project")
	}
	if selected.ID != "p1" {
		t.Errorf("expected p1, got %s", selected.ID)
	}
}

func TestNavigationDownUp(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	projects := []domain.Project{
		{ID: "p1", Name: "Project 1", Status: "active", TaskCount: 5},
		{ID: "p2", Name: "Project 2", Status: "on-hold", TaskCount: 3},
		{ID: "p3", Name: "Project 3", Status: "active", TaskCount: 1},
	}
	m = m.SetProjects(projects)

	// Initial cursor should be 0
	if m.cursor != 0 {
		t.Errorf("expected cursor 0, got %d", m.cursor)
	}

	// Press down
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", m.cursor)
	}

	// Press down again
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 2 {
		t.Errorf("expected cursor 2, got %d", m.cursor)
	}

	// Press down - should wrap to 0
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.cursor != 0 {
		t.Errorf("expected cursor to wrap to 0, got %d", m.cursor)
	}

	// Press up - should wrap to last
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 2 {
		t.Errorf("expected cursor to wrap to 2, got %d", m.cursor)
	}

	// Press up
	m, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", m.cursor)
	}
}

func TestViewLoading(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	m = m.SetLoading(true)
	m.height = 20

	view := m.View()

	if !strings.Contains(view, "Loading projects") {
		t.Error("should show loading message")
	}
}

func TestViewEmpty(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	m.height = 20

	view := m.View()

	if !strings.Contains(view, "No projects") {
		t.Error("should show empty message")
	}
}

func TestViewWithProjects(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	projects := []domain.Project{
		{ID: "p1", Name: "Work Project", Status: "active", TaskCount: 5},
		{ID: "p2", Name: "Personal", Status: "on-hold", TaskCount: 3},
	}
	m = m.SetProjects(projects)
	m.width = 80

	view := m.View()

	if !strings.Contains(view, "Work Project") {
		t.Error("should show project name")
	}
	if !strings.Contains(view, "(5)") {
		t.Error("should show task count")
	}
}

func TestFormatProjectLine_StatusIcons(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)
	m.width = 80

	tests := []struct {
		name         string
		status       string
		expectedIcon string
	}{
		{"Active project", "active", FolderIcon},
		{"Completed project (done)", "done", CheckIcon},
		{"Completed project (completed)", "completed", CheckIcon},
		{"Dropped project", "dropped", DropIcon},
		{"On hold project (space)", "on hold", PauseIcon},
		{"On hold project (hyphen)", "on-hold", PauseIcon},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := domain.Project{
				ID:        "p1",
				Name:      "Test Project",
				Status:    tt.status,
				TaskCount: 5,
			}

			line := m.formatProjectLine(project, false)

			// Verify the line is not empty
			if len(line) == 0 {
				t.Error("expected non-empty line")
			}

			// Verify the expected icon appears in the line
			if !strings.Contains(line, tt.expectedIcon) {
				t.Errorf("expected icon %s in line, got: %s", tt.expectedIcon, line)
			}
		})
	}
}

func TestWindowSizeMsg(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	m, _ = m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})

	if m.width != 100 {
		t.Errorf("expected width 100, got %d", m.width)
	}
	if m.height != 50 {
		t.Errorf("expected height 50, got %d", m.height)
	}
}

func TestProjects_ReturnsAllProjects(t *testing.T) {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()
	m := New(styles, keys)

	projects := []domain.Project{
		{ID: "p1", Name: "Project 1", Status: "active", TaskCount: 5},
		{ID: "p2", Name: "Project 2", Status: "on-hold", TaskCount: 3},
	}
	m = m.SetProjects(projects)

	allProjects := m.Projects()
	if len(allProjects) != 2 {
		t.Errorf("expected 2 projects, got %d", len(allProjects))
	}
}
