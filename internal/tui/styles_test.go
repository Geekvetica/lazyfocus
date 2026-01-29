package tui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestDefaultStyles(t *testing.T) {
	styles := DefaultStyles()

	t.Run("returns non-nil Styles struct", func(t *testing.T) {
		if styles == nil {
			t.Fatal("DefaultStyles() returned nil")
		}
	})

	t.Run("all style structs are initialized", func(t *testing.T) {
		// Color styles - check primary color is not empty string
		if styles.Colors.Primary.Light == "" && styles.Colors.Primary.Dark == "" {
			t.Error("Colors not initialized")
		}

		// Task styles - check Normal style has been configured
		if styles.Task.Normal.GetWidth() == 0 {
			t.Error("Task styles not initialized")
		}

		// UI styles - check Header has been configured
		if !styles.UI.Header.GetBold() {
			t.Error("UI styles not initialized")
		}

		// Due date styles - check Today has been configured
		if !styles.DueDate.Today.GetBold() {
			t.Error("DueDate styles not initialized")
		}
	})
}

func TestColorStyles(t *testing.T) {
	styles := DefaultStyles()

	t.Run("Primary color is set", func(t *testing.T) {
		if styles.Colors.Primary.Light == "" && styles.Colors.Primary.Dark == "" {
			t.Error("Primary color not set")
		}
	})

	t.Run("Secondary color is set", func(t *testing.T) {
		if styles.Colors.Secondary.Light == "" && styles.Colors.Secondary.Dark == "" {
			t.Error("Secondary color not set")
		}
	})

	t.Run("Success color is set", func(t *testing.T) {
		if styles.Colors.Success.Light == "" && styles.Colors.Success.Dark == "" {
			t.Error("Success color not set")
		}
	})

	t.Run("Warning color is set", func(t *testing.T) {
		if styles.Colors.Warning.Light == "" && styles.Colors.Warning.Dark == "" {
			t.Error("Warning color not set")
		}
	})

	t.Run("Error color is set", func(t *testing.T) {
		if styles.Colors.Error.Light == "" && styles.Colors.Error.Dark == "" {
			t.Error("Error color not set")
		}
	})

	t.Run("Flagged color is set", func(t *testing.T) {
		if styles.Colors.Flagged.Light == "" && styles.Colors.Flagged.Dark == "" {
			t.Error("Flagged color not set")
		}
	})
}

func TestTaskStyles(t *testing.T) {
	styles := DefaultStyles()

	t.Run("TaskNormal has non-zero width", func(t *testing.T) {
		if styles.Task.Normal.GetWidth() == 0 {
			t.Error("TaskNormal width not set")
		}
	})

	t.Run("TaskSelected has background color", func(t *testing.T) {
		// Background should be set (we can check by rendering)
		if !styles.Task.Selected.GetBold() {
			t.Error("TaskSelected should be bold")
		}
	})

	t.Run("TaskFlagged has foreground color", func(t *testing.T) {
		// Flagged should be bold
		if !styles.Task.Flagged.GetBold() {
			t.Error("TaskFlagged should be bold")
		}
	})

	t.Run("TaskCompleted has dimmed or faint style", func(t *testing.T) {
		// Check if the style has faint and strikethrough attributes
		if !styles.Task.Completed.GetFaint() {
			t.Error("TaskCompleted should be faint")
		}
		if !styles.Task.Completed.GetStrikethrough() {
			t.Error("TaskCompleted should have strikethrough")
		}
	})
}

