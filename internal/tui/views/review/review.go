// Package review provides the review view for the TUI.
package review

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/tasklist"
)

// Model represents the review view state
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

// New creates a new review view
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

// Init initializes the review view
func (m Model) Init() tea.Cmd {
	return m.loadFlaggedTasks()
}

func (m Model) loadFlaggedTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.service.GetFlaggedTasks()
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TasksLoadedMsg{Tasks: tasks}
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tui.TasksLoadedMsg:
		m.taskList = m.taskList.SetTasks(msg.Tasks)
		m.taskCount = len(msg.Tasks)
		m.loaded = true
		m.err = nil
		return m, nil

	case tui.ErrorMsg:
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 3 // Header + subtext
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
		var cmd tea.Cmd
		m.taskList, cmd = m.taskList.Update(msg)
		return m, cmd
	}
}

// View renders the review view
func (m Model) View() string {
	if m.err != nil {
		return m.renderError()
	}

	header := m.renderHeader()
	content := m.taskList.View()

	return header + "\n" + content
}

func (m Model) renderHeader() string {
	headerText := fmt.Sprintf("REVIEW - Flagged Tasks (%d)", m.taskCount)
	styled := m.styles.UI.Header.Render(headerText)

	// Add subtext
	subtext := m.styles.UI.Help.Render("Review flagged tasks: [c]omplete, [d]elete, [f]unflag")

	return styled + "\n" + subtext
}

func (m Model) renderError() string {
	header := m.styles.UI.Header.Render("REVIEW")
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

// SelectedTask returns the currently selected task
func (m Model) SelectedTask() *domain.Task {
	return m.taskList.SelectedTask()
}

// TaskCount returns the number of flagged tasks
func (m Model) TaskCount() int {
	return m.taskCount
}

// Refresh reloads flagged tasks
func (m Model) Refresh() tea.Cmd {
	return m.loadFlaggedTasks()
}
