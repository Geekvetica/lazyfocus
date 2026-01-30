// Package app provides the main TUI application model and orchestration.
package app

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/command"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/commandinput"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/confirm"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/quickadd"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/searchinput"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/taskdetail"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/taskedit"
	"github.com/pwojciechowski/lazyfocus/internal/tui/filter"
	"github.com/pwojciechowski/lazyfocus/internal/tui/overlay"
	"github.com/pwojciechowski/lazyfocus/internal/tui/views/forecast"
	"github.com/pwojciechowski/lazyfocus/internal/tui/views/inbox"
	"github.com/pwojciechowski/lazyfocus/internal/tui/views/projects"
	"github.com/pwojciechowski/lazyfocus/internal/tui/views/review"
	"github.com/pwojciechowski/lazyfocus/internal/tui/views/tags"
)

// DeleteContext stores context for delete confirmation
type DeleteContext struct {
	TaskID   string
	TaskName string
}

// Model represents the main TUI application state
type Model struct {
	// Views
	inboxView    inbox.Model
	projectsView projects.Model
	tagsView     tags.Model
	forecastView forecast.Model
	reviewView   review.Model
	currentView  int // tui.ViewInbox, tui.ViewProjects, etc from messages.go

	// Overlays
	quickAdd     quickadd.Model
	taskDetail   taskdetail.Model
	taskEdit     taskedit.Model
	confirmModal confirm.Model
	searchInput  searchinput.Model
	commandInput commandinput.Model
	showHelp     bool
	compositor   *overlay.Compositor

	// State
	filterState filter.State
	service     service.OmniFocusService
	styles      *tui.Styles
	keys        tui.KeyMap
	width       int
	height      int
	err         error
	ready       bool // true after first WindowSizeMsg
}

// NewApp creates a new TUI application instance
func NewApp(svc service.OmniFocusService) Model {
	styles := tui.DefaultStyles()
	keys := tui.DefaultKeyMap()

	return Model{
		// Views
		inboxView:    inbox.New(styles, keys, svc),
		projectsView: projects.New(styles, keys, svc),
		tagsView:     tags.New(styles, keys, svc),
		forecastView: forecast.New(styles, keys, svc),
		reviewView:   review.New(styles, keys, svc),
		currentView:  tui.ViewInbox,

		// Overlays
		quickAdd:     quickadd.New(styles, svc),
		taskDetail:   taskdetail.New(styles, keys),
		taskEdit:     taskedit.New(styles),
		confirmModal: confirm.New(styles),
		searchInput:  searchinput.New(styles),
		commandInput: commandinput.New(styles),
		showHelp:     false,
		compositor:   overlay.New(styles.UI.OverlayBackdrop),

		// State
		filterState: filter.State{},
		service:     svc,
		styles:      styles,
		keys:        keys,
		ready:       false,
	}
}

// Init initializes the application
func (m Model) Init() tea.Cmd {
	return m.initCurrentView()
}

// initCurrentView initializes the current view
func (m Model) initCurrentView() tea.Cmd {
	switch m.currentView {
	case tui.ViewInbox:
		return m.inboxView.Init()
	case tui.ViewProjects:
		return m.projectsView.Init()
	case tui.ViewTags:
		return m.tagsView.Init()
	case tui.ViewForecast:
		return m.forecastView.Init()
	case tui.ViewReview:
		return m.reviewView.Init()
	default:
		return nil
	}
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
		return m.handleWindowResize(msg)
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

	// Handle task detail action messages before overlay delegation
	// These are emitted by taskdetail component and must be handled at app level
	if newModel, cmd, handled := m.handleTaskDetailMessages(msg); handled {
		return newModel, cmd
	}

	// Handle task edit messages before overlay delegation
	if newModel, cmd, handled := m.handleTaskEditMessages(msg); handled {
		return newModel, cmd
	}

	// Handle overlays in priority order (highest to lowest)
	if newModel, cmd, handled := m.handleOverlays(msg); handled {
		return newModel, cmd
	}

	// Handle custom messages
	if newModel, cmd, handled := m.handleCustomMessages(msg); handled {
		return newModel, cmd
	}

	// Handle global keys when overlay is not visible
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		return m.handleKeyMsg(keyMsg)
	}

	// Delegate to current view
	return m.delegateToCurrentView(msg)
}

