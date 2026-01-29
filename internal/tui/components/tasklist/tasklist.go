package tasklist

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// Icons for task display
const (
	CheckboxEmpty   = "☐"
	CheckboxChecked = "☑"
	FlagIcon        = "🚩"
	CalendarIcon    = "📅"
)

// Model represents the task list component state
type Model struct {
	tasks   []domain.Task
	cursor  int
	width   int
	height  int
	styles  *tui.Styles
	keys    tui.KeyMap
	loading bool
	empty   bool
}

// New creates a new task list component
func New(styles *tui.Styles, keys tui.KeyMap) Model {
	return Model{
		tasks:   []domain.Task{},
		cursor:  0,
		styles:  styles,
		keys:    keys,
		loading: false,
		empty:   true,
	}
}

// Init initializes the component
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

// handleKeyPress processes keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if len(m.tasks) == 0 {
		return m, nil
	}

	// Handle down navigation
	if key.Matches(msg, m.keys.Down) {
		m.cursor++
		if m.cursor >= len(m.tasks) {
			m.cursor = 0
		}
		return m, nil
	}

	// Handle up navigation
	if key.Matches(msg, m.keys.Up) {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.tasks) - 1
		}
		return m, nil
	}

	return m, nil
}

// View renders the component
func (m Model) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.empty {
		return m.renderEmpty()
	}

	return m.renderTasks()
}

// renderLoading renders the loading state
func (m Model) renderLoading() string {
	if m.height == 0 {
		return "Loading..."
	}

	// Center the loading message
	padding := strings.Repeat("\n", m.height/2)
	return padding + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "Loading...")
}

// renderEmpty renders the empty state
func (m Model) renderEmpty() string {
	if m.height == 0 {
		return "No tasks"
	}

	// Center the empty message
	padding := strings.Repeat("\n", m.height/2)
	return padding + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "No tasks")
}

// renderTasks renders the task list
func (m Model) renderTasks() string {
	var b strings.Builder

	for i, task := range m.tasks {
		line := m.formatTaskLine(task, i == m.cursor)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

// formatTaskLine formats a single task line
func (m Model) formatTaskLine(task domain.Task, selected bool) string {
	// Status icon
	statusIcon := CheckboxEmpty
	if task.Completed {
		statusIcon = CheckboxChecked
	}

	// Build the left side (status icon + task name)
	leftSide := fmt.Sprintf("%s %s", statusIcon, task.Name)

	// Build the right side (due date or flag)
	var rightSide string
	if task.DueDate != nil {
		rightSide = fmt.Sprintf("%s %s", CalendarIcon, formatDate(*task.DueDate))
	} else if task.Flagged {
		rightSide = FlagIcon
	}

	// Calculate spacing to align right side
	// Account for the width of emoji characters (they take up more space visually)
	contentWidth := m.width
	if contentWidth == 0 {
		contentWidth = 80
	}

	// Calculate actual content length
	leftLen := len(statusIcon) + 1 + len(task.Name)
	rightLen := len(rightSide)

	// Add visual compensation for emoji width
	emojiCount := strings.Count(statusIcon, "☐") + strings.Count(statusIcon, "☑")
	if task.Flagged {
		emojiCount++
	}
	if task.DueDate != nil {
		emojiCount++
	}

	// Emoji characters typically take 2 display cells but count as 1-4 bytes
	// We need to compensate for the visual width
	visualAdjustment := emojiCount

	spacing := contentWidth - leftLen - rightLen - visualAdjustment - 2
	if spacing < 0 {
		spacing = 1
	}

	line := leftSide
	if rightSide != "" {
		line = leftSide + strings.Repeat(" ", spacing) + rightSide
	}

	// Apply styles
	if selected {
		return m.styles.Task.Selected.Render(line)
	}

	if task.Completed {
		return m.styles.Task.Completed.Render(line)
	}

	return m.styles.Task.Normal.Render(line)
}

// formatDate formats a time.Time into a human-readable string
func formatDate(t time.Time) string {
	now := time.Now()

	// Check if it's today
	if isSameDay(t, now) {
		return "Today"
	}

	// Check if it's tomorrow
	tomorrow := now.AddDate(0, 0, 1)
	if isSameDay(t, tomorrow) {
		return "Tomorrow"
	}

	// Check if it's yesterday
	yesterday := now.AddDate(0, 0, -1)
	if isSameDay(t, yesterday) {
		return "Yesterday"
	}

	// Check if it's within the same year
	if t.Year() == now.Year() {
		return t.Format("Jan 2")
	}

	return t.Format("Jan 2, 2006")
}

// isSameDay checks if two times are on the same calendar day
func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// SetTasks updates the task list
func (m Model) SetTasks(tasks []domain.Task) Model {
	m.tasks = tasks
	m.empty = len(tasks) == 0
	m.loading = false

	// Clamp cursor to valid range
	if m.cursor >= len(m.tasks) {
		if len(m.tasks) > 0 {
			m.cursor = len(m.tasks) - 1
		} else {
			m.cursor = 0
		}
	}

	return m
}

// SetLoading sets the loading state
func (m Model) SetLoading(loading bool) Model {
	m.loading = loading
	return m
}

// SelectedTask returns the currently selected task
func (m Model) SelectedTask() *domain.Task {
	if len(m.tasks) == 0 || m.cursor >= len(m.tasks) {
		return nil
	}

	return &m.tasks[m.cursor]
}

// SelectedIndex returns the current cursor position
func (m Model) SelectedIndex() int {
	return m.cursor
}
