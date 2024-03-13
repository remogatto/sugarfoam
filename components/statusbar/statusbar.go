package statusbar

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	foam "github.com/remogatto/sugarfoam"
)

type Option func(*Model)

func WithStyles(styles *foam.Styles) Option {
	return func(sb *Model) {
		sb.Common.SetStyles(styles)
	}
}

func WithContent(left, center, right string) Option {
	return func(sb *Model) {
		sb.left, sb.center, sb.right = left, center, right
	}
}

type Model struct {
	foam.Common

	left, center, right string
	bindings            help.KeyMap
}

func New(bindings help.KeyMap, opts ...Option) *Model {
	m := &Model{}

	m.Common.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) SetContent(left, center, right string) {
	m.left, m.center, m.right = left, center, right
}

func (m *Model) View() string {
	doc := strings.Builder{}

	w := lipgloss.Width

	frameSizes := leftStyle.GetHorizontalFrameSize() +
		centerStyle.GetHorizontalFrameSize() +
		rightStyle.GetHorizontalBorderSize() +
		m.GetStyles().NoBorder.GetHorizontalFrameSize()

	statusKey := statusStyle.Render(m.left)
	right := rightStyle.Render(m.right)
	statusVal := centerStyle.Copy().
		Width(m.GetWidth() - w(m.left) - w(m.right) - frameSizes).
		Render(m.center)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		right,
	)

	doc.WriteString(centerStyle.Render(bar))

	return m.GetStyles().NoBorder.Render(doc.String())
}

func (m *Model) CanGrow() bool {
	return false
}

func (m *Model) GetHeight() int {
	return lipgloss.Height(m.View())
}
