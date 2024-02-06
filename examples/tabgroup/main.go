package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	foam "github.com/remogatto/sugarfoam"
	"github.com/remogatto/sugarfoam/components/group"
	"github.com/remogatto/sugarfoam/components/tabgroup"
	"github.com/remogatto/sugarfoam/components/table"
	"github.com/remogatto/sugarfoam/components/textinput"
	"github.com/remogatto/sugarfoam/layout/tiled"

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
		textinput.WithStyles(&foam.Styles{
			Blurred: blurredTextInputStyle,
			Focused: focusedTextInputStyle,
		}),
		textinput.WithPlaceholder("Text input goes here..."),
	)

	table1 := table.New(
		table.WithStyles(&foam.Styles{
			Blurred: blurredTableStyle,
			Focused: focusedTableStyle,
		}),
	)

	table2 := table.New(
		table.WithStyles(&foam.Styles{
			Blurred: blurredTableStyle,
			Focused: focusedTableStyle,
		}),
	)

	table3 := table.New(
		table.WithStyles(&foam.Styles{
			Blurred: blurredTableStyle,
			Focused: focusedTableStyle,
		}),
	)

	group1 := group.New(
		group.WithItems(textinput, table1, table2),
		group.WithLayout(foam.NewLayout(80, 25).AddItem(textinput).AddItem(tiled.New(80, 25, table1, table2))),
	)

	group2 := group.New(
		group.WithItems(table3),
		group.WithLayout(foam.NewLayout(80, 25).AddItem(table3)),
	)

	tabgroup := tabgroup.New().
		AddItem(&tabgroup.TabItem{Title: "Group 1", Group: group1}).
		AddItem(&tabgroup.TabItem{Title: "Group 2", Group: group2})

	return model{
		tabgroup: tabgroup,
		bindings: newBindings(),
	}
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
		m.tabgroup.SetSize(msg.Width, msg.Height)

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
