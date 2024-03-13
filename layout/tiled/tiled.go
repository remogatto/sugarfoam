package tiled

// Package tiled provides utilities for creating and managing
// horizontal tiles in a UI.  It leverages the Lipgloss library for
// styling and the Sugarfoam layout management system for positioning
// and sizing components. The package is designed to facilitate the
// creation of horizontally aligned layouts with a focus on simplicity
// and ease of use.
import (
	"github.com/charmbracelet/lipgloss"
	"github.com/remogatto/sugarfoam/layout"
)

// DefaultWidth and DefaultHeight define the default dimensions for a
// horizontal tiled layout.
var (
	DefaultWidth  = 80
	DefaultHeight = 25
)

// HorizontalTile represents a container for horizontally aligned
// placeable items.
type HorizontalTile struct {
	width, height int
	items         []layout.Placeable
}

// View returns the string representation of the horizontal tile, with
// all its items rendered horizontally. The width of each item is
// evenly distributed, with the last item potentially taking up the
// remaining space if the total width is not evenly divisible by the
// number of items.
func (ht *HorizontalTile) View() string {
	strs := make([]string, 0)

	for _, item := range ht.items {
		strs = append(strs, item.View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, strs...)
}

// SetSize sets the dimensions of the horizontal tile.
func (ht *HorizontalTile) SetSize(width int, height int) {
	ht.width = width
	ht.height = height

	w := ht.width / len(ht.items)
	dw := ht.width - w*len(ht.items)

	for i, item := range ht.items {
		if i == len(ht.items)-1 {
			item.SetSize(w+dw, height)
		} else {
			item.SetSize(w, height)
		}
	}

}

// GetWidth returns the current width of the horizontal tile.
func (ht *HorizontalTile) GetWidth() int { return ht.width }

// GetHeight returns the current height of the horizontal tile.
func (ht *HorizontalTile) GetHeight() int { return ht.height }

func (ht *HorizontalTile) SetWidth(width int) { ht.width = width }

func (ht *HorizontalTile) SetHeight(height int) { ht.height = height }

func (ht *HorizontalTile) CanGrow() bool {
	return true
}

// New creates a new HorizontalTile with the specified items, using
// the default dimensions defined by DefaultWidth and DefaultHeight.
func New(items ...layout.Placeable) *HorizontalTile {
	return &HorizontalTile{DefaultWidth, DefaultHeight, items}
}
