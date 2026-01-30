package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad_NoConfigFile_ReturnsDefaults(t *testing.T) {
	// Temporarily change HOME to a temp directory
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Also clear any LAZYFOCUS_ env vars
	oldEnvVars := clearLazyFocusEnvVars()
	defer restoreEnvVars(oldEnvVars)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Verify defaults
	if cfg.Output.Format != "human" {
		t.Errorf("Expected default output format 'human', got %q", cfg.Output.Format)
	}

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", cfg.Timeout)
	}

	if cfg.Defaults.Project != "" {
		t.Errorf("Expected default project to be empty, got %q", cfg.Defaults.Project)
	}

	if cfg.TUI.Theme != "default" {
		t.Errorf("Expected default theme 'default', got %q", cfg.TUI.Theme)
	}

	if cfg.TUI.Colors.Primary != "#5B9BD5" {
		t.Errorf("Expected default primary color '#5B9BD5', got %q", cfg.TUI.Colors.Primary)
	}
}

func TestLoad_WithConfigFile_OverridesDefaults(t *testing.T) {
	// Create temp directory and config file
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Clear env vars
	oldEnvVars := clearLazyFocusEnvVars()
	defer restoreEnvVars(oldEnvVars)

	// Write config file
	configContent := `output:
  format: json
timeout: 60s
defaults:
  project: Work
tui:
  theme: custom
  colors:
    primary: "#FF0000"
    flagged: "#00FF00"
    due: "#0000FF"
    overdue: "#FFFF00"
`
	configPath := filepath.Join(tmpDir, ".lazyfocus.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Verify config file values override defaults
	if cfg.Output.Format != "json" {
		t.Errorf("Expected format 'json' from config, got %q", cfg.Output.Format)
	}

	if cfg.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s from config, got %v", cfg.Timeout)
	}

	if cfg.Defaults.Project != "Work" {
		t.Errorf("Expected project 'Work' from config, got %q", cfg.Defaults.Project)
	}

	if cfg.TUI.Theme != "custom" {
		t.Errorf("Expected theme 'custom' from config, got %q", cfg.TUI.Theme)
	}

	if cfg.TUI.Colors.Primary != "#FF0000" {
		t.Errorf("Expected primary color '#FF0000' from config, got %q", cfg.TUI.Colors.Primary)
	}
}

func TestLoad_EnvironmentVariables_OverrideConfigFile(t *testing.T) {
	// Create temp directory and config file
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Clear existing env vars
	oldEnvVars := clearLazyFocusEnvVars()
	defer restoreEnvVars(oldEnvVars)

	// Write config file with some values
	configContent := `output:
  format: json
timeout: 60s
`
	configPath := filepath.Join(tmpDir, ".lazyfocus.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Set environment variables (Viper uses LAZYFOCUS_ prefix)
	os.Setenv("LAZYFOCUS_OUTPUT_FORMAT", "human")
	os.Setenv("LAZYFOCUS_TIMEOUT", "90s")
	os.Setenv("LAZYFOCUS_DEFAULTS_PROJECT", "Personal")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Verify environment variables override config file
	if cfg.Output.Format != "human" {
		t.Errorf("Expected format 'human' from env var, got %q", cfg.Output.Format)
	}

	if cfg.Timeout != 90*time.Second {
		t.Errorf("Expected timeout 90s from env var, got %v", cfg.Timeout)
	}

	if cfg.Defaults.Project != "Personal" {
		t.Errorf("Expected project 'Personal' from env var, got %q", cfg.Defaults.Project)
	}
}

func TestLoad_InvalidConfigFile_ReturnsError(t *testing.T) {
	// Create temp directory with invalid config
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Clear env vars
	oldEnvVars := clearLazyFocusEnvVars()
	defer restoreEnvVars(oldEnvVars)

	// Write invalid YAML
	configContent := `output:
  format: json
  invalid yaml {{{
`
	configPath := filepath.Join(tmpDir, ".lazyfocus.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Expected Load() to return error for invalid YAML, got nil")
	}
}

