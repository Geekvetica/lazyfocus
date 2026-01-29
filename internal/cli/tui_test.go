package cli

import (
	"strings"
	"testing"
)

func TestTUICommand_IsRegistered(t *testing.T) {
	// Create root command and add TUI command
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTUICommand())

	// Find the tui command
	tuiCmd, _, err := rootCmd.Find([]string{"tui"})

	// Assert command was found
	if err != nil {
		t.Fatalf("Expected no error finding tui command, got: %v", err)
	}

	if tuiCmd == nil {
		t.Fatal("Expected tui command to be registered, got nil")
	}

	if tuiCmd.Use != "tui" {
		t.Errorf("Expected Use to be 'tui', got: %s", tuiCmd.Use)
	}
}

func TestTUICommand_HasCorrectMetadata(t *testing.T) {
	// Create root command and add TUI command
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTUICommand())

	// Find the tui command
	tuiCmd, _, err := rootCmd.Find([]string{"tui"})
	if err != nil {
		t.Fatalf("Expected no error finding tui command, got: %v", err)
	}

	// Assert metadata
	if tuiCmd.Use != "tui" {
		t.Errorf("Expected Use to be 'tui', got: %s", tuiCmd.Use)
	}

	expectedShort := "Launch the interactive TUI"
	if tuiCmd.Short != expectedShort {
		t.Errorf("Expected Short to be %q, got: %s", expectedShort, tuiCmd.Short)
	}

	if !strings.Contains(tuiCmd.Long, "interactive terminal user interface") {
		t.Errorf("Expected Long to contain 'interactive terminal user interface', got: %s", tuiCmd.Long)
	}

	if !strings.Contains(tuiCmd.Long, "OmniFocus") {
		t.Errorf("Expected Long to contain 'OmniFocus', got: %s", tuiCmd.Long)
	}
}

func TestTUICommand_HasRunFunction(t *testing.T) {
	// Create root command and add TUI command
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTUICommand())

	// Find the tui command
	tuiCmd, _, err := rootCmd.Find([]string{"tui"})
	if err != nil {
		t.Fatalf("Expected no error finding tui command, got: %v", err)
	}

	// Assert RunE is set (not nil)
	if tuiCmd.RunE == nil {
		t.Error("Expected RunE to be set, got nil")
	}
}

func TestTUICommand_SkipsServiceSetup(t *testing.T) {
	// Create root command and add TUI command
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTUICommand())

	// Find the tui command
	tuiCmd, _, err := rootCmd.Find([]string{"tui"})
	if err != nil {
		t.Fatalf("Expected no error finding tui command, got: %v", err)
	}

	// Assert skipServiceSetup annotation is set
	// The TUI command should skip service setup because it creates its own service
	if tuiCmd.Annotations == nil {
		t.Error("Expected Annotations to be set, got nil")
	}

	if value, exists := tuiCmd.Annotations["skipServiceSetup"]; !exists || value != "true" {
		t.Errorf("Expected skipServiceSetup annotation to be 'true', got: %s (exists: %v)", value, exists)
	}
}
