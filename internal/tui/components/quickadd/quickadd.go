// Package quickadd provides a quick add overlay component for creating tasks in the TUI.
package quickadd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/cli/taskparse"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// Model represents the quick add overlay component state
type Model struct {
	textInput textinput.Model
	visible   bool
	width     int
	height    int
	styles    *tui.Styles
	err       error
	service   service.OmniFocusService
}

// New creates a new quick add overlay component
func New(styles *tui.Styles, svc service.OmniFocusService) Model {
	ti := textinput.New()
	ti.Placeholder = "Add task (e.g., Buy milk #groceries due:tomorrow)"
	ti.CharLimit = 256
	ti.Width = 60

	return Model{
		textInput: ti,
		visible:   false,
		styles:    styles,
		service:   svc,
	}
}

// Init initializes the component (Bubble Tea interface)
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and returns updated model (Bubble Tea interface)
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			// Cancel and close
			return m.Hide(), nil

		case tea.KeyEnter:
			// Submit task
			return m.submitTask()

		default:
			// Pass through to text input
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	}

	// Pass through to text input for other messages
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the component (Bubble Tea interface)
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal dimensions
	modalWidth := min(70, m.width-4)
	modalHeight := 8

	// Build content
	var content string

	// Title
	title := m.styles.UI.Header.
		Width(modalWidth - 4).
		Align(lipgloss.Center).
		Render("Quick Add Task")
	content += title + "\n\n"

	// Input field with border
	inputStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(m.styles.Colors.Primary).
		Padding(0, 1).
		Width(modalWidth - 4)

	input := inputStyle.Render(m.textInput.View())
	content += input + "\n"

	// Error display (if any)
	if m.err != nil {
		errorStyle := lipgloss.NewStyle().
			Foreground(m.styles.Colors.Error).
			Width(modalWidth - 4)
		content += errorStyle.Render(fmt.Sprintf("Error: %s", m.err.Error())) + "\n"
	} else {
		content += "\n"
	}

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary).
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	help := helpStyle.Render("Enter: add • Escape: cancel")
	content += help

	// Wrap in overlay style
	overlay := m.styles.UI.Overlay.
		Width(modalWidth).
		Height(modalHeight).
		Render(content)

	// Center the modal
	return m.centerModal(overlay)
}

// Show makes the component visible and focuses the input
func (m Model) Show() Model {
	m.visible = true
	m.err = nil
	m.textInput.Focus()
	return m
}

// Hide makes the component invisible and clears the input
func (m Model) Hide() Model {
	m.visible = false
	m.err = nil
	m.textInput.SetValue("")
	m.textInput.Blur()
	return m
}

// IsVisible returns whether the component is currently visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the component's dimensions for layout calculations
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// submitTask parses the input and creates a task
func (m Model) submitTask() (Model, tea.Cmd) {
	input := m.textInput.Value()

	// Parse the input using natural syntax parser
	taskInput, err := taskparse.Parse(input)
	if err != nil {
		m.err = err
		return m, func() tea.Msg {
			return tui.ErrorMsg{Err: err}
		}
	}

	// Resolve project name to ID if specified
	if taskInput.ProjectName != "" {
		projectID, err := m.service.ResolveProjectName(taskInput.ProjectName)
		if err != nil {
			m.err = err
			return m, func() tea.Msg {
				return tui.ErrorMsg{Err: err}
			}
		}
		taskInput.ProjectID = projectID
	}

	// Create the task
	task, err := m.service.CreateTask(taskInput)
	if err != nil {
		m.err = err
		return m, func() tea.Msg {
			return tui.ErrorMsg{Err: err}
		}
	}

	// Success - hide the overlay and return success message
	m = m.Hide()
	return m, func() tea.Msg {
		return tui.TaskCreatedMsg{Task: *task}
	}
}

// centerModal centers the modal content on the screen
func (m Model) centerModal(content string) string {
	// Calculate vertical padding
	lines := lipgloss.Height(content)
	verticalPad := (m.height - lines) / 2
	if verticalPad < 0 {
		verticalPad = 0
	}

	// Calculate horizontal padding
	contentWidth := lipgloss.Width(content)
	horizontalPad := (m.width - contentWidth) / 2
	if horizontalPad < 0 {
		horizontalPad = 0
	}

	// Apply padding
	return lipgloss.NewStyle().
		Padding(verticalPad, 0, 0, horizontalPad).
		Render(content)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
