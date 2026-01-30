// Package forecast provides the forecast view for the TUI.
package forecast

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/filter"
)

// DueGroup represents a group of tasks by due date category
type DueGroup int

// DueGroup constants for grouping tasks by due date category.
const (
	GroupOverdue DueGroup = iota
	GroupToday
	GroupTomorrow
	GroupThisWeek
	GroupLater
	GroupNoDue
)

// GroupedTask wraps a task with its group info
type GroupedTask struct {
	Task     domain.Task
	Group    DueGroup
	IsHeader bool // True if this is a group header, not a task
}

// Model represents the forecast view state
type Model struct {
	items     []GroupedTask
	cursor    int
	service   service.OmniFocusService
	styles    *tui.Styles
	keys      tui.KeyMap
	filter    filter.State
	width     int
	height    int
	err       error
	loaded    bool
	collapsed map[DueGroup]bool // Track collapsed groups
	allTasks  []domain.Task     // Store all tasks for filtering
}

// New creates a new forecast view
func New(styles *tui.Styles, keys tui.KeyMap, svc service.OmniFocusService) Model {
	return Model{
		items:     []GroupedTask{},
		cursor:    0,
		service:   svc,
		styles:    styles,
		keys:      keys,
		collapsed: make(map[DueGroup]bool),
		loaded:    false,
	}
}

// Init initializes the forecast view
func (m Model) Init() tea.Cmd {
	return m.loadTasks()
}

func (m Model) loadTasks() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.service.GetAllTasks(service.TaskFilters{})
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
		// Store all tasks and apply filter
		m.allTasks = msg.Tasks
		filteredTasks := m.applyFilter(msg.Tasks)
		m.items = m.groupTasks(filteredTasks)
		m.loaded = true
		m.err = nil
		// Move cursor to first task (skip header)
		if len(m.items) > 0 && m.items[0].IsHeader && len(m.items) > 1 {
			m.cursor = 1
		}
		return m, nil

	case tui.ErrorMsg:
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if len(m.items) == 0 {
		return m, nil
	}

	// Navigation
	if key.Matches(msg, m.keys.Down) {
		m.cursor = m.nextSelectableIndex(m.cursor, 1)
		return m, nil
	}
	if key.Matches(msg, m.keys.Up) {
		m.cursor = m.nextSelectableIndex(m.cursor, -1)
		return m, nil
	}

	// Toggle group collapse on Enter when on header
	if key.Matches(msg, enterKey) {
		if m.cursor < len(m.items) && m.items[m.cursor].IsHeader {
			group := m.items[m.cursor].Group
			m.collapsed[group] = !m.collapsed[group]
			m.items = m.rebuildItems()
			return m, nil
		}
	}

	return m, nil
}

// nextSelectableIndex finds the next selectable item (skips headers optionally)
func (m Model) nextSelectableIndex(current, direction int) int {
	next := current + direction
	if next < 0 {
		next = len(m.items) - 1
	} else if next >= len(m.items) {
		next = 0
	}
	// Allow selecting headers for collapse toggle
	return next
}

func (m Model) groupTasks(tasks []domain.Task) []GroupedTask {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.AddDate(0, 0, 1)
	weekEnd := today.AddDate(0, 0, 7)

	groups := map[DueGroup][]domain.Task{
		GroupOverdue:  {},
		GroupToday:    {},
		GroupTomorrow: {},
		GroupThisWeek: {},
		GroupLater:    {},
		GroupNoDue:    {},
	}

	for _, task := range tasks {
		if task.Completed {
			continue
		}

		group := m.categorizeTask(task, today, tomorrow, weekEnd)
		groups[group] = append(groups[group], task)
	}

	return m.buildGroupedItems(groups)
}

func (m Model) categorizeTask(task domain.Task, today, tomorrow, weekEnd time.Time) DueGroup {
	if task.DueDate == nil {
		return GroupNoDue
	}

	due := *task.DueDate
	if due.Before(today) {
		return GroupOverdue
	}
	if due.Before(tomorrow) {
		return GroupToday
	}
	dayAfterTomorrow := tomorrow.AddDate(0, 0, 1)
	if due.Before(dayAfterTomorrow) {
		return GroupTomorrow
	}
	if due.Before(weekEnd) {
		return GroupThisWeek
	}
	return GroupLater
}

func (m Model) buildGroupedItems(groups map[DueGroup][]domain.Task) []GroupedTask {
	var items []GroupedTask

	groupOrder := []DueGroup{GroupOverdue, GroupToday, GroupTomorrow, GroupThisWeek, GroupLater, GroupNoDue}

	for _, group := range groupOrder {
		tasks := groups[group]
		if len(tasks) == 0 {
			continue
		}

		// Add header
		items = append(items, GroupedTask{
			Group:    group,
			IsHeader: true,
		})

		// Add tasks if not collapsed
		if !m.collapsed[group] {
			for _, task := range tasks {
				items = append(items, GroupedTask{
					Task:     task,
					Group:    group,
					IsHeader: false,
				})
			}
		}
	}

	return items
}

