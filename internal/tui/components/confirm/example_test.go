package confirm_test

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/tui"
	"github.com/pwojciechowski/lazyfocus/internal/tui/components/confirm"
)

// Example usage of the confirmation modal component
func ExampleModel_basic() {
	// Create a confirmation modal
	styles := tui.DefaultStyles()
	modal := confirm.New(styles)

	// Show the modal with a title and message
	modal = modal.Show("Delete Task", "Are you sure you want to delete this task?")

	// In your main app's Update function, handle the messages:
	// case confirm.ConfirmedMsg:
	//     // User confirmed - proceed with deletion
	// case confirm.CancelledMsg:
	//     // User cancelled - do nothing
}

// Example showing how to pass context through the modal
func ExampleModel_withContext() {
	styles := tui.DefaultStyles()
	modal := confirm.New(styles)

	// Pass a task ID as context
	taskID := "task123"
	modal = modal.ShowWithContext(
		"Delete Task",
		"Are you sure?",
		taskID,
	)

	// In your Update function:
	// case confirm.ConfirmedMsg:
	//     taskID := msg.Context.(string)
	//     // Delete the task with this ID
}

// Example integration in a Bubble Tea application
type exampleApp struct {
	confirmModal confirm.Model
	taskToDelete string
}

func (a exampleApp) Init() tea.Cmd {
	return nil
}

func (a exampleApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// When user presses 'd' to delete
		if msg.String() == "d" && !a.confirmModal.IsVisible() {
			a.taskToDelete = "some-task-id"
			a.confirmModal = a.confirmModal.ShowWithContext(
				"Delete Task",
				"Are you sure you want to delete this task?",
				a.taskToDelete,
			)
			return a, nil
		}

	case confirm.ConfirmedMsg:
		// User confirmed deletion
		taskID := msg.Context.(string)
		fmt.Printf("Deleting task: %s\n", taskID)
		// Perform actual deletion here
		return a, nil

	case confirm.CancelledMsg:
		// User cancelled
		return a, nil
	}

	// Update the modal
	var cmd tea.Cmd
	a.confirmModal, cmd = a.confirmModal.Update(msg)
	return a, cmd
}

func (a exampleApp) View() string {
	view := "Main app view\n"

	// Render the modal overlay if visible
	if a.confirmModal.IsVisible() {
		view += a.confirmModal.View()
	}

	return view
}
