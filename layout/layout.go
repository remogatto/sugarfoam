package layout

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	DefaultWidth  = 80
	DefaultHeight = 25
)

type Option func(*Layout)

func WithItem(i Placeable) Option {
	return func(l *Layout) {
		l.items = append(l.items, i)
	}
}

func WithStyles(styles *Styles) Option {
	return func(l *Layout) {
		l.styles = styles
	}
}

type Sizeable interface {
	SetSize(width int, height int)
	GetWidth() int
	GetHeight() int
}

type Placeable interface {
	Sizeable

	View() string
}

type Styles struct {
	Container lipgloss.Style
}

type Layout struct {
	width, height int
	items         []Placeable
	styles        *Styles
}

func New(opts ...Option) *Layout {
	l := &Layout{width: DefaultWidth, height: DefaultHeight}

	for _, opt := range opts {
		opt(l)
	}

	return l
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
		item.SetSize(width-l.styles.Container.GetHorizontalFrameSize(), height-l.styles.Container.GetVerticalFrameSize())
	}
}

func (l *Layout) View() string {
	views := []string{}

	for _, item := range l.Items() {
		views = append(views, item.View())
	}

	return l.styles.Container.Render(lipgloss.JoinVertical(lipgloss.Top, views...))
}
