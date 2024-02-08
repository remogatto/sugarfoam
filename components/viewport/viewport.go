package viewport

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Model struct {
	foam.Common
	*viewport.Model

	focused bool
}

func New(width, heigth int, opts ...Option) *Model {
	vp := viewport.New(width, heigth)

	v := &Model{
		Model: &vp,
	}

	v.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(v)
	}

	return v

}

func WithStyles(styles *foam.Styles) Option {
	return func(m *Model) {
		m.SetStyles(styles)
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	t, cmd := m.Model.Update(msg)
	m.Model = &t

	return m, cmd
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Blur() {
	m.focused = false
}

func (m *Model) Focus() tea.Cmd {
	m.focused = true

	return nil
}

// func (m *Model) SetSize(width int, height int) {
// 	m.Common.SetSize(width, height)
// 	m.Model.Height = m.GetHeight()
// 	log.Println(m.Height)
// }

func (m *Model) View() string {
	if m.Focused() {
		return m.GetStyles().Focused.Render(m.Model.View())
	}
	return m.GetStyles().Blurred.Render(m.Model.View())
}
