package launcher

import (
	"image/color"
	"iter"
	"testing"
	"time"

	"github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/sym"
	"github.com/erik-adelbert/firework/internal/vec"
)

type benchmarkSpawner struct {
	c color.RGBA
}

func (s benchmarkSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return helper.EmptySeq[*particle.Particle]()
}

func (s benchmarkSpawner) Gradient(p *particle.Particle) color.RGBA {
	return s.c
}

func (s benchmarkSpawner) Reset() {}

func benchmarkLauncher(particleCount, trailLen, h, w int) *Launcher {
	spawner := benchmarkSpawner{c: color.RGBA{R: 255, G: 180, B: 32, A: 255}}
	forces := firework.MkForces(0.2, 0.01, 40, nil)
	shell := firework.MkShell(time.Millisecond, spawner)
	timing := firework.MkTiming(5*time.Second, 0)

	fw := firework.New(vec.Vec{X: 0, Y: 0}, forces, shell, timing)
	fw.State = firework.Spawning
	fw.Particles = make([]*particle.Particle, 0, particleCount)

	for i := range particleCount {
		p := particle.New(
			vec.Vec{X: float64(i % w), Y: float64(i % h)},
			vec.Vec{X: 2 + 0.1*float64(i%7), Y: 1.5 + 0.1*float64(i%5)},
			vec.Vec{},
			color.RGBA{R: 255, G: 160, B: 32, A: 255},
			10*time.Second,
			trailLen,
		)

		for range trailLen {
			p.Update(30*time.Millisecond, 0.1, 0.1, 10, nil)
		}

		fw.Particles = append(fw.Particles, p)
	}

	launcher := new(Launcher)
	launcher.Init([]*firework.Firework{fw}, true)

	return launcher
}

func benchmarkRender(b *testing.B, h, w, particleCount, trailLen int) {
	// renderCache = make(map[sym.Symbol]string)

	l := benchmarkLauncher(particleCount, trailLen, h, w)
	screen := make([]sym.Symbol, h*w)

	for b.Loop() {
		_ = l.Render(screen, h, w)
	}
}

func BenchmarkRender(b *testing.B) {
	cases := []struct {
		name          string
		h, w          int
		particleCount int
		trailLen      int
	}{
		{name: "22x80-120-24", h: 22, w: 80, particleCount: 120, trailLen: 24},
		{name: "22x80-240-24", h: 22, w: 80, particleCount: 240, trailLen: 24},
		{name: "22x80-480-24", h: 22, w: 80, particleCount: 480, trailLen: 24},
		{name: "22x80-120-48", h: 22, w: 80, particleCount: 120, trailLen: 48},
		{name: "22x80-120-96", h: 22, w: 80, particleCount: 120, trailLen: 96},
		{name: "44x160-120-24", h: 44, w: 160, particleCount: 120, trailLen: 24},
		{name: "44x160-220-30", h: 44, w: 160, particleCount: 220, trailLen: 30},
		{name: "88x320-400-40", h: 88, w: 320, particleCount: 400, trailLen: 40},
		{name: "88x320-800-33", h: 88, w: 320, particleCount: 800, trailLen: 33},
		{name: "88x320-1200-30", h: 88, w: 320, particleCount: 1200, trailLen: 30},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			benchmarkRender(b, tc.h, tc.w, tc.particleCount, tc.trailLen)
		})
	}
}

func TestRender(t *testing.T) {
	h, w := 88, 320

	l := benchmarkLauncher(1200, 30, h, w)
	screen := make([]sym.Symbol, h*w)

	l.Update(time.Now(), 30*time.Millisecond)

	s := l.Render(screen, h, w)
	t.Log(s)
}
