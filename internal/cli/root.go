// Package cli provides the command-line interface for LazyFocus.
package cli

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	quietMode  bool
	timeout    time.Duration
)

// NewRootCommand creates the root cobra command for lazyfocus
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lazyfocus",
		Short: "CLI interface for OmniFocus",
		Long: `LazyFocus (lf) is a CLI and TUI tool for interacting with OmniFocus on macOS.

It provides both human-readable output for terminal use and JSON output for
scripting and AI agent integration.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	cmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	cmd.PersistentFlags().BoolVar(&quietMode, "quiet", false, "Suppress output, exit codes only")
	cmd.PersistentFlags().DurationVar(&timeout, "timeout", 30*time.Second, "Timeout for OmniFocus operations")

	return cmd
}

// GetJSONFlag returns the value of the --json flag
func GetJSONFlag() bool {
	return jsonOutput
}

// GetQuietFlag returns the value of the --quiet flag
func GetQuietFlag() bool {
	return quietMode
}

// GetTimeoutFlag returns the value of the --timeout flag
func GetTimeoutFlag() time.Duration {
	return timeout
}
