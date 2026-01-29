package cli

import (
	"github.com/pwojciechowski/lazyfocus/internal/cli/output"
	"github.com/spf13/cobra"
)

// NewTagsCommand creates the tags command
func NewTagsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List tags from OmniFocus",
		Long: `List tags from OmniFocus with optional hierarchy and task counts.

By default, shows tags with hierarchy. Use --flat to show tags in a flat list.
Use --with-counts to include task counts for each tag.`,
		RunE: runTags,
	}

	cmd.Flags().Bool("flat", false, "Show tags in flat list (no hierarchy)")
	cmd.Flags().Bool("with-counts", false, "Show task count per tag")

	return cmd
}

func runTags(cmd *cobra.Command, args []string) error {
	// Get flag values
	flatFlag, _ := cmd.Flags().GetBool("flat")
	withCountsFlag, _ := cmd.Flags().GetBool("with-counts")

	// Get service
	svc, err := getServiceFromCmd(cmd)
	if err != nil {
		return handleError(cmd, err)
	}

	// Get tags from service
	tags, getErr := svc.GetTags()
	if getErr != nil {
		return handleError(cmd, getErr)
	}

	// Get tag counts if requested
	if withCountsFlag {
		_, countErr := svc.GetTagCounts()
		if countErr != nil {
			return handleError(cmd, countErr)
		}
		// TODO: Pass counts to formatter when implementing count display
	}

	// Format and output results
	if GetQuietFlag() {
		// Quiet mode: no output, just exit code
		return nil
	}

	formatOptions := output.TagFormatOptions{
		Flat:       flatFlag,
		ShowCounts: withCountsFlag,
	}

	formatter := getFormatter()
	outputStr := formatter.FormatTags(tags, formatOptions)
	cmd.Print(outputStr)

	return nil
}
