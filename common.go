package foam

import (
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Focused  lipgloss.Style
	Blurred  lipgloss.Style
	NoBorder lipgloss.Style
}

type Common struct {
	width, height int
	styles        *Styles
}

func DefaultStyles() *Styles {
	return &Styles{
		Focused:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")),
		Blurred:  lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
		NoBorder: lipgloss.NewStyle(),
	}
}

func (c *Common) GetWidth() int  { return c.width }
func (c *Common) GetHeight() int { return c.height }

func (c *Common) SetWidth(width int)   { c.width = width }
func (c *Common) SetHeight(height int) { c.height = height }

func (c *Common) SetSize(width int, height int) {
	c.width = width - c.styles.Focused.GetHorizontalFrameSize()
	c.height = height

	c.styles.Focused = c.styles.Focused.Width(c.width)
	c.styles.Blurred = c.styles.Blurred.Width(c.width)
}

func (c *Common) GetStyles() *Styles { return c.styles }
func (c *Common) SetStyles(styles *Styles) {
	c.styles = styles
}
