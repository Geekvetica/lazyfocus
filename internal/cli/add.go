package cli

import (
	"fmt"
	"strings"

	"github.com/pwojciechowski/lazyfocus/internal/cli/dateparse"
	"github.com/pwojciechowski/lazyfocus/internal/cli/taskparse"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

// NewAddCommand creates the add command
func NewAddCommand() *cobra.Command {
	var (
		projectFlag string
		tagFlags    []string
		dueFlag     string
		deferFlag   string
		flaggedFlag bool
		noteFlag    string
	)

	cmd := &cobra.Command{
		Use:   "add <task description>",
		Short: "Create a new task in OmniFocus",
		Long: `Create a new task in OmniFocus with natural syntax or flags.

Natural syntax in description:
  #tag        Add tag
  @project    Add to project (by name)
  due:xxx     Set due date
  defer:xxx   Set defer date
  !           Mark flagged

Command-line flags override natural syntax when both are present.

Examples:
  lazyfocus add "Buy milk #groceries"
  lazyfocus add "Call dentist" --due tomorrow
  lazyfocus add "Review PR @Work due:friday !"
  lazyfocus add "Meeting prep" --project Work --flagged --note "Prepare slides"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd, args, projectFlag, tagFlags, dueFlag, deferFlag, flaggedFlag, noteFlag)
		},
	}

	cmd.Flags().StringVarP(&projectFlag, "project", "p", "", "Project name or ID")
	cmd.Flags().StringSliceVarP(&tagFlags, "tag", "t", []string{}, "Tags (repeatable)")
	cmd.Flags().StringVarP(&dueFlag, "due", "d", "", "Due date")
	cmd.Flags().StringVar(&deferFlag, "defer", "", "Defer date")
	cmd.Flags().BoolVarP(&flaggedFlag, "flagged", "f", false, "Mark flagged")
	cmd.Flags().StringVarP(&noteFlag, "note", "n", "", "Task note")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string, projectFlag string, tagFlags []string, dueFlag, deferFlag string, flaggedFlag bool, noteFlag string) error {
	// Combine all args into a single task description
	taskDescription := strings.Join(args, " ")

	// Parse the task description with natural syntax
	taskInput, err := taskparse.Parse(taskDescription)
	if err != nil {
		return handleError(cmd, fmt.Errorf("failed to parse task: %w", err))
	}

	// Apply command-line flags (flags take precedence over natural syntax)
	if err := applyAddFlags(cmd, &taskInput, projectFlag, tagFlags, dueFlag, deferFlag, flaggedFlag, noteFlag); err != nil {
		return handleError(cmd, err)
	}

	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Resolve project name to ID if needed
	if taskInput.ProjectName != "" && taskInput.ProjectID == "" {
		projectID, err := svc.ResolveProjectName(taskInput.ProjectName)
		if err != nil {
			return handleError(cmd, fmt.Errorf("failed to resolve project: %w", err))
		}
		taskInput.ProjectID = projectID
	}

	// Create the task
	task, err := svc.CreateTask(taskInput)
	if err != nil {
		return handleError(cmd, fmt.Errorf("failed to create task: %w", err))
	}

	// Format and output results
	if GetQuietFlag() {
		return nil
	}

	formatter := getFormatter()
	outputStr := formatter.FormatCreatedTask(*task)
	cmd.Print(outputStr)

	return nil
}

// applyAddFlags applies command-line flags to TaskInput, overriding natural syntax values.
func applyAddFlags(cmd *cobra.Command, taskInput *domain.TaskInput, projectFlag string, tagFlags []string, dueFlag, deferFlag string, flaggedFlag bool, noteFlag string) error {
	if noteFlag != "" {
		taskInput.Note = noteFlag
	}

	if projectFlag != "" {
		taskInput.ProjectName = projectFlag
		taskInput.ProjectID = "" // Will be resolved later
	}

	if len(tagFlags) > 0 {
		taskInput.TagNames = tagFlags
	}

	if dueFlag != "" {
		dueDate, err := dateparse.Parse(dueFlag)
		if err != nil {
			return fmt.Errorf("invalid due date: %w", err)
		}
		taskInput.DueDate = &dueDate
	}

	if deferFlag != "" {
		deferDate, err := dateparse.Parse(deferFlag)
		if err != nil {
			return fmt.Errorf("invalid defer date: %w", err)
		}
		taskInput.DeferDate = &deferDate
	}

	// Handle flagged flag (only override if explicitly set)
	if cmd.Flags().Changed("flagged") {
		taskInput.Flagged = &flaggedFlag
	}

	return nil
}