// handleWindowResize handles tea.WindowSizeMsg
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	m.ready = true

	// Update compositor dimensions
	m.compositor.SetSize(msg.Width, msg.Height)

	// Update all overlays
	m.quickAdd = m.quickAdd.SetSize(msg.Width, msg.Height)
	m.taskDetail = m.taskDetail.SetSize(msg.Width, msg.Height)
	m.taskEdit = m.taskEdit.SetSize(msg.Width, msg.Height)
	m.confirmModal = m.confirmModal.SetSize(msg.Width, msg.Height)
	m.searchInput = m.searchInput.SetWidth(msg.Width)
	m.commandInput = m.commandInput.SetWidth(msg.Width)

	// Pass resize to all views
	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.inboxView, cmd = m.inboxView.Update(msg)
	cmds = append(cmds, cmd)
	m.projectsView, cmd = m.projectsView.Update(msg)
	cmds = append(cmds, cmd)
	m.tagsView, cmd = m.tagsView.Update(msg)
	cmds = append(cmds, cmd)
	m.forecastView, cmd = m.forecastView.Update(msg)
	cmds = append(cmds, cmd)
	m.reviewView, cmd = m.reviewView.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// handleOverlays delegates messages to visible overlays
// Returns the updated model, command, and true if an overlay handled the message
func (m Model) handleOverlays(msg tea.Msg) (Model, tea.Cmd, bool) {
	// 1. Confirm modal (highest - blocking)
	if m.confirmModal.IsVisible() {
		var cmd tea.Cmd
		m.confirmModal, cmd = m.confirmModal.Update(msg)
		return m, cmd, true
	}

	// 2. Task edit overlay
	if m.taskEdit.IsVisible() {
		var cmd tea.Cmd
		m.taskEdit, cmd = m.taskEdit.Update(msg)
		return m, cmd, true
	}

	// 3. Task detail overlay
	if m.taskDetail.IsVisible() {
		var cmd tea.Cmd
		m.taskDetail, cmd = m.taskDetail.Update(msg)
		return m, cmd, true
	}

	// 4. Quick add overlay
	if m.quickAdd.IsVisible() {
		var cmd tea.Cmd
		m.quickAdd, cmd = m.quickAdd.Update(msg)
		return m, cmd, true
	}

	// 5. Search input
	if m.searchInput.IsVisible() {
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd, true
	}

	// 6. Command input
	if m.commandInput.IsVisible() {
		var cmd tea.Cmd
		m.commandInput, cmd = m.commandInput.Update(msg)
		return m, cmd, true
	}

	return m, nil, false
}

// handleCustomMessages handles custom message types from components
// Returns the updated model, command, and true if message was handled
func (m Model) handleCustomMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	// Handle search input messages
	if newModel, cmd, handled := m.handleSearchInputMessages(msg); handled {
		return newModel, cmd, true
	}

	// Handle command input messages
	if newModel, cmd, handled := m.handleCommandInputMessages(msg); handled {
		return newModel, cmd, true
	}

	// Handle confirm messages
	if newModel, cmd, handled := m.handleConfirmMessages(msg); handled {
		return newModel, cmd, true
	}

	// Handle task operation messages
	if newModel, cmd, handled := m.handleTaskOperationMessages(msg); handled {
		return newModel, cmd, true
	}

	return m, nil, false
}

// handleTaskDetailMessages handles task detail related messages
func (m Model) handleTaskDetailMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	if _, ok := msg.(taskdetail.CloseMsg); ok {
		m.taskDetail = m.taskDetail.Hide()
		return m, nil, true
	}

	if editMsg, ok := msg.(taskdetail.EditRequestedMsg); ok {
		m.taskDetail = m.taskDetail.Hide()
		m.taskEdit = m.taskEdit.Show(&editMsg.Task)
		return m, nil, true
	}

	if completeMsg, ok := msg.(taskdetail.CompleteRequestedMsg); ok {
		m.taskDetail = m.taskDetail.Hide()
		return m, m.completeTask(completeMsg.TaskID), true
	}

	if deleteMsg, ok := msg.(taskdetail.DeleteRequestedMsg); ok {
		m.taskDetail = m.taskDetail.Hide()
		ctx := DeleteContext{TaskID: deleteMsg.TaskID, TaskName: deleteMsg.TaskName}
		m.confirmModal = m.confirmModal.ShowWithContext(
			"Delete Task",
			fmt.Sprintf("Delete \"%s\"?", deleteMsg.TaskName),
			ctx,
		)
		return m, nil, true
	}

	if _, ok := msg.(taskdetail.FlagRequestedMsg); ok {
		task := m.taskDetail.Task()
		m.taskDetail = m.taskDetail.Hide()
		if task != nil {
			return m, m.toggleTaskFlag(task), true
		}
		return m, nil, true
	}

	return m, nil, false
}

