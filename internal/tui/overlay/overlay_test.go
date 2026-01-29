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

	// Both empty - should not panic
	_ = c.Compose("", "", false)
}

// TestComposeEmptyOverlayWithDimming verifies that empty overlay with dimming applies dim to base
func TestComposeEmptyOverlayWithDimming(t *testing.T) {
	c := New()
	c.SetSize(80, 24)

	base := "Base content"

	// Empty overlay with dimming should apply dim to base
	result := c.Compose(base, "", true)

	// Should contain base content
	if !strings.Contains(result, "Base") {
		t.Error("should contain base content")
	}

	// Should have faint ANSI code applied (line 61 coverage)
	if !strings.Contains(result, "\x1b[2m") {
		t.Error("should apply faint styling when dim=true and overlay is empty")
	}
}

// TestPlaceWithNegativeDimensions verifies graceful handling of negative viewport dimensions
func TestPlaceWithNegativeDimensions(t *testing.T) {
	c := New()
	c.SetSize(-10, -5)

	content := "test content"
	result := c.Place(content)

	// Should return content unchanged (not panic)
	if result != content {
		t.Errorf("expected content unchanged for negative dimensions, got %q", result)
	}
}

// TestPadToWidth verifies the padding function handles various width scenarios
func TestPadToWidth(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		width         int
		expectedWidth int
	}{
		{
			name:          "shorter than width gets padded",
			input:         "Hi",
			width:         10,
			expectedWidth: 10,
		},
		{
			name:          "equal to width unchanged",
			input:         "0123456789",
			width:         10,
			expectedWidth: 10,
		},
		{
			name:          "longer than width unchanged",
			input:         "This is longer",
			width:         5,
			expectedWidth: 14, // unchanged - content is 14 chars wide
		},
		{
			name:          "empty string gets padded",
			input:         "",
			width:         5,
			expectedWidth: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padToWidth(tt.input, tt.width)
			gotWidth := lipgloss.Width(result)
			if gotWidth != tt.expectedWidth {
				t.Errorf("expected width %d, got %d", tt.expectedWidth, gotWidth)
			}
		})
	}
}

// TestFindContentBounds verifies finding content boundaries in lines
func TestFindContentBounds(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLeft  int
		wantRight int
	}{
		{
			name:      "content at start",
			input:     "Hello     ",
			wantLeft:  0,
			wantRight: 4,
		},
		{
			name:      "content at end",
			input:     "     World",
			wantLeft:  5,
			wantRight: 9,
		},
		{
			name:      "content in middle",
			input:     "   XYZ   ",
			wantLeft:  3,
			wantRight: 5,
		},
		{
			name:      "all spaces",
			input:     "          ",
			wantLeft:  -1,
			wantRight: -1,
		},
		{
			name:      "empty string",
			input:     "",
			wantLeft:  -1,
			wantRight: -1,
		},
		{
			name:      "single char",
			input:     "  X  ",
			wantLeft:  2,
			wantRight: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, right := findContentBounds(tt.input)
			if left != tt.wantLeft {
				t.Errorf("left: expected %d, got %d", tt.wantLeft, left)
			}
			if right != tt.wantRight {
				t.Errorf("right: expected %d, got %d", tt.wantRight, right)
			}
		})
	}
}

// TestCompositeLineCharLevel verifies character-level compositing
func TestCompositeLineCharLevel(t *testing.T) {
	c := New()
	c.SetSize(20, 1)

	// Base: "│ABCDEFGHIJKLMNOP│  " (border chars at positions 0 and 17)
	// Overlay: "     [OVERLAY]     " (content from 5 to 13)
	// Expected: "│ABCD[OVERLAY]OP│  " (base on sides, overlay in middle)

	baseLine := "|ABCDEFGHIJKLMNOP|  "
	overlayLine := "     [OVERLAY]      "

	result := c.compositeLineCharLevel(baseLine, overlayLine)

	// Should have base content on left
	if !strings.HasPrefix(result, "|ABCD") {
		t.Errorf("expected base content on left, got: %q", result)
	}

	// Should have overlay content in middle
	if !strings.Contains(result, "[OVERLAY]") {
		t.Errorf("expected overlay content in middle, got: %q", result)
	}

	// Should have base content on right
	if !strings.Contains(result, "OP|") {
		t.Errorf("expected base content on right, got: %q", result)
	}
}

// TestComposePreservesBaseBorders verifies that borders are preserved during compositing
func TestComposePreservesBaseBorders(t *testing.T) {
	c := New()
	c.SetSize(40, 10)

	// Create base with borders
	base := strings.Repeat("│content here content│\n", 10)

	// Create small centered overlay
	overlay := "MODAL"

	result := c.Compose(base, overlay, true)

	// The lines without overlay should have borders preserved
	lines := strings.Split(result, "\n")
	foundBorderLine := false
	for _, line := range lines {
		// Skip overlay lines
		if strings.Contains(line, "MODAL") {
			continue
		}
		// Non-overlay lines should have border character
		if strings.Contains(line, "│") {
			foundBorderLine = true
			break
		}
	}

	if !foundBorderLine {
		t.Error("expected border characters to be preserved in non-overlay lines")
	}
}

// TestCompositeLineCharLevelWithEmptyOverlay verifies handling of empty overlay
func TestCompositeLineCharLevelWithEmptyOverlay(t *testing.T) {
	c := New()
	c.SetSize(20, 1)

	baseLine := "base content here"
	overlayLine := "                    " // all spaces

	result := c.compositeLineCharLevel(baseLine, overlayLine)

	// Should return padded base line since overlay is empty
	if !strings.Contains(result, "base content here") {
		t.Errorf("expected base content for empty overlay, got: %q", result)
	}
}

// TestCompositeLineCharLevelOverlayAtEdges verifies overlay at viewport edges
func TestCompositeLineCharLevelOverlayAtEdges(t *testing.T) {
	c := New()
	c.SetSize(20, 1)

	// Overlay content at the start
	baseLine := "BBBBBBBBBBBBBBBBBBBB"
	overlayLine := "OOOO                "

	result := c.compositeLineCharLevel(baseLine, overlayLine)

	// Should have overlay at start and base at end
	if !strings.HasPrefix(result, "OOOO") {
		t.Errorf("expected overlay at start, got: %q", result)
	}
	if !strings.HasSuffix(result, "BBBBBBBBBBBBBBBB") {
		t.Errorf("expected base at end, got: %q", result)
	}
}
