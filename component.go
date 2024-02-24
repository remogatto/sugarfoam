package foam

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remogatto/sugarfoam/layout"
)

type Focusable interface {
	tea.Model
	layout.Placeable

	Focus() tea.Cmd
	Blur()
}

type Groupable interface {
	Focusable

	Current() Focusable
}

type Tabbable interface {
	Groupable

	Title() string
	Active() bool
}