// handleTaskEditMessages handles task edit related messages
func (m Model) handleTaskEditMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	if saveMsg, ok := msg.(taskedit.SaveMsg); ok {
		m.taskEdit = m.taskEdit.Hide()
		return m, m.modifyTask(saveMsg.TaskID, saveMsg.Modification), true
	}

	if _, ok := msg.(taskedit.CancelMsg); ok {
		m.taskEdit = m.taskEdit.Hide()
		return m, nil, true
	}

	return m, nil, false
}

// handleSearchInputMessages handles search input related messages
func (m Model) handleSearchInputMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	if searchMsg, ok := msg.(searchinput.SearchChangedMsg); ok {
		m.filterState = m.filterState.WithSearchText(searchMsg.Text)
		m = m.applyFilterToCurrentView()
		return m, nil, true
	}

	if _, ok := msg.(searchinput.SearchClearedMsg); ok {
		m.filterState = m.filterState.Clear()
		m = m.applyFilterToCurrentView()
		return m, nil, true
	}

	if searchMsg, ok := msg.(searchinput.SearchConfirmedMsg); ok {
		m.filterState = m.filterState.WithSearchText(searchMsg.Text)
		m = m.applyFilterToCurrentView()
		return m, nil, true
	}

	return m, nil, false
}

// handleCommandInputMessages handles command input related messages
func (m Model) handleCommandInputMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	if cmdMsg, ok := msg.(commandinput.CommandExecutedMsg); ok {
		newModel, cmd := m.executeCommand(cmdMsg.Command)
		return newModel, cmd, true
	}

	if _, ok := msg.(commandinput.CommandCancelledMsg); ok {
		return m, nil, true
	}

	if errMsg, ok := msg.(commandinput.CommandErrorMsg); ok {
		m.err = fmt.Errorf("%s", errMsg.Error)
		return m, nil, true
	}

	return m, nil, false
}

// handleConfirmMessages handles confirmation modal messages
func (m Model) handleConfirmMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	if msg, ok := msg.(confirm.ConfirmedMsg); ok {
		if ctx, ok := msg.Context.(DeleteContext); ok {
			return m, m.deleteTask(ctx.TaskID), true
		}
		return m, nil, true
	}

	return m, nil, false
}

// handleTaskOperationMessages handles task operation result messages
func (m Model) handleTaskOperationMessages(msg tea.Msg) (Model, tea.Cmd, bool) {
	if _, ok := msg.(tui.TaskCompletedMsg); ok {
		return m, m.refreshCurrentView(), true
	}

	if _, ok := msg.(tui.TaskDeletedMsg); ok {
		return m, m.refreshCurrentView(), true
	}

	if _, ok := msg.(tui.TaskModifiedMsg); ok {
		return m, m.refreshCurrentView(), true
	}

	return m, nil, false
}

