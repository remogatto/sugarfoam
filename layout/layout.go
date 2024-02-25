package layout

// Package layout provides utilities for managing the layout of UI components.
// It offers a flexible way to organize components into a container, with support
// for styling and resizing. The package is designed to work with the Lipgloss library
// for styling, providing a visually appealing interface.
import (
	"github.com/charmbracelet/lipgloss"
)

// DefaultWidth and DefaultHeight define the default dimensions for a layout.
var (
	DefaultWidth  = 80
	DefaultHeight = 25
)

// Option is a type for functions that modify a Layout.
type Option func(*Layout)

// WithItem adds a Placeable item to the layout.
func WithItem(i Placeable) Option {
	return func(l *Layout) {
		l.items = append(l.items, i)
	}
}

// WithStyles sets the styles for the layout container.
func WithStyles(styles *Styles) Option {
	return func(l *Layout) {
		l.styles = styles
	}
}

// Sizeable is an interface for components that can be resized.
type Sizeable interface {
	SetSize(width int, height int)
	GetWidth() int
	GetHeight() int
}

// Placeable is an interface for components that can be placed within a layout.
// It extends the Sizeable interface, indicating that placeable items must
// be resizable.
type Placeable interface {
	Sizeable

	// View returns the string representation of the component.
	View() string
}

// Styles defines the styles for the layout container.
type Styles struct {
	Container lipgloss.Style
}

// Layout represents a layout container for UI components.
type Layout struct {
	width, height int
	items         []Placeable
	styles        *Styles
}

// New creates a new Layout with optional configurations.
func New(opts ...Option) *Layout {
	l := &Layout{width: DefaultWidth, height: DefaultHeight}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// AddItem adds a Placeable item to the layout.
func (l *Layout) AddItem(i Placeable) *Layout {
	l.items = append(l.items, i)

	return l
}

// Items returns the items currently in the layout.
func (l *Layout) Items() []Placeable {
	return l.items
}

// SetSize resizes the layout and its items, adjusting for container styling.
func (l *Layout) SetSize(width int, height int) {
	for _, item := range l.items {
		item.SetSize(width-l.styles.Container.GetHorizontalFrameSize(), height-l.styles.Container.GetVerticalFrameSize())
	}
}

// View returns the string representation of the layout, rendered with its container styles.
func (l *Layout) View() string {
	views := []string{}

	for _, item := range l.Items() {
		views = append(views, item.View())
	}

	return l.styles.Container.Render(lipgloss.JoinVertical(lipgloss.Top, views...))
}
