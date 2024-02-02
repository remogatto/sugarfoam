package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/remogatto/bubbletea-app/ui"
	"github.com/remogatto/bubbletea-app/ui/components/table"
	"github.com/remogatto/bubbletea-app/ui/components/textinput"
	"github.com/remogatto/bubbletea-app/ui/layout/tiled"

	tea "github.com/charmbracelet/bubbletea"
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
	layout     *ui.Layout
	components []ui.Component
	bindings   keyBindings

	currFocus int
}

func initialModel() model {
	textinput := textinput.New(
		textinput.WithStyles(&textinput.Styles{
			BlurredBorder: blurredTextInputStyle,
			FocusedBorder: focusedTextInputStyle,
		}),
		textinput.WithPlaceholder("Text input goes here..."),
	)

	table1 := table.New(
		table.WithStyles(&table.Styles{
			BlurredBorder: blurredTableStyle,
			FocusedBorder: focusedTableStyle,
		}),
	)

	table2 := table.New(
		table.WithStyles(&table.Styles{
			BlurredBorder: blurredTableStyle,
			FocusedBorder: focusedTableStyle,
		}),
	)

	table3 := table.New(
		table.WithStyles(&table.Styles{
			BlurredBorder: blurredTableStyle,
			FocusedBorder: focusedTableStyle,
		}),
	)

	return model{
		layout: ui.NewLayout(80, 25).
			AddItem(textinput).
			AddItem(tiled.New(80, 25, table1, table2, table3)),
		components: []ui.Component{textinput, table1, table2},
		bindings:   newBindings(),
	}
}

func (m *model) setSize(msg tea.WindowSizeMsg) {
	m.layout.SetSize(msg.Width, msg.Height)
}

func (m *model) updateComponents(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0)

	for _, c := range m.components {
		_, cmd := c.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (m *model) nextFocus() {
	m.components[m.currFocus].Blur()
	m.currFocus = (m.currFocus + 1) % len(m.components)
	m.components[m.currFocus].Focus()
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
		case key.Matches(msg, m.bindings.Tab):
			m.nextFocus()

		case key.Matches(msg, m.bindings.Quit):
			return m, tea.Quit

		}
	}

	cmds = append(cmds, m.updateComponents(msg)...)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.layout.String()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
