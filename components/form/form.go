package form

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Model struct {
	foam.Common
	*huh.Form

	focused bool
}

func New(opts ...Option) *Model {
	f := huh.NewForm()

	form := &Model{
		Form: f,
	}

	form.Common.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(form)
	}

	return form
}

func WithGroups(groups ...*huh.Group) Option {
	return func(form *Model) {
		form.Form = huh.NewForm(groups...)
	}

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

func (m *Model) CanGrow() bool {
	return false
}

func (m *Model) GetHeight() int {
	return lipgloss.Height(m.View())
}

func (t *Model) SetWidth(width int) {
	t.Common.SetWidth(width)
	t.Form.WithWidth(width)

	ww := lipgloss.Width(t.Form.View()) - width
	availableW := width - ww

	t.Form.WithWidth(availableW)
}

func (t *Model) SetHeight(h int) {
	t.Common.SetHeight(h)
	t.Form.WithHeight(h)

	hh := lipgloss.Height(t.Form.View()) - h

	t.Form.WithHeight(h - hh)
}

func (t *Model) SetSize(w, h int) {
	t.Common.SetSize(w, h)

	t.SetWidth(w)
	t.SetHeight(h)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		if !m.focused {
			return m, nil
		}
	}

	f, cmd := m.Form.Update(msg)
	m.Form = f.(*huh.Form)

	return m, cmd
}

func (t *Model) View() string {
	if t.Focused() {
		return t.GetStyles().Focused.Render(t.Form.View())
	}
	return t.GetStyles().Blurred.Render(t.Form.View())
}
