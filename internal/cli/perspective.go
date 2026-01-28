package cli

import (
	"github.com/pwojciechowski/lazyfocus/internal/cli/output"
	"github.com/spf13/cobra"
)

// NewPerspectiveCommand creates the perspective command
func NewPerspectiveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "perspective <name>",
		Short: "Show tasks from a perspective",
		Long: `Show tasks from a named OmniFocus perspective.

Note: Custom perspectives require OmniFocus Pro.`,
		Args: cobra.ExactArgs(1),
		RunE: runPerspective,
	}

	return cmd
}

func runPerspective(cmd *cobra.Command, args []string) error {
	perspectiveName := args[0]

	svc := getService()

	tasks, err := svc.GetPerspectiveTasks(perspectiveName)
	if err != nil {
		return handleError(cmd, err)
	}

	if GetQuietFlag() {
		return nil
	}

	formatter := getFormatter()
	options := output.TaskFormatOptions{
		ShowProject: true,
		ShowTags:    true,
	}

	cmd.Print(formatter.FormatTasks(tasks, options))
	return nil
}
