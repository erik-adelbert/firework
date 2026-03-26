package tui

import (
	"math/rand/v2"
	"time"

	"github.com/erik-adelbert/firework/fireworks"
	"github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/sym"
	"github.com/erik-adelbert/firework/internal/vec"
)

const (
	H = 88
	W = 320
)

type Model struct {
	Launcher
	*fps
	origins []vec.Vec
	screen  []sym.Symbol
	h       int
	w       int
	mode    showMode
	paused  bool
}

type fps struct {
	frames int
	last   time.Time
	cur    float64
}

// FPS returns the current frames per second.
func (f *fps) FPS() float64 {
	return f.cur
}

func (f *fps) sample(now time.Time) {
	f.frames++

	if f.last.IsZero() {
		f.last = now
	}

	dt := now.Sub(f.last)

	if dt >= time.Second {
		f.cur = float64(f.frames) / dt.Seconds()
		f.frames = 0
		f.last = now
	}
}

func NewModel(show Launcher) *Model {
	origins := make([]vec.Vec, 0, show.Len())
	for fw := range show.AllFireworks() {
		origins = append(origins, fw.Center())
	}

	return &Model{
		fps:      &fps{},
		h:        H,
		w:        W,
		mode:     modeMixed,
		origins:  origins,
		Launcher: show,
		screen:   make([]sym.Symbol, H*W),
	}
}

func (m *Model) SetDemo() {
	m.mode = modeDemo
}

func (m *Model) H() int {
	return m.h
}

func (m *Model) W() int {
	return m.w
}

func (m *Model) Size() int {
	return m.h * m.w
}

func (m *Model) step(t time.Time) {
	if m.paused {
		return
	}

	dt := 30 * time.Millisecond

	m.Launcher.Update(t, dt)
}

func (m *Model) Render() string {
	return m.Launcher.Render(m.screen, m.h, m.w)
}

func (m *Model) cycleMode() {
	m.mode = (m.mode + 1) % 15
	m.applyMode(time.Now())
}

func (m *Model) togglePause() {
	m.paused = !m.paused
}

func (m *Model) resetNow() {
	now := time.Now()
	m.Reset()
	m.Trigger(now)
}

var news = []firework.NewFirework{
	modePeony:   fireworks.NewPeony,
	modeChrys:   fireworks.NewChrysanthemum,
	modeBrocade: fireworks.NewBrocade,
	modePalm:    fireworks.NewPalm,
	modeWillow:  fireworks.NewWillow,
	modeFish:    fireworks.NewFish,
	modeGlitter: fireworks.NewGlitter,
	modeSaturn:  fireworks.NewSaturn,
	modeSphere:  fireworks.NewSphere,
	modeSun:     fireworks.NewSun,
	modeKamuro:  fireworks.NewKamuro,
	modeComet:   fireworks.NewComet,
	modeLaser:   fireworks.NewLaser,
	modeFeather: fireworks.NewFeather,
	modeDigit: func(o vec.Vec) *firework.Firework {
		d := rand.IntN(4)
		return fireworks.NewDigit(o, d)
	},
}

func (m *Model) applyMode(now time.Time) {
	fws := make([]*firework.Firework, 0, len(m.origins))

	for i, o := range m.origins {
		switch {
		case m.mode == modeMixed:
			ii := 1 + i%14
			fws = append(fws, news[ii](o))
		case m.mode < modeDemo:
			fws = append(fws, news[m.mode](o))
		}
	}

	m.Launcher.Init(fws, true)
	m.Reset()
	m.Trigger(now)
}

type showMode int

const (
	modeMixed showMode = iota
	modePeony
	modeChrys
	modeBrocade
	modePalm
	modeWillow
	modeFish
	modeGlitter
	modeSaturn
	modeSphere
	modeSun
	modeKamuro
	modeComet
	modeLaser
	modeFeather
	modeDigit
	modeDemo
)

func (m showMode) String() string {
	if m < 0 || m > modeDemo {
		return "Unknown"
	}

	return []string{
		modeMixed:   "Mixed",
		modePeony:   "Peony",
		modeChrys:   "Chrysanthemum",
		modeBrocade: "Brocade",
		modePalm:    "Palm",
		modeWillow:  "Willow",
		modeFish:    "Fish",
		modeGlitter: "Glitter",
		modeSaturn:  "Saturn Missile",
		modeSphere:  "Sphere",
		modeSun:     "Sun",
		modeKamuro:  "Kamuro",
		modeComet:   "Comet",
		modeLaser:   "Laser",
		modeFeather: "Feather",
		modeDigit:   "Digit",
		modeDemo:    "Demo",
	}[m]
}
