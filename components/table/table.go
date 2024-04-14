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

	RelWidths []int
}

func New(opts ...Option) *Model {
	t := table.New()

	relWidths := make([]int, len(t.Columns()))

	for i := range t.Columns() {
		relWidths[i] = 100
	}

	ti := &Model{
		Model:     &t,
		RelWidths: relWidths,
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

func WithRelWidths(percentages ...int) Option {
	return func(m *Model) {
		m.RelWidths = make([]int, 0)
		m.RelWidths = append(m.RelWidths, percentages...)
	}
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

func (t *Model) SetWidth(width int) {
	t.Model.SetWidth(width)

	ww := lipgloss.Width(t.Model.View()) - width
	availableW := width - ww

	t.Model.SetWidth(availableW)

	cols := make([]table.Column, 0)

	for i, col := range t.Model.Columns() {
		colW := availableW * t.RelWidths[i] / 100
		col.Width = colW - table.DefaultStyles().Cell.GetHorizontalFrameSize() - 1
		cols = append(cols, col)
	}

	t.Model.SetColumns(cols)
}

func (t *Model) SetHeight(h int) {
	t.Model.SetHeight(h)

	hh := lipgloss.Height(t.Model.View()) - h

	t.Model.SetHeight(h - hh)
}

func (t *Model) SetSize(w, h int) {
	t.Common.SetSize(w, h)

	t.SetWidth(w)
	t.SetHeight(h)
}

func (m *Model) SetRelWidths(percentages ...int) {
	m.RelWidths = make([]int, 0)
	m.RelWidths = append(m.RelWidths, percentages...)

	m.SetWidth(m.GetWidth())
}

func (m *Model) GetHeight() int {
	return lipgloss.Height(m.View())
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

func (m *Model) CanGrow() bool {
	return true
}
