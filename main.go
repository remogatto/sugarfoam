package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/remogatto/bubbletea-app/ui"
	"github.com/remogatto/bubbletea-app/ui/components/group"
	"github.com/remogatto/bubbletea-app/ui/components/tabgroup"
	"github.com/remogatto/bubbletea-app/ui/components/table"
	"github.com/remogatto/bubbletea-app/ui/components/textinput"
	"github.com/remogatto/bubbletea-app/ui/layout/tiled"

	tea "github.com/charmbracelet/bubbletea"
)

type keyBindings struct {
	Quit key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyBindings) ShortHelp() []key.Binding {
	return []key.Binding{
		k.Quit,
	}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.Quit,
		},
	}
}

func newBindings() keyBindings {
	return keyBindings{
		key.NewBinding(
			key.WithKeys("esc"), key.WithHelp("esc", "Quit app"),
		),
	}
}

type model struct {
	tabgroup *tabgroup.TabGroup
	bindings keyBindings
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

	group := group.New(
		group.WithItems(textinput, table1, table2),
		group.WithLayout(ui.NewLayout(80, 25).AddItem(textinput).AddItem(tiled.New(80, 25, table1, table2))),
	)

	tabgroup := tabgroup.New().
		AddItem(&tabgroup.TabItem{Title: "Tab1", Groupable: group}).
		AddItem(&tabgroup.TabItem{Title: "Tab2", Groupable: table3})

	return model{
		tabgroup: tabgroup,
		bindings: newBindings(),
	}
}

func (m *model) setSize(msg tea.WindowSizeMsg) {
	m.tabgroup.SetSize(msg.Width, msg.Height)
}

func (m model) Init() tea.Cmd {
	var cmds []tea.Cmd

	cmds = append(cmds, m.tabgroup.Init())

	m.tabgroup.Focus()

	return tea.Batch(cmds...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.setSize(msg)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.bindings.Quit):
			return m, tea.Quit

		}
	}

	_, cmd := m.tabgroup.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return m.tabgroup.View()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
