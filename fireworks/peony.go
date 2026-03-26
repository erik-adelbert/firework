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

func NewPeony(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.65, 0.28, 0, nil),
		fw.MkShell(0, NewPeonySpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type peonySpawner struct {
	LUT
}

func NewPeonySpawner() *peonySpawner {
	return &peonySpawner{
		LUT: PickLut(),
	}
}

func (s *peonySpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		ρ := helper.JitterFloat(245, 1./5)
		n := helper.JitterInt(40, 1./2)

		for v := range AllNormalDisk9Balanced(ρ, n) {
			ttl := time.Duration(helper.JitterInt(2000, .1)) * time.Millisecond

			p := particle.New(o, v, vec.Vec{}, s.PickColor(), ttl, helper.JitterInt(27, 2./27))

			if !yield(p) {
				return
			}
		}

	}
}

func (s *peonySpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}