// handleKeyMsg handles global key messages
func (m Model) handleKeyMsg(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
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

	// Show task detail on Enter
	if keyMsg.String() == "enter" {
		task := m.getSelectedTask()
		if task != nil {
			m.taskDetail = m.taskDetail.Show(task)
			return m, nil
		}
	}

	// Show edit task overlay
	if key.Matches(keyMsg, m.keys.Edit) {
		task := m.getSelectedTask()
		if task != nil {
			m.taskEdit = m.taskEdit.Show(task)
			return m, nil
		}
		return m, nil
	}

	// Delete task - show confirmation
	if key.Matches(keyMsg, m.keys.Delete) {
		task := m.getSelectedTask()
		if task != nil {
			ctx := DeleteContext{TaskID: task.ID, TaskName: task.Name}
			m.confirmModal = m.confirmModal.ShowWithContext(
				"Delete Task",
				fmt.Sprintf("Delete \"%s\"?", task.Name),
				ctx,
			)
		}
		return m, nil
	}

	// Toggle flag - immediate action (no confirmation)
	if key.Matches(keyMsg, m.keys.Flag) {
		task := m.getSelectedTask()
		if task != nil {
			return m, m.toggleTaskFlag(task)
		}
		return m, nil
	}

	// Show search input
	if keyMsg.String() == "/" {
		m.searchInput = m.searchInput.Show()
		return m, nil
	}

	// Show command input
	if keyMsg.String() == ":" {
		m.commandInput = m.commandInput.Show()
		return m, nil
	}

	// Handle view switching
	return m.handleViewSwitching(keyMsg)
}

// handleViewSwitching handles view switching key presses
func (m Model) handleViewSwitching(keyMsg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Matches(keyMsg, m.keys.View1) {
		if m.currentView != tui.ViewInbox {
			m.currentView = tui.ViewInbox
			return m, m.inboxView.Init()
		}
		return m, nil
	}
	if key.Matches(keyMsg, m.keys.View2) {
		if m.currentView != tui.ViewProjects {
			m.currentView = tui.ViewProjects
			return m, m.projectsView.Init()
		}
		return m, nil
	}
	if key.Matches(keyMsg, m.keys.View3) {
		if m.currentView != tui.ViewTags {
			m.currentView = tui.ViewTags
			return m, m.tagsView.Init()
		}
		return m, nil
	}
	if key.Matches(keyMsg, m.keys.View4) {
		if m.currentView != tui.ViewForecast {
			m.currentView = tui.ViewForecast
			return m, m.forecastView.Init()
		}
		return m, nil
	}
	if key.Matches(keyMsg, m.keys.View5) {
		if m.currentView != tui.ViewReview {
			m.currentView = tui.ViewReview
			return m, m.reviewView.Init()
		}
		return m, nil
	}
	return m, nil
}

// delegateToCurrentView delegates messages to the current view
func (m Model) delegateToCurrentView(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentView {
	case tui.ViewInbox:
		m.inboxView, cmd = m.inboxView.Update(msg)
	case tui.ViewProjects:
		m.projectsView, cmd = m.projectsView.Update(msg)
	case tui.ViewTags:
		m.tagsView, cmd = m.tagsView.Update(msg)
	case tui.ViewForecast:
		m.forecastView, cmd = m.forecastView.Update(msg)
	case tui.ViewReview:
		m.reviewView, cmd = m.reviewView.Update(msg)
	}
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
	case tui.ViewProjects:
		view = m.projectsView.View()
	case tui.ViewTags:
		view = m.tagsView.View()
	case tui.ViewForecast:
		view = m.forecastView.View()
	case tui.ViewReview:
		view = m.reviewView.View()
	default:
		view = "View not implemented"
	}

	// Layer overlays from lowest to highest priority
	// Bottom bar overlays (search, command)
	if m.searchInput.IsVisible() {
		view = m.renderWithBottomBar(view, m.searchInput.View())
	}

	if m.commandInput.IsVisible() {
		view = m.renderWithBottomBar(view, m.commandInput.View())
	}

	// Center overlays
	if m.quickAdd.IsVisible() {
		view = m.layerOverlay(view, m.quickAdd.View())
	}

	if m.taskDetail.IsVisible() {
		view = m.layerOverlay(view, m.taskDetail.View())
	}

	if m.taskEdit.IsVisible() {
		view = m.layerOverlay(view, m.taskEdit.View())
	}

	// Top priority overlays
	if m.confirmModal.IsVisible() {
		view = m.layerOverlay(view, m.confirmModal.View())
	}

	if m.showHelp {
		view = m.layerOverlay(view, m.renderHelp())
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

	return overlay
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
	return m.compositor.Compose(base, overlay, true)
}

