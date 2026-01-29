package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		name        string
		binding     key.Binding
		wantKeys    []string
		wantHelp    string
		wantEnabled bool
	}{
		// Navigation
		{
			name:        "Up binding",
			binding:     km.Up,
			wantKeys:    []string{"k", "up"},
			wantHelp:    "k/↑",
			wantEnabled: true,
		},
		{
			name:        "Down binding",
			binding:     km.Down,
			wantKeys:    []string{"j", "down"},
			wantHelp:    "j/↓",
			wantEnabled: true,
		},
		{
			name:        "Left binding",
			binding:     km.Left,
			wantKeys:    []string{"h", "left"},
			wantHelp:    "h/←",
			wantEnabled: true,
		},
		{
			name:        "Right binding",
			binding:     km.Right,
			wantKeys:    []string{"l", "right"},
			wantHelp:    "l/→",
			wantEnabled: true,
		},
		// View Switching
		{
			name:        "View1 binding",
			binding:     km.View1,
			wantKeys:    []string{"1"},
			wantHelp:    "1",
			wantEnabled: true,
		},
		{
			name:        "View2 binding",
			binding:     km.View2,
			wantKeys:    []string{"2"},
			wantHelp:    "2",
			wantEnabled: true,
		},
		{
			name:        "View3 binding",
			binding:     km.View3,
			wantKeys:    []string{"3"},
			wantHelp:    "3",
			wantEnabled: true,
		},
		{
			name:        "View4 binding",
			binding:     km.View4,
			wantKeys:    []string{"4"},
			wantHelp:    "4",
			wantEnabled: true,
		},
		{
			name:        "View5 binding",
			binding:     km.View5,
			wantKeys:    []string{"5"},
			wantHelp:    "5",
			wantEnabled: true,
		},
		// Actions
		{
			name:        "QuickAdd binding",
			binding:     km.QuickAdd,
			wantKeys:    []string{"a"},
			wantHelp:    "a",
			wantEnabled: true,
		},
		{
			name:        "Complete binding",
			binding:     km.Complete,
			wantKeys:    []string{"c"},
			wantHelp:    "c",
			wantEnabled: true,
		},
		{
			name:        "Edit binding",
			binding:     km.Edit,
			wantKeys:    []string{"e"},
			wantHelp:    "e",
			wantEnabled: true,
		},
		{
			name:        "Delete binding",
			binding:     km.Delete,
			wantKeys:    []string{"d"},
			wantHelp:    "d",
			wantEnabled: true,
		},
		{
			name:        "Flag binding",
			binding:     km.Flag,
			wantKeys:    []string{"f"},
			wantHelp:    "f",
			wantEnabled: true,
		},
		// Global
		{
			name:        "Quit binding",
			binding:     km.Quit,
			wantKeys:    []string{"q"},
			wantHelp:    "q",
			wantEnabled: true,
		},
		{
			name:        "Help binding",
			binding:     km.Help,
			wantKeys:    []string{"?"},
			wantHelp:    "?",
			wantEnabled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if keys are set
			if len(tt.binding.Keys()) != len(tt.wantKeys) {
				t.Errorf("got %d keys, want %d keys", len(tt.binding.Keys()), len(tt.wantKeys))
			}

			// Verify each key
			gotKeys := tt.binding.Keys()
			for i, wantKey := range tt.wantKeys {
				if i >= len(gotKeys) {
					t.Errorf("missing key at index %d: want %q", i, wantKey)
					continue
				}
				if gotKeys[i] != wantKey {
					t.Errorf("key[%d] = %q, want %q", i, gotKeys[i], wantKey)
				}
			}

			// Check help key
			helpKey := tt.binding.Help().Key
			if helpKey != tt.wantHelp {
				t.Errorf("help key = %q, want %q", helpKey, tt.wantHelp)
			}

			// Check if enabled
			if tt.binding.Enabled() != tt.wantEnabled {
				t.Errorf("enabled = %v, want %v", tt.binding.Enabled(), tt.wantEnabled)
			}
		})
	}
}

func TestKeyBindingMatches(t *testing.T) {
	km := DefaultKeyMap()

	tests := []struct {
		name    string
		binding key.Binding
		testKey string
		want    bool
	}{
		// Navigation
		{"Up with k", km.Up, "k", true},
		{"Up with up arrow", km.Up, "up", true},
		{"Up with wrong key", km.Up, "j", false},
		{"Down with j", km.Down, "j", true},
		{"Down with down arrow", km.Down, "down", true},
		{"Left with h", km.Left, "h", true},
		{"Left with left arrow", km.Left, "left", true},
		{"Right with l", km.Right, "l", true},
		{"Right with right arrow", km.Right, "right", true},
		// View Switching
		{"View1 with 1", km.View1, "1", true},
		{"View2 with 2", km.View2, "2", true},
		{"View3 with 3", km.View3, "3", true},
		{"View4 with 4", km.View4, "4", true},
		{"View5 with 5", km.View5, "5", true},
		{"View1 with wrong key", km.View1, "6", false},
		// Actions
		{"QuickAdd with a", km.QuickAdd, "a", true},
		{"Complete with c", km.Complete, "c", true},
		{"Edit with e", km.Edit, "e", true},
		{"Delete with d", km.Delete, "d", true},
		{"Flag with f", km.Flag, "f", true},
		{"QuickAdd with wrong key", km.QuickAdd, "b", false},
		// Global
		{"Quit with q", km.Quit, "q", true},
		{"Help with ?", km.Help, "?", true},
		{"Quit with wrong key", km.Quit, "x", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check if the test key is in the binding's keys list
			keys := tt.binding.Keys()
			found := false
			for _, k := range keys {
				if k == tt.testKey {
					found = true
					break
				}
			}
			if found != tt.want {
				t.Errorf("key %q in binding = %v, want %v", tt.testKey, found, tt.want)
			}
		})
	}
}

func TestKeyMapStructHasAllFields(t *testing.T) {
	km := DefaultKeyMap()

	// Use reflection-like approach via field access
	// This ensures all fields are present and of type key.Binding
	_ = km.Up
	_ = km.Down
	_ = km.Left
	_ = km.Right
	_ = km.View1
	_ = km.View2
	_ = km.View3
	_ = km.View4
	_ = km.View5
	_ = km.QuickAdd
	_ = km.Complete
	_ = km.Edit
	_ = km.Delete
	_ = km.Flag
	_ = km.Quit
	_ = km.Help

	// If we get here without compilation errors, all fields exist
	t.Log("All required key binding fields are present in KeyMap")
}
