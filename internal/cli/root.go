// Package cli provides the command-line interface for LazyFocus.
package cli

import (
	"context"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/bridge"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
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
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip service setup for commands that don't need it (like version, help)
			if cmd.Name() == "version" || cmd.Name() == "help" {
				return nil
			}

			// Get current context, use background if nil
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}

			// Check if service is already in context (e.g., from tests)
			if _, err := ServiceFromContext(ctx); err == nil {
				// Service already exists, skip setup
				return nil
			}

			// Create executor and service
			executor := bridge.NewOSAScriptExecutor()
			svc := service.NewOmniFocusService(executor, GetTimeoutFlag())

			// Inject service into context
			ctx = ContextWithService(ctx, svc)
			cmd.SetContext(ctx)

			return nil
		},
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
