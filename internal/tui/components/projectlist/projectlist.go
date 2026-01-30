// Package projectlist provides the project list TUI component.
package projectlist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// Icons for project display
const (
	FolderIcon     = "📁"
	FolderOpenIcon = "📂"
	CheckIcon      = "✓"
	PauseIcon      = "⏸"
	DropIcon       = "✗"
)

// Model represents the project list component state
type Model struct {
	projects []domain.Project
	cursor   int
	width    int
	height   int
	styles   *tui.Styles
	keys     tui.KeyMap
	loading  bool
	empty    bool
}

// New creates a new project list component
func New(styles *tui.Styles, keys tui.KeyMap) Model {
	return Model{
		projects: []domain.Project{},
		cursor:   0,
		styles:   styles,
		keys:     keys,
		loading:  false,
		empty:    true,
	}
}

// Init initializes the component
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if len(m.projects) == 0 {
		return m, nil
	}

	if key.Matches(msg, m.keys.Down) {
		m.cursor++
		if m.cursor >= len(m.projects) {
			m.cursor = 0
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Up) {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.projects) - 1
		}
		return m, nil
	}

	return m, nil
}

// View renders the component
func (m Model) View() string {
	if m.loading {
		return m.renderLoading()
	}
	if m.empty {
		return m.renderEmpty()
	}
	return m.renderProjects()
}

func (m Model) renderLoading() string {
	if m.height == 0 {
		return "Loading..."
	}
	padding := strings.Repeat("\n", m.height/2)
	return padding + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "Loading projects...")
}

func (m Model) renderEmpty() string {
	if m.height == 0 {
		return "No projects"
	}
	padding := strings.Repeat("\n", m.height/2)
	return padding + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "No projects")
}

func (m Model) renderProjects() string {
	var b strings.Builder

	for i, project := range m.projects {
		line := m.formatProjectLine(project, i == m.cursor)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) formatProjectLine(project domain.Project, selected bool) string {
	// Status icon based on project status
	statusIcon := FolderIcon
	switch project.Status {
	case "done", "completed":
		statusIcon = CheckIcon
	case "dropped":
		statusIcon = DropIcon
	case "on hold", "on-hold":
		statusIcon = PauseIcon
	}

	// Build left side
	leftSide := fmt.Sprintf("%s %s", statusIcon, project.Name)

	// Build right side (task count)
	rightSide := fmt.Sprintf("(%d)", project.TaskCount)

	// Calculate spacing
	contentWidth := m.width
	if contentWidth == 0 {
		contentWidth = 80
	}

	leftLen := runewidth.StringWidth(statusIcon) + 1 + runewidth.StringWidth(project.Name)
	rightLen := runewidth.StringWidth(rightSide)
	spacing := contentWidth - leftLen - rightLen - 2
	if spacing < 0 {
		spacing = 1
	}

	line := leftSide + strings.Repeat(" ", spacing) + rightSide

	// Apply styles based on status and selection
	if selected {
		return m.styles.Task.Selected.Render(line)
	}

	switch project.Status {
	case "done", "completed":
		return m.styles.Project.Completed.Render(line)
	case "dropped":
		return m.styles.Project.Dropped.Render(line)
	case "on hold", "on-hold":
		return m.styles.Project.OnHold.Render(line)
	default:
		return m.styles.Project.Active.Render(line)
	}
}

// SetProjects updates the project list
func (m Model) SetProjects(projects []domain.Project) Model {
	m.projects = projects
	m.empty = len(projects) == 0
	m.loading = false
	if m.cursor >= len(m.projects) {
		if len(m.projects) > 0 {
			m.cursor = len(m.projects) - 1
		} else {
			m.cursor = 0
		}
	}
	return m
}

// SetLoading sets the loading state
func (m Model) SetLoading(loading bool) Model {
	m.loading = loading
	return m
}

// SelectedProject returns the currently selected project
func (m Model) SelectedProject() *domain.Project {
	if len(m.projects) == 0 || m.cursor >= len(m.projects) {
		return nil
	}
	return &m.projects[m.cursor]
}

// SelectedIndex returns the current cursor position
func (m Model) SelectedIndex() int {
	return m.cursor
}

// Projects returns all projects (needed for count display)
func (m Model) Projects() []domain.Project {
	return m.projects
}
