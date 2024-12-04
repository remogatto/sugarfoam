package main

import (
	"fmt"
	"log"
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

const (
	DraculaStyle = "dracula"
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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleWindowSize(msg)
	case tea.KeyMsg:
		cmds = m.handleKeyMsg(msg, cmds)
	case checkConnectionMsg:
		cmds = m.handleCheckConnectionMsg(msg, cmds)
		m.state = DownloadingState
	case charactersResponseMsg:
		m.handleCharactersResponseMsg(msg)
		m.state = BrowseState
	default:
		if m.state == DownloadingState {
			return m.handleSpinner(msg)
		}

	}

	_, cmd = m.group.Update(msg)
	cmds = append(cmds, cmd)

	if m.state == BrowseState {
		currRow := m.table.Model.Cursor() + 1
		m.statusBar.SetContent(formats[BrowseState][0],
			fmt.Sprintf(formats[BrowseState][1], currRow, len(m.characters)),
			formats[BrowseState][2],
		)

		cmd = m.updateViewport()
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *model) handleSpinner(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	m.statusBar.SetContent(
		fmt.Sprintf(formats[DownloadingState][0], m.spinner.View()),
		formats[DownloadingState][1],
		formats[DownloadingState][2],
	)

	return m, cmd

}

func (m *model) handleWindowSize(msg tea.WindowSizeMsg) {
	m.group.SetSize(msg.Width, msg.Height)
	m.document.SetSize(msg.Width, msg.Height)
	m.renderer = m.createRenderer(msg.Width)
	m.api.limit = m.group.GetHeight() * 2
}

func (m *model) createRenderer(width int) *glamour.TermRenderer {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle(DraculaStyle),
		glamour.WithWordWrap(m.viewport.GetWidth()),
	)
	if err != nil {
		// Log the error instead of panicking
		log.Println("Error creating renderer:", err)
		return nil
	}
	return renderer
}

func (m *model) handleKeyMsg(msg tea.KeyMsg, cmds []tea.Cmd) []tea.Cmd {
	if key.Matches(msg, m.bindings.quit) {
		return append(cmds, tea.Quit)
	}
	return cmds
}

func (m *model) handleCheckConnectionMsg(msg checkConnectionMsg, cmds []tea.Cmd) []tea.Cmd {
	if msg {
		cmds = append(cmds, m.api.getCharacters, m.spinner.Tick)
	}

	return cmds
}

func (m *model) handleCharactersResponseMsg(msg charactersResponseMsg) {
	m.updateTableRows(msg.Data)
	m.updateViewport()
}

// View renders the application UI.
func (m model) View() string {
	return m.document.View()
}

// updateViewport updates the viewport with the selected character's details.
func (m model) updateViewport() tea.Cmd {
	if m.table.Cursor() >= 0 {
		if m.table.Cursor() >= len(m.characters)-1 {
			m.api.page++
			return m.api.checkConnection
		}

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

	return nil
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
