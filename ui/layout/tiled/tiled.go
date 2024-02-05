package tiled

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/bubbletea-app/ui"
)

type HorizontalTile struct {
	width, height int
	items         []ui.LayoutItem
}

func (ht *HorizontalTile) View() string {
	w := ht.width / len(ht.items)
	dw := ht.width - w*len(ht.items)

	strs := make([]string, 0)

	for i, item := range ht.items {
		if i == len(ht.items)-1 {
			item.SetSize(w+dw, item.Height())
		} else {
			item.SetSize(w, item.Height())
		}
		strs = append(strs, item.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

func (ht *HorizontalTile) SetSize(width int, height int) {
	ht.width = width
	ht.height = height
}

func (ht *HorizontalTile) Width() int  { return ht.width }
func (ht *HorizontalTile) Height() int { return ht.height }

func New(width int, height int, items ...ui.LayoutItem) *HorizontalTile {
	return &HorizontalTile{width, height, items}
}
