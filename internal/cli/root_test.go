package cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

func TestRootCommand_Flags(t *testing.T) {
	cmd := NewRootCommand()

	tests := []struct {
		name     string
		flagName string
		flagType string
	}{
		{
			name:     "json flag exists",
			flagName: "json",
			flagType: "bool",
		},
		{
			name:     "quiet flag exists",
			flagName: "quiet",
			flagType: "bool",
		},
		{
			name:     "timeout flag exists",
			flagName: "timeout",
			flagType: "duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.PersistentFlags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %q not found", tt.flagName)
				return
			}

			if flag.Value.Type() != tt.flagType {
				t.Errorf("Flag %q type = %v, want %v", tt.flagName, flag.Value.Type(), tt.flagType)
			}
		})
	}
}

func TestRootCommand_JSONFlagDefault(t *testing.T) {
	cmd := NewRootCommand()

	flag := cmd.PersistentFlags().Lookup("json")
	if flag == nil {
		t.Fatal("json flag not found")
	}

	if flag.DefValue != "false" {
		t.Errorf("json flag default = %v, want false", flag.DefValue)
	}
}

func TestRootCommand_QuietFlagDefault(t *testing.T) {
	cmd := NewRootCommand()

	flag := cmd.PersistentFlags().Lookup("quiet")
	if flag == nil {
		t.Fatal("quiet flag not found")
	}

	if flag.DefValue != "false" {
		t.Errorf("quiet flag default = %v, want false", flag.DefValue)
	}
}

func TestRootCommand_TimeoutFlagDefault(t *testing.T) {
	cmd := NewRootCommand()

	flag := cmd.PersistentFlags().Lookup("timeout")
	if flag == nil {
		t.Fatal("timeout flag not found")
	}

	// Default should be 30s
	if flag.DefValue != "30s" {
		t.Errorf("timeout flag default = %v, want 30s", flag.DefValue)
	}
}

func TestRootCommand_Metadata(t *testing.T) {
	cmd := NewRootCommand()

	if cmd.Use != "lazyfocus" && cmd.Use != "lf" {
		t.Errorf("Command Use = %v, want 'lazyfocus' or 'lf'", cmd.Use)
	}

	if cmd.Short == "" {
		t.Error("Command Short description is empty")
	}

	if cmd.Long == "" {
		t.Error("Command Long description is empty")
	}
}

func TestGetJSONFlag(t *testing.T) {
	cmd := NewRootCommand()

	// Test default value
	if GetJSONFlag() {
		t.Error("GetJSONFlag() = true, want false (default)")
	}

	// Set flag and test
	cmd.PersistentFlags().Set("json", "true")
	if !GetJSONFlag() {
		t.Error("GetJSONFlag() = false, want true after setting flag")
	}
}

func TestGetQuietFlag(t *testing.T) {
	cmd := NewRootCommand()

	// Test default value
	if GetQuietFlag() {
		t.Error("GetQuietFlag() = true, want false (default)")
	}

	// Set flag and test
	cmd.PersistentFlags().Set("quiet", "true")
	if !GetQuietFlag() {
		t.Error("GetQuietFlag() = false, want true after setting flag")
	}
}

func TestGetTimeoutFlag(t *testing.T) {
	cmd := NewRootCommand()

	// Test default value
	timeout := GetTimeoutFlag()
	if timeout != 30*time.Second {
		t.Errorf("GetTimeoutFlag() = %v, want 30s (default)", timeout)
	}

	// Set flag and test
	cmd.PersistentFlags().Set("timeout", "1m")
	timeout = GetTimeoutFlag()
	if timeout != 60*time.Second {
		t.Errorf("GetTimeoutFlag() = %v, want 1m after setting flag", timeout)
	}
}

func TestRootCommand_PersistentPreRunE_SkipsForVersionCommand(t *testing.T) {
	rootCmd := NewRootCommand()

	// Add the actual version command (which should have skipServiceSetup annotation)
	rootCmd.AddCommand(NewVersionCommand())

	// Execute the version command - PersistentPreRunE should skip service setup
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"version"})
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRootCommand_PersistentPreRunE_SkipsForCommandsWithAnnotation(t *testing.T) {
	rootCmd := NewRootCommand()

	// Create a test command with skipServiceSetup annotation
	testCmd := &cobra.Command{
		Use:   "testskip",
		Short: "Test skip command",
		Annotations: map[string]string{
			"skipServiceSetup": "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// Verify service was NOT set in context
			_, err := ServiceFromContext(cmd.Context())
			if err == nil {
				t.Error("Expected service to NOT be in context for commands with skipServiceSetup annotation")
			}
			return nil
		},
	}
	rootCmd.AddCommand(testCmd)

	// Execute the test command - should skip service setup
	rootCmd.SetArgs([]string{"testskip"})
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRootCommand_PersistentPreRunE_SkipsForHelpCommand(t *testing.T) {
	rootCmd := NewRootCommand()

	// Create a help subcommand (simulating the built-in help)
	helpCmd := &cobra.Command{
		Use:   "help",
		Short: "Help about any command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	rootCmd.AddCommand(helpCmd)

	// Execute the help command - PersistentPreRunE should skip service setup
	rootCmd.SetArgs([]string{"help"})
	err := rootCmd.Execute()

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestRootCommand_PersistentPreRunE_HandlesNilContext(t *testing.T) {
	rootCmd := NewRootCommand()

	// Create a test subcommand that will trigger PersistentPreRunE
	// but we won't provide a context (Execute without ExecuteContext)
	testCmd := &cobra.Command{
		Use:   "testcmd",
		Short: "Test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Verify that context was set (should be Background context)
			ctx := cmd.Context()
			if ctx == nil {
				t.Error("Expected context to be set, got nil")
			}
			return nil
		},
	}
	rootCmd.AddCommand(testCmd)

	// Execute without providing context - PersistentPreRunE should use Background context
	// Note: This will attempt to create a real executor, but that's OK for this test
	// since we're just testing that nil context doesn't cause a panic
	rootCmd.SetArgs([]string{"testcmd"})

	// We expect this may fail due to OmniFocus not being available, but it should
	// NOT panic due to nil context
	_ = rootCmd.Execute()
}

func TestRootCommand_PersistentPreRunE_UsesExistingServiceFromContext(t *testing.T) {
	rootCmd := NewRootCommand()

	// Create a mock service
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test task"},
		},
	}

	// Add tasks command for testing
	rootCmd.AddCommand(NewTasksCommand())

	// Track if service was replaced
	var serviceFromCommand service.OmniFocusService

	// Create a wrapper command to capture the service after PersistentPreRunE
	tasksCmd, _, _ := rootCmd.Find([]string{"tasks"})
	originalRunE := tasksCmd.RunE
	tasksCmd.RunE = func(cmd *cobra.Command, args []string) error {
		var err error
		serviceFromCommand, err = ServiceFromContext(cmd.Context())
		if err != nil {
			return err
		}
		return originalRunE(cmd, args)
	}

	// Set output buffer
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Execute with context that already has a service
	ctx := ContextWithService(context.Background(), mockService)
	rootCmd.SetArgs([]string{"tasks"})
	err := rootCmd.ExecuteContext(ctx)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify that the original mock service was used (not replaced)
	if serviceFromCommand != mockService {
		t.Error("Expected PersistentPreRunE to use existing service from context, but it was replaced")
	}
}
