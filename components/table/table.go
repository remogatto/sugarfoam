package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Table)

type Table struct {
	foam.Common

	*table.Model
}

func New(opts ...Option) *Table {
	t := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "Col1", Width: 10},
				{Title: "Col2", Width: 10},
				{Title: "Col3", Width: 10},
			},
		),
		table.WithRows(
			[]table.Row{
				[]string{"Cell 11", "Cell 12", "Cell 13"},
				[]string{"Cell 21", "Cell 22", "Cell 23"},
				[]string{"Cell 31", "Cell 32", "Cell 33"},
				[]string{"Cell 41", "Cell 42", "Cell 43"},
			},
		),
	)

	ti := &Table{
		Model: &t,
	}

	ti.Common.SetStyles(foam.DefaultStyles())

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	ti.Model.SetStyles(s)

	for _, opt := range opts {
		opt(ti)
	}

	return ti
}

func WithStyles(styles *foam.Styles) Option {
	return func(ti *Table) {
		ti.Common.SetStyles(styles)
	}
}

func (t *Table) Focus() tea.Cmd {
	t.Model.Focus()

	return nil
}

func (t *Table) Blur() {
	t.Model.Blur()
}

func (t *Table) Init() tea.Cmd {
	return nil
}

func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	table, cmd := t.Model.Update(msg)

	t.Model = &table

	return t, cmd
}

func (t *Table) View() string {
	if t.Focused() {
		return t.GetStyles().Focused.Render(t.Model.View())
	}
	return t.GetStyles().Blurred.Render(t.Model.View())
}

func (t *Table) String() string { return t.View() }
