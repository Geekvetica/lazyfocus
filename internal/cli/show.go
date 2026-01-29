package cli

import (
	"fmt"

	"github.com/pwojciechowski/lazyfocus/internal/cli/output"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

// ItemNotFoundError is returned when an item (task, project, or tag) is not found
type ItemNotFoundError struct {
	ID string
}

func (e *ItemNotFoundError) Error() string {
	return fmt.Sprintf("item not found: %s", e.ID)
}

// ExitCode returns the exit code for this error
func (e *ItemNotFoundError) ExitCode() int {
	return output.ExitItemNotFound
}

// NewShowCommand creates the show command
func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show <id>",
		Short: "Show details for a task, project, or tag",
		Long: `Show detailed information for a specific item by its ID.

The command will attempt to auto-detect the type of item (task, project, or tag)
unless you specify the type explicitly with --type flag.

Examples:
  lazyfocus show abc123              # Auto-detect type
  lazyfocus show abc123 --type task  # Show as task
  lazyfocus show abc123 --json       # Output as JSON`,
		Args: cobra.ExactArgs(1),
		RunE: runShow,
	}

	cmd.Flags().String("type", "", "Item type: task, project, or tag (auto-detect if not specified)")

	return cmd
}

func runShow(cmd *cobra.Command, args []string) error {
	id := args[0]
	itemType, _ := cmd.Flags().GetString("type")

	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}
	formatter := getFormatter()

	switch itemType {
	case "task":
		return showTask(cmd, svc, formatter, id)
	case "project":
		return showProject(cmd, svc, formatter, id)
	case "tag":
		return showTag(cmd, svc, formatter, id)
	case "":
		return autoDetectAndShow(cmd, svc, formatter, id)
	default:
		return fmt.Errorf("unknown type: %s (valid types: task, project, tag)", itemType)
	}
}

func showTask(cmd *cobra.Command, svc service.OmniFocusService, formatter output.Formatter, id string) error {
	task, err := svc.GetTaskByID(id)
	if err != nil {
		return handleError(cmd, err)
	}

	if task == nil {
		return handleError(cmd, &ItemNotFoundError{ID: id})
	}

	return outputItem(cmd, formatter, *task)
}

func showProject(cmd *cobra.Command, svc service.OmniFocusService, formatter output.Formatter, id string) error {
	project, err := svc.GetProjectByID(id)
	if err != nil {
		return handleError(cmd, err)
	}

	if project == nil {
		return handleError(cmd, &ItemNotFoundError{ID: id})
	}

	return outputItem(cmd, formatter, *project)
}

func showTag(cmd *cobra.Command, svc service.OmniFocusService, formatter output.Formatter, id string) error {
	tag, err := svc.GetTagByID(id)
	if err != nil {
		return handleError(cmd, err)
	}

	if tag == nil {
		return handleError(cmd, &ItemNotFoundError{ID: id})
	}

	return outputItem(cmd, formatter, *tag)
}

func autoDetectAndShow(cmd *cobra.Command, svc service.OmniFocusService, formatter output.Formatter, id string) error {
	// Try task first
	task, err := svc.GetTaskByID(id)
	if err == nil && task != nil {
		return outputItem(cmd, formatter, *task)
	}

	// Try project
	project, err := svc.GetProjectByID(id)
	if err == nil && project != nil {
		return outputItem(cmd, formatter, *project)
	}

	// Try tag
	tag, err := svc.GetTagByID(id)
	if err == nil && tag != nil {
		return outputItem(cmd, formatter, *tag)
	}

	// Not found in any category
	return handleError(cmd, &ItemNotFoundError{ID: id})
}

func outputItem(cmd *cobra.Command, formatter output.Formatter, item interface{}) error {
	if GetQuietFlag() {
		// Quiet mode: no output, just exit code
		return nil
	}

	var outputStr string
	switch v := item.(type) {
	case domain.Task:
		outputStr = formatter.FormatTask(v)
	case domain.Project:
		outputStr = formatter.FormatProject(v)
	case domain.Tag:
		outputStr = formatter.FormatTag(v)
	default:
		return fmt.Errorf("unsupported item type: %T", item)
	}

	cmd.Print(outputStr)
	return nil
}
