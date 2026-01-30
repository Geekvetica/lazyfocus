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

// ProjectStyles defines styles for project display
type ProjectStyles struct {
	Active    lipgloss.Style
	OnHold    lipgloss.Style
	Completed lipgloss.Style
	Dropped   lipgloss.Style
}

// ForecastStyles defines styles for forecast view groups
type ForecastStyles struct {
	Overdue     lipgloss.Style
	Today       lipgloss.Style
	Tomorrow    lipgloss.Style
	Later       lipgloss.Style
	GroupHeader lipgloss.Style
}

// SearchStyles defines styles for search highlighting
type SearchStyles struct {
	Highlight lipgloss.Style
	Input     lipgloss.Style
}

// TagStyles defines styles for tag display
type TagStyles struct {
	Badge    lipgloss.Style
	Selected lipgloss.Style
}

// Styles contains all organized style groups
type Styles struct {
	Colors   ColorStyles
	Task     TaskStyles
	Project  ProjectStyles
	Forecast ForecastStyles
	Search   SearchStyles
	Tag      TagStyles
	UI       UIStyles
	DueDate  DueDateStyles
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

	// Project styles
	projectStyles := ProjectStyles{
		Active: lipgloss.NewStyle().
			Foreground(colors.Primary),
		OnHold: lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Faint(true),
		Completed: lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Faint(true).
			Strikethrough(true),
		Dropped: lipgloss.NewStyle().
			Foreground(colors.Error).
			Faint(true).
			Strikethrough(true),
	}

	// Forecast styles
	forecastStyles := ForecastStyles{
		Overdue: lipgloss.NewStyle().
			Foreground(colors.Error).
			Bold(true),
		Today: lipgloss.NewStyle().
			Foreground(colors.Warning).
			Bold(true),
		Tomorrow: lipgloss.NewStyle().
			Foreground(colors.Primary),
		Later: lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Faint(true),
		GroupHeader: lipgloss.NewStyle().
			Bold(true).
			Underline(true).
			Foreground(colors.Primary),
	}

	// Search styles
	searchStyles := SearchStyles{
		Highlight: lipgloss.NewStyle().
			Background(colors.Warning).
			Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#000000"}),
		Input: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Primary).
			Padding(0, 1),
	}

	// Tag styles
	tagStyles := TagStyles{
		Badge: lipgloss.NewStyle().
			Foreground(colors.Secondary).
			Background(lipgloss.AdaptiveColor{Light: "#E8E8E8", Dark: "#3A3A3A"}).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Secondary),
		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#FFFFFF", Dark: "#000000"}).
			Background(colors.Primary).
			Padding(0, 1).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Primary).
			Bold(true),
	}

	return &Styles{
		Colors:   colors,
		Task:     taskStyles,
		Project:  projectStyles,
		Forecast: forecastStyles,
		Search:   searchStyles,
		Tag:      tagStyles,
		UI:       uiStyles,
		DueDate:  dueDateStyles,
	}
}
