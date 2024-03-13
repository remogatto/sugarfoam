package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	btTable "github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/sugarfoam/components/group"
	"github.com/remogatto/sugarfoam/components/header"
	"github.com/remogatto/sugarfoam/components/statusbar"
	"github.com/remogatto/sugarfoam/components/table"
	"github.com/remogatto/sugarfoam/components/viewport"
	"github.com/remogatto/sugarfoam/layout"
	"github.com/remogatto/sugarfoam/layout/tiled"
)

const (
	GotConnectionState = iota
	GotResponseState
	BeginDownloadState
	DownloadingState
	BrowseState
)

const characterTpl = `
* **ID**: %s
* **Name**: %s

# Description

%s
`

type model struct {
	group     *group.Model
	statusBar *statusbar.Model
	table     *table.Model
	viewport  *viewport.Model
	spinner   spinner.Model
	renderer  *glamour.TermRenderer
	document  *layout.Layout

	bindings *keyBindings

	characters map[string]character

	api   *swDbApi
	state int
}

type keyBindings struct {
	group *group.Model
	quit  key.Binding
}

func (k *keyBindings) ShortHelp() []key.Binding {
	keys := make([]key.Binding, 0)

	switch currItem := k.group.Current().(type) {
	case *table.Model:
		keys = append(
			keys,
			currItem.KeyMap.LineDown,
			currItem.KeyMap.LineUp,
		)
	}

	keys = append(
		keys,
		k.group.KeyMap.FocusNext,
		k.group.KeyMap.FocusPrev,
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

func newBindings(group *group.Model) *keyBindings {
	return &keyBindings{
		group: group,
		quit: key.NewBinding(
			key.WithKeys("esc"), key.WithHelp("esc", "Quit app"),
		),
	}
}

func initialModel() model {
	viewport := viewport.New()

	table := table.New(table.WithRelWidths(30, 70))
	table.Model.SetColumns([]btTable.Column{
		{Title: "ID", Width: 20},
		{Title: "Name", Width: 10},
	})

	group := group.New(
		group.WithItems(table, viewport),
		group.WithLayout(
			layout.New(
				layout.WithStyles(&layout.Styles{Container: lipgloss.NewStyle().Padding(1, 1)}),
				layout.WithItem(tiled.New(table, viewport)),
			),
		),
	)

	bindings := newBindings(group)
	statusBar := statusbar.New(bindings, statusbar.WithContent("Idle", "", "API 🔴"))

	header := header.New(
		header.WithContent(
			lipgloss.NewStyle().Bold(true).Border(lipgloss.NormalBorder(), false, false, true, false).Render("⭐🌌 Star Wars characters browser 🌌⭐"),
		),
	)

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dracula"),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		panic(err)
	}

	document := layout.New(
		layout.WithStyles(&layout.Styles{Container: lipgloss.NewStyle().Margin(1)}),
		layout.WithItem(header),
		layout.WithItem(group),
		layout.WithItem(statusBar),
	)

	s := spinner.New(
		spinner.WithStyle(
			lipgloss.NewStyle().Foreground(lipgloss.Color("265"))),
	)
	s.Spinner = spinner.Dot

	return model{
		group:      group,
		table:      table,
		viewport:   viewport,
		statusBar:  statusBar,
		spinner:    s,
		document:   document,
		renderer:   renderer,
		bindings:   bindings,
		characters: make(map[string]character),
		api:        &swDbApi{1, 20},
	}
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.group.Init(), m.api.ping)

	m.group.Focus()

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.group.SetSize(msg.Width, msg.Height)
		m.document.SetSize(msg.Width, msg.Height)

		renderer, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle("dracula"),
			glamour.WithWordWrap(m.viewport.GetWidth()),
		)
		if err != nil {
			panic(err)
		}

		m.renderer = renderer

		m.api.limit = m.group.GetHeight() * 2

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.quit):
			return m, tea.Quit

		}
	case onlineMsg:
		cmds = append(cmds, m.api.getCharacters, m.spinner.Tick)
		m.state = DownloadingState

	case charactersResponseMsg:
		m.state = GotResponseState
	}

	return m, tea.Batch(append(cmds, m.handleState(msg, cmds)...)...)
}

func (m model) View() string {
	return m.document.View()
}

func (m *model) handleState(msg tea.Msg, cmds []tea.Cmd) []tea.Cmd {
	_, cmd := m.group.Update(msg)
	switch m.state {
	case BrowseState:
		m.updateViewport()

	case GotResponseState:
		m.setTableRows(msg.(charactersResponseMsg).Data)
		m.updateViewport()
		m.statusBar.SetContent("Browse 📖", "Browse the results using arrow keys", "API 🟢")
		m.state = BrowseState

	case DownloadingState:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		m.statusBar.SetContent(fmt.Sprintf("Loading %s", m.spinner.View()), "Fetching results from the endpoint", "API 🟢")
	}

	return append(cmds, cmd)
}

func (m model) updateViewport() {
	if currentRow := m.table.SelectedRow(); currentRow != nil {
		currentId := currentRow[0]
		character := m.characters[currentId]

		// FIXME: The use of a standard '-' character causes rendering
		// issues within the viewport. Further investigation is
		// required to resolve this issue.
		description := strings.Replace(character.Description, "-", "–", -1)

		md, _ := m.renderer.Render(fmt.Sprintf(characterTpl, character.ID, character.Name, description))
		m.viewport.SetContent(md)
	}
}

func (m model) setTableRows(data []character) {
	rows := make([]btTable.Row, 0)

	for _, c := range data {
		rows = append(rows, btTable.Row{c.ID, c.Name})
		m.characters[c.ID] = c
	}

	m.table.Model.SetRows(rows)
}

func main() {
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	if _, err := tea.NewProgram(initialModel() /*, tea.WithAltScreen()*/).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
