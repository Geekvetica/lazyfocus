// Package commandinput provides a command input component for the TUI.
package commandinput

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/command"
)

// CommandExecutedMsg is sent when a command is executed
type CommandExecutedMsg struct {
	Command *command.Command
}

// CommandCancelledMsg is sent when command input is cancelled
type CommandCancelledMsg struct{}

// CommandErrorMsg is sent when a command parsing error occurs
type CommandErrorMsg struct {
	Error string
}

// Model represents the command input state
type Model struct {
	input       textinput.Model
	parser      *command.Parser
	visible     bool
	styles      *tui.Styles
	history     []string
	historyIdx  int
	width       int
	completions []string
	compIdx     int
}

// New creates a new command input
func New(styles *tui.Styles) Model {
	ti := textinput.New()
	ti.Prompt = ":"
	ti.CharLimit = 200

	return Model{
		input:      ti,
		parser:     command.NewParser(),
		visible:    false,
		styles:     styles,
		history:    []string{},
		historyIdx: -1,
	}
}

// Show makes the command input visible
func (m Model) Show() Model {
	m.visible = true
	m.input.Focus()
	m.input.SetValue("")
	m.historyIdx = -1
	m.completions = nil
	m.compIdx = 0
	return m
}

// Hide hides the command input
func (m Model) Hide() Model {
	m.visible = false
	m.input.Blur()
	return m
}

// IsVisible returns true if the input is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetWidth sets the width for the input
func (m Model) SetWidth(width int) Model {
	m.width = width
	m.input.Width = width - 4
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
			m = m.Hide()
			return m, func() tea.Msg { return CommandCancelledMsg{} }

		case key.Matches(msg, enterKey):
			input := m.input.Value()
			if input != "" {
				// Add to history
				m.history = append(m.history, input)

				// Parse and execute
				cmd, err := m.parser.Parse(input)
				if err != nil {
					m = m.Hide()
					errStr := err.Error()
					return m, func() tea.Msg { return CommandErrorMsg{Error: errStr} }
				}

				m = m.Hide()
				return m, func() tea.Msg { return CommandExecutedMsg{Command: cmd} }
			}
			m = m.Hide()
			return m, func() tea.Msg { return CommandCancelledMsg{} }

		case key.Matches(msg, upKey):
			// Navigate history backward
			if len(m.history) > 0 {
				if m.historyIdx < 0 {
					m.historyIdx = len(m.history) - 1
				} else if m.historyIdx > 0 {
					m.historyIdx--
				}
				m.input.SetValue(m.history[m.historyIdx])
				m.input.CursorEnd()
			}
			return m, nil

		case key.Matches(msg, downKey):
			// Navigate history forward
			if m.historyIdx >= 0 {
				m.historyIdx++
				if m.historyIdx >= len(m.history) {
					m.historyIdx = -1
					m.input.SetValue("")
				} else {
					m.input.SetValue(m.history[m.historyIdx])
					m.input.CursorEnd()
				}
			}
			return m, nil

		case key.Matches(msg, tabKey):
			// Tab completion
			text := m.input.Value()
			completions := m.parser.GetCompletions(text)
			if len(completions) == 1 {
				m.input.SetValue(completions[0])
				m.input.CursorEnd()
			} else if len(completions) > 1 {
				// Cycle through completions
				m.completions = completions
				m.compIdx = (m.compIdx + 1) % len(completions)
				m.input.SetValue(completions[m.compIdx])
				m.input.CursorEnd()
			}
			return m, nil
		}
	}

	// Update text input
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View renders the command input
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	inputStyle := lipgloss.NewStyle().
		Background(m.styles.Colors.Primary).
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Width(m.width)

	// Show completions hint if available
	var content strings.Builder
	content.WriteString(m.input.View())

	if len(m.completions) > 1 {
		hint := strings.Join(m.completions, " | ")
		hintStyle := lipgloss.NewStyle().
			Foreground(m.styles.Colors.Secondary).
			Faint(true)
		content.WriteString("\n")
		content.WriteString(hintStyle.Render(hint))
	}

	return inputStyle.Render(content.String())
}

var (
	escapeKey = key.NewBinding(key.WithKeys("esc", "escape"))
	enterKey  = key.NewBinding(key.WithKeys("enter"))
	upKey     = key.NewBinding(key.WithKeys("up"))
	downKey   = key.NewBinding(key.WithKeys("down"))
	tabKey    = key.NewBinding(key.WithKeys("tab"))
)
