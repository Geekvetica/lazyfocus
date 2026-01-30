// Package config provides configuration management for LazyFocus.
package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// contextKey is used to store config in context
type contextKey string

const configKey contextKey = "lazyfocus-config"

// ErrConfigNotFound is returned when config is not found in context
var ErrConfigNotFound = errors.New("config not found in context")

// Config holds the application configuration
type Config struct {
	Output   OutputConfig   `mapstructure:"output"`
	Timeout  time.Duration  `mapstructure:"timeout"`
	Defaults DefaultsConfig `mapstructure:"defaults"`
	TUI      TUIConfig      `mapstructure:"tui"`
}

// OutputConfig holds output-related configuration
type OutputConfig struct {
	Format string `mapstructure:"format"` // "human" or "json"
}

// DefaultsConfig holds default values for commands
type DefaultsConfig struct {
	Project string `mapstructure:"project"` // Default project name
}

// TUIConfig holds TUI-related configuration
type TUIConfig struct {
	Theme  string      `mapstructure:"theme"` // "default" or custom
	Colors ColorConfig `mapstructure:"colors"`
}

// ColorConfig holds color configuration for TUI
type ColorConfig struct {
	Primary string `mapstructure:"primary"` // Primary accent color
	Flagged string `mapstructure:"flagged"` // Color for flagged items
	Due     string `mapstructure:"due"`     // Color for due items
	Overdue string `mapstructure:"overdue"` // Color for overdue items
}

// Load loads configuration from file and environment
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Config file settings
	v.SetConfigName(".lazyfocus")
	v.SetConfigType("yaml")

	// Look in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		v.AddConfigPath(home)
	}

	// Also look in current directory
	v.AddConfigPath(".")

	// Environment variables
	v.SetEnvPrefix("LAZYFOCUS")
	v.AutomaticEnv()

	// Bind environment variables to config keys explicitly
	// This is needed for nested keys to work properly
	_ = v.BindEnv("output.format", "LAZYFOCUS_OUTPUT_FORMAT")
	_ = v.BindEnv("timeout", "LAZYFOCUS_TIMEOUT")
	_ = v.BindEnv("defaults.project", "LAZYFOCUS_DEFAULTS_PROJECT")
	_ = v.BindEnv("tui.theme", "LAZYFOCUS_TUI_THEME")
	_ = v.BindEnv("tui.colors.primary", "LAZYFOCUS_TUI_COLORS_PRIMARY")
	_ = v.BindEnv("tui.colors.flagged", "LAZYFOCUS_TUI_COLORS_FLAGGED")
	_ = v.BindEnv("tui.colors.due", "LAZYFOCUS_TUI_COLORS_DUE")
	_ = v.BindEnv("tui.colors.overdue", "LAZYFOCUS_TUI_COLORS_OVERDUE")

	// Read config file (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// FilePath returns the path to the config file
func FilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".lazyfocus.yaml"
	}
	return filepath.Join(home, ".lazyfocus.yaml")
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("output.format", "human")
	v.SetDefault("timeout", "30s")
	v.SetDefault("defaults.project", "")
	v.SetDefault("tui.theme", "default")
	v.SetDefault("tui.colors.primary", "#5B9BD5")
	v.SetDefault("tui.colors.flagged", "#ED7D31")
	v.SetDefault("tui.colors.due", "#70AD47")
	v.SetDefault("tui.colors.overdue", "#FF6B6B")
}

// FromContext extracts the Config from the context.
// Returns ErrConfigNotFound if the context is nil or config is not present.
func FromContext(ctx context.Context) (*Config, error) {
	if ctx == nil {
		return nil, ErrConfigNotFound
	}
	cfg, ok := ctx.Value(configKey).(*Config)
	if !ok || cfg == nil {
		return nil, ErrConfigNotFound
	}
	return cfg, nil
}

// ContextWithConfig returns a new context with the config attached.
// If ctx is nil, context.Background() is used as the parent context.
func ContextWithConfig(ctx context.Context, cfg *Config) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, configKey, cfg)
}
