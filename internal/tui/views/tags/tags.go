// Package tags provides the tags view for the TUI.
package tags

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/taglist"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/tasklist"
)

// ViewMode represents whether we're viewing tags or a tag's tasks
type ViewMode int

const (
	ModeTagList ViewMode = iota
	ModeTagTasks
)

// TagsAndCountsLoadedMsg is sent when tags and counts are loaded
type TagsAndCountsLoadedMsg struct {
	Tags   []domain.Tag
	Counts map[string]int
}

// Model represents the tags view state
type Model struct {
	tagList    taglist.Model
	taskList   tasklist.Model
	service    service.OmniFocusService
	styles     *tui.Styles
	keys       tui.KeyMap
	mode       ViewMode
	currentTag *domain.Tag
	width      int
	height     int
	err        error
	loaded     bool
}

// New creates a new tags view
func New(styles *tui.Styles, keys tui.KeyMap, svc service.OmniFocusService) Model {
	return Model{
		tagList:  taglist.New(styles, keys),
		taskList: tasklist.New(styles, keys),
		service:  svc,
		styles:   styles,
		keys:     keys,
		mode:     ModeTagList,
		loaded:   false,
	}
}

// Init initializes the tags view
func (m Model) Init() tea.Cmd {
	return m.loadTagsAndCounts()
}

func (m Model) loadTagsAndCounts() tea.Cmd {
	return func() tea.Msg {
		tags, err := m.service.GetTags()
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		counts, err := m.service.GetTagCounts()
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return TagsAndCountsLoadedMsg{Tags: tags, Counts: counts}
	}
}

func (m Model) loadTagTasks(tagID string) tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.service.GetTasksByTag(tagID)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TasksLoadedMsg{Tasks: tasks}
	}
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TagsAndCountsLoadedMsg:
		m.tagList = m.tagList.SetTags(msg.Tags, msg.Counts)
		m.loaded = true
		m.err = nil
		return m, nil

	case tui.TagsLoadedMsg:
		// If we receive just tags, load counts separately
		return m, m.loadTagsAndCounts()

	case tui.TasksLoadedMsg:
		m.taskList = m.taskList.SetTasks(msg.Tasks)
		return m, nil

	case tui.ErrorMsg:
		m.err = msg.Err
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		headerHeight := 2
		availableHeight := msg.Height - headerHeight
		if availableHeight < 0 {
			availableHeight = 0
		}

		subMsg := tea.WindowSizeMsg{Width: msg.Width, Height: availableHeight}

		var cmd tea.Cmd
		if m.mode == ModeTagList {
			m.tagList, cmd = m.tagList.Update(subMsg)
		} else {
			m.taskList, cmd = m.taskList.Update(subMsg)
		}
		return m, cmd

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	default:
		return m.delegateToCurrentList(msg)
	}
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	// Handle drill-down with Enter
	if key.Matches(msg, enterKey) {
		if m.mode == ModeTagList {
			tag := m.tagList.SelectedTag()
			if tag != nil {
				m.mode = ModeTagTasks
				m.currentTag = tag
				return m, m.loadTagTasks(tag.ID)
			}
		}
		return m, nil
	}

	// Handle back navigation
	if key.Matches(msg, backKey) || key.Matches(msg, escapeKey) {
		if m.mode == ModeTagTasks {
			m.mode = ModeTagList
			m.currentTag = nil
			return m, nil
		}
		return m, nil
	}

	return m.delegateToCurrentList(msg)
}

func (m Model) delegateToCurrentList(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.mode == ModeTagList {
		m.tagList, cmd = m.tagList.Update(msg)
	} else {
		m.taskList, cmd = m.taskList.Update(msg)
	}
	return m, cmd
}

// View renders the tags view
func (m Model) View() string {
	if m.err != nil {
		return m.renderError()
	}

	header := m.renderHeader()

	var content string
	if m.mode == ModeTagList {
		content = m.tagList.View()
	} else {
		content = m.taskList.View()
	}

	return header + "\n" + content
}

func (m Model) renderHeader() string {
	var headerText string
	if m.mode == ModeTagList {
		count := len(m.tagList.Tags())
		headerText = fmt.Sprintf("TAGS (%d)", count)
	} else if m.currentTag != nil {
		headerText = fmt.Sprintf("🏷 %s", m.currentTag.Name)
	} else {
		headerText = "TAG TASKS"
	}

	styled := m.styles.UI.Header.Render(headerText)

	if m.mode == ModeTagTasks {
		hint := m.styles.UI.Help.Render("  [h/Esc] back")
		styled += hint
	}

	return styled
}

func (m Model) renderError() string {
	header := m.styles.UI.Header.Render("TAGS")
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

// SelectedTask returns the currently selected task (when in task mode)
func (m Model) SelectedTask() *domain.Task {
	if m.mode == ModeTagTasks {
		return m.taskList.SelectedTask()
	}
	return nil
}

// Refresh reloads tags
func (m Model) Refresh() tea.Cmd {
	if m.mode == ModeTagTasks && m.currentTag != nil {
		return m.loadTagTasks(m.currentTag.ID)
	}
	return m.loadTagsAndCounts()
}

// Mode returns the current view mode
func (m Model) Mode() ViewMode {
	return m.mode
}

var (
	enterKey  = key.NewBinding(key.WithKeys("enter"))
	backKey   = key.NewBinding(key.WithKeys("h", "left"))
	escapeKey = key.NewBinding(key.WithKeys("esc", "escape"))
)
