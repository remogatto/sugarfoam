package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Groupable interface {
	tea.Model
	LayoutItem

	Focus() tea.Cmd
	Blur()
}

type LayoutItem interface {
	View() string

	SetSize(width int, height int)
	Width() int
	Height() int
}

type Layout struct {
	width, height int
	items         []LayoutItem
}

func NewLayout(width, height int) *Layout {
	return &Layout{width: width, height: height}
}

func (l *Layout) AddItem(i LayoutItem) *Layout {
	l.items = append(l.items, i)

	return l
}

func (l *Layout) Items() []LayoutItem {
	return l.items
}

func (l *Layout) SetSize(width int, height int) {
	for _, item := range l.items {
		item.SetSize(width, height)
	}
}

func (l *Layout) View() string {
	views := []string{}

	for _, item := range l.Items() {
		views = append(views, item.View())
	}

	return lipgloss.JoinVertical(lipgloss.Top, views...)
}
