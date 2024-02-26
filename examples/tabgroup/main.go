package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"

	"github.com/remogatto/sugarfoam/components/group"
	"github.com/remogatto/sugarfoam/components/header"
	"github.com/remogatto/sugarfoam/components/statusbar"
	"github.com/remogatto/sugarfoam/components/tabgroup"
	"github.com/remogatto/sugarfoam/components/tabgroup/tabitem"
	"github.com/remogatto/sugarfoam/components/table"
	"github.com/remogatto/sugarfoam/components/textinput"
	"github.com/remogatto/sugarfoam/components/viewport"
	"github.com/remogatto/sugarfoam/layout"
	"github.com/remogatto/sugarfoam/layout/tiled"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tabGroup  *tabgroup.Model
	help      help.Model
	statusBar *statusbar.Model
	spinner   spinner.Model

	document *layout.Layout

	bindings *keyBindings

	loadCompleted bool
}

type keyBindings struct {
	tabGroup *tabgroup.Model

	quit key.Binding
}

func (k *keyBindings) ShortHelp() []key.Binding {
	keys := make([]key.Binding, 0)

	currentTabItem := k.tabGroup.Current()

	switch tabItem := currentTabItem.(type) {
	case *tabitem.Model:
		keys = append(
			keys,
			tabItem.KeyMap.FocusNext,
			tabItem.KeyMap.FocusPrev,
		)

		switch model := tabItem.Current().(type) {
		case *table.Model:
			keys = append(
				keys,
				model.KeyMap.LineUp,
				model.KeyMap.LineDown,
			)
		}
	}

	keys = append(
		keys,
		k.tabGroup.KeyMap.TabNext,
		k.tabGroup.KeyMap.TabPrev,
		k.quit,
	)

	return keys
}

func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.quit,
		},
	}
}

func newBindings(tg *tabgroup.Model) *keyBindings {
	return &keyBindings{
		tabGroup: tg,
		quit: key.NewBinding(
			key.WithKeys("esc"), key.WithHelp("esc", "Quit app"),
		),
	}
}

func initialModel() model {
	textinput := textinput.New(
		textinput.WithPlaceholder("Text input goes here..."),
	)

	table1 := table.New()
	table1.SetHeight(25)

	viewport := viewport.New()
	table2 := table.New()
	help := help.New()

	group1 := group.New(
		group.WithItems(textinput, viewport, table1),
		group.WithLayout(
			layout.New(
				layout.WithStyles(&layout.Styles{Container: lipgloss.NewStyle().Padding(1, 0, 1, 0)}),
				layout.WithItem(textinput),
				layout.WithItem(tiled.New(viewport, table1)),
			),
		),
	)

	group2 := group.New(
		group.WithItems(table2),
		group.WithLayout(
			layout.New(
				layout.WithStyles(&layout.Styles{Container: lipgloss.NewStyle().Padding(1, 0, 1, 0)}),
				layout.WithItem(table2),
			),
		),
	)

	tabGroup := tabgroup.New(
		tabgroup.WithItems(
			tabitem.New(group1, tabitem.WithTitle("Tiled layout"), tabitem.WithActive(true)),
			tabitem.New(group2, tabitem.WithTitle("Single layout")),
		),
	)

	bindings := newBindings(tabGroup)
	statusBar := statusbar.New(bindings)

	header := header.New(
		header.WithContent(
			lipgloss.NewStyle().Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).Render("🧋Sugarfoam TabGroup Example🧋"),
		),
	)

	document := layout.New(
		layout.WithStyles(&layout.Styles{Container: docStyle}),
		layout.WithItem(header),
		layout.WithItem(tabGroup),
		layout.WithItem(statusBar),
	)

	return model{
		tabGroup:  tabGroup,
		statusBar: statusBar,
		document:  document,
		bindings:  bindings,
		help:      help,
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.tabGroup.Init())

	m.tabGroup.Focus()

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.document.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.quit):
			return m, tea.Quit

		}
	}

	_, cmd := m.tabGroup.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.document.View()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
