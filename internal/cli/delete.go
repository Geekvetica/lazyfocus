package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewDeleteCommand creates the delete command
func NewDeleteCommand() *cobra.Command {
	var forceFlag bool

	cmd := &cobra.Command{
		Use:   "delete <task-id> [task-id...] [flags]",
		Short: "Delete tasks from OmniFocus",
		Long: `Delete one or more tasks from OmniFocus.

Accepts one or more task IDs as arguments. By default, prompts for confirmation
before deleting. Use --force to skip confirmation.

In JSON mode, confirmation is automatically skipped.

Examples:
  lazyfocus delete abc123 --force
  lazyfocus delete task1 task2 task3 --force
  lazyfocus delete abc123 --json`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(cmd, args, forceFlag)
		},
	}

	cmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "Skip confirmation")

	return cmd
}

func runDelete(cmd *cobra.Command, args []string, forceFlag bool) error {
	// Skip confirmation in JSON mode or quiet mode
	skipConfirmation := forceFlag || GetJSONFlag() || GetQuietFlag()

	// If not skipping, we would prompt here
	// For now, we require --force for non-interactive mode
	// In a real implementation, we'd use a prompt library
	if !skipConfirmation {
		// In a real CLI, we'd prompt the user here
		// For testing purposes and batch operations, we require --force
		return fmt.Errorf("confirmation required: use --force to delete without confirmation")
	}

	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Track if any errors occurred
	var lastError error
	successCount := 0

	// Attempt to delete each task
	for _, taskID := range args {
		result, err := svc.DeleteTask(taskID)
		if err != nil {
			lastError = err
			// In non-quiet mode, show the error
			if !GetQuietFlag() {
				formatter := getFormatter()
				cmd.Print(formatter.FormatError(fmt.Errorf("failed to delete %s: %w", taskID, err)))
			}
			continue
		}

		successCount++

		// Format and output result
		if !GetQuietFlag() {
			formatter := getFormatter()
			outputStr := formatter.FormatDeletedTask(*result)
			cmd.Print(outputStr)
		}
	}

	// If all tasks failed, return the last error
	if successCount == 0 && lastError != nil {
		return lastError
	}

	return nil
}