func (m Model) rebuildItems() []GroupedTask {
	// Re-group with current collapse state
	var allTasks []domain.Task
	for _, item := range m.items {
		if !item.IsHeader {
			allTasks = append(allTasks, item.Task)
		}
	}
	return m.groupTasks(allTasks)
}

// View renders the forecast view
func (m Model) View() string {
	if m.err != nil {
		return m.renderError()
	}

	header := m.renderHeader()
	content := m.renderContent()

	return header + "\n" + content
}

func (m Model) renderHeader() string {
	taskCount := 0
	for _, item := range m.items {
		if !item.IsHeader {
			taskCount++
		}
	}
	headerText := fmt.Sprintf("FORECAST (%d tasks)", taskCount)
	return m.styles.UI.Header.Render(headerText)
}

func (m Model) renderContent() string {
	if !m.loaded {
		return "Loading..."
	}
	if len(m.items) == 0 {
		return "No tasks"
	}

	var b strings.Builder

	for i, item := range m.items {
		selected := i == m.cursor
		if item.IsHeader {
			b.WriteString(m.renderGroupHeader(item.Group, selected))
		} else {
			b.WriteString(m.renderTask(item.Task, item.Group, selected))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) renderGroupHeader(group DueGroup, selected bool) string {
	name := groupName(group)
	icon := "▼" // Expanded state - down arrow means "can collapse"
	if m.collapsed[group] {
		icon = "▶" // Collapsed state - right arrow means "can expand"
	}

	header := fmt.Sprintf("%s %s", icon, name)

	// Apply group-specific styling
	var style lipgloss.Style
	switch group {
	case GroupOverdue:
		style = m.styles.Forecast.Overdue
	case GroupToday:
		style = m.styles.Forecast.Today
	case GroupTomorrow:
		style = m.styles.Forecast.Tomorrow
	default:
		style = m.styles.Forecast.Later
	}

	if selected {
		style = style.Background(m.styles.Colors.Primary).Foreground(lipgloss.Color("#FFFFFF"))
	}

	return style.Bold(true).Render(header)
}

func (m Model) renderTask(task domain.Task, _ DueGroup, selected bool) string {
	statusIcon := "☐"
	if task.Completed {
		statusIcon = "☑"
	}

	flagIcon := ""
	if task.Flagged {
		flagIcon = " 🚩"
	}

	line := fmt.Sprintf("  %s %s%s", statusIcon, task.Name, flagIcon)

	if selected {
		return m.styles.Task.Selected.Render(line)
	}
	return m.styles.Task.Normal.Render(line)
}

func (m Model) renderError() string {
	header := m.styles.UI.Header.Render("FORECAST")
	errorText := fmt.Sprintf("Error: %v", m.err)
	errorStyle := m.styles.UI.Help.Foreground(m.styles.Colors.Error)
	return header + "\n" + errorStyle.Render(errorText)
}

// SelectedTask returns the currently selected task
func (m Model) SelectedTask() *domain.Task {
	if m.cursor >= len(m.items) || m.items[m.cursor].IsHeader {
		return nil
	}
	return &m.items[m.cursor].Task
}

// Refresh reloads tasks
func (m Model) Refresh() tea.Cmd {
	return m.loadTasks()
}

// SetFilter sets the filter state and applies it to tasks
func (m Model) SetFilter(f filter.State) Model {
	m.filter = f
	// Re-apply filter to existing tasks
	filteredTasks := m.applyFilter(m.allTasks)
	m.items = m.groupTasks(filteredTasks)
	// Reset cursor to first valid position
	if len(m.items) > 0 && m.items[0].IsHeader && len(m.items) > 1 {
		m.cursor = 1
	} else if len(m.items) > 0 {
		m.cursor = 0
	}
	return m
}

// applyFilter filters tasks based on current filter state
func (m Model) applyFilter(tasks []domain.Task) []domain.Task {
	if !m.filter.IsActive() {
		return tasks
	}
	matcher := filter.NewMatcher(m.filter)
	return matcher.FilterTasks(tasks)
}

func groupName(g DueGroup) string {
	switch g {
	case GroupOverdue:
		return "Overdue"
	case GroupToday:
		return "Today"
	case GroupTomorrow:
		return "Tomorrow"
	case GroupThisWeek:
		return "This Week"
	case GroupLater:
		return "Later"
	case GroupNoDue:
		return "No Due Date"
	default:
		return "Unknown"
	}
}

var enterKey = key.NewBinding(key.WithKeys("enter"))
