package group

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

type Modelable interface {
	tea.Model
	foam.Placeable

	Focus() tea.Cmd
	Blur()
}

type Model struct {
	KeyMap KeyMap

	layout *foam.Layout
	items  []Modelable

	width, height int
	styles        *Styles
	focused       bool
	currFocus     int
}

type Styles struct {
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style
}

func New(opts ...Option) *Model {
	group := new(Model)

	for _, opt := range opts {
		opt(group)
	}

	group.KeyMap = DefaultKeyMap()
	group.styles = DefaultStyles()

	return group
}

type KeyMap struct {
	FocusNext key.Binding
	FocusPrev key.Binding
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		FocusNext: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "focus next"),
		),
		FocusPrev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "focus prev")),
	}
}

func DefaultStyles() *Styles {
	return &Styles{
		FocusedBorder: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")),
		BlurredBorder: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
	}
}

func (g *Model) SetSize(width int, height int) {
	g.width = width - g.styles.FocusedBorder.GetHorizontalFrameSize()
	g.height = height

	g.styles.FocusedBorder = g.styles.FocusedBorder.Width(g.width)
	g.styles.BlurredBorder = g.styles.BlurredBorder.Width(g.width)

	g.layout.SetSize(width, height)
}

func (g *Model) Width() int  { return g.width }
func (g *Model) Height() int { return g.height }

func (g *Model) Current() Modelable {
	return g.items[g.currFocus]
}

func (g *Model) Blur() {
	g.focused = false
}

func (g *Model) Focus() tea.Cmd {
	g.focused = true
	return g.Current().Focus()
}

func (g *Model) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, c := range g.items {
		cmds = append(cmds, c.Init())
	}

	g.Current().Focus()

	return tea.Batch(cmds...)
}

func (g *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, g.KeyMap.FocusNext):
			cmds = append(cmds, g.nextFocus())

		}
	}

	cmds = append(cmds, g.updateComponents(msg)...)

	return g, tea.Batch(cmds...)
}

func (g *Model) View() string {
	return g.layout.View()
}

func WithKeyMap(km KeyMap) Option {
	return func(g *Model) {
		g.KeyMap = km
	}
}
func WithItems(items ...Modelable) Option {
	return func(g *Model) {
		g.items = items
	}
}

func WithLayout(layout *foam.Layout) Option {
	return func(g *Model) {
		g.layout = layout
	}
}

func (g *Model) nextFocus() tea.Cmd {
	g.Current().Blur()

	g.currFocus = (g.currFocus + 1) % len(g.items)

	return g.Current().Focus()
}

func (g *Model) updateComponents(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, item := range g.items {
		_, cmd := item.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}
