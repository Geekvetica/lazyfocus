package cli

import (
	"fmt"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/output"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

// NewTasksCommand creates the tasks command
func NewTasksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "List tasks from OmniFocus",
		Long: `List tasks from OmniFocus with various filtering options.

By default, shows inbox tasks. Use flags to filter by project, tag, due date, etc.`,
		RunE: runTasks,
	}

	cmd.Flags().Bool("inbox", false, "Show inbox tasks only")
	cmd.Flags().Bool("all", false, "Show all tasks")
	cmd.Flags().String("project", "", "Filter by project ID")
	cmd.Flags().String("tag", "", "Filter by tag ID")
	cmd.Flags().Bool("flagged", false, "Show flagged tasks only")
	cmd.Flags().String("due", "", "Show tasks due on/before date (supports 'today', 'tomorrow', or YYYY-MM-DD)")
	cmd.Flags().Bool("completed", false, "Include completed tasks")

	return cmd
}

func runTasks(cmd *cobra.Command, args []string) error {
	// Get flag values
	allFlag, _ := cmd.Flags().GetBool("all")
	projectFlag, _ := cmd.Flags().GetString("project")
	tagFlag, _ := cmd.Flags().GetString("tag")
	flaggedFlag, _ := cmd.Flags().GetBool("flagged")
	dueFlag, _ := cmd.Flags().GetString("due")
	completedFlag, _ := cmd.Flags().GetBool("completed")

	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Determine which service method to call based on flags
	var tasks []domain.Task

	switch {
	case flaggedFlag:
		tasks, err = svc.GetFlaggedTasks()
	case projectFlag != "":
		tasks, err = svc.GetTasksByProject(projectFlag)
	case tagFlag != "":
		tasks, err = svc.GetTasksByTag(tagFlag)
	case allFlag:
		filters := service.TaskFilters{
			Completed: completedFlag,
		}
		tasks, err = svc.GetAllTasks(filters)
	default:
		// Default to inbox (inbox flag is redundant with default behavior)
		tasks, err = svc.GetInboxTasks()
	}

	if err != nil {
		return handleError(cmd, err)
	}

	// Apply due date filter if specified
	if dueFlag != "" {
		tasks, err = filterTasksByDueDate(tasks, dueFlag)
		if err != nil {
			return handleError(cmd, err)
		}
	}

	// Format and output results
	if GetQuietFlag() {
		// Quiet mode: no output, just exit code
		return nil
	}

	formatOptions := output.TaskFormatOptions{
		ShowCompleted: completedFlag,
		ShowProject:   true,
		ShowTags:      true,
	}

	formatter := getFormatter()
	outputStr := formatter.FormatTasks(tasks, formatOptions)
	cmd.Print(outputStr)

	return nil
}

// getServiceFromCmd retrieves the service from the command context.
// Returns an error if the service is not found in context.
func getServiceFromCmd(cmd *cobra.Command) (service.OmniFocusService, error) {
	return ServiceFromContext(cmd.Context())
}

// getFormatter returns the appropriate formatter based on the --json flag
func getFormatter() output.Formatter {
	if GetJSONFlag() {
		return output.NewJSONFormatter()
	}
	return output.NewHumanFormatter()
}

// handleError handles errors and formats them appropriately
func handleError(cmd *cobra.Command, err error) error {
	if GetQuietFlag() {
		// In quiet mode, just return the error for exit code
		return err
	}

	formatter := getFormatter()
	cmd.Print(formatter.FormatError(err))

	return err
}

// filterTasksByDueDate filters tasks by due date
// Tasks with due dates on or before the specified date are included.
// Timezone handling: dates from OmniFocus come as UTC ISO strings and are
// compared correctly with local timezone dates using Go's time.Time comparison.
func filterTasksByDueDate(tasks []domain.Task, dueStr string) ([]domain.Task, error) {
	dueDate, err := parseDueDate(dueStr)
	if err != nil {
		return nil, fmt.Errorf("invalid due date format: %w", err)
	}

	var filtered []domain.Task
	for _, task := range tasks {
		if task.DueDate != nil && !task.DueDate.After(dueDate) {
			filtered = append(filtered, task)
		}
	}

	return filtered, nil
}

// parseDueDate parses a due date string (today, tomorrow, or YYYY-MM-DD)
// Returns a time at 23:59:59 in the local timezone to include all tasks due on that day
func parseDueDate(dueStr string) (time.Time, error) {
	now := time.Now()
	loc := now.Location()

	switch dueStr {
	case "today":
		return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc), nil
	case "tomorrow":
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 23, 59, 59, 0, loc), nil
	default:
		// Try to parse as YYYY-MM-DD in local timezone
		parsed, err := time.ParseInLocation("2006-01-02", dueStr, loc)
		if err != nil {
			return time.Time{}, fmt.Errorf("expected 'today', 'tomorrow', or YYYY-MM-DD format: %w", err)
		}
		return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, loc), nil
	}
}
