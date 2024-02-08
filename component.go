package foam

import tea "github.com/charmbracelet/bubbletea"

type Focusable interface {
	tea.Model
	Placeable

	Focus() tea.Cmd
	Blur()
}