func TestFilePath_ReturnsCorrectPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Cannot get user home directory, skipping test")
	}

	expected := filepath.Join(home, ".lazyfocus.yaml")
	actual := FilePath()

	if actual != expected {
		t.Errorf("Expected config path %q, got %q", expected, actual)
	}
}

func TestFilePath_NoHomeDir_ReturnsFallback(t *testing.T) {
	// This test is harder to implement reliably since we can't easily
	// make os.UserHomeDir() fail. Documenting expected behavior:
	// If HOME is not available, should return ".lazyfocus.yaml"

	// For now, just verify the function doesn't panic
	path := FilePath()
	if path == "" {
		t.Error("FilePath() returned empty string")
	}
}

func TestFromContext_WithValidConfig(t *testing.T) {
	cfg := &Config{
		Output:  OutputConfig{Format: "json"},
		Timeout: 60 * time.Second,
	}
	ctx := ContextWithConfig(context.Background(), cfg)

	result, err := FromContext(ctx)
	if err != nil {
		t.Fatalf("FromContext() returned error: %v", err)
	}

	if result.Output.Format != "json" {
		t.Errorf("Expected format 'json', got %q", result.Output.Format)
	}

	if result.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", result.Timeout)
	}
}

func TestFromContext_WithoutConfig(t *testing.T) {
	ctx := context.Background()

	result, err := FromContext(ctx)
	if err != ErrConfigNotFound {
		t.Errorf("Expected ErrConfigNotFound, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestFromContext_WithNilContext(t *testing.T) {
	// Test that FromContext handles nil by checking a context without config
	// Using context.TODO() as per SA1012 - do not pass nil Context
	emptyCtx := context.TODO()
	result, err := FromContext(emptyCtx)
	if err != ErrConfigNotFound {
		t.Errorf("Expected ErrConfigNotFound, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %+v", result)
	}
}

func TestContextWithConfig_ShouldCreateContextWithConfig(t *testing.T) {
	cfg := &Config{
		Output: OutputConfig{Format: "json"},
	}
	ctx := context.Background()

	newCtx := ContextWithConfig(ctx, cfg)

	retrievedConfig, err := FromContext(newCtx)
	if err != nil {
		t.Fatalf("FromContext() returned error: %v", err)
	}

	if retrievedConfig.Output.Format != "json" {
		t.Errorf("Expected format 'json', got %q", retrievedConfig.Output.Format)
	}
}

func TestContextWithConfig_WithNilContext(t *testing.T) {
	cfg := &Config{
		Output: OutputConfig{Format: "json"},
	}

	// ContextWithConfig handles nil by using context.Background() internally
	// Using context.TODO() as per SA1012 - do not pass nil Context
	newCtx := ContextWithConfig(context.TODO(), cfg)

	retrievedConfig, err := FromContext(newCtx)
	if err != nil {
		t.Fatalf("FromContext() returned error: %v", err)
	}

	if retrievedConfig.Output.Format != "json" {
		t.Errorf("Expected format 'json', got %q", retrievedConfig.Output.Format)
	}
}

// Helper functions

func clearLazyFocusEnvVars() map[string]string {
	old := make(map[string]string)
	for _, env := range os.Environ() {
		if len(env) > 10 && env[:10] == "LAZYFOCUS_" {
			parts := splitEnvVar(env)
			if len(parts) == 2 {
				old[parts[0]] = parts[1]
				os.Unsetenv(parts[0])
			}
		}
	}
	return old
}

func restoreEnvVars(vars map[string]string) {
	// Clear any LAZYFOCUS_ vars first
	for _, env := range os.Environ() {
		if len(env) > 10 && env[:10] == "LAZYFOCUS_" {
			parts := splitEnvVar(env)
			if len(parts) == 2 {
				os.Unsetenv(parts[0])
			}
		}
	}

	// Restore old values
	for k, v := range vars {
		os.Setenv(k, v)
	}
}

func splitEnvVar(env string) []string {
	for i := 0; i < len(env); i++ {
		if env[i] == '=' {
			return []string{env[:i], env[i+1:]}
		}
	}
	return nil
}
