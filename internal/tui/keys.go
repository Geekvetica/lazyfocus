// Package tui provides shared types for the TUI layer.
package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings for the TUI
type KeyMap struct {
	// Navigation
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding

	// View Switching (1-5)
	View1 key.Binding
	View2 key.Binding
	View3 key.Binding
	View4 key.Binding
	View5 key.Binding

	// Actions
	QuickAdd key.Binding
	Complete key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Flag     key.Binding

	// Global
	Quit key.Binding
	Help key.Binding
}

// DefaultKeyMap returns the default key bindings for the TUI
func DefaultKeyMap() KeyMap {
	return KeyMap{
		// Navigation
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("k/↑", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("j/↓", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("h/←", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("l/→", "move right"),
		),

		// View Switching
		View1: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "inbox view"),
		),
		View2: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "projects view"),
		),
		View3: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "tags view"),
		),
		View4: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "forecast view"),
		),
		View5: key.NewBinding(
			key.WithKeys("5"),
			key.WithHelp("5", "review view"),
		),

		// Actions
		QuickAdd: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "quick add task"),
		),
		Complete: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "complete task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete task"),
		),
		Flag: key.NewBinding(
			key.WithKeys("f"),
			key.WithHelp("f", "toggle flag"),
		),

		// Global
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}
}
