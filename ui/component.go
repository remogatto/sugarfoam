package ui

import tea "github.com/charmbracelet/bubbletea"

type Component interface {
	tea.Model

	SetSize(tea.WindowSizeMsg)
	Focus() tea.Cmd
}
