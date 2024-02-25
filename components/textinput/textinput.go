package textinput

// Package textinput provides a wrapper around the Bubble Tea
// textinput component, enhancing it with additional styling options
// and customization. It integrates with the Sugarfoam UI framework
// for a more flexible text input experience.
import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	foam "github.com/remogatto/sugarfoam"
)

// Option is a type for functions that modify a text input model.
type Option func(*Model)

// Model wraps the Bubble Tea textinput model with additional
// functionality and styling options provided by the Sugarfoam
// framework.
type Model struct {
	foam.Common

	*textinput.Model
}

// New creates a new text input model with optional configurations.
func New(opts ...Option) *Model {
	t := textinput.New()
	t.Placeholder = "Text here..."

	ti := &Model{
		Model: &t,
	}

	ti.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(ti)
	}

	return ti
}

// WithPlaceholder sets the placeholder text for the text input.
func WithPlaceholder(placeholder string) Option {
	return func(ti *Model) {
		ti.Model.Placeholder = placeholder
	}
}

// WithStyles sets the styles for the text input using the provided
// styles.
func WithStyles(styles *foam.Styles) Option {
	return func(ti *Model) {
		ti.SetStyles(styles)
	}
}

// Init initializes the text input model with a blink command.
func (ti *Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update updates the text input model based on the received message.
func (ti *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	t, cmd := ti.Model.Update(msg)
	ti.Model = &t

	return ti, cmd
}

// View renders the text input model, applying the appropriate style
// based on focus state.
func (ti *Model) View() string {
	if ti.Focused() {
		return ti.GetStyles().Focused.Render(ti.Model.View())
	}
	return ti.GetStyles().Blurred.Render(ti.Model.View())
}
