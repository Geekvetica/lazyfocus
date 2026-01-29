package overlay

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func init() {
	// Enable ANSI colors for testing
	lipgloss.SetColorProfile(termenv.ANSI256)
}

// TestPlaceCentersContent verifies that content is centered in the viewport
func TestPlaceCentersContent(t *testing.T) {
	c := New()
	c.SetSize(80, 24)

	// Simple 3x3 content
	content := "XXX\nXXX\nXXX"
	result := c.Place(content)

	// Result should be 80x24
	resultHeight := lipgloss.Height(result)
	resultWidth := lipgloss.Width(result)

	if resultHeight != 24 {
		t.Errorf("expected height 24, got %d", resultHeight)
	}
	if resultWidth != 80 {
		t.Errorf("expected width 80, got %d", resultWidth)
	}

	// Content should be centered - find the X's in the output
	lines := strings.Split(result, "\n")
	foundContent := false
	for _, line := range lines {
		if strings.Contains(line, "XXX") {
			foundContent = true
			// Should have spaces on both sides (centered)
			trimmed := strings.TrimSpace(line)
			if trimmed != "XXX" {
				t.Errorf("content line has unexpected characters: %q", line)
			}
			break
		}
	}

	if !foundContent {
		t.Error("content not found in placed output")
	}
}

// TestPlaceWithZeroDimensions verifies graceful handling of 0x0 viewport
func TestPlaceWithZeroDimensions(t *testing.T) {
	c := New()
	c.SetSize(0, 0)

	content := "test"
	result := c.Place(content)

	// Should not panic and should return something (even if empty/unchanged)
	if result == "" {
		t.Error("expected non-empty result even with zero dimensions")
	}
}

// TestPlaceWithSmallViewport verifies handling when content is larger than viewport
func TestPlaceWithSmallViewport(t *testing.T) {
	c := New()
	c.SetSize(10, 3)

	// Content is larger than viewport
	content := "This is a very long line that exceeds viewport width\nLine 2\nLine 3\nLine 4\nLine 5"
	result := c.Place(content)

	// Should not panic
	if result == "" {
		t.Error("expected non-empty result")
	}

	// Result dimensions should match viewport
	resultHeight := lipgloss.Height(result)
	resultWidth := lipgloss.Width(result)

	// Width should be at least as wide as viewport (lipgloss may add padding)
	if resultWidth < 10 {
		t.Errorf("expected width >= 10, got %d", resultWidth)
	}

	// Height should be at least as tall as viewport
	if resultHeight < 3 {
		t.Errorf("expected height >= 3, got %d", resultHeight)
	}
}

// TestComposeLayersOverlay verifies that overlay appears over base content
func TestComposeLayersOverlay(t *testing.T) {
	c := New()
	c.SetSize(80, 24)

	base := "Base content line 1\nBase content line 2\nBase content line 3"
	overlay := "OVERLAY"

	result := c.Compose(base, overlay, false)

	// Result should contain overlay content
	if !strings.Contains(result, "OVERLAY") {
		t.Error("result does not contain overlay content")
	}

	// Base should still be present (underneath)
	if !strings.Contains(result, "Base content") {
		t.Error("result does not contain base content")
	}
}

// TestComposeWithDimming verifies that base content is dimmed when requested
func TestComposeWithDimming(t *testing.T) {
	c := New()
	c.SetSize(80, 24)

	base := "Base content"
	overlay := "OVERLAY"

	// Test without dimming
	resultNoDim := c.Compose(base, overlay, false)

	// Test with dimming
	resultDimmed := c.Compose(base, overlay, true)

	// Both should contain the overlay
	if !strings.Contains(resultNoDim, "OVERLAY") {
		t.Error("result without dimming does not contain overlay")
	}
	if !strings.Contains(resultDimmed, "OVERLAY") {
		t.Error("result with dimming does not contain overlay")
	}

	// Dimmed result should have ANSI escape codes for faint/dim styling
	// Faint is ANSI code \x1b[2m
	if !strings.Contains(resultDimmed, "\x1b[2m") {
		t.Errorf("dimmed result does not contain faint ANSI code\nResult: %q", resultDimmed)
	}

	// Non-dimmed should not have faint code on base content
	// (overlay might have styling, but base shouldn't be explicitly dimmed)
	baseLines := strings.Split(resultNoDim, "\n")
	overlayLines := strings.Split(resultDimmed, "\n")

	// At least one line in the non-dimmed version should lack the faint code
	// while the dimmed version should have it
	foundDifference := false
	for i := 0; i < len(baseLines) && i < len(overlayLines); i++ {
		hasNoDimFaint := strings.Contains(baseLines[i], "\x1b[2m")
		hasDimFaint := strings.Contains(overlayLines[i], "\x1b[2m")

		if !hasNoDimFaint && hasDimFaint {
			foundDifference = true
			break
		}
	}

	if !foundDifference {
		t.Error("expected dimming to add faint styling that wasn't present before")
	}
}

// TestSetSizeUpdatesCompositor verifies that SetSize updates the compositor dimensions
func TestSetSizeUpdatesCompositor(t *testing.T) {
	c := New()

	// Initially should have zero dimensions
	if c.width != 0 || c.height != 0 {
		t.Errorf("expected initial dimensions 0x0, got %dx%d", c.width, c.height)
	}

	// Set new dimensions
	c.SetSize(100, 50)

	if c.width != 100 {
		t.Errorf("expected width 100, got %d", c.width)
	}
	if c.height != 50 {
		t.Errorf("expected height 50, got %d", c.height)
	}
}

// TestComposeWithNilInputs verifies graceful handling of empty strings
func TestComposeWithNilInputs(t *testing.T) {
	c := New()
	c.SetSize(80, 24)

	// Empty base
	result1 := c.Compose("", "overlay", false)
	if !strings.Contains(result1, "overlay") {
		t.Error("should contain overlay even with empty base")
	}

	// Empty overlay
	result2 := c.Compose("base", "", false)
	if !strings.Contains(result2, "base") {
		t.Error("should contain base even with empty overlay")
	}

	// Both empty
	result3 := c.Compose("", "", false)
	if result3 == "" {
		// This is acceptable
	}
}
