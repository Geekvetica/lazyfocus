// Package taglist provides the tag list TUI component.
package taglist

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// TagIcon for tag display
const TagIcon = "🏷"

// TagWithCount represents a tag with its task count
type TagWithCount struct {
	Tag   domain.Tag
	Count int
	Depth int // For hierarchical display
}

// Model represents the tag list component state
type Model struct {
	tags    []TagWithCount
	cursor  int
	width   int
	height  int
	styles  *tui.Styles
	keys    tui.KeyMap
	loading bool
	empty   bool
}

// New creates a new tag list component
func New(styles *tui.Styles, keys tui.KeyMap) Model {
	return Model{
		tags:    []TagWithCount{},
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

// Update handles messages
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

func (m Model) handleKeyPress(msg tea.KeyMsg) (Model, tea.Cmd) {
	if len(m.tags) == 0 {
		return m, nil
	}

	if key.Matches(msg, m.keys.Down) {
		m.cursor++
		if m.cursor >= len(m.tags) {
			m.cursor = 0
		}
		return m, nil
	}

	if key.Matches(msg, m.keys.Up) {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.tags) - 1
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
	return m.renderTags()
}

func (m Model) renderLoading() string {
	if m.height == 0 {
		return "Loading..."
	}
	padding := strings.Repeat("\n", m.height/2)
	return padding + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "Loading tags...")
}

func (m Model) renderEmpty() string {
	if m.height == 0 {
		return "No tags"
	}
	padding := strings.Repeat("\n", m.height/2)
	return padding + lipgloss.PlaceHorizontal(m.width, lipgloss.Center, "No tags")
}

func (m Model) renderTags() string {
	var b strings.Builder

	for i, tagWithCount := range m.tags {
		line := m.formatTagLine(tagWithCount, i == m.cursor)
		b.WriteString(line)
		b.WriteString("\n")
	}

	return b.String()
}

func (m Model) formatTagLine(twc TagWithCount, selected bool) string {
	// Indentation for hierarchy
	indent := strings.Repeat("  ", twc.Depth)

	// Build left side with tag icon and name
	leftSide := fmt.Sprintf("%s%s %s", indent, TagIcon, twc.Tag.Name)

	// Build right side (task count)
	rightSide := fmt.Sprintf("(%d)", twc.Count)

	// Calculate spacing
	contentWidth := m.width
	if contentWidth == 0 {
		contentWidth = 80
	}

	leftLen := len(indent) + runewidth.StringWidth(TagIcon) + 1 + runewidth.StringWidth(twc.Tag.Name)
	rightLen := runewidth.StringWidth(rightSide)
	spacing := contentWidth - leftLen - rightLen - 2
	if spacing < 0 {
		spacing = 1
	}

	line := leftSide + strings.Repeat(" ", spacing) + rightSide

	// Apply styles
	if selected {
		return m.styles.Task.Selected.Render(line)
	}

	return m.styles.Tag.Badge.Render(line)
}

// SetTags updates the tag list with counts
func (m Model) SetTags(tags []domain.Tag, counts map[string]int) Model {
	m.tags = m.flattenTags(tags, counts, 0)
	m.empty = len(m.tags) == 0
	m.loading = false
	if m.cursor >= len(m.tags) {
		if len(m.tags) > 0 {
			m.cursor = len(m.tags) - 1
		} else {
			m.cursor = 0
		}
	}
	return m
}

// flattenTags converts hierarchical tags to flat list with depth info
func (m Model) flattenTags(tags []domain.Tag, counts map[string]int, depth int) []TagWithCount {
	var result []TagWithCount
	for _, tag := range tags {
		count := counts[tag.ID]
		result = append(result, TagWithCount{
			Tag:   tag,
			Count: count,
			Depth: depth,
		})
		// Recursively add children
		if len(tag.Children) > 0 {
			children := m.flattenTags(tag.Children, counts, depth+1)
			result = append(result, children...)
		}
	}
	return result
}

// SetLoading sets the loading state
func (m Model) SetLoading(loading bool) Model {
	m.loading = loading
	return m
}

// SelectedTag returns the currently selected tag
func (m Model) SelectedTag() *domain.Tag {
	if len(m.tags) == 0 || m.cursor >= len(m.tags) {
		return nil
	}
	return &m.tags[m.cursor].Tag
}

// SelectedIndex returns the current cursor position
func (m Model) SelectedIndex() int {
	return m.cursor
}

// Tags returns the current tags for header display
func (m Model) Tags() []TagWithCount {
	return m.tags
}
