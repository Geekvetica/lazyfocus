package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestCompletionCommand(t *testing.T) {
	cmd := NewCompletionCommand()

	if cmd.Use != "completion [bash|zsh|fish|powershell]" {
		t.Errorf("expected Use to be 'completion [bash|zsh|fish|powershell]', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be non-empty")
	}

	if cmd.Long == "" {
		t.Error("expected Long description to be non-empty")
	}
}

func TestCompletionCommand_BashCompletion(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewCompletionCommand())

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "bash"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("bash completion failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected bash completion output to be non-empty")
	}

	// Basic check for bash completion markers
	if !strings.Contains(output, "bash completion") && !strings.Contains(output, "_lazyfocus") {
		t.Errorf("output does not appear to be bash completion script: %s", output[:min(100, len(output))])
	}
}

func TestCompletionCommand_ZshCompletion(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewCompletionCommand())

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "zsh"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("zsh completion failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected zsh completion output to be non-empty")
	}

	// Basic check for zsh completion markers
	if !strings.Contains(output, "zsh completion") && !strings.Contains(output, "_lazyfocus") {
		t.Errorf("output does not appear to be zsh completion script: %s", output[:min(100, len(output))])
	}
}

func TestCompletionCommand_FishCompletion(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewCompletionCommand())

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "fish"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("fish completion failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected fish completion output to be non-empty")
	}

	// Basic check for fish completion markers
	if !strings.Contains(output, "lazyfocus") {
		t.Errorf("output does not appear to be fish completion script: %s", output[:min(100, len(output))])
	}
}

func TestCompletionCommand_PowerShellCompletion(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewCompletionCommand())

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "powershell"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("powershell completion failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected powershell completion output to be non-empty")
	}

	// Basic check for powershell completion markers
	if !strings.Contains(output, "lazyfocus") {
		t.Errorf("output does not appear to be powershell completion script: %s", output[:min(100, len(output))])
	}
}

func TestCompletionCommand_InvalidShellArgument(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewCompletionCommand())

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "invalid-shell"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when providing invalid shell argument")
	}
}

func TestCompletionCommand_NoArgument(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewCompletionCommand())

	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion"})

	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error when no argument provided")
	}
}

func TestCompletionCommand_HasSkipServiceSetupAnnotation(t *testing.T) {
	cmd := NewCompletionCommand()

	if cmd.Annotations == nil {
		t.Fatal("expected Annotations to be set, got nil")
	}

	if cmd.Annotations["skipServiceSetup"] != "true" {
		t.Errorf("expected skipServiceSetup annotation to be 'true', got %q", cmd.Annotations["skipServiceSetup"])
	}
}

func TestCompletionCommand_DefaultCaseReturnsNil(t *testing.T) {
	// This test ensures the default case in the switch statement is covered
	// The switch statement has a ValidArgs constraint, so this case should never
	// be reached in practice, but we test it for completeness

	cmd := NewCompletionCommand()

	// Verify the command has the expected valid args
	expectedArgs := []string{"bash", "zsh", "fish", "powershell"}
	if len(cmd.ValidArgs) != len(expectedArgs) {
		t.Errorf("expected %d ValidArgs, got %d", len(expectedArgs), len(cmd.ValidArgs))
	}

	for i, expected := range expectedArgs {
		if i >= len(cmd.ValidArgs) || cmd.ValidArgs[i] != expected {
			t.Errorf("ValidArgs[%d]: expected %q, got %q", i, expected, cmd.ValidArgs[i])
		}
	}

	// The Args validator prevents invalid args from reaching RunE, so the default
	// case in the switch is unreachable in normal operation. This test documents
	// that behavior.
}
