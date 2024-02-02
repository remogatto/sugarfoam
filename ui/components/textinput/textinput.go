package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Option func(*TextInput)

type Styles struct {
	FocusedBorder lipgloss.Style
	BlurredBorder lipgloss.Style
}

type TextInput struct {
	m *textinput.Model

	width, height int

	styles *Styles
}

func New(opts ...Option) *TextInput {
	t := textinput.New()
	t.Placeholder = "Text here..."

	ti := &TextInput{
		m: &t,
	}

	ti.styles = DefaultStyles()

	for _, opt := range opts {
		opt(ti)
	}

	return ti
}

func DefaultStyles() *Styles {
	return &Styles{
		FocusedBorder: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")),
		BlurredBorder: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()),
	}
}

func WithPlaceholder(placeholder string) Option {
	return func(ti *TextInput) {
		ti.m.Placeholder = placeholder
	}
}

func WithStyles(styles *Styles) Option {
	return func(ti *TextInput) {
		ti.styles = styles
	}
}

func (ti *TextInput) SetSize(width int, height int) {

	ti.width = width - ti.styles.FocusedBorder.GetHorizontalFrameSize()
	ti.height = height

	ti.styles.FocusedBorder = ti.styles.FocusedBorder.Width(ti.width)
	ti.styles.BlurredBorder = ti.styles.BlurredBorder.Width(ti.width)

}

func (ti *TextInput) Width() int  { return ti.width }
func (ti *TextInput) Height() int { return ti.height }

func (ti *TextInput) Focus() tea.Cmd {
	return ti.m.Focus()
}

func (ti *TextInput) Blur() tea.Cmd {
	ti.m.Blur()

	return nil
}

func (ti *TextInput) Init() tea.Cmd {
	return textinput.Blink
}

func (ti *TextInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	t, cmd := ti.m.Update(msg)
	ti.m = &t

	return ti, cmd
}

func (ti *TextInput) View() string {
	if ti.m.Focused() {
		return ti.styles.FocusedBorder.Render(ti.m.View())
	}
	return ti.styles.BlurredBorder.Render(ti.m.View())
}

func (ti *TextInput) String() string { return ti.View() }
