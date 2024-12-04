package textarea

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Model struct {
	foam.Common

	*textarea.Model

	heightCorrection int
}

func New(opts ...Option) *Model {
	t := textarea.New()
	t.Placeholder = ""

	m := &Model{
		Model: &t,
	}

	m.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// WithPlaceholder sets the placeholder text for the text input.
func WithPlaceholder(placeholder string) Option {
	return func(m *Model) {
		m.Model.Placeholder = placeholder
	}
}

// WithStyles sets the styles for the text input using the provided
// styles.
func WithStyles(styles *foam.Styles) Option {
	return func(m *Model) {
		m.SetStyles(styles)
	}
}

// Init initializes the text input model with a blink command.
func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

// Update updates the text input model based on the received message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	t, cmd := m.Model.Update(msg)
	m.Model = &t

	diff := m.LineCount() - m.Line()
	if diff > 1 {
		m.heightCorrection = 1
	}

	return m, cmd
}

// View renders the text input model, applying the appropriate style
// based on focus state.
func (m *Model) View() string {
	if m.Focused() {
		return m.GetStyles().Focused.Render(m.Model.View())
	}
	return m.GetStyles().Blurred.Render(m.Model.View())
}

func (m *Model) CanGrow() bool {
	return true
}

// func (m *Model) GetHeight() int {
// 	return lipgloss.Height(m.View())
// }

func (m *Model) SetWidth(width int) {
	m.Model.SetWidth(width)
	ww := lipgloss.Width(m.Model.View()) - width + 1
	m.Model.SetWidth(width - ww)
}

func (m *Model) SetHeight(height int) {
	m.Model.SetHeight(height)
	hh := lipgloss.Height(m.Model.View()) - height + 1
	m.Model.SetHeight(height - hh)
}

func (t *Model) SetSize(w, h int) {
	t.Common.SetSize(w, h)

	t.SetWidth(w)
	t.SetHeight(h)
}
