package foam

import (
	"github.com/charmbracelet/lipgloss"
)

type Sizeable interface {
	SetSize(width int, height int)
	GetWidth() int
	GetHeight() int
}

type Placeable interface {
	Sizeable

	View() string
}

type Layout struct {
	width, height int
	items         []Placeable
}

func NewLayout(width, height int) *Layout {
	return &Layout{width: width, height: height}
}

func (l *Layout) AddItem(i Placeable) *Layout {
	l.items = append(l.items, i)

	return l
}

func (l *Layout) Items() []Placeable {
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
