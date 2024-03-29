package tabgroup

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

func WithItems(items ...foam.Tabbable) Option {
	return func(m *Model) {
		m.items = items
	}
}

type KeyMap struct {
	TabNext key.Binding
	TabPrev key.Binding
}

type Styles struct {
	Navbar                lipgloss.Style
	NavbarTitleUnselected lipgloss.Style
	NavbarTitleSelected   lipgloss.Style
}

type Model struct {
	foam.Common

	KeyMap KeyMap

	items         []foam.Tabbable
	currItemIndex int

	focused bool
	styles  *Styles
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		TabNext: key.NewBinding(
			key.WithKeys("alt+right"),
			key.WithHelp("alt+right", "Next tab"),
		),
		TabPrev: key.NewBinding(
			key.WithKeys("alt+left"),
			key.WithHelp("alt+left", "Prev tab")),
	}
}

func New(opts ...Option) *Model {
	tg := &Model{
		items: make([]foam.Tabbable, 0),
	}

	tg.KeyMap = DefaultKeyMap()

	tg.Common.SetStyles(foam.DefaultStyles())
	tg.styles = DefaultStyles()

	for _, opt := range opts {
		opt(tg)
	}

	return tg
}

func (tg *Model) Items() []foam.Tabbable {
	return tg.items
}

func (tg *Model) AddItem(item foam.Tabbable) *Model {
	tg.items = append(tg.items, item)

	return tg
}

func DefaultStyles() *Styles {
	return &Styles{
		Navbar: lipgloss.NewStyle().Padding(1, 1, 0, 1),
		NavbarTitleUnselected: lipgloss.NewStyle().
			Background(lipgloss.Color("#373B41")).
			Foreground(lipgloss.Color("240")).
			Padding(0, 2, 0, 2),
		NavbarTitleSelected: lipgloss.NewStyle().
			Background(lipgloss.Color("5")).
			Foreground(lipgloss.Color("#ffffff")).
			Padding(0, 2, 0, 2),
	}
}

func (tg *Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	for _, item := range tg.items {
		cmds = append(cmds, item.Init())
	}

	return tea.Batch(cmds...)
}

func (tg *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, tg.KeyMap.TabNext):
			tg.nextTab()

		case key.Matches(msg, tg.KeyMap.TabPrev):
			tg.prevTab()
		}
	}

	cmds = append(cmds, tg.updateTabItems(msg)...)

	return tg, tea.Batch(cmds...)

}

func (tg *Model) Focus() tea.Cmd {
	tg.focused = true
	return nil
}

func (tg *Model) Blur() {
	tg.focused = false
}

func (tg *Model) SetSize(width int, height int) {
	tg.Common.SetSize(width, height)

	tg.styles.Navbar = tg.styles.Navbar.Width(tg.GetWidth())

	for _, item := range tg.items {
		item.SetSize(tg.GetWidth(), tg.GetHeight())
	}
}

func (tg *Model) View() string {
	var navbar string

	for i, item := range tg.Items() {
		tabTitle := tg.styles.NavbarTitleUnselected.Render(item.Title())
		if tg.currItemIndex == i {
			tabTitle = tg.styles.NavbarTitleSelected.Render(item.Title())
		}
		navbar += fmt.Sprintf("%s • ", tabTitle)
	}

	navbar = tg.styles.Navbar.Render(strings.TrimRight(navbar, "• "))

	if len(tg.items) > 0 {
		return tg.Common.GetStyles().NoBorder.Render(lipgloss.JoinVertical(lipgloss.Top, navbar, tg.items[tg.currItemIndex].View()))
	}

	return tg.Common.GetStyles().Focused.Render(navbar)

}

func (m *Model) Current() foam.Tabbable {
	return m.items[m.currItemIndex]
}

func (m *Model) CanGrow() bool {
	return true
}

func (tg *Model) nextTab() {
	tg.currItemIndex = (tg.currItemIndex + 1) % len(tg.items)

	tg.Current().SetSize(tg.GetWidth(), tg.GetHeight())
}

func (tg *Model) prevTab() {
	tg.currItemIndex = (tg.currItemIndex - 1) % len(tg.items)

	if tg.currItemIndex < 0 {
		tg.currItemIndex = len(tg.items) - 1
	}
	tg.Current().SetSize(tg.GetWidth(), tg.GetHeight())
}

func (tg *Model) updateTabItems(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, item := range tg.items {
		_, cmd := item.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}
