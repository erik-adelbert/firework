package tui

import (
	"iter"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/sym"
)

type Launcher interface {
	Init(state0 []*fw.Firework, looping bool)
	Len() int
	Reset()
	Trigger(now time.Time)
	Update(now time.Time, dt time.Duration)
	Render(screen []sym.Symbol, h, w int) string
	AllFireworks() iter.Seq[*fw.Firework]
}
