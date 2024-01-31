package main

import "github.com/charmbracelet/lipgloss"

var (
	blurTextInputStyle    = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Margin(1)
	focusedTextInputStyle = blurTextInputStyle.Copy().BorderForeground(lipgloss.Color("5"))

	tableStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	headerTableStyle = lipgloss.NewStyle().Bold(true)

	jsonViewportStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())

	infoStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("222"))
	helpStyle = lipgloss.NewStyle().Padding(1, 1)
)
