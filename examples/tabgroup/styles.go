package main

import (
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

var (
	blurredTextInputStyle = foam.DefaultStyles().Blurred.Margin(1, 1, 0, 1)
	focusedTextInputStyle = foam.DefaultStyles().Focused.Margin(1, 1, 0, 1)

	blurredTableStyle = foam.DefaultStyles().Blurred.Margin(1)
	focusedTableStyle = foam.DefaultStyles().Focused.Margin(1)

	blurredViewportStyle = foam.DefaultStyles().Blurred.Margin(1)
	focusedViewportStyle = foam.DefaultStyles().Focused.Margin(1)

	headerTableStyle = lipgloss.NewStyle().Bold(true)

	infoStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("222"))
	helpStyle = lipgloss.NewStyle().Padding(1, 1)
)
