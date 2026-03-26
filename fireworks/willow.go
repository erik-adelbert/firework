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
	willowSpawningDuration = 170 * time.Millisecond
	willowActiveDuration   = 600 * time.Millisecond
)

func NewWillow(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./20)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(50, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(1.25, 0.22, 0, nil),
		fw.MkShell(0, NewWillowSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type willowSpawner struct {
	LUT
}

func NewWillowSpawner() *willowSpawner {
	return &willowSpawner{LUT: PickWillowLut()}
}

func (s *willowSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		strands := helper.JitterInt(24, 1./4)
		base := helper.JitterFloat(88, 3./22)

		for i := range strands {
			θ := 2 * math.Pi * (float64(i) + rand.Float64()) / float64(strands)
			cluster := 4 + rand.IntN(3)

			for range cluster {
				δ := (rand.Float64() - 0.5) * 0.03 * math.Pi
				φ := θ + δ

				dir := vec.Vec{X: math.Cos(φ), Y: -math.Sin(φ)}
				speed := base * (0.90 + 0.16*rand.Float64())
				v := dir.Scale(speed)

				ttl := time.Duration(helper.JitterInt(3300, 5./33)) * time.Millisecond
				trail := helper.JitterInt(55, 1./11)
				phases := particle.PhaseTiming{
					SpawningEnd: float64(willowSpawningDuration) / float64(ttl),
					ActiveEnd:   float64(willowActiveDuration) / float64(ttl),
				}

				p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)

				if !yield(p) {
					return
				}
			}
		}
	}
}

func (s *willowSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient2(p.Life()) * 1.12

	return scaleColor(c, g)
}
