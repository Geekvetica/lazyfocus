// Package app provides the main TUI application model and orchestration.
package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/quickadd"
	"github.com/pwojciechowski/lazyfocus/internal/tui/views/inbox"
)

// Model represents the main TUI application state
type Model struct {
	// Views
	inboxView   inbox.Model
	currentView int // tui.ViewInbox, tui.ViewProjects, etc from messages.go

	// Overlays
	quickAdd quickadd.Model
	showHelp bool

	// State
	service service.OmniFocusService
	styles  *tui.Styles
	keys    tui.KeyMap
	width   int
	height  int
	err     error
	ready   bool // true after first WindowSizeMsg
}

// NewApp creates a new TUI application instance
func NewApp(svc service.OmniFocusService) Model {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	return Model{
		inboxView:   inbox.New(styles, keys, svc),
		currentView: tui.ViewInbox,
		quickAdd:    quickadd.New(styles, svc),
		showHelp:    false,
		service:     svc,
		styles:      styles,
		keys:        keys,
		ready:       false,
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return m.inboxView.Init()
}

// Update handles messages and updates the application state
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit immediately
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if key.Matches(keyMsg, m.keys.Quit) {
			return m, tea.Quit
		}
	}

	// Handle window resize
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

		// Update quick add size
		m.quickAdd = m.quickAdd.SetSize(msg.Width, msg.Height)

		// Pass resize to current view
		var cmd tea.Cmd
		m.inboxView, cmd = m.inboxView.Update(msg)
		return m, cmd
	}

	// Handle TaskCreatedMsg - hide quick add and refresh view
	// Must come before quick add delegation since quick add emits this message
	if msg, ok := msg.(tui.TaskCreatedMsg); ok {
		_ = msg // Task created successfully
		m.quickAdd = m.quickAdd.Hide()
		// Refresh the current view
		return m, m.inboxView.Refresh()
	}

	// Handle ErrorMsg
	if msg, ok := msg.(tui.ErrorMsg); ok {
		m.err = msg.Err
		return m, nil
	}

	// If quick add is visible, delegate to it
	if m.quickAdd.IsVisible() {
		var cmd tea.Cmd
		m.quickAdd, cmd = m.quickAdd.Update(msg)
		return m, cmd
	}

	// Handle global keys when overlay is not visible
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		// Toggle help
		if key.Matches(keyMsg, m.keys.Help) {
			m.showHelp = !m.showHelp
			return m, nil
		}

		// Show quick add
		if key.Matches(keyMsg, m.keys.QuickAdd) {
			m.quickAdd = m.quickAdd.Show()
			return m, nil
		}

		// View switching (only 1 is functional in Phase 4)
		if key.Matches(keyMsg, m.keys.View1) {
			m.currentView = tui.ViewInbox
			return m, nil
		}
		// Views 2-5 not implemented yet, stay on inbox
		if key.Matches(keyMsg, m.keys.View2) || key.Matches(keyMsg, m.keys.View3) ||
			key.Matches(keyMsg, m.keys.View4) || key.Matches(keyMsg, m.keys.View5) {
			// No-op in Phase 4
			return m, nil
		}
	}

	// Delegate to current view
	var cmd tea.Cmd
	m.inboxView, cmd = m.inboxView.Update(msg)
	return m, cmd
}

// View renders the application
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// Render current view
	var view string
	switch m.currentView {
	case tui.ViewInbox:
		view = m.inboxView.View()
	default:
		view = "View not implemented"
	}

	// Overlay quick add if visible
	if m.quickAdd.IsVisible() {
		quickAddView := m.quickAdd.View()
		// Layer quick add over the main view
		view = m.layerOverlay(view, quickAddView)
	}

	// Overlay help if visible
	if m.showHelp {
		helpView := m.renderHelp()
		view = m.layerOverlay(view, helpView)
	}

	return view
}

// CurrentViewName returns the name of the current view
func (m Model) CurrentViewName() string {
	switch m.currentView {
	case tui.ViewInbox:
		return "Inbox"
	case tui.ViewProjects:
		return "Projects"
	case tui.ViewTags:
		return "Tags"
	case tui.ViewForecast:
		return "Forecast"
	case tui.ViewReview:
		return "Review"
	default:
		return "Unknown"
	}
}

// renderHelp renders the help overlay
func (m Model) renderHelp() string {
	// Calculate modal dimensions
	modalWidth := min(60, m.width-4)

	// Build help content
	var content strings.Builder

	// Title
	title := m.styles.UI.Header.
		Width(modalWidth - 4).
		Align(lipgloss.Center).
		Render("lazyfocus - Keyboard Shortcuts")
	content.WriteString(title)
	content.WriteString("\n\n")

	// Navigation section
	content.WriteString(m.styles.UI.Header.
		Width(modalWidth - 4).
		Render("Navigation"))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Down.Help().Key, m.keys.Down.Help().Desc))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Up.Help().Key, m.keys.Up.Help().Desc))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine("1-5", "switch views"))
	content.WriteString("\n\n")

	// Actions section
	content.WriteString(m.styles.UI.Header.
		Width(modalWidth - 4).
		Render("Actions"))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.QuickAdd.Help().Key, m.keys.QuickAdd.Help().Desc))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Complete.Help().Key, m.keys.Complete.Help().Desc))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Delete.Help().Key, m.keys.Delete.Help().Desc))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Flag.Help().Key, m.keys.Flag.Help().Desc))
	content.WriteString("\n\n")

	// General section
	content.WriteString(m.styles.UI.Header.
		Width(modalWidth - 4).
		Render("General"))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Help.Help().Key, m.keys.Help.Help().Desc))
	content.WriteString("\n")
	content.WriteString(m.formatHelpLine(m.keys.Quit.Help().Key, m.keys.Quit.Help().Desc))
	content.WriteString("\n")

	// Wrap in overlay style
	overlay := m.styles.UI.Overlay.
		Width(modalWidth).
		Render(content.String())

	return m.centerOverlay(overlay)
}

// formatHelpLine formats a help line with key and description
func (m Model) formatHelpLine(key, desc string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Primary).
		Bold(true).
		Width(10)
	descStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary)

	return "  " + keyStyle.Render(key) + " " + descStyle.Render(desc)
}

// layerOverlay layers an overlay on top of the base view
func (m Model) layerOverlay(base, overlay string) string {
	// Simply append the overlay after the base for now
	// In a more sophisticated implementation, we would composite them
	return base + "\n" + overlay
}

// centerOverlay centers an overlay on the screen
func (m Model) centerOverlay(content string) string {
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
