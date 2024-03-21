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
	"github.com/remogatto/sugarfoam/components/image"
	"github.com/remogatto/sugarfoam/components/statusbar"
	"github.com/remogatto/sugarfoam/components/table"
	"github.com/remogatto/sugarfoam/components/viewport"
	"github.com/remogatto/sugarfoam/layout"
	"github.com/remogatto/sugarfoam/layout/tiled"
)

// Application states.
const (
	CheckConnectionState = iota
	GotConnectionState
	GotResponseState
	DownloadingState
	BrowseState
)

// characterTpl is a template for displaying character information.
const characterTpl = `
* **ID**: %s
* **Name**: %s

# Description

%s
`

type loadImgMsg struct {
	url string
}

// model represents the application state.
type model struct {
	group     *group.Model
	table     *table.Model
	viewport  *viewport.Model
	image     *image.Model
	statusBar *statusbar.Model
	spinner   spinner.Model
	renderer  *glamour.TermRenderer
	document  *layout.Layout

	bindings *keyBindings

	characters []character

	api   *swDbApi
	state int
}

// keyBindings holds the key bindings for the application.
type keyBindings struct {
	group *group.Model
	quit  key.Binding
}

// ShortHelp returns a list of key bindings for the current state.
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

// FullHelp returns a list of all key bindings.
func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.quit,
		},
	}
}

// newBindings creates a new set of key bindings for the application.
func newBindings(group *group.Model) *keyBindings {
	return &keyBindings{
		group: group,
		quit: key.NewBinding(
			key.WithKeys("esc"), key.WithHelp("esc", "Quit app"),
		),
	}
}

// initialModel initializes the application model.
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
	statusBar := statusbar.New(bindings, statusbar.WithContent(formats[CheckConnectionState]...))

	header := header.New(
		header.WithContent(
			lipgloss.NewStyle().Bold(true).Border(
				lipgloss.NormalBorder(),
				false,
				false,
				true,
				false).Render("⭐🌌 Star Wars characters browser 🌌⭐"),
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
		characters: make([]character, 0),
		api:        &swDbApi{1, 20},
	}
}

// Init initializes the application.
func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.group.Init(), m.api.checkConnection)

	m.group.Focus()

	return tea.Batch(cmds...)
}

// Update updates the application state based on messages.
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

	case checkConnectionMsg:
		m.state = CheckConnectionState

		if msg {
			cmds = append(cmds, m.api.getCharacters, m.spinner.Tick)
			m.state = DownloadingState
		}

	case charactersResponseMsg:
		m.state = GotResponseState

	}

	return m, tea.Batch(m.handleState(msg, cmds)...)
}

// View renders the application UI.
func (m model) View() string {
	return m.document.View()
}

// handleState updates the application state based on the current state.
func (m *model) handleState(msg tea.Msg, cmds []tea.Cmd) []tea.Cmd {
	_, cmd := m.group.Update(msg)

	switch m.state {

	case CheckConnectionState:
		m.statusBar.SetContent(formats[CheckConnectionState]...)

	case BrowseState:
		m.updateViewport()
		currRow := m.table.Model.Cursor() + 1

		if currRow < len(m.characters) {
			m.statusBar.SetContent(formats[BrowseState][0],
				fmt.Sprintf(formats[BrowseState][1], currRow, len(m.characters)),
				formats[BrowseState][2],
			)
		} else {
			m.api.page++
			cmds = append(cmds, m.api.checkConnection)
		}

	case GotResponseState:
		m.updateTableRows(msg.(charactersResponseMsg).Data)
		m.updateViewport()

		m.statusBar.SetContent(formats[BrowseState]...)
		m.state = BrowseState

	case DownloadingState:
		var cmd tea.Cmd

		m.spinner, cmd = m.spinner.Update(msg)
		m.statusBar.SetContent(fmt.Sprintf(formats[DownloadingState][0], m.spinner.View()), formats[DownloadingState][1], formats[DownloadingState][2])

		cmds = append(cmds, cmd)
	}

	return append(cmds, cmd)
}

func (m model) updateImage() tea.Msg {
	character := m.characters[m.table.Cursor()]
	return loadImgMsg{character.ImageURL}
}

// updateViewport updates the viewport with the selected character's details.
func (m model) updateViewport() {
	character := m.characters[m.table.Cursor()]
	md, _ := m.renderer.Render(
		fmt.Sprintf(
			characterTpl,
			character.ID,
			sanitize(character.Name),
			sanitize(character.Description),
		),
	)
	m.viewport.SetContent(md)
}

// setTableRows sets the table rows with the provided character data.
func (m *model) updateTableRows(data []character) {
	for _, c := range data {
		m.characters = append(m.characters, c)
	}

	rows := make([]btTable.Row, 0)

	for _, character := range m.characters {
		rows = append(rows, btTable.Row{character.ID, sanitize(character.Name)})
	}

	m.table.Model.SetRows(rows)
}

func sanitize(text string) string {
	// FIXME: The use of a standard '-' character causes rendering
	// issues within the viewport. Further investigation is
	// required to resolve this problem.
	return strings.Replace(text, "-", "–", -1)
}

// main is the entry point of the application.
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
