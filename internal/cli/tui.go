package cli

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pwojciechowski/lazyfocus/internal/app"
	"github.com/pwojciechowski/lazyfocus/internal/bridge"
	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/spf13/cobra"
)

// NewTUICommand creates the tui command
func NewTUICommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tui",
		Short: "Launch the interactive TUI",
		Long:  `Launch the interactive terminal user interface for managing OmniFocus tasks.`,
		RunE:  runTUI,
		Annotations: map[string]string{
			"skipServiceSetup": "true",
		},
	}

	return cmd
}

func runTUI(cmd *cobra.Command, args []string) error {
	// Create executor and service
	executor := bridge.NewOSAScriptExecutor()
	svc := service.NewOmniFocusService(executor, 30*time.Second)

	// Create app model
	model := app.NewApp(svc)

	// Create and run Bubble Tea program with alt screen
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}
