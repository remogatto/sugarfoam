package header

import (
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Model struct {
	foam.Common

	content string
}

func WithContent(text string) Option {
	return func(m *Model) {
		m.content = text
	}
}

func New(opts ...Option) *Model {
	m := &Model{}

	m.Common.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) View() string {
	return m.Common.GetStyles().NoBorder.Render(m.content)
}

func (m *Model) CanGrow() bool {
	return false
}

func (m *Model) GetHeight() int {
	return lipgloss.Height(m.View())
}
