package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd.Use != "version" {
		t.Errorf("expected Use to be 'version', got %q", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("expected Short description to be non-empty")
	}
}

func TestVersionCommandOutput(t *testing.T) {
	cmd := NewVersionCommand()

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "lazyfocus") {
		t.Errorf("expected output to contain 'lazyfocus', got: %s", output)
	}

	if !strings.Contains(output, "version") {
		t.Errorf("expected output to contain 'version', got: %s", output)
	}
}

func TestVersionCommandNoArgs(t *testing.T) {
	cmd := NewVersionCommand()

	// Version command should not accept any arguments
	cmd.SetArgs([]string{"unexpected-arg"})

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when providing arguments to version command")
	}
}

func TestVersionCommand_HasSkipServiceSetupAnnotation(t *testing.T) {
	cmd := NewVersionCommand()

	if cmd.Annotations == nil {
		t.Fatal("expected Annotations to be set, got nil")
	}

	if cmd.Annotations["skipServiceSetup"] != "true" {
		t.Errorf("expected skipServiceSetup annotation to be 'true', got %q", cmd.Annotations["skipServiceSetup"])
	}
}
