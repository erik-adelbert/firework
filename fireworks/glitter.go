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
	glitterSpawningDuration = 90 * time.Millisecond
	glitterActiveDuration   = 240 * time.Millisecond
)

func NewGlitter(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.85, 0.30, 0, nil),
		fw.MkShell(0, NewGlitterSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type glitterSpawner struct {
	LUT
}

func NewGlitterSpawner() *glitterSpawner {
	return &glitterSpawner{LUT: PickGlitterLut()}
}

func (s *glitterSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		n := helper.JitterInt(110, 3./11)

		for dir := range AllStratifiedCircle(1, n) {
			speed := helper.JitterFloat(107, 15./107)
			v := dir.Scale(speed)

			ttl := time.Duration(helper.JitterInt(950, 5./19)) * time.Millisecond
			phases := particle.PhaseTiming{
				SpawningEnd: float64(glitterSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(glitterActiveDuration) / float64(ttl),
			}

			p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, 1, phases)

			// trail length 1 keeps only the current point for point-like glitter sparks.
			if !yield(p) {
				return
			}
		}
	}
}

func (s *glitterSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := helper.JitterFloat(PhaseGradient5(p.Life()), .21)

	return scaleColor(c, g)
}
