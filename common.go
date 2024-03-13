package foam

// Package foam provides utilities for styling and managing UI components.
// It uses the Lipgloss library for styling and offers a structured way to
// handle styles for focused, blurred, and no-border states of UI elements.
import (
	"github.com/charmbracelet/lipgloss"
)

// Styles defines the styles for focused, blurred, and no-border
// states of UI elements.
type Styles struct {
	Focused  lipgloss.Style
	Blurred  lipgloss.Style
	NoBorder lipgloss.Style
}

// Common encapsulates common properties for UI components, including
// width, height, and styles.
type Common struct {
	width, height int
	styles        *Styles
}

// DefaultStyles returns a new Styles struct with default styles for
// focused, blurred, and no-border states, leveraging Lipgloss for
// styling.
func DefaultStyles() *Styles {
	return &Styles{
		Focused:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")),
		Blurred:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
		NoBorder: lipgloss.NewStyle(),
	}
}

// GetWidth returns the current width of the UI component.
func (c *Common) GetWidth() int { return c.width }

// GetHeight returns the current height of the UI component.
func (c *Common) GetHeight() int { return c.height }

// SetWidth sets the width of the UI component.
func (c *Common) SetWidth(width int) {
	c.width = width - c.styles.Focused.GetHorizontalFrameSize()
	c.styles.Focused = c.styles.Focused.Width(c.width)
	c.styles.Blurred = c.styles.Blurred.Width(c.width)
}

// SetHeight sets the height of the UI component.
func (c *Common) SetHeight(height int) {
	// c.height = height
	c.height = height - c.styles.Focused.GetHorizontalFrameSize()
	c.styles.Focused = c.styles.Focused.Height(c.height)
	c.styles.Blurred = c.styles.Blurred.Height(c.height)

}

// SetSize sets the width and height of the UI component, adjusting
// the focused style width to account for the border size.
func (c *Common) SetSize(width int, height int) {
	c.width = width - c.styles.Focused.GetHorizontalFrameSize()
	c.height = height - c.styles.Focused.GetVerticalFrameSize()

	c.styles.Focused = c.styles.Focused.Width(c.width)
	c.styles.Focused = c.styles.Focused.Height(c.height)

	c.styles.Blurred = c.styles.Blurred.Width(c.width)
	c.styles.Blurred = c.styles.Blurred.Height(c.height)
}

// GetStyles returns the current styles of the UI component.
func (c *Common) GetStyles() *Styles { return c.styles }

// SetStyles sets the styles of the UI component.
func (c *Common) SetStyles(styles *Styles) {
	c.styles = styles
}
