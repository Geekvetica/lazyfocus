package cli

import (
	"github.com/pwojciechowski/lazyfocus/internal/cli/output"
	"github.com/spf13/cobra"
)

// NewProjectsCommand creates the projects command
func NewProjectsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "projects",
		Short: "List projects from OmniFocus",
		Long: `List projects from OmniFocus with filtering options.

By default, shows active projects. Use --status flag to filter by status.`,
		RunE: runProjects,
	}

	cmd.Flags().String("status", "active", "Filter by status (active, on-hold, completed, dropped, all)")
	cmd.Flags().Bool("with-tasks", false, "Include nested tasks")

	return cmd
}

func runProjects(cmd *cobra.Command, args []string) error {
	// Get flag values
	statusFlag, _ := cmd.Flags().GetString("status")
	withTasksFlag, _ := cmd.Flags().GetBool("with-tasks")

	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Get projects from service
	projects, getErr := svc.GetProjects(statusFlag)
	if getErr != nil {
		return handleError(cmd, getErr)
	}

	// Format and output results
	if GetQuietFlag() {
		// Quiet mode: no output, just exit code
		return nil
	}

	formatOptions := output.ProjectFormatOptions{
		ShowTasks: withTasksFlag,
		ShowNotes: false,
	}

	formatter := getFormatter()
	outputStr := formatter.FormatProjects(projects, formatOptions)
	cmd.Print(outputStr)

	return nil
}
