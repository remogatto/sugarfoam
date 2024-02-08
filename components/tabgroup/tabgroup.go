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

type KeyMap struct {
	TabNext key.Binding
	TabPrev key.Binding
}

type Styles struct {
	Navbar                lipgloss.Style
	NavbarTitleUnselected lipgloss.Style
	NavbarTitleSelected   lipgloss.Style
}

type TabItem struct {
	foam.Focusable

	Title  string
	Active bool
}

type Model struct {
	foam.Common

	KeyMap KeyMap

	items         []*TabItem
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
		items: make([]*TabItem, 0),
	}

	tg.KeyMap = DefaultKeyMap()

	tg.Common.SetStyles(foam.DefaultStyles())
	tg.styles = DefaultStyles()

	return tg
}

func (tg *Model) Items() []*TabItem {
	return tg.items
}

func (tg *Model) AddItem(item *TabItem) *Model {
	tg.items = append(tg.items, item)

	return tg
}

func DefaultStyles() *Styles {
	return &Styles{
		Navbar: lipgloss.NewStyle().Padding(0, 1, 0, 1),
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
		cmds = append(cmds, item.Focusable.Init())
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
		item.Focusable.SetSize(tg.GetWidth(), tg.GetHeight())
	}
}

func (tg *Model) View() string {
	var navbar string

	for i, item := range tg.Items() {
		tabTitle := tg.styles.NavbarTitleUnselected.Render(item.Title)
		if tg.currItemIndex == i {
			tabTitle = tg.styles.NavbarTitleSelected.Render(item.Title)
		}
		navbar += fmt.Sprintf("%s • ", tabTitle)
	}

	navbar = tg.styles.Navbar.Render(strings.TrimRight(navbar, "• "))

	if len(tg.items) > 0 {
		return tg.Common.GetStyles().Focused.Render(lipgloss.JoinVertical(lipgloss.Top, navbar, tg.items[tg.currItemIndex].Focusable.View()))
	}

	return tg.Common.GetStyles().Focused.Render(navbar)

}

func (tg *Model) nextTab() {
	tg.currItemIndex = (tg.currItemIndex + 1) % len(tg.items)
}

func (tg *Model) prevTab() {
	tg.currItemIndex = (tg.currItemIndex - 1) % len(tg.items)

	if tg.currItemIndex < 0 {
		tg.currItemIndex = len(tg.items) - 1
	}
}

func (tg *Model) updateTabItems(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, item := range tg.items {
		_, cmd := item.Focusable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}
