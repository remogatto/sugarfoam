package tiled

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/sugarfoam/layout"
)

var (
	DefaultWidth  = 80
	DefaultHeight = 25
)

type HorizontalTile struct {
	width, height int
	items         []layout.Placeable
}

func (ht *HorizontalTile) View() string {
	w := ht.width / len(ht.items)
	dw := ht.width - w*len(ht.items)

	strs := make([]string, 0)

	for i, item := range ht.items {
		if i == len(ht.items)-1 {
			item.SetSize(w+dw, item.GetHeight())
		} else {
			item.SetSize(w, item.GetHeight())
		}
		strs = append(strs, item.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

func (ht *HorizontalTile) SetSize(width int, height int) {
	ht.width = width
	ht.height = height
}

func (ht *HorizontalTile) GetWidth() int  { return ht.width }
func (ht *HorizontalTile) GetHeight() int { return ht.height }

func New(items ...layout.Placeable) *HorizontalTile {
	return &HorizontalTile{DefaultWidth, DefaultHeight, items}
}
