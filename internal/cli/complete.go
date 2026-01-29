package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCompleteCommand creates the complete command
func NewCompleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "complete <task-id> [task-id...]",
		Short: "Mark tasks as complete in OmniFocus",
		Long: `Mark one or more tasks as complete in OmniFocus.

Accepts one or more task IDs as arguments. The command will attempt to
complete all specified tasks, continuing even if some fail.

Examples:
  lazyfocus complete abc123
  lazyfocus complete abc123 def456
  lazyfocus complete task1 task2 task3 --json`,
		Args: cobra.MinimumNArgs(1),
		RunE: runComplete,
	}

	return cmd
}

func runComplete(cmd *cobra.Command, args []string) error {
	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Track if any errors occurred
	var lastError error
	successCount := 0

	// Attempt to complete each task
	for _, taskID := range args {
		result, err := svc.CompleteTask(taskID)
		if err != nil {
			lastError = err
			// In non-quiet mode, show the error
			if !GetQuietFlag() {
				formatter := getFormatter()
				cmd.Print(formatter.FormatError(fmt.Errorf("failed to complete %s: %w", taskID, err)))
			}
			continue
		}

		successCount++

		// Format and output result
		if !GetQuietFlag() {
			formatter := getFormatter()
			outputStr := formatter.FormatCompletedTask(*result)
			cmd.Print(outputStr)
		}
	}

	// If all tasks failed, return the last error
	if successCount == 0 && lastError != nil {
		return lastError
	}

	return nil
}
