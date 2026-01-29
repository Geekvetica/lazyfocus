// Package inbox provides the inbox view for the TUI.
package inbox

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/tasklist"
)

// Model represents the inbox view state
type Model struct {
	taskList  tasklist.Model
	service   service.OmniFocusService
	styles    *tui.Styles
	keys      tui.KeyMap
	width     int
	height    int
	err       error
	loaded    bool
	taskCount int
}

// New creates a new inbox view
func New(styles *tui.Styles, keys tui.KeyMap, svc service.OmniFocusService) Model {
	return Model{
		taskList:  tasklist.New(styles, keys),
		service:   svc,
		styles:    styles,
		keys:      keys,
		loaded:    false,
		taskCount: 0,
	}
}

// Init initializes the inbox view
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.taskList.Init(),
		m.loadTasks(),
	)
}

// loadTasks loads tasks from the OmniFocus service
func (m Model) loadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.service.GetInboxTasks()
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TasksLoadedMsg{Tasks: tasks}
	}
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.TasksLoadedMsg:
		// Update task list with loaded tasks
		m.taskList = m.taskList.SetTasks(msg.Tasks)
		m.taskCount = len(msg.Tasks)
		m.loaded = true
		m.err = nil
		return m, nil

	case tui.ErrorMsg:
		// Store error for display
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		// Update dimensions
		m.width = msg.Width
		m.height = msg.Height

		// Pass resize to task list
		// Calculate available height for task list (subtract header height)
		headerHeight := 2 // Header + border
		availableHeight := msg.Height - headerHeight
		if availableHeight < 0 {
			availableHeight = 0
		}

		taskListMsg := tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: availableHeight,
		}
		var cmd tea.Cmd
		m.taskList, cmd = m.taskList.Update(taskListMsg)
		return m, cmd

	default:
		// Delegate other messages to task list
		var cmd tea.Cmd
		m.taskList, cmd = m.taskList.Update(msg)
		return m, cmd
	}
}

// View renders the inbox view
func (m Model) View() string {
	// Show error if present
	if m.err != nil {
		return m.renderError()
	}

	// Render header
	header := m.renderHeader()

	// Render task list
	taskListView := m.taskList.View()

	return header + "\n" + taskListView
}

// renderHeader renders the inbox header with task count
func (m Model) renderHeader() string {
	taskCount := m.TaskCount()
	headerText := fmt.Sprintf("INBOX (%d tasks)", taskCount)

	// Apply header style
	styled := m.styles.UI.Header.Render(headerText)

	return styled
}

// renderError renders the error view
func (m Model) renderError() string {
	header := m.styles.UI.Header.Render("INBOX")

	// Calculate separator width (default to 40 if width not set)
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

// TaskCount returns the number of tasks in the inbox
func (m Model) TaskCount() int {
	return m.taskCount
}

// SelectedTask returns the currently selected task
func (m Model) SelectedTask() *domain.Task {
	return m.taskList.SelectedTask()
}

// Refresh reloads tasks from the service
func (m Model) Refresh() tea.Cmd {
	return m.loadTasks()
}
