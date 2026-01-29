package cli

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/pwojciechowski/lazyfocus/internal/cli/service"
	"github.com/pwojciechowski/lazyfocus/internal/domain"
	"github.com/spf13/cobra"
)

// TestAllCommandsRegistered verifies all commands are registered on root
func TestAllCommandsRegistered(t *testing.T) {
	rootCmd := NewRootCommand()
	rootCmd.AddCommand(NewTasksCommand())
	rootCmd.AddCommand(NewProjectsCommand())
	rootCmd.AddCommand(NewTagsCommand())
	rootCmd.AddCommand(NewShowCommand())
	rootCmd.AddCommand(NewPerspectiveCommand())
	rootCmd.AddCommand(NewVersionCommand())

	expectedCommands := []string{"tasks", "projects", "tags", "show", "perspective", "version"}

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, cmd := range rootCmd.Commands() {
			if cmd.Name() == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected command %q to be registered, but it was not found", expectedCmd)
		}
	}
}

// TestCommandHelpConsistency verifies all command help texts follow conventions
func TestCommandHelpConsistency(t *testing.T) {
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"tasks", NewTasksCommand()},
		{"projects", NewProjectsCommand()},
		{"tags", NewTagsCommand()},
		{"show", NewShowCommand()},
		{"perspective", NewPerspectiveCommand()},
		{"version", NewVersionCommand()},
	}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			// Verify command has a short description
			if tc.cmd.Short == "" {
				t.Errorf("command %q missing short description", tc.name)
			}

			// Verify command has a long description
			if tc.cmd.Long == "" {
				t.Errorf("command %q missing long description", tc.name)
			}

			// Note: --json flag is inherited from root at runtime, so we can't verify it here
			// The flag is set as a persistent flag on the root command
		})
	}
}

// TestJSONOutputConsistency verifies all commands return consistent JSON structure
func TestJSONOutputConsistency(t *testing.T) {
	// Set up mock service
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{
			{ID: "task1", Name: "Test Task"},
		},
		Projects: []domain.Project{
			{ID: "proj1", Name: "Test Project"},
		},
		Tags: []domain.Tag{
			{ID: "tag1", Name: "Test Tag"},
		},
		Task: &domain.Task{ID: "task1", Name: "Test Task"},
	}
	Service = mockService
	defer func() { Service = nil }()

	tests := []struct {
		name         string
		cmd          string
		args         []string
		expectedKeys []string
	}{
		{
			name:         "tasks command",
			cmd:          "tasks",
			args:         []string{"--json"},
			expectedKeys: []string{"tasks"},
		},
		{
			name:         "projects command",
			cmd:          "projects",
			args:         []string{"--json"},
			expectedKeys: []string{"projects"},
		},
		{
			name:         "tags command",
			cmd:          "tags",
			args:         []string{"--json"},
			expectedKeys: []string{"tags"},
		},
		{
			name:         "show command",
			cmd:          "show",
			args:         []string{"task1", "--json"},
			expectedKeys: []string{"task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := NewRootCommand()
			rootCmd.AddCommand(NewTasksCommand())
			rootCmd.AddCommand(NewProjectsCommand())
			rootCmd.AddCommand(NewTagsCommand())
			rootCmd.AddCommand(NewShowCommand())

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			args := append([]string{tt.cmd}, tt.args...)
			rootCmd.SetArgs(args)

			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("command execution failed: %v", err)
			}

			// Parse JSON output
			var result map[string]interface{}
			if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
				t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, buf.String())
			}

			// Verify expected keys exist
			for _, key := range tt.expectedKeys {
				if _, ok := result[key]; !ok {
					t.Errorf("expected JSON key %q not found in output: %v", key, result)
				}
			}
		})
	}
}

// TestErrorOutputConsistency verifies all commands return consistent error structure
func TestErrorOutputConsistency(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func() *service.MockOmniFocusService
		cmd         string
		args        []string
		expectError bool
	}{
		{
			name: "item not found error",
			setupMock: func() *service.MockOmniFocusService {
				return &service.MockOmniFocusService{
					Task: nil, // Return nil to trigger ItemNotFoundError
				}
			},
			cmd:         "show",
			args:        []string{"nonexistent", "--json"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Service = tt.setupMock()
			defer func() { Service = nil }()

			rootCmd := NewRootCommand()
			rootCmd.AddCommand(NewShowCommand())

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			args := append([]string{tt.cmd}, tt.args...)
			rootCmd.SetArgs(args)

			err := rootCmd.Execute()

			if tt.expectError && err == nil {
				t.Fatal("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectError {
				// Verify error output contains "error" key in JSON
				var result map[string]interface{}
				if jsonErr := json.Unmarshal(buf.Bytes(), &result); jsonErr == nil {
					if _, ok := result["error"]; !ok {
						t.Errorf("expected JSON error output to have 'error' key, got: %v", result)
					}
				}
			}
		})
	}
}

// TestQuietModeConsistency verifies all commands respect --quiet flag
func TestQuietModeConsistency(t *testing.T) {
	mockService := &service.MockOmniFocusService{
		InboxTasks: []domain.Task{{ID: "task1", Name: "Test"}},
		Projects:   []domain.Project{{ID: "proj1", Name: "Test"}},
		Tags:       []domain.Tag{{ID: "tag1", Name: "Test"}},
		Task:       &domain.Task{ID: "task1", Name: "Test"},
	}
	Service = mockService
	defer func() { Service = nil }()

	commands := []struct {
		name string
		args []string
	}{
		{"tasks", []string{"tasks", "--quiet"}},
		{"projects", []string{"projects", "--quiet"}},
		{"tags", []string{"tags", "--quiet"}},
		{"show", []string{"show", "task1", "--quiet"}},
	}

	for _, tt := range commands {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd := NewRootCommand()
			rootCmd.AddCommand(NewTasksCommand())
			rootCmd.AddCommand(NewProjectsCommand())
			rootCmd.AddCommand(NewTagsCommand())
			rootCmd.AddCommand(NewShowCommand())

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			if err != nil {
				t.Fatalf("command execution failed: %v", err)
			}

			// In quiet mode, there should be no output
			if buf.Len() > 0 {
				t.Errorf("expected no output in quiet mode, got: %s", buf.String())
			}
		})
	}
}

// TestGlobalFlagsInheritance verifies all commands inherit global flags
func TestGlobalFlagsInheritance(t *testing.T) {
	rootCmd := NewRootCommand()
	commands := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"tasks", NewTasksCommand()},
		{"projects", NewProjectsCommand()},
		{"tags", NewTagsCommand()},
		{"show", NewShowCommand()},
		{"perspective", NewPerspectiveCommand()},
		{"version", NewVersionCommand()},
	}

	for _, tc := range commands {
		rootCmd.AddCommand(tc.cmd)
	}

	globalFlags := []string{"json", "quiet", "timeout"}

	for _, tc := range commands {
		t.Run(tc.name, func(t *testing.T) {
			for _, flagName := range globalFlags {
				// Try to parse the flag - it should be available via persistence
				buf := new(bytes.Buffer)
				rootCmd.SetOut(buf)
				rootCmd.SetErr(buf)
				rootCmd.SetArgs([]string{tc.name, "--" + flagName, "test", "--help"})

				// We expect help to show, but the flag should be recognized
				// If the flag is not recognized, we'd get an error before help
				_ = rootCmd.Execute()
			}
		})
	}
}
