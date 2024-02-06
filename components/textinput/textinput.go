package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Model struct {
	foam.Common

	*textinput.Model
}

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

func WithPlaceholder(placeholder string) Option {
	return func(ti *Model) {
		ti.Model.Placeholder = placeholder
	}
}

func WithStyles(styles *foam.Styles) Option {
	return func(ti *Model) {
		ti.SetStyles(styles)
	}
}

func (ti *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (ti *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	t, cmd := ti.Model.Update(msg)
	ti.Model = &t

	return ti, cmd
}

func (ti *Model) View() string {
	if ti.Focused() {
		return ti.GetStyles().Focused.Render(ti.Model.View())
	}
	return ti.GetStyles().Blurred.Render(ti.Model.View())
}
