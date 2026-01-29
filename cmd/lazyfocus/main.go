// Package main provides the entry point for the lazyfocus CLI application.
package main

import (
	"context"
	"os"

	"github.com/pwojciechowski/lazyfocus/internal/cli"
	"github.com/pwojciechowski/lazyfocus/internal/cli/output"
)

func main() {
	rootCmd := cli.NewRootCommand()

	// Add commands
	rootCmd.AddCommand(cli.NewTasksCommand())
	rootCmd.AddCommand(cli.NewProjectsCommand())
	rootCmd.AddCommand(cli.NewTagsCommand())
	rootCmd.AddCommand(cli.NewShowCommand())
	rootCmd.AddCommand(cli.NewPerspectiveCommand())
	rootCmd.AddCommand(cli.NewVersionCommand())

	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		// Determine exit code based on error type
		exitCode := output.ExitGeneralError

		// Check for specific error types
		if itemNotFoundErr, ok := err.(*cli.ItemNotFoundError); ok {
			exitCode = itemNotFoundErr.ExitCode()
		}

		os.Exit(exitCode)
	}
}
