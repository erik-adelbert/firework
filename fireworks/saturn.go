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
	saturnSpawningDuration = 120 * time.Millisecond
	saturnActiveDuration   = 380 * time.Millisecond
)

func NewSaturn(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.92, 0.18, 0, nil),
		fw.MkShell(0, NewSaturnSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type saturnSpawner struct {
	LUT
}

func NewSaturnSpawner() *saturnSpawner {
	return &saturnSpawner{LUT: PickSaturnLut()}
}

func (s *saturnSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		head := 2 * math.Pi * rand.Float64()

		forward := vec.Vec{
			X: math.Cos(head), Y: -math.Sin(head),
		}

		side := vec.Vec{
			X: -forward.Y, Y: forward.X,
		}

		coreCount := helper.JitterInt(23, 5./23)
		for range coreCount {
			spread := (rand.Float64() - 0.5) * 0.10 * math.Pi
			speed := helper.JitterFloat(162, 2./27)

			v := rotate(forward, spread).Scale(speed)

			ttl := time.Duration(helper.JitterInt(1500, 2./15)) * time.Millisecond
			trail := helper.JitterInt(29, 5./29)
			phases := particle.PhaseTiming{
				SpawningEnd: float64(saturnSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(saturnActiveDuration) / float64(ttl),
			}

			if !yield(particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)) {
				return
			}
		}

		ringCount := helper.JitterInt(27, 5./27)
		for i := range ringCount {
			θ := 2 * math.Pi * (float64(i) + rand.Float64()) / float64(ringCount)

			x := 0.92 * math.Cos(θ)
			y := 0.52 * math.Sin(θ)

			// Build a tilted local ellipse: X along side axis, Y along head axis.
			dir := side.Scale(x).Add(forward.Scale(y)).Normalize()
			if dir.Length2() == 0 {
				continue
			}

			speed := rand.Float64()*22 + 78
			v := dir.Scale(speed)

			ttl := time.Duration(helper.JitterInt(2400, 5./24)) * time.Millisecond
			trail := helper.JitterInt(35, 1./7)
			phases := particle.PhaseTiming{
				SpawningEnd: float64(saturnSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(saturnActiveDuration) / float64(ttl),
			}

			p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)

			if !yield(p) {
				return
			}
		}
	}
}

func (s *saturnSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient5(p.Life())

	return scaleColor(c, g)
}
