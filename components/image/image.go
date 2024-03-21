package image

import (
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
	foam "github.com/remogatto/sugarfoam"
)

type DoneMsg struct {
	done bool
}

type errMsg struct{ error }

type loadMsg struct {
	io.ReadCloser
}

type redrawMsg struct {
	width  uint
	height uint
	url    string
}

type Option func(*Model)

type Model struct {
	foam.Common

	textImage string
	image     image.Image
	url       string
	focused   bool
}

func New(opts ...Option) *Model {
	img := new(Model)

	img.Common.SetStyles(foam.DefaultStyles())

	for _, opt := range opts {
		opt(img)
	}

	return img
}

func WithStyles(styles *foam.Styles) Option {
	return func(ti *Model) {
		ti.Common.SetStyles(styles)
	}
}

// Blur removes focus from the viewport.
func (m *Model) Blur() {
	m.focused = false
}

// Focus sets the viewport to be focused.
func (m *Model) Focus() tea.Cmd {
	m.focused = true

	return nil
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) SetWidth(w int) {
	m.Common.SetHeight(w)

	if m.image != nil {
		m.imageToString()
	}
}

func (m *Model) SetHeight(h int) {
	m.Common.SetHeight(h)

	if m.image != nil {
		m.imageToString()
	}
}

func (m *Model) SetSize(w, h int) {
	m.Common.SetSize(w, h)

	if m.image != nil {
		m.imageToString()
	}
}

func (t *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		// m.err = msg
		return m, nil
	case redrawMsg:
		return m, m.LoadURL(m.url)

	case loadMsg:
		return m.handleLoadMsg(msg)

	}
	return m, nil
}

func (m *Model) View() string {
	if m.Focused() {
		return m.GetStyles().Focused.Render(m.textImage)
	}
	return m.GetStyles().Blurred.Render(m.textImage)

}

func (m *Model) Redraw() tea.Cmd {
	return func() tea.Msg {
		return redrawMsg{
			width:  uint(m.GetWidth()),
			height: uint(m.GetHeight()),
			url:    m.url,
		}
	}
}

func (m *Model) CanGrow() bool {
	return true
}

func (m *Model) LoadURL(url string) tea.Cmd {
	var r io.ReadCloser
	var err error

	if strings.HasPrefix(m.url, "http") {
		var resp *http.Response
		resp, err = http.Get(m.url)
		r = resp.Body
	} else {
		r, err = os.Open(m.url)
	}

	if err != nil {
		return func() tea.Msg {
			return errMsg{err}
		}
	}

	return load(r)
}

func (m *Model) SetURL(url string) {
	m.url = url
}

func load(r io.ReadCloser) tea.Cmd {
	return func() tea.Msg {
		return loadMsg{r}
	}
}

func (m *Model) handleLoadMsg(msg loadMsg) (*Model, tea.Cmd) {
	// blank out image so it says "loading..."
	m.textImage = ""

	return m.handleLoadMsgStatic(msg)
}

func (m *Model) handleLoadMsgStatic(msg loadMsg) (*Model, tea.Cmd) {
	defer msg.Close()

	img, err := m.readerToimage(msg)
	if err != nil {
		return m, func() tea.Msg { return errMsg{err} }
	}

	m.textImage = img

	return m, func() tea.Msg { return DoneMsg{true} }
}

func (m *Model) imageToString() (string, error) {
	m.image = resize.Thumbnail(uint(m.GetWidth()), uint(m.GetHeight()*2-4), m.image, resize.Lanczos3)
	b := m.image.Bounds()
	w := b.Max.X
	h := b.Max.Y
	p := termenv.ColorProfile()
	str := strings.Builder{}
	for y := 0; y < h; y += 2 {
		for x := w; x < int(m.GetWidth()); x = x + 2 {
			str.WriteString(" ")
		}
		for x := 0; x < w; x++ {
			c1, _ := colorful.MakeColor(m.image.At(x, y))
			color1 := p.Color(c1.Hex())
			c2, _ := colorful.MakeColor(m.image.At(x, y+1))
			color2 := p.Color(c2.Hex())
			str.WriteString(termenv.String("▀").
				Foreground(color1).
				Background(color2).
				String())
		}
		str.WriteString("\n")
	}
	return str.String(), nil
}

func (m *Model) readerToimage(r io.Reader) (string, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return "", err
	}

	m.image = img

	return m.imageToString()
}
