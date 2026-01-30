// Package confirm provides a reusable confirmation modal component.
package confirm

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// ConfirmedMsg is sent when the user confirms the action
type ConfirmedMsg struct {
	Context interface{} // Optional context passed through from Show()
}

// CancelledMsg is sent when the user cancels the action
type CancelledMsg struct{}

// Model represents the confirmation modal state
type Model struct {
	title   string
	message string
	context interface{}
	visible bool
	styles  *tui.Styles
	width   int
	height  int
}

// New creates a new confirmation modal
func New(styles *tui.Styles) Model {
	return Model{
		styles:  styles,
		visible: false,
	}
}

// Show makes the modal visible with the given title and message
func (m Model) Show(title, message string) Model {
	m.title = title
	m.message = message
	m.visible = true
	m.context = nil
	return m
}

// ShowWithContext makes the modal visible with context that will be passed to ConfirmedMsg
func (m Model) ShowWithContext(title, message string, context interface{}) Model {
	m.title = title
	m.message = message
	m.visible = true
	m.context = context
	return m
}

// Hide hides the modal
func (m Model) Hide() Model {
	m.visible = false
	return m
}

// IsVisible returns true if the modal is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the dimensions for the modal
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// Init initializes the component
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, confirmKey):
			m.visible = false
			return m, func() tea.Msg {
				return ConfirmedMsg{Context: m.context}
			}
		case key.Matches(msg, cancelKey):
			m.visible = false
			return m, func() tea.Msg {
				return CancelledMsg{}
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the modal
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal width
	modalWidth := min(50, m.width-4)
	if modalWidth < 20 {
		modalWidth = 20
	}

	// Build modal content
	var content string

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Colors.Warning).
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	content = titleStyle.Render(m.title) + "\n\n"

	// Message
	messageStyle := lipgloss.NewStyle().
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	content += messageStyle.Render(m.message) + "\n\n"

	// Buttons hint
	hintStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary).
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	content += hintStyle.Render("[y/Enter] Confirm  [n/Esc] Cancel")

	// Wrap in overlay style
	return m.styles.UI.Overlay.
		Width(modalWidth).
		Render(content)
}

// Key bindings for the confirmation modal
var (
	confirmKey = key.NewBinding(
		key.WithKeys("y", "Y", "enter"),
	)
	cancelKey = key.NewBinding(
		key.WithKeys("n", "N", "esc"),
	)
)
