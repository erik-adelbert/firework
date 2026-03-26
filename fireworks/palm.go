package fireworks

import (
	"image/color"
	"iter"
	"math"
	"math/rand/v2"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

const (
	palmSpawningDuration = 160 * time.Millisecond
	palmActiveDuration   = 500 * time.Millisecond
)

func NewPalm(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(1.15, 0.10, 0, nil),
		fw.MkShell(0, NewPalmSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type palmSpawner struct {
	LUT
}

func NewPalmSpawner() *palmSpawner {
	return &palmSpawner{LUT: PickPalmLut()}
}

func (s *palmSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		fronds := helper.JitterInt(12, 1./12)
		base := helper.JitterFloat(95, 16./19)

		for i := range fronds {
			θ := 2 * math.Pi * (float64(i) + rand.Float64()) / float64(fronds)
			clusters := 4 + rand.IntN(3)

			for range clusters {
				δ := (rand.Float64() - 0.5) * 0.035 * math.Pi
				φ := θ + δ

				dir := vec.Vec{X: math.Cos(φ), Y: -math.Sin(φ)}
				speed := base * (0.92 + 0.12*rand.Float64())
				v := dir.Scale(speed)

				ttl := time.Duration(helper.JitterInt(3600, 1./6)) * time.Millisecond
				trail := helper.JitterInt(61, 11./61)

				phases := particle.PhaseTiming{
					SpawningEnd: float64(palmSpawningDuration) / float64(ttl),
					ActiveEnd:   float64(palmActiveDuration) / float64(ttl),
				}

				p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)

				if !yield(p) {
					return
				}
			}
		}
	}
}

func (s *palmSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}
