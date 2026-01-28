package cli

import (
	"testing"
	"time"
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