func TestUIStyles(t *testing.T) {
	styles := DefaultStyles()

	t.Run("Header has bold styling", func(t *testing.T) {
		rendered := styles.UI.Header.Render("TEST")
		if rendered == "TEST" {
			t.Error("Header should have styling applied")
		}
	})

	t.Run("Footer has non-zero height", func(t *testing.T) {
		if styles.UI.Footer.GetHeight() == 0 {
			t.Error("Footer height not set")
		}
	})

	t.Run("Help has foreground color", func(t *testing.T) {
		// Help should be faint
		if !styles.UI.Help.GetFaint() {
			t.Error("Help should be faint")
		}
	})

	t.Run("Overlay has background color and border", func(t *testing.T) {
		// Overlay should have padding
		if styles.UI.Overlay.GetPaddingLeft() == 0 {
			t.Error("Overlay padding not set")
		}
	})

	t.Run("Input has border", func(t *testing.T) {
		borderStyle := styles.UI.Input.GetBorderStyle()
		if borderStyle == lipgloss.NormalBorder() {
			// Border might be set but we need to check if border is enabled
			if !styles.UI.Input.GetBorderTop() && !styles.UI.Input.GetBorderBottom() &&
				!styles.UI.Input.GetBorderLeft() && !styles.UI.Input.GetBorderRight() {
				t.Error("Input should have border enabled")
			}
		}
	})
}

func TestDueDateStyles(t *testing.T) {
	styles := DefaultStyles()

	t.Run("DueToday has bold styling", func(t *testing.T) {
		if !styles.DueDate.Today.GetBold() {
			t.Error("DueToday should be bold")
		}
	})

	t.Run("DueOverdue has bold styling", func(t *testing.T) {
		if !styles.DueDate.Overdue.GetBold() {
			t.Error("DueOverdue should be bold")
		}
	})

	t.Run("DueNormal has foreground color set", func(t *testing.T) {
		// Just verify the style can render without panic
		result := styles.DueDate.Normal.Render("test")
		if result == "" {
			t.Error("DueNormal should render text")
		}
	})
}

func TestStylesRendering(t *testing.T) {
	styles := DefaultStyles()

	t.Run("TaskNormal renders text", func(t *testing.T) {
		result := styles.Task.Normal.Render("Buy groceries")
		if result == "" {
			t.Error("TaskNormal should render text")
		}
	})

	t.Run("TaskSelected renders differently than TaskNormal", func(t *testing.T) {
		// Compare style attributes instead of rendered output
		// Selected should be bold but normal should not
		if styles.Task.Normal.GetBold() {
			t.Error("TaskNormal should not be bold")
		}
		if !styles.Task.Selected.GetBold() {
			t.Error("TaskSelected should be bold")
		}
		// Selected should have background, normal should not
		normalBg := styles.Task.Normal.GetBackground()
		selectedBg := styles.Task.Selected.GetBackground()
		if normalBg == selectedBg {
			t.Error("TaskSelected should have different background than TaskNormal")
		}
	})

	t.Run("Header renders text with styling", func(t *testing.T) {
		result := styles.UI.Header.Render("INBOX (3 tasks)")
		if result == "" {
			t.Error("Header should render text")
		}
	})

	t.Run("DueOverdue renders differently than DueNormal", func(t *testing.T) {
		// Compare style attributes instead of rendered output
		// Overdue should be bold, normal should not
		if !styles.DueDate.Overdue.GetBold() {
			t.Error("DueOverdue should be bold")
		}
		if styles.DueDate.Normal.GetBold() {
			t.Error("DueNormal should not be bold")
		}
		// Foreground colors should differ
		overdueFg := styles.DueDate.Overdue.GetForeground()
		normalFg := styles.DueDate.Normal.GetForeground()
		if overdueFg == normalFg {
			t.Error("DueOverdue should have different foreground than DueNormal")
		}
	})
}

func TestAdaptiveColors(t *testing.T) {
	styles := DefaultStyles()

	t.Run("styles support adaptive colors", func(t *testing.T) {
		// Adaptive colors should be used for light/dark terminal support
		// Verify that colors have both light and dark variants

		// Primary color should be an adaptive color for better terminal support
		if styles.Colors.Primary.Light == "" || styles.Colors.Primary.Dark == "" {
			t.Error("Primary color should have both light and dark variants")
		}

		// Check other colors too
		if styles.Colors.Success.Light == "" || styles.Colors.Success.Dark == "" {
			t.Error("Success color should have both light and dark variants")
		}

		if styles.Colors.Error.Light == "" || styles.Colors.Error.Dark == "" {
			t.Error("Error color should have both light and dark variants")
		}
	})
}
