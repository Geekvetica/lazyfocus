// Package projects provides the projects view for the TUI.
package projects

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/projectlist"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/tasklist"
)

// ViewMode represents whether we're viewing projects or a project's tasks
type ViewMode int

const (
	ModeProjectList ViewMode = iota
	ModeProjectTasks
)

// Model represents the projects view state
type Model struct {
	projectList    projectlist.Model
	taskList       tasklist.Model
	service        service.OmniFocusService
	styles         *tui.Styles
	keys           tui.KeyMap
	mode           ViewMode
	currentProject *domain.Project
	width          int
	height         int
	err            error
	loaded         bool
}

// New creates a new projects view
func New(styles *tui.Styles, keys tui.KeyMap, svc service.OmniFocusService) Model {
	return Model{
		projectList: projectlist.New(styles, keys),
		taskList:    tasklist.New(styles, keys),
		service:     svc,
		styles:      styles,
		keys:        keys,
		mode:        ModeProjectList,
		loaded:      false,
	}
}

// Init initializes the projects view
func (m Model) Init() tea.Cmd {
	return m.loadProjects()
}

func (m Model) loadProjects() tea.Cmd {
	return func() tea.Msg {
		projects, err := m.service.GetProjects("")
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.ProjectsLoadedMsg{Projects: projects}
	}
}

func (m Model) loadProjectTasks(projectID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.service.GetTasksByProject(projectID)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TasksLoadedMsg{Tasks: tasks}
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.ProjectsLoadedMsg:
		m.projectList = m.projectList.SetProjects(msg.Projects)
		m.loaded = true
		m.err = nil
		return m, nil

	case tui.TasksLoadedMsg:
		m.taskList = m.taskList.SetTasks(msg.Tasks)
		return m, nil

	case tui.ErrorMsg:
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 2
		availableHeight := msg.Height - headerHeight
		if availableHeight < 0 {
			availableHeight = 0
		}

		subMsg := tea.WindowSizeMsg{Width: msg.Width, Height: availableHeight}

		var cmd tea.Cmd
		if m.mode == ModeProjectList {
			m.projectList, cmd = m.projectList.Update(subMsg)
		} else {
			m.taskList, cmd = m.taskList.Update(subMsg)
		}
		return m, cmd

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	default:
		return m.delegateToCurrentList(msg)
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle drill-down with Enter
	if key.Matches(msg, enterKey) {
		if m.mode == ModeProjectList {
			project := m.projectList.SelectedProject()
			if project != nil {
				m.mode = ModeProjectTasks
				m.currentProject = project
				return m, m.loadProjectTasks(project.ID)
			}
		}
		return m, nil
	}

	// Handle back navigation with h or Escape
	if key.Matches(msg, backKey) || key.Matches(msg, escapeKey) {
		if m.mode == ModeProjectTasks {
			m.mode = ModeProjectList
			m.currentProject = nil
			return m, nil
		}
		return m, nil
	}

	// Delegate to current list
	return m.delegateToCurrentList(msg)
}

func (m Model) delegateToCurrentList(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.mode == ModeProjectList {
		m.projectList, cmd = m.projectList.Update(msg)
	} else {
		m.taskList, cmd = m.taskList.Update(msg)
	}
	return m, cmd
}

// View renders the projects view
func (m Model) View() string {
	if m.err != nil {
		return m.renderError()
	}

	header := m.renderHeader()

	var content string
	if m.mode == ModeProjectList {
		content = m.projectList.View()
	} else {
		content = m.taskList.View()
	}

	return header + "\n" + content
}

func (m Model) renderHeader() string {
	var headerText string
	if m.mode == ModeProjectList {
		count := len(m.projectList.Projects())
		headerText = fmt.Sprintf("PROJECTS (%d)", count)
	} else if m.currentProject != nil {
		headerText = fmt.Sprintf("📁 %s", m.currentProject.Name)
	} else {
		headerText = "PROJECT TASKS"
	}

	styled := m.styles.UI.Header.Render(headerText)

	// Add back hint when in drill-down mode
	if m.mode == ModeProjectTasks {
		hint := m.styles.UI.Help.Render("  [h/Esc] back")
		styled += hint
	}

	return styled
}

func (m Model) renderError() string {
	header := m.styles.UI.Header.Render("PROJECTS")
	separatorWidth := m.width
	if separatorWidth == 0 {
		separatorWidth = 40
	}
	separator := strings.Repeat("─", separatorWidth)
	errorText := fmt.Sprintf("Error: %v", m.err)
	errorStyle := m.styles.UI.Help.Foreground(m.styles.Colors.Error)
	errorStyled := errorStyle.Render(errorText)
	return header + "\n" + separator + "\n" + errorStyled
}

// SelectedTask returns the currently selected task (when in task mode)
func (m Model) SelectedTask() *domain.Task {
	if m.mode == ModeProjectTasks {
		return m.taskList.SelectedTask()
	}
	return nil
}

// Refresh reloads projects
func (m Model) Refresh() tea.Cmd {
	if m.mode == ModeProjectTasks && m.currentProject != nil {
		return m.loadProjectTasks(m.currentProject.ID)
	}
	return m.loadProjects()
}

// Mode returns the current view mode
func (m Model) Mode() ViewMode {
	return m.mode
}

var (
	enterKey  = key.NewBinding(key.WithKeys("enter"))
	backKey   = key.NewBinding(key.WithKeys("h", "left"))
	escapeKey = key.NewBinding(key.WithKeys("esc", "escape"))
)
