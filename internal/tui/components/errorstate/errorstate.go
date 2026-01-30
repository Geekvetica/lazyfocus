// Package errorstate provides a reusable error state component.
package errorstate

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// ErrorDismissedMsg indicates the user dismissed an error
type ErrorDismissedMsg struct{}

// Model represents the error state component
type Model struct {
	err      error
	retryCmd tea.Cmd
	visible  bool
	styles   Styles
	width    int
	height   int
}

// Styles for the error state
type Styles struct {
	Container lipgloss.Style
	Title     lipgloss.Style
	Message   lipgloss.Style
	Hint      lipgloss.Style
}

// DefaultStyles returns default error state styles
func DefaultStyles() Styles {
	errorColor := lipgloss.AdaptiveColor{Light: "#C00000", Dark: "#FF6B6B"}
	secondaryColor := lipgloss.AdaptiveColor{Light: "#808080", Dark: "#A0A0A0"}

	return Styles{
		Container: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#2A2A2A"}),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(errorColor),
		Message: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}),
		Hint: lipgloss.NewStyle().
			Foreground(secondaryColor).
			Faint(true),
	}
}

// New creates a new error state model
func New() Model {
	return Model{
		styles:  DefaultStyles(),
		visible: false,
	}
}

// NewWithStyles creates a new error state model with custom styles
func NewWithStyles(styles *tui.Styles) Model {
	return Model{
		styles: Styles{
			Container: styles.UI.Overlay,
			Title: lipgloss.NewStyle().
				Bold(true).
				Foreground(styles.Colors.Error),
			Message: lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#FFFFFF"}),
			Hint: lipgloss.NewStyle().
				Foreground(styles.Colors.Secondary).
				Faint(true),
		},
		visible: false,
	}
}

// Show displays the error with optional retry command
func (m Model) Show(err error, retryCmd tea.Cmd) Model {
	m.err = err
	m.retryCmd = retryCmd
	m.visible = true
	return m
}

// Hide hides the error state
func (m Model) Hide() Model {
	m.visible = false
	return m
}

// IsVisible returns whether the error is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the dimensions for the error state
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// Init initializes the component
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles key events
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		switch keyMsg.Type {
		case tea.KeyEsc, tea.KeyEnter:
			m.visible = false
			return m, func() tea.Msg {
				return ErrorDismissedMsg{}
			}
		case tea.KeyRunes:
			if len(keyMsg.Runes) > 0 && keyMsg.Runes[0] == 'r' {
				if m.retryCmd != nil {
					m.visible = false
					return m, m.retryCmd
				}
			}
		}
	}

	return m, nil
}

// View renders the error state
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Calculate modal width
	modalWidth := min(60, m.width-4)
	if modalWidth < 30 {
		modalWidth = 30
	}

	// Build content
	var content string

	// Title
	titleStyle := m.styles.Title.
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	content = titleStyle.Render("Error") + "\n\n"

	// Error message
	messageStyle := m.styles.Message.
		Width(modalWidth - 4).
		Align(lipgloss.Left)
	errorMessage := "An unknown error occurred"
	if m.err != nil {
		errorMessage = m.err.Error()
	}
	content += messageStyle.Render(errorMessage) + "\n\n"

	// Build hints
	var hints string
	if m.retryCmd != nil {
		hints = "[r] Retry  "
	}
	hints += "[Enter/Esc] Dismiss"

	hintStyle := m.styles.Hint.
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	content += hintStyle.Render(hints)

	// Wrap in container
	return m.styles.Container.
		Width(modalWidth).
		Render(content)
}
