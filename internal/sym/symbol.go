package sym

import (
	"image/color"
	"sync"

	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

type Symbol struct {
	rune
	color.RGBA
}

func MkSymbol(r rune, c color.RGBA) Symbol {
	return Symbol{rune: r, RGBA: c}
}

func (s Symbol) IsZero() bool {
	return s.rune == 0
}

func (s Symbol) String() string {
	if s.rune == 0 {
		return "  " // two spaces to preserve cell width for empty symbols
	}

	str, ok := renderCache[s]
	if !ok {
		// If the rune is narrow, add a space to make it take up the full cell width.
		text := string(s.rune)
		if runewidth.RuneWidth(s.rune) < 2 {
			text += " "
		}

		style := getStyle()
		str = style.Foreground(s.RGBA).Render(text)
		putStyle(style)

		renderCache[s] = str
	}

	return str
}

var renderCache = make(map[Symbol]string)

var stylePool = sync.Pool{
	New: func() any {
		st := lipgloss.NewStyle()
		return &st
	},
}

func getStyle() *lipgloss.Style {
	return stylePool.Get().(*lipgloss.Style)
}

func putStyle(s *lipgloss.Style) {
	stylePool.Put(s)
}
