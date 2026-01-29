// Package overlay provides utilities for compositing overlay content on top of base views.
package overlay

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Compositor handles the positioning and compositing of overlay content.
type Compositor struct {
	width         int
	height        int
	backdropStyle lipgloss.Style
}

// New creates a new Compositor with zero dimensions.
// The backdropStyle is used to style the base content when dimming is enabled.
func New(backdropStyle lipgloss.Style) *Compositor {
	return &Compositor{
		width:         0,
		height:        0,
		backdropStyle: backdropStyle,
	}
}

// SetSize updates the compositor's viewport dimensions.
func (c *Compositor) SetSize(width, height int) {
	c.width = width
	c.height = height
}

// Place centers the given content within the compositor's viewport dimensions.
// Uses lipgloss.Place to handle the centering and padding.
func (c *Compositor) Place(content string) string {
	// Handle zero dimensions gracefully
	if c.width <= 0 || c.height <= 0 {
		// Return content as-is if dimensions are invalid
		return content
	}

	// Use lipgloss.Place to center the content
	return lipgloss.Place(
		c.width,
		c.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}

// Compose layers an overlay on top of a base view.
// If dim is true, the base content is rendered with a faint style.
func (c *Compositor) Compose(base, overlay string, dim bool) string {
	// Handle empty inputs
	if base == "" && overlay == "" {
		return ""
	}
	if base == "" {
		return c.Place(overlay)
	}
	if overlay == "" {
		if dim {
			return c.applyDim(base)
		}
		return base
	}

	// Apply dimming to base if requested
	processedBase := base
	if dim {
		processedBase = c.applyDim(base)
	}

	// Center the overlay - this creates a viewport-sized string
	// with the overlay content centered and spaces elsewhere
	centeredOverlay := c.Place(overlay)

	// Simple approach: place the base first, then overlay on top
	// Since lipgloss.Place fills the entire viewport with spaces where
	// there's no content, we just need to combine them with the overlay
	// taking precedence where it has content.
	return c.layerContentWithCharacterPrecision(processedBase, centeredOverlay)
}

// applyDim applies the backdrop style to make content appear dimmed.
func (c *Compositor) applyDim(content string) string {
	return c.backdropStyle.Render(content)
}

// layerContentWithCharacterPrecision places overlay content on top of base content.
// Uses ANSI-aware character-level compositing to preserve base content on the sides
// of the overlay, properly handling styled text via ansi.Cut, ansi.Truncate, and
// ansi.TruncateLeft for left/middle/right reconstruction.
func (c *Compositor) layerContentWithCharacterPrecision(base, overlay string) string {
	// Split into lines for compositing
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	// Ensure we have enough base lines to match viewport height
	for len(baseLines) < c.height {
		baseLines = append(baseLines, "")
	}

	// Ensure overlay has enough lines too
	for len(overlayLines) < c.height {
		overlayLines = append(overlayLines, "")
	}

	result := make([]string, c.height)

	for i := 0; i < c.height; i++ {
		baseLine := baseLines[i]
		overlayLine := overlayLines[i]

		// If overlay line is all spaces, keep base line
		trimmedOverlay := strings.TrimSpace(overlayLine)
		if trimmedOverlay == "" {
			result[i] = padToWidth(baseLine, c.width)
		} else {
			// Character-level compositing: base + overlay + base
			result[i] = c.compositeLineCharLevel(baseLine, overlayLine)
		}
	}

	return strings.Join(result, "\n")
}

// compositeLineCharLevel composites a single line with character-level precision.
// It preserves base content on the left and right of the overlay content.
func (c *Compositor) compositeLineCharLevel(baseLine, overlayLine string) string {
	// Find the bounds of actual content in the overlay line (non-space characters)
	left, right := findContentBounds(overlayLine)

	if left == -1 {
		// No content in overlay, return base
		return padToWidth(baseLine, c.width)
	}

	// Ensure base line is wide enough
	baseLine = padToWidth(baseLine, c.width)

	// Build the composite line:
	// 1. Left part from base (columns 0 to left-1)
	// 2. Middle part from overlay (columns left to right inclusive)
	// 3. Right part from base (columns right+1 to end)

	var result strings.Builder

	// Left part from base (0 to left-1)
	if left > 0 {
		leftPart := ansi.Truncate(baseLine, left, "")
		result.WriteString(leftPart)
	}

	// Middle part from overlay using ansi.Cut (left to right inclusive)
	overlayMiddle := ansi.Cut(overlayLine, left, right+1)
	result.WriteString(overlayMiddle)

	// Right part from base (right+1 to end)
	if right < c.width-1 {
		rightPart := ansi.TruncateLeft(baseLine, right+1, "")
		result.WriteString(rightPart)
	}

	return result.String()
}

// findContentBounds finds the leftmost and rightmost non-space character positions.
// Returns (-1, -1) if the line has no non-space content.
func findContentBounds(line string) (left, right int) {
	// We need to work with display positions, ignoring ANSI codes
	stripped := ansi.Strip(line)

	left = -1
	right = -1

	for i, r := range stripped {
		if r != ' ' {
			if left == -1 {
				left = i
			}
			right = i
		}
	}

	return left, right
}

// padToWidth pads a string to the specified width with spaces.
// Preserves ANSI codes by adding spaces at the end.
func padToWidth(s string, width int) string {
	// Get display width (without ANSI codes)
	displayWidth := lipgloss.Width(s)

	if displayWidth >= width {
		return s
	}

	// Add padding
	padding := strings.Repeat(" ", width-displayWidth)
	return s + padding
}
