// Package taskdetail provides a task detail view component.
package taskdetail

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// CloseMsg signals the task detail view should be closed.
type CloseMsg struct{}

// EditRequestedMsg signals the user wants to edit the task.
type EditRequestedMsg struct{ Task domain.Task }

// CompleteRequestedMsg signals the user wants to complete the task.
type CompleteRequestedMsg struct{ TaskID string }

// DeleteRequestedMsg signals the user wants to delete the task.
type DeleteRequestedMsg struct{ TaskID, TaskName string }

// FlagRequestedMsg signals the user wants to toggle the task flag.
type FlagRequestedMsg struct {
	TaskID  string
	Flagged bool
}

// Model represents the task detail view state
type Model struct {
	task     *domain.Task
	visible  bool
	styles   *tui.Styles
	keys     tui.KeyMap
	viewport viewport.Model
	width    int
	height   int
	ready    bool
}

// New creates a new task detail view
func New(styles *tui.Styles, keys tui.KeyMap) Model {
	return Model{
		styles:  styles,
		keys:    keys,
		visible: false,
	}
}

// Show displays the task detail view with the given task
func (m Model) Show(task *domain.Task) Model {
	m.task = task
	m.visible = true
	m.ready = false
	return m
}

// Hide closes the task detail view
func (m Model) Hide() Model {
	m.visible = false
	m.task = nil
	return m
}

// IsVisible returns true if the view is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// Task returns the current task being displayed
func (m Model) Task() *domain.Task {
	return m.task
}

// SetSize updates the dimensions
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height

	// Update viewport size
	modalWidth := min(70, width-4)
	modalHeight := min(20, height-4)

	if m.visible && m.ready {
		m.viewport.Width = modalWidth - 4
		m.viewport.Height = modalHeight - 6 // Account for header and footer
	}

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
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	// Pass to viewport for scrolling
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.task == nil {
		return m, nil
	}

	switch {
	// Close on Escape
	case key.Matches(msg, escapeKey):
		m.visible = false
		return m, func() tea.Msg { return CloseMsg{} }

	// Edit task
	case key.Matches(msg, m.keys.Edit):
		return m, func() tea.Msg { return EditRequestedMsg{Task: *m.task} }

	// Complete task
	case key.Matches(msg, m.keys.Complete):
		return m, func() tea.Msg { return CompleteRequestedMsg{TaskID: m.task.ID} }

	// Delete task
	case key.Matches(msg, m.keys.Delete):
		return m, func() tea.Msg {
			return DeleteRequestedMsg{TaskID: m.task.ID, TaskName: m.task.Name}
		}

	// Toggle flag
	case key.Matches(msg, m.keys.Flag):
		return m, func() tea.Msg {
			return FlagRequestedMsg{TaskID: m.task.ID, Flagged: !m.task.Flagged}
		}

	// Scroll down
	case key.Matches(msg, m.keys.Down):
		m.viewport.ScrollDown(1)
		return m, nil

	// Scroll up
	case key.Matches(msg, m.keys.Up):
		m.viewport.ScrollUp(1)
		return m, nil
	}

	return m, nil
}

// View renders the task detail view
func (m Model) View() string {
	if !m.visible || m.task == nil {
		return ""
	}

	modalWidth := min(70, m.width-4)
	if modalWidth < 30 {
		modalWidth = 30
	}

	// Build content
	content := m.buildContent(modalWidth - 4)

	// Initialize viewport if not ready
	if !m.ready {
		modalHeight := min(20, m.height-4)
		m.viewport = viewport.New(modalWidth-4, modalHeight-6)
		m.viewport.SetContent(content)
		m.ready = true
	} else {
		m.viewport.SetContent(content)
	}

	// Header
	header := m.renderHeader(modalWidth - 4)

	// Footer with actions
	footer := m.renderFooter(modalWidth - 4)

	// Combine
	fullContent := header + "\n" + m.viewport.View() + "\n" + footer

	return m.styles.UI.Overlay.
		Width(modalWidth).
		Render(fullContent)
}

func (m Model) renderHeader(width int) string {
	// Status icon
	statusIcon := "☐"
	if m.task.Completed {
		statusIcon = "☑"
	}

	// Flag icon
	flagIcon := ""
	if m.task.Flagged {
		flagIcon = " 🚩"
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Colors.Primary).
		Width(width)

	title := fmt.Sprintf("%s %s%s", statusIcon, m.task.Name, flagIcon)

	return titleStyle.Render(title)
}

func (m Model) buildContent(width int) string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary).
		Width(12)
	valueStyle := lipgloss.NewStyle().
		Width(width - 14)

	// Project
	if m.task.ProjectName != "" {
		b.WriteString(labelStyle.Render("Project:"))
		b.WriteString(valueStyle.Render(m.task.ProjectName))
		b.WriteString("\n")
	}

	// Tags
	if len(m.task.Tags) > 0 {
		b.WriteString(labelStyle.Render("Tags:"))
		b.WriteString(valueStyle.Render(strings.Join(m.task.Tags, ", ")))
		b.WriteString("\n")
	}

	// Due Date
	if m.task.DueDate != nil {
		b.WriteString(labelStyle.Render("Due:"))
		b.WriteString(m.formatDueDate(*m.task.DueDate, valueStyle))
		b.WriteString("\n")
	}

	// Defer Date
	if m.task.DeferDate != nil {
		b.WriteString(labelStyle.Render("Defer:"))
		b.WriteString(valueStyle.Render(formatDateTime(*m.task.DeferDate)))
		b.WriteString("\n")
	}

	// Completed Date
	if m.task.Completed && m.task.CompletedDate != nil {
		b.WriteString(labelStyle.Render("Completed:"))
		b.WriteString(valueStyle.Render(formatDateTime(*m.task.CompletedDate)))
		b.WriteString("\n")
	}

	// Note
	if m.task.Note != "" {
		b.WriteString("\n")
		b.WriteString(labelStyle.Render("Note:"))
		b.WriteString("\n")
		noteStyle := lipgloss.NewStyle().
			Width(width).
			Foreground(m.styles.Colors.Secondary)
		b.WriteString(noteStyle.Render(m.task.Note))
	}

	return b.String()
}

func (m Model) formatDueDate(t time.Time, style lipgloss.Style) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	dateStr := formatDateTime(t)

	if t.Before(today) {
		return m.styles.DueDate.Overdue.Render(dateStr + " (Overdue)")
	}

	tomorrow := today.AddDate(0, 0, 1)
	if t.Before(tomorrow) {
		return m.styles.DueDate.Today.Render(dateStr + " (Today)")
	}

	return style.Render(dateStr)
}

func (m Model) renderFooter(width int) string {
	hintStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary).
		Width(width).
		Align(lipgloss.Center)

	hints := "[e]dit  [c]omplete  [d]elete  [f]lag  [Esc] close"
	return hintStyle.Render(hints)
}

// Helper function
func formatDateTime(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() {
		return t.Format("Jan 2 at 3:04 PM")
	}
	return t.Format("Jan 2, 2006 at 3:04 PM")
}

var escapeKey = key.NewBinding(key.WithKeys("esc", "escape"))
