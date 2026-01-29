package cli

import (
	"fmt"
	"strconv"

	"github.com/pwojciechowski/lazyfocus/internal/cli/dateparse"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

// NewModifyCommand creates the modify command
func NewModifyCommand() *cobra.Command {
	var (
		nameFlag       string
		noteFlag       string
		projectFlag    string
		addTagFlags    []string
		removeTagFlag  []string
		dueFlag        string
		deferFlag      string
		flaggedFlag    string
		clearDueFlag   bool
		clearDeferFlag bool
	)

	cmd := &cobra.Command{
		Use:   "modify <task-id> [flags]",
		Short: "Modify an existing task in OmniFocus",
		Long: `Modify an existing task in OmniFocus.

Requires exactly one task ID as argument. Use flags to specify which
fields to modify. At least one modification flag is required.

Note: Due to OmniFocus automation limitations, only the first tag specified
with --add-tag will be applied to the task (as the primary tag). Multiple
tags can be specified but only the first will be used. Using --remove-tag
will only remove the primary tag if it matches.

Examples:
  lazyfocus modify task123 --name "New name"
  lazyfocus modify task123 --due tomorrow --flagged true
  lazyfocus modify task123 --add-tag urgent --remove-tag low
  lazyfocus modify task123 --clear-due
  lazyfocus modify task123 --project Work --note "Updated note"`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runModify(cmd, args, nameFlag, noteFlag, projectFlag, addTagFlags, removeTagFlag,
				dueFlag, deferFlag, flaggedFlag, clearDueFlag, clearDeferFlag)
		},
	}

	cmd.Flags().StringVar(&nameFlag, "name", "", "New name")
	cmd.Flags().StringVar(&noteFlag, "note", "", "New note")
	cmd.Flags().StringVar(&projectFlag, "project", "", "Move to project (name or ID)")
	cmd.Flags().StringSliceVar(&addTagFlags, "add-tag", []string{}, "Add tag (repeatable)")
	cmd.Flags().StringSliceVar(&removeTagFlag, "remove-tag", []string{}, "Remove tag (repeatable)")
	cmd.Flags().StringVar(&dueFlag, "due", "", "Set due date")
	cmd.Flags().StringVar(&deferFlag, "defer", "", "Set defer date")
	cmd.Flags().StringVar(&flaggedFlag, "flagged", "", "Set flagged (true/false)")
	cmd.Flags().BoolVar(&clearDueFlag, "clear-due", false, "Clear due date")
	cmd.Flags().BoolVar(&clearDeferFlag, "clear-defer", false, "Clear defer date")

	return cmd
}

func runModify(cmd *cobra.Command, args []string, nameFlag, noteFlag, projectFlag string,
	addTagFlags, removeTagFlags []string, dueFlag, deferFlag, flaggedFlag string,
	clearDueFlag, clearDeferFlag bool) error {

	taskID := args[0]

	// Build TaskModification from flags
	mod, err := buildModificationFromFlags(nameFlag, noteFlag, projectFlag, addTagFlags, removeTagFlags,
		dueFlag, deferFlag, flaggedFlag, clearDueFlag, clearDeferFlag)
	if err != nil {
		return handleError(cmd, err)
	}

	// Check that at least one modification is specified
	if mod.IsEmpty() {
		return handleError(cmd, fmt.Errorf("no modifications specified"))
	}

	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Resolve project name to ID if needed
	if mod.ProjectID != nil && *mod.ProjectID != "" {
		projectID, err := svc.ResolveProjectName(*mod.ProjectID)
		if err != nil {
			return handleError(cmd, fmt.Errorf("failed to resolve project: %w", err))
		}
		mod.ProjectID = &projectID
	}

	// Modify the task
	task, err := svc.ModifyTask(taskID, mod)
	if err != nil {
		return handleError(cmd, fmt.Errorf("failed to modify task: %w", err))
	}

	// Format and output results
	if GetQuietFlag() {
		return nil
	}

	formatter := getFormatter()
	outputStr := formatter.FormatModifiedTask(*task)
	cmd.Print(outputStr)

	return nil
}

// buildModificationFromFlags constructs a TaskModification from command-line flags.
func buildModificationFromFlags(nameFlag, noteFlag, projectFlag string,
	addTagFlags, removeTagFlags []string, dueFlag, deferFlag, flaggedFlag string,
	clearDueFlag, clearDeferFlag bool) (domain.TaskModification, error) {

	mod := domain.TaskModification{
		AddTags:    addTagFlags,
		RemoveTags: removeTagFlags,
		ClearDue:   clearDueFlag,
		ClearDefer: clearDeferFlag,
	}

	if nameFlag != "" {
		mod.Name = &nameFlag
	}

	if noteFlag != "" {
		mod.Note = &noteFlag
	}

	if projectFlag != "" {
		mod.ProjectID = &projectFlag
	}

	if dueFlag != "" {
		dueDate, err := dateparse.Parse(dueFlag)
		if err != nil {
			return domain.TaskModification{}, fmt.Errorf("invalid due date: %w", err)
		}
		mod.DueDate = &dueDate
	}

	if deferFlag != "" {
		deferDate, err := dateparse.Parse(deferFlag)
		if err != nil {
			return domain.TaskModification{}, fmt.Errorf("invalid defer date: %w", err)
		}
		mod.DeferDate = &deferDate
	}

	if flaggedFlag != "" {
		flaggedBool, err := strconv.ParseBool(flaggedFlag)
		if err != nil {
			return domain.TaskModification{}, fmt.Errorf("invalid flagged value (use true/false): %w", err)
		}
		mod.Flagged = &flaggedBool
	}

	return mod, nil
}
