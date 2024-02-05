package main

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/bubbletea-app/ui/components/table"
	"github.com/remogatto/bubbletea-app/ui/components/textinput"
)

var (
	blurredTextInputStyle = textinput.DefaultStyles().BlurredBorder.Margin(1, 1, 0, 1)
	focusedTextInputStyle = textinput.DefaultStyles().FocusedBorder.Margin(1, 1, 0, 1)

	blurredTableStyle = table.DefaultStyles().BlurredBorder.Margin(1)
	focusedTableStyle = table.DefaultStyles().FocusedBorder.Margin(1)

	headerTableStyle = lipgloss.NewStyle().Bold(true)

	jsonViewportStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())

	infoStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("222"))
	helpStyle = lipgloss.NewStyle().Padding(1, 1)
)
