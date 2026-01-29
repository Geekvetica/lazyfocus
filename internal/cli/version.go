package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - can be set at build time with -ldflags
var (
	Version   = "0.1.0"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  `Print version information for lazyfocus.`,
		Args:  cobra.NoArgs,
		Annotations: map[string]string{
			"skipServiceSetup": "true",
		},
		Run: runVersion,
	}

	return cmd
}

func runVersion(cmd *cobra.Command, args []string) {
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "lazyfocus version %s\n", Version)
	if BuildDate != "unknown" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Build date: %s\n", BuildDate)
	}
	if GitCommit != "unknown" {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Git commit: %s\n", GitCommit)
	}
}
