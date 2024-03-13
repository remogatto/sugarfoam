package viewport

// Package viewport provides a viewport model for managing scrollable
// content within a UI.  It integrates with the Bubble Tea framework
// and the SugarFoam UI framework to offer a flexible and customizable
// scrolling experience. The package supports styling and key bindings
// to enhance the user interaction with the viewport.
import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

// DefaultWidth and DefaultHeight define the default dimensions for
// the viewport.
var (
	DefaultWidth  = 80
	DefaultHeight = 25
)

// Option is a type for functions that modify a viewport model.
type Option func(*Model)

// Model wraps the Bubble Tea viewport model with additional
// functionality and styling options provided by the Sugarfoam
// framework.
type Model struct {
	foam.Common
	*viewport.Model

	focused bool
}

// New creates a new viewport model with optional configurations.
func New(opts ...Option) *Model {
	vp := viewport.New(DefaultWidth, DefaultHeight)

	v := &Model{
		Model: &vp,
	}

	v.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(v)
	}

	return v
}

// WithStyles sets the styles for the viewport using the provided
// styles.
func WithStyles(styles *foam.Styles) Option {
	return func(m *Model) {
		m.SetStyles(styles)
	}
}

// Init initializes the viewport model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update updates the viewport model based on the received message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	t, cmd := m.Model.Update(msg)
	m.Model = &t

	return m, cmd
}

// Focused returns the focus state of the viewport.
func (m *Model) Focused() bool {
	return m.focused
}

// Blur removes focus from the viewport.
func (m *Model) Blur() {
	m.focused = false
}

// Focus sets the viewport to be focused.
func (m *Model) Focus() tea.Cmd {
	m.focused = true

	return nil
}

func (m *Model) CanGrow() bool {
	return true
}

func (m *Model) GetHeight() int {
	return lipgloss.Height(m.View())
}

func (t *Model) SetWidth(width int) {
	t.Model.Width = width
	ww := lipgloss.Width(t.Model.View()) - width
	t.Model.Width = width - ww
}

func (t *Model) SetHeight(height int) {
	t.Model.Height = height

	hh := lipgloss.Height(t.Model.View()) - height

	t.Model.Height = height - hh
}

func (t *Model) SetSize(w, h int) {
	t.Common.SetSize(w, h)

	t.SetWidth(w)
	t.SetHeight(h)
}

// View renders the viewport model, applying the appropriate style
// based on focus state.
func (m *Model) View() string {
	if m.Focused() {
		return m.GetStyles().Focused.Render(m.Model.View())
	}
	return m.GetStyles().Blurred.Render(m.Model.View())
}
