// Package taskedit provides an edit task overlay component.
package taskedit

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/cli/dateparse"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
)

// Field indices
const (
	FieldName = iota
	FieldNote
	FieldProject
	FieldTags
	FieldDueDate
	FieldDeferDate
	FieldFlagged
	NumFields
)

// SaveMsg is sent when the user saves changes
type SaveMsg struct {
	TaskID       string
	Modification domain.TaskModification
}

// CancelMsg is sent when the user cancels editing
type CancelMsg struct{}

// Model represents the edit task overlay state
type Model struct {
	task       *domain.Task
	visible    bool
	styles     *tui.Styles
	inputs     []textinput.Model
	focusIndex int
	flagged    bool
	width      int
	height     int
	err        string
}

// New creates a new edit task overlay
func New(styles *tui.Styles) Model {
	inputs := make([]textinput.Model, NumFields)

	// Name field
	inputs[FieldName] = textinput.New()
	inputs[FieldName].Placeholder = "Task name (required)"
	inputs[FieldName].CharLimit = 200

	// Note field
	inputs[FieldNote] = textinput.New()
	inputs[FieldNote].Placeholder = "Note"
	inputs[FieldNote].CharLimit = 1000

	// Project field
	inputs[FieldProject] = textinput.New()
	inputs[FieldProject].Placeholder = "Project name"
	inputs[FieldProject].CharLimit = 100

	// Tags field
	inputs[FieldTags] = textinput.New()
	inputs[FieldTags].Placeholder = "Tags (comma-separated)"
	inputs[FieldTags].CharLimit = 200

	// Due date field
	inputs[FieldDueDate] = textinput.New()
	inputs[FieldDueDate].Placeholder = "Due date (e.g., tomorrow, next monday)"
	inputs[FieldDueDate].CharLimit = 50

	// Defer date field
	inputs[FieldDeferDate] = textinput.New()
	inputs[FieldDeferDate].Placeholder = "Defer date"
	inputs[FieldDeferDate].CharLimit = 50

	// Flagged is a toggle, not a text input - index 6
	inputs[FieldFlagged] = textinput.New()
	inputs[FieldFlagged].Placeholder = "[Press Enter to toggle]"

	return Model{
		styles:     styles,
		inputs:     inputs,
		focusIndex: 0,
		visible:    false,
	}
}

// Show makes the overlay visible with the task to edit
func (m Model) Show(task *domain.Task) Model {
	m.task = task
	m.visible = true
	m.focusIndex = 0
	m.err = ""

	// Populate fields with current values
	m.inputs[FieldName].SetValue(task.Name)
	m.inputs[FieldNote].SetValue(task.Note)
	m.inputs[FieldProject].SetValue(task.ProjectName)

	// Tags as comma-separated
	if len(task.Tags) > 0 {
		m.inputs[FieldTags].SetValue(strings.Join(task.Tags, ", "))
	} else {
		m.inputs[FieldTags].SetValue("")
	}

	// Due date
	if task.DueDate != nil {
		m.inputs[FieldDueDate].SetValue(task.DueDate.Format("2006-01-02"))
	} else {
		m.inputs[FieldDueDate].SetValue("")
	}

	// Defer date
	if task.DeferDate != nil {
		m.inputs[FieldDeferDate].SetValue(task.DeferDate.Format("2006-01-02"))
	} else {
		m.inputs[FieldDeferDate].SetValue("")
	}

	m.flagged = task.Flagged

	// Focus first input
	m.inputs[m.focusIndex].Focus()

	return m
}

// Hide closes the overlay
func (m Model) Hide() Model {
	m.visible = false
	m.task = nil
	return m
}

// IsVisible returns true if the overlay is visible
func (m Model) IsVisible() bool {
	return m.visible
}

// SetSize updates the dimensions
func (m Model) SetSize(width, height int) Model {
	m.width = width
	m.height = height
	return m
}

