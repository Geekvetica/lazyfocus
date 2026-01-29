// Package overlay provides utilities for compositing overlay content on top of base views.
package overlay

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Compositor handles the positioning and compositing of overlay content.
type Compositor struct {
	width  int
	height int
}

// New creates a new Compositor with zero dimensions.
func New() *Compositor {
	return &Compositor{
		width:  0,
		height: 0,
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
	return c.simpleLayerContent(processedBase, centeredOverlay)
}

// applyDim applies a faint style to make content appear dimmed.
func (c *Compositor) applyDim(content string) string {
	dimStyle := lipgloss.NewStyle().Faint(true)
	return dimStyle.Render(content)
}

// simpleLayerContent places overlay content on top of base content.
// This preserves ANSI codes in the base by not re-processing it through Place.
func (c *Compositor) simpleLayerContent(base, overlay string) string {
	// We need to preserve ANSI codes in base content
	// The overlay is already centered and sized to viewport by Place()
	// We just need to ensure base content appears behind overlay content

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
		// Get lines (guaranteed to exist now)
		baseLine := baseLines[i]
		overlayLine := overlayLines[i]

		// If overlay line is all spaces (empty after trim), keep base line
		trimmedOverlay := strings.TrimSpace(overlayLine)
		if trimmedOverlay == "" {
			// Ensure base line is padded to viewport width
			baseLine = padToWidth(baseLine, c.width)
			result[i] = baseLine
		} else {
			// Overlay has content, use the overlay line as-is
			// (it's already sized and positioned by Place)
			result[i] = overlayLine
		}
	}

	return strings.Join(result, "\n")
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
