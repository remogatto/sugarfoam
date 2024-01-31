package textinput

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Option func(*TextInput)

type Styles struct {
	FocusedBorder lipgloss.Style
	BlurBorder    lipgloss.Style
}

type TextInput struct {
	m      *textinput.Model
	styles *Styles
}

func New(opts ...Option) *TextInput {
	t := textinput.New()
	t.Placeholder = "Text here..."

	ti := &TextInput{
		m: &t,
	}

	for _, opt := range opts {
		opt(ti)
	}

	return ti
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

func (ti *TextInput) SetSize(msg tea.WindowSizeMsg) {
	ti.styles.FocusedBorder = ti.styles.FocusedBorder.Width(msg.Width - len(ti.m.Prompt) - ti.styles.FocusedBorder.GetHorizontalFrameSize())
}

func (ti *TextInput) Focus() tea.Cmd {
	return ti.m.Focus()
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
	return ti.styles.FocusedBorder.Render(ti.m.View())
}
