package table

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Option func(*Table)

type Styles struct {
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style
}

type Table struct {
	m             *table.Model
	width, height int
	styles        *Styles
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
		m: &t,
	}

	ti.styles = DefaultStyles()

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
	ti.m.SetStyles(s)

	for _, opt := range opts {
		opt(ti)
	}

	return ti
}

func DefaultStyles() *Styles {
	return &Styles{
		FocusedBorder: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")),
		BlurredBorder: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
	}
}

func WithStyles(styles *Styles) Option {
	return func(ti *Table) {
		ti.styles = styles
	}
}

func (t *Table) SetSize(width int, height int) {
	t.width = width - t.styles.FocusedBorder.GetHorizontalFrameSize()
	t.height = height

	t.styles.FocusedBorder = t.styles.FocusedBorder.Width(t.width)
	t.styles.BlurredBorder = t.styles.BlurredBorder.Width(t.width)
}

func (t *Table) Width() int  { return t.width }
func (t *Table) Height() int { return t.height }

func (t *Table) Focus() tea.Cmd {
	t.m.Focus()

	return nil
}

func (t *Table) Blur() tea.Cmd {
	t.m.Blur()

	return nil
}

func (t *Table) Init() tea.Cmd {
	return nil
}

func (t *Table) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	table, cmd := t.m.Update(msg)
	t.m = &table

	return t, cmd
}

func (t *Table) View() string {
	if t.m.Focused() {
		return t.styles.FocusedBorder.Render(t.m.View())
	}
	return t.styles.BlurredBorder.Render(t.m.View())
}

func (t *Table) String() string { return t.View() }
