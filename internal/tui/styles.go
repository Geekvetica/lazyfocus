package tui

import (
	"github.com/charmbracelet/lipgloss"
)

// ColorStyles defines the color palette for the TUI
type ColorStyles struct {
	Primary   lipgloss.AdaptiveColor
	Secondary lipgloss.AdaptiveColor
	Success   lipgloss.AdaptiveColor
	Warning   lipgloss.AdaptiveColor
	Error     lipgloss.AdaptiveColor
	Flagged   lipgloss.AdaptiveColor
}

// TaskStyles defines styles for task display
type TaskStyles struct {
	Normal    lipgloss.Style
	Selected  lipgloss.Style
	Flagged   lipgloss.Style
	Completed lipgloss.Style
}

// UIStyles defines styles for UI elements
type UIStyles struct {
	Header          lipgloss.Style
	Footer          lipgloss.Style
	Help            lipgloss.Style
	Overlay         lipgloss.Style
	OverlayBackdrop lipgloss.Style
	Input           lipgloss.Style
}

// DueDateStyles defines styles for due date display
type DueDateStyles struct {
	Today   lipgloss.Style
	Overdue lipgloss.Style
	Normal  lipgloss.Style
}

// Styles contains all organized style groups
type Styles struct {
	Colors  ColorStyles
	Task    TaskStyles
	UI      UIStyles
	DueDate DueDateStyles
}

// DefaultStyles returns the default style configuration
func DefaultStyles() *Styles {
	// Define color palette
	colors := ColorStyles{
		Primary: lipgloss.AdaptiveColor{
			Light: "#5B9BD5", // Blue
			Dark:  "#7FB3D5",
		},
		Secondary: lipgloss.AdaptiveColor{
			Light: "#808080", // Gray
			Dark:  "#A0A0A0",
		},
		Success: lipgloss.AdaptiveColor{
			Light: "#70AD47", // Green
			Dark:  "#90CD67",
		},
		Warning: lipgloss.AdaptiveColor{
			Light: "#FFC000", // Orange/Yellow
			Dark:  "#FFD666",
		},
		Error: lipgloss.AdaptiveColor{
			Light: "#C00000", // Red
			Dark:  "#FF6B6B",
		},
		Flagged: lipgloss.AdaptiveColor{
			Light: "#ED7D31", // Orange/Red
			Dark:  "#FF9F66",
		},
	}

	// Task styles
	taskStyles := TaskStyles{
		Normal: lipgloss.NewStyle().
			Width(80).
			PaddingLeft(1),
		Selected: lipgloss.NewStyle().
			Width(80).
			PaddingLeft(1).
			Background(colors.Primary).
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}).
			Bold(true),
		Flagged: lipgloss.NewStyle().
			Foreground(colors.Flagged).
			Bold(true),
		Completed: lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Faint(true).
			Strikethrough(true),
	}

	// UI styles
	uiStyles := UIStyles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Primary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(colors.Secondary).
			PaddingLeft(1).
			PaddingRight(1),
		Footer: lipgloss.NewStyle().
			Height(1).
			Foreground(colors.Secondary).
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(colors.Secondary).
			PaddingLeft(1).
			PaddingRight(1),
		Help: lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Faint(true),
		Overlay: lipgloss.NewStyle().
			Background(lipgloss.AdaptiveColor{Light: "#F0F0F0", Dark: "#2A2A2A"}).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Primary).
			Padding(1, 2),
		OverlayBackdrop: lipgloss.NewStyle().
			Faint(true),
		Input: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Primary).
			Padding(0, 1),
	}

	// Due date styles
	dueDateStyles := DueDateStyles{
		Today: lipgloss.NewStyle().
			Foreground(colors.Warning).
			Bold(true),
		Overdue: lipgloss.NewStyle().
			Foreground(colors.Error).
			Bold(true),
		Normal: lipgloss.NewStyle().
			Foreground(colors.Secondary),
	}

	return &Styles{
		Colors:  colors,
		Task:    taskStyles,
		UI:      uiStyles,
		DueDate: dueDateStyles,
	}
}
