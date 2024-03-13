package help

import (
	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

func WithStyles(styles *foam.Styles) Option {
	return func(ti *Model) {
		ti.Common.SetStyles(styles)
	}
}

type Model struct {
	foam.Common
	*help.Model

	bindings help.KeyMap
}

func New(bindings help.KeyMap, opts ...Option) *Model {
	h := help.New()

	m := &Model{
		Model:    &h,
		bindings: bindings,
	}

	m.Common.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	help, cmd := m.Model.Update(msg)

	m.Model = &help

	return m, cmd
}

func (m *Model) View() string {
	return m.GetStyles().NoBorder.Render(m.Model.View(m.bindings))
}

func (m *Model) CanGrow() bool {
	return false
}

func (m *Model) GetHeight() int {
	return lipgloss.Height(m.View())
}
