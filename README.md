# SugarFoam

SugarFoam is a library written in Go, based on the delightful [Bubble
Tea](https://github.com/charmbracelet/bubbletea), aiming to provide a higher-level API compared to the latter. With SugarFoam, it is possible to develop TUI applications that present organized components in layouts. Compared to BubbleTea, it tries to save a bit of boilerplate code, especially if you want to develop applications that present articulated components such as tabgroup, statusbar, etc.

# Project Status

SugarFoam is in a pre-alpha development stage. APIs are subject to change and not recommended for production use. We welcome feedback and contributions for early development.

# Features

- High-level API based on [Bubble
Tea](https://github.com/charmbracelet/bubbletea) for building TUI applications.
- Idiomatic approach to building applications that preserves the architecture inspired by [Elm](https://guide.elm-lang.org/architecture/).
- Customizable layouts to easily define layouts for rendering components.

# Quickstart

To get started with SugarFoam, follow these steps:

1. **Clone the Repository**:

```bash
git clone github.com/remogatto/sugarfoam
```

2. **Explore the examples/ folder**:

```bash
cd sugarfoam/examples/tabgroup
go run .
```

# Components and layout

SugarFoam includes a small library of components that serve as containers for [Bubbles](https://github.com/charmbracelet/bubbletea) models. These components can be organized into groups, and within these groups, a layout can be defined to distribute the components in space.

# Example

```go
table1 := table.New()
viewport := viewport.New()

// Create a new group with text input, viewport, and table1 elements.
group1 := group.New(
	group.WithItems(textinput, viewport, table1),
	group.WithLayout(
		layout.New(
			layout.WithStyles(&layout.Styles{Container: lipgloss.NewStyle().Padding(1, 0, 1, 0)}),
			layout.WithItem(tiled.New(viewport, table1)),
		),
	),
)

// Create another group, this time with only table2 element.
group2 := group.New(
	group.WithItems(table2),
	group.WithLayout(
		layout.New(
			layout.WithStyles(&layout.Styles{Container: lipgloss.NewStyle().Padding(1, 0, 1, 0)}),
			layout.WithItem(table2),
		),
	),
)

// Organize the groups into a tab group with two tabs.
tabGroup := tabgroup.New(
	tabgroup.WithItems(
		tabitem.New(group1, tabitem.WithTitle("Tiled layout"), tabitem.WithActive(true)),
		tabitem.New(group2, tabitem.WithTitle("Single layout")),
	),
)
```

# License

Copyright © 2024 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

