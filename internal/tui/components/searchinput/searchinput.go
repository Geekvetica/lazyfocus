// Package searchinput provides a search input component for the TUI.
package searchinput

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// SearchChangedMsg is sent when search text changes
type SearchChangedMsg struct {
	Text string
}

// SearchClearedMsg is sent when search is cleared
type SearchClearedMsg struct{}

// SearchConfirmedMsg is sent when search is confirmed (Enter)
type SearchConfirmedMsg struct {
	Text string
}

// Model represents the search input state
type Model struct {
	input   textinput.Model
	visible bool
	styles  *tui.Styles
	width   int
}

// New creates a new search input
func New(styles *tui.Styles) Model {
	ti := textinput.New()
	ti.Placeholder = "Search tasks..."
	ti.Prompt = "/ "
	ti.CharLimit = 100

	return Model{
		input:   ti,
		visible: false,
		styles:  styles,
	}
}

// Show makes the search input visible
func (m Model) Show() Model {
	m.visible = true
	m.input.Focus()
	m.input.SetValue("")
	return m
}

// Hide hides the search input
func (m Model) Hide() Model {
	m.visible = false
	m.input.Blur()
	return m
}

// IsVisible returns true if the input is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// Value returns the current search text
func (m Model) Value() string {
	return m.input.Value()
}

// SetWidth sets the width for the input
func (m Model) SetWidth(width int) Model {
	m.width = width
	m.input.Width = width - 4 // Account for prompt and padding
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
		case key.Matches(msg, escapeKey):
			m.visible = false
			m.input.Blur()
			m.input.SetValue("")
			return m, func() tea.Msg { return SearchClearedMsg{} }

		case key.Matches(msg, enterKey):
			m.visible = false
			m.input.Blur()
			text := m.input.Value()
			return m, func() tea.Msg { return SearchConfirmedMsg{Text: text} }
		}
	}

	// Update text input
	prevValue := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Emit change event if text changed
	newValue := m.input.Value()
	if newValue != prevValue {
		return m, tea.Batch(cmd, func() tea.Msg {
			return SearchChangedMsg{Text: newValue}
		})
	}

	return m, cmd
}

// View renders the search input
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	// Render at bottom of screen
	inputStyle := lipgloss.NewStyle().
		Background(m.styles.Colors.Primary).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width(m.width)

	return inputStyle.Render(m.input.View())
}

var (
	escapeKey = key.NewBinding(key.WithKeys("esc", "escape"))
	enterKey  = key.NewBinding(key.WithKeys("enter"))
)