// renderWithBottomBar renders a bottom bar overlay (search, command)
func (m Model) renderWithBottomBar(base, bottomBar string) string {
	// Split base into lines
	baseLines := strings.Split(base, "\n")
	if len(baseLines) == 0 {
		return bottomBar
	}

	// Replace last line with bottom bar
	if len(baseLines) > 0 {
		baseLines[len(baseLines)-1] = bottomBar
	}

	return strings.Join(baseLines, "\n")
}

// getSelectedTask returns the currently selected task from the current view
func (m Model) getSelectedTask() *domain.Task {
	switch m.currentView {
	case tui.ViewInbox:
		return m.inboxView.SelectedTask()
	case tui.ViewProjects:
		return m.projectsView.SelectedTask()
	case tui.ViewTags:
		return m.tagsView.SelectedTask()
	case tui.ViewForecast:
		return m.forecastView.SelectedTask()
	case tui.ViewReview:
		return m.reviewView.SelectedTask()
	default:
		return nil
	}
}

// deleteTask creates a command to delete a task
func (m Model) deleteTask(taskID string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.service.DeleteTask(taskID)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TaskDeletedMsg{
			TaskID:   result.ID,
			TaskName: result.Message, // Message contains task name
		}
	}
}

// toggleTaskFlag creates a command to toggle a task's flag status
func (m Model) toggleTaskFlag(task *domain.Task) tea.Cmd {
	return func() tea.Msg {
		newFlagged := !task.Flagged
		mod := domain.TaskModification{
			Flagged: &newFlagged,
		}
		result, err := m.service.ModifyTask(task.ID, mod)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TaskModifiedMsg{Task: *result}
	}
}

// completeTask creates a command to complete a task
func (m Model) completeTask(taskID string) tea.Cmd {
	return func() tea.Msg {
		result, err := m.service.CompleteTask(taskID)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TaskCompletedMsg{
			TaskID:   result.ID,
			TaskName: result.Message,
		}
	}
}

// modifyTask creates a command to modify a task
func (m Model) modifyTask(taskID string, mod domain.TaskModification) tea.Cmd {
	return func() tea.Msg {
		result, err := m.service.ModifyTask(taskID, mod)
		if err != nil {
			return tui.ErrorMsg{Err: err}
		}
		return tui.TaskModifiedMsg{Task: *result}
	}
}

// refreshCurrentView creates a command to refresh the current view
func (m Model) refreshCurrentView() tea.Cmd {
	switch m.currentView {
	case tui.ViewInbox:
		return m.inboxView.Refresh()
	case tui.ViewProjects:
		return m.projectsView.Refresh()
	case tui.ViewTags:
		return m.tagsView.Refresh()
	case tui.ViewForecast:
		return m.forecastView.Refresh()
	case tui.ViewReview:
		return m.reviewView.Refresh()
	default:
		return nil
	}
}

// executeCommand handles command execution
func (m Model) executeCommand(cmd *command.Command) (Model, tea.Cmd) {
	if cmd == nil {
		return m, nil
	}

	switch cmd.Name {
	case "quit":
		return m, tea.Quit
	case "refresh":
		return m, m.refreshCurrentView()
	case "add":
		return m.executeAddCommand(cmd)
	case "complete":
		return m.executeCompleteCommand()
	case "delete":
		return m.executeDeleteCommand()
	case "project":
		return m.executeProjectCommand(cmd)
	case "tag":
		return m.executeTagCommand(cmd)
	case "due":
		return m.executeDueCommand(cmd)
	case "flagged":
		return m.executeFlaggedCommand()
	case "clear":
		return m.executeClearCommand()
	case "help":
		m.showHelp = !m.showHelp
		return m, nil
	default:
		return m, nil
	}
}

// executeAddCommand handles the "add" command
func (m Model) executeAddCommand(cmd *command.Command) (Model, tea.Cmd) {
	// Open quick add with args if provided
	if len(cmd.Args) > 0 {
		_ = strings.Join(cmd.Args, " ") // taskText for future pre-fill feature
		m.quickAdd = m.quickAdd.Show()
		// TODO: Pre-fill quick add with taskText
		// This would require adding a method to quickadd component
	} else {
		m.quickAdd = m.quickAdd.Show()
	}
	return m, nil
}

