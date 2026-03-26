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

func NewKamuro(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.7, 0.18, 0, nil),
		fw.MkShell(0, NewKamuroSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type kamuroSpawner struct {
	LUT
}

func NewKamuroSpawner() *kamuroSpawner {
	return &kamuroSpawner{LUT: PickKamuroLut()}
}

func (s *kamuroSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		n := helper.JitterInt(135, 24./27)

		for v := range AllNormalDisk9(350, n) {
			ttl := time.Duration(helper.JitterInt(4250, 3./17)) * time.Millisecond
			trail := helper.JitterInt(33, 23./33)

			p := particle.New(o, v, vec.Vec{}, s.PickColor(), ttl, trail)

			if !yield(p) {
				return
			}
		}
	}
}

func (s *kamuroSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}
