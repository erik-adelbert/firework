package fireworks

import (
	"image/color"
	"iter"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

var cometColors = []color.RGBA{
	{R: 0xF2, G: 0xE9, B: 0xBE, A: 0xFF}, // Cream
	{R: 0xE2, G: 0xC4, B: 0x88, A: 0xFF}, // Tan
	{R: 0xFF, G: 0xF8, B: 0xFD, A: 0xFF}, // Near-white
}

func NewComet(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.3, 0.18, 0, nil),
		fw.MkShell(0, NewCometSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type cometSpawner struct{}

func NewCometSpawner() *cometSpawner {
	return &cometSpawner{}
}

func (s *cometSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		n := helper.JitterInt(26, 11./13)

		for v := range AllNormalDisk9(350, n) {
			c, _ := helper.Pick(cometColors)
			ttl := time.Duration(helper.JitterInt(4250, 1./5)) * time.Millisecond
			trail := helper.JitterInt(35, 3./7)

			if !yield(particle.New(o, v, vec.Vec{}, c, ttl, trail)) {
				return
			}
		}
	}
}

func (s *cometSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}