// executeCompleteCommand handles the "complete" command
func (m Model) executeCompleteCommand() (Model, tea.Cmd) {
	task := m.getSelectedTask()
	if task != nil {
		return m, m.completeTask(task.ID)
	}
	return m, nil
}

// executeDeleteCommand handles the "delete" command
func (m Model) executeDeleteCommand() (Model, tea.Cmd) {
	task := m.getSelectedTask()
	if task != nil {
		ctx := DeleteContext{TaskID: task.ID, TaskName: task.Name}
		m.confirmModal = m.confirmModal.ShowWithContext(
			"Delete Task",
			fmt.Sprintf("Delete \"%s\"?", task.Name),
			ctx,
		)
	}
	return m, nil
}

// executeProjectCommand handles the "project" command
func (m Model) executeProjectCommand(cmd *command.Command) (Model, tea.Cmd) {
	if len(cmd.Args) > 0 {
		projectName := strings.Join(cmd.Args, " ")
		// Resolve project name to ID
		projects, err := m.service.GetProjects("")
		if err != nil {
			m.err = fmt.Errorf("failed to get projects: %w", err)
			return m, nil
		}

		// Find project by name (case-insensitive)
		var projectID string
		projectNameLower := strings.ToLower(projectName)
		for _, proj := range projects {
			if strings.ToLower(proj.Name) == projectNameLower {
				projectID = proj.ID
				break
			}
		}

		if projectID == "" {
			m.err = fmt.Errorf("project not found: %s", projectName)
			return m, nil
		}

		m.filterState = m.filterState.WithProject(projectID)
		m = m.applyFilterToCurrentView()
	}
	return m, nil
}

// executeTagCommand handles the "tag" command
func (m Model) executeTagCommand(cmd *command.Command) (Model, tea.Cmd) {
	if len(cmd.Args) > 0 {
		tagName := strings.Join(cmd.Args, " ")
		// Resolve tag name to ID
		tags, err := m.service.GetTags()
		if err != nil {
			m.err = fmt.Errorf("failed to get tags: %w", err)
			return m, nil
		}

		// Find tag by name (case-insensitive)
		var tagID string
		tagNameLower := strings.ToLower(tagName)
		for _, tag := range tags {
			if strings.ToLower(tag.Name) == tagNameLower {
				tagID = tag.ID
				break
			}
		}

		if tagID == "" {
			m.err = fmt.Errorf("tag not found: %s", tagName)
			return m, nil
		}

		m.filterState = m.filterState.WithTag(tagID)
		m = m.applyFilterToCurrentView()
	}
	return m, nil
}

// executeDueCommand handles the "due" command
func (m Model) executeDueCommand(cmd *command.Command) (Model, tea.Cmd) {
	if len(cmd.Args) > 0 {
		dueFilter := cmd.Args[0]
		var df filter.DueFilter
		switch strings.ToLower(dueFilter) {
		case "today":
			df = filter.DueToday
		case "tomorrow":
			df = filter.DueTomorrow
		case "week":
			df = filter.DueWeek
		case "overdue":
			df = filter.DueOverdue
		default:
			df = filter.DueNone
		}
		m.filterState = m.filterState.WithDueFilter(df)
		m = m.applyFilterToCurrentView()
	}
	return m, nil
}

// executeFlaggedCommand handles the "flagged" command
func (m Model) executeFlaggedCommand() (Model, tea.Cmd) {
	m.filterState = m.filterState.WithFlaggedOnly(true)
	m = m.applyFilterToCurrentView()
	return m, nil
}

// executeClearCommand handles the "clear" command
func (m Model) executeClearCommand() (Model, tea.Cmd) {
	m.filterState = m.filterState.Clear()
	m = m.applyFilterToCurrentView()
	return m, nil
}

// applyFilterToCurrentView applies the current filter state to the active view
func (m Model) applyFilterToCurrentView() Model {
	switch m.currentView {
	case tui.ViewInbox:
		m.inboxView = m.inboxView.SetFilter(m.filterState)
	case tui.ViewForecast:
		m.forecastView = m.forecastView.SetFilter(m.filterState)
	case tui.ViewReview:
		m.reviewView = m.reviewView.SetFilter(m.filterState)
		// Projects and Tags views don't support filtering (they have their own navigation)
	}
	return m
}
