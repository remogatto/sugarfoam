package main

import (
	"fmt"
	"os"

	"git.andreafazzi.eu/andrea/probo/lib/store/file"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/bubbletea-app/ui"
	"github.com/remogatto/bubbletea-app/ui/components/table"
	"github.com/remogatto/bubbletea-app/ui/components/textinput"

	tea "github.com/charmbracelet/bubbletea"
)

type focusState int

const (
	inputFocus focusState = iota
	jsonViewportFocus
	tableFocus

	tableAndViewportHeight = 20
)

type keyBindings struct {
	Quit key.Binding
	Tab  key.Binding

	// Enter key.Binding
	// table.KeyMap
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyBindings) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
		k.Tab,
		// k.Enter,
		// k.LineUp,
		// k.LineDown,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Quit,
			k.Tab,
			// k.Enter,
			// k.LineUp,
			// k.LineDown,
		},
	}
}

func newBindings() keyBindings {
	return keyBindings{
		key.NewBinding(
			key.WithKeys("esc"), key.WithHelp("esc", "Quit app"),
		),
		key.NewBinding(
			key.WithKeys("tab"), key.WithHelp("tab", "Switch table/input"),
		),
		// key.NewBinding(
		// 	key.WithKeys("enter"), key.WithHelp("enter", "Submit query"),
		// ),

		//		table.KeyMap,
	}
}

type model struct {
	components []ui.Component
	bindings   keyBindings

	pStore       *file.ParticipantFileStore
	pJson        []byte
	participants []any
}

func initialModel() model {
	components := []ui.Component{
		textinput.New(
			textinput.WithStyles(&textinput.Styles{
				BlurBorder:    blurTextInputStyle,
				FocusedBorder: focusedTextInputStyle,
			}),
			textinput.WithPlaceholder("Insert a jq filter..."),
		),
		table.New(),
	}

	return model{
		components: components,
		bindings:   newBindings(),
	}
}

func (m *model) setSize(msg tea.WindowSizeMsg) {
	for _, c := range m.components {
		c.SetSize(msg)
	}
}

func (m *model) updateComponents(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, c := range m.components {
		_, cmd := c.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m model) Init() tea.Cmd {
	return m.components[0].Focus()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.setSize(msg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.Quit):
			return m, tea.Quit

		}
	}

	cmds = append(cmds, m.updateComponents(msg)...)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	views := []string{}

	for _, c := range m.components {
		views = append(views, c.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, views...)
	// for _, c := range m.components {
	// 	content += c.View()
	// }

	// return content
}

func main() {
	if _, err := tea.NewProgram(initialModel() /*tea.WithAltScreen()*/).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
