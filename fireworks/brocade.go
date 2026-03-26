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

const (
	brocadeSpawningDuration = 220 * time.Millisecond
	brocadeActiveDuration   = 700 * time.Millisecond
)

func NewBrocade(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.56, 0.20, 0, nil),
		fw.MkShell(0, NewBrocadeSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type brocadeSpawner struct {
	LUT
}

func NewBrocadeSpawner() *brocadeSpawner {
	return &brocadeSpawner{LUT: PickLut()}
}

func (s *brocadeSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		n := helper.JitterInt(24, 1./3)
		base := helper.JitterFloat(100, 1./17)

		for dir := range AllStratifiedCircle(1, n) {
			speed := base * helper.JitterFloat(0.86, 0.25)
			v := dir.Scale(speed)

			ttl := time.Duration(helper.JitterInt(3000, .2)) * time.Millisecond
			size := helper.JitterInt(36, 1./6)
			phases := particle.PhaseTiming{
				SpawningEnd: float64(brocadeSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(brocadeActiveDuration) / float64(ttl),
			}

			p := particle.NewWithPhaseTiming(
				o, v, vec.Vec{},
				s.PickColor(),
				ttl,
				size,
				phases,
			)

			if !yield(p) {
				return
			}
		}
	}
}

func (s *brocadeSpawner) Gradient(p *particle.Particle) color.RGBA {
	g := PhaseGradient2(p.Life())
	c := p.Color()

	return scaleColor(c, g)
}
