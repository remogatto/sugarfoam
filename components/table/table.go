package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Model struct {
	foam.Common

	*table.Model
}

func New(opts ...Option) *Model {
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

	ti := &Model{
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
	return func(ti *Model) {
		ti.Common.SetStyles(styles)
	}
}

func (t *Model) Focus() tea.Cmd {
	t.Model.Focus()

	return nil
}

func (t *Model) Blur() {
	t.Model.Blur()
}

func (t *Model) SetHeight(h int) {
	t.Model.SetHeight(h)

	hh := lipgloss.Height(t.Model.View()) - h

	t.Model.SetHeight(h - hh)
}

func (t *Model) Init() tea.Cmd {
	return nil
}

func (t *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	table, cmd := t.Model.Update(msg)

	t.Model = &table

	return t, cmd
}

func (t *Model) View() string {
	if t.Focused() {
		return t.GetStyles().Focused.Render(t.Model.View())
	}
	return t.GetStyles().Blurred.Render(t.Model.View())
}

func (t *Model) String() string { return t.View() }
