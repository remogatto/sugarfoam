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

func WithContent(text string) Option {
	return func(sb *Model) {
		sb.content = text
	}
}

type Model struct {
	foam.Common

	content  string
	bindings help.KeyMap
}

func New(bindings help.KeyMap, opts ...Option) *Model {
	m := &Model{}

	m.Common.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) SetContent(text string) {
	m.content = text
}

func (m *Model) View() string {
	doc := strings.Builder{}
	w := lipgloss.Width

	statusKey := statusStyle.Render("STATUS")
	encoding := encodingStyle.Render("UTF-8")
	fishCake := fishCakeStyle.Render("🍥 Fish Cake")
	statusVal := statusText.Copy().
		Width(m.GetWidth() - w(statusKey) - w(encoding) - w(fishCake) - m.GetStyles().NoBorder.GetHorizontalFrameSize()).
		Render(m.content)

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	)

	doc.WriteString(statusBarStyle.Render(bar))
	return m.GetStyles().NoBorder.Render(doc.String())
}
