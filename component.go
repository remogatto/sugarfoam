package foam

// Package foam provides interfaces and types for managing focus, grouping, and
// tabbing behavior in a user interface.
//
// The package is designed to work with the Bubble Tea framework and Sugarfoam layout
// management system, providing a structured way to handle UI interactions such as
// focusing on elements, grouping them, and managing tab focus.
import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/remogatto/sugarfoam/layout"
)

// Focusable is an interface that extends the tea.Model and layout.Placeable
// interfaces, adding methods for managing focus state.
//
// Implementations of Focusable are expected to define how an element can be
// focused and blurred.
type Focusable interface {
	tea.Model
	layout.Placeable

	// Focus returns a command that focuses the element.
	Focus() tea.Cmd

	// Blur removes focus from the element.
	Blur()
}

// Groupable extends the Focusable interface, adding methods for
// managing a group of focusable elements.
//
// Implementations of Groupable are expected to define the current
// focused element within the group.
type Groupable interface {
	Focusable

	// Current returns the currently focused element within the group.
	Current() Focusable
}

// Tabbable extends the Groupable interface, adding methods for
// managing tab navigation and active state.
//
// Implementations of Tabbable are expected to define a title for the
// tab, and whether the tab is currently active.
type Tabbable interface {
	Groupable

	// Title returns the title of the tab.
	Title() string

	// Active returns true if the tab is currently active.
	Active() bool
}