// Init initializes the component
func (m Model) Init() tea.Cmd {
	return textinput.Blink
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
			return m, func() tea.Msg { return CancelMsg{} }

		case key.Matches(msg, submitKey):
			// On flagged field, toggle instead of submit
			if m.focusIndex == FieldFlagged {
				m.flagged = !m.flagged
				return m, nil
			}

			// Validate and submit
			if err := m.validate(); err != "" {
				m.err = err
				return m, nil
			}

			mod := m.buildModification()
			m.visible = false
			return m, func() tea.Msg {
				return SaveMsg{
					TaskID:       m.task.ID,
					Modification: mod,
				}
			}

		case key.Matches(msg, tabKey):
			m = m.nextField()
			return m, nil

		case key.Matches(msg, shiftTabKey):
			m = m.prevField()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	// Update the focused input
	if m.focusIndex < FieldFlagged {
		var cmd tea.Cmd
		m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) nextField() Model {
	// Blur current
	if m.focusIndex < FieldFlagged {
		m.inputs[m.focusIndex].Blur()
	}

	// Move to next
	m.focusIndex++
	if m.focusIndex >= NumFields {
		m.focusIndex = 0
	}

	// Focus new
	if m.focusIndex < FieldFlagged {
		m.inputs[m.focusIndex].Focus()
	}

	return m
}

func (m Model) prevField() Model {
	// Blur current
	if m.focusIndex < FieldFlagged {
		m.inputs[m.focusIndex].Blur()
	}

	// Move to previous
	m.focusIndex--
	if m.focusIndex < 0 {
		m.focusIndex = NumFields - 1
	}

	// Focus new
	if m.focusIndex < FieldFlagged {
		m.inputs[m.focusIndex].Focus()
	}

	return m
}

func (m Model) validate() string {
	name := strings.TrimSpace(m.inputs[FieldName].Value())
	if name == "" {
		return "Task name is required"
	}

	// Validate due date if provided
	dueStr := strings.TrimSpace(m.inputs[FieldDueDate].Value())
	if dueStr != "" {
		if _, err := dateparse.Parse(dueStr); err != nil {
			return "Invalid due date format"
		}
	}

	// Validate defer date if provided
	deferStr := strings.TrimSpace(m.inputs[FieldDeferDate].Value())
	if deferStr != "" {
		if _, err := dateparse.Parse(deferStr); err != nil {
			return "Invalid defer date format"
		}
	}

	return ""
}

func (m Model) buildModification() domain.TaskModification {
	mod := domain.TaskModification{}

	// Name
	newName := strings.TrimSpace(m.inputs[FieldName].Value())
	if newName != m.task.Name {
		mod.Name = &newName
	}

	// Note
	newNote := strings.TrimSpace(m.inputs[FieldNote].Value())
	if newNote != m.task.Note {
		mod.Note = &newNote
	}

	// Project
	newProject := strings.TrimSpace(m.inputs[FieldProject].Value())
	if newProject != m.task.ProjectName {
		if newProject == "" {
			// Clear project
			empty := ""
			mod.ProjectID = &empty
		} else {
			// Note: app.go will need to resolve project name to ID
			mod.ProjectID = &newProject
		}
	}

	// Tags - compare and build add/remove lists
	currentTagNames := make(map[string]bool)
	for _, tagName := range m.task.Tags {
		currentTagNames[strings.ToLower(tagName)] = true
	}

	newTagsStr := strings.TrimSpace(m.inputs[FieldTags].Value())
	newTagNames := make(map[string]bool)
	if newTagsStr != "" {
		for _, tagName := range strings.Split(newTagsStr, ",") {
			trimmed := strings.TrimSpace(tagName)
			if trimmed != "" {
				newTagNames[strings.ToLower(trimmed)] = true
			}
		}
	}

	// Find tags to add
	for tagName := range newTagNames {
		if !currentTagNames[tagName] {
			mod.AddTags = append(mod.AddTags, tagName)
		}
	}

	// Find tags to remove
	for tagName := range currentTagNames {
		if !newTagNames[tagName] {
			mod.RemoveTags = append(mod.RemoveTags, tagName)
		}
	}

	// Due date
	dueStr := strings.TrimSpace(m.inputs[FieldDueDate].Value())
	if dueStr == "" && m.task.DueDate != nil {
		mod.ClearDue = true
	} else if dueStr != "" {
		if dueDate, err := dateparse.Parse(dueStr); err == nil {
			mod.DueDate = &dueDate
		}
	}

	// Defer date
	deferStr := strings.TrimSpace(m.inputs[FieldDeferDate].Value())
	if deferStr == "" && m.task.DeferDate != nil {
		mod.ClearDefer = true
	} else if deferStr != "" {
		if deferDate, err := dateparse.Parse(deferStr); err == nil {
			mod.DeferDate = &deferDate
		}
	}

	// Flagged
	if m.flagged != m.task.Flagged {
		mod.Flagged = &m.flagged
	}

	return mod
}

// View renders the overlay
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	modalWidth := min(60, m.width-4)
	if modalWidth < 30 {
		modalWidth = 30
	}

	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.styles.Colors.Primary).
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render("Edit Task"))
	b.WriteString("\n\n")

	// Error message if any
	if m.err != "" {
		errStyle := lipgloss.NewStyle().
			Foreground(m.styles.Colors.Error).
			Width(modalWidth - 4)
		b.WriteString(errStyle.Render(m.err))
		b.WriteString("\n\n")
	}

	// Fields
	labels := []string{"Name:", "Note:", "Project:", "Tags:", "Due:", "Defer:", "Flagged:"}

	labelStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary).
		Width(10)

	inputWidth := modalWidth - 16

	for i := 0; i < NumFields; i++ {
		// Label
		b.WriteString(labelStyle.Render(labels[i]))

		if i == FieldFlagged {
			// Flagged toggle
			flagText := "[ ] No"
			if m.flagged {
				flagText = "[✓] Yes"
			}

			var style lipgloss.Style
			if i == m.focusIndex {
				style = lipgloss.NewStyle().
					Background(m.styles.Colors.Primary).
					Foreground(lipgloss.Color("#FFFFFF")).
					Width(inputWidth)
			} else {
				style = lipgloss.NewStyle().Width(inputWidth)
			}
			b.WriteString(style.Render(flagText))
		} else {
			// Text input
			m.inputs[i].Width = inputWidth
			b.WriteString(m.inputs[i].View())
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Footer with hints
	hintStyle := lipgloss.NewStyle().
		Foreground(m.styles.Colors.Secondary).
		Width(modalWidth - 4).
		Align(lipgloss.Center)
	b.WriteString(hintStyle.Render("Tab/Shift+Tab: Navigate  Enter: Save  Esc: Cancel"))

	return m.styles.UI.Overlay.
		Width(modalWidth).
		Render(b.String())
}

// Key bindings
var (
	escapeKey   = key.NewBinding(key.WithKeys("esc", "escape"))
	submitKey   = key.NewBinding(key.WithKeys("enter"))
	tabKey      = key.NewBinding(key.WithKeys("tab"))
	shiftTabKey = key.NewBinding(key.WithKeys("shift+tab"))
)
