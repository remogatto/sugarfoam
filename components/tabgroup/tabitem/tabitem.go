package tabitem

import (
	"github.com/remogatto/sugarfoam/components/group"
)

type Option func(*Model)

type Model struct {
	*group.Model

	title  string
	active bool
}

func (m *Model) Active() bool {
	return m.active
}

func (m *Model) Title() string {
	return m.title
}

func New(group *group.Model, opts ...Option) *Model {
	m := &Model{
		Model:  group,
		title:  "Tab Item",
		active: false,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func WithActive(active bool) Option {
	return func(m *Model) {
		m.active = active
	}
}

func WithTitle(title string) Option {
	return func(m *Model) {
		m.title = title
	}
}

func WithGroup(group *group.Model) Option {
	return func(m *Model) {
		m.Model = group
	}
}
