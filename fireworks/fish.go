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
	fishSpawningDuration = 140 * time.Millisecond
	fishActiveDuration   = 480 * time.Millisecond
)

func NewFish(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.72, 0.18, 0, nil),
		fw.MkShell(0, NewFishSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type fishSpawner struct {
	LUT
}

func NewFishSpawner() *fishSpawner {
	return &fishSpawner{LUT: PickFishLut()}
}

func (s *fishSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		heading := 2 * math.Pi * rand.Float64()

		bodyCount := helper.JitterInt(31, 13./31)
		bodySpeed := helper.JitterFloat(105, .1)
		for range bodyCount {
			// Local fish body: tapered ellipse shifted forward on heading axis.
			x := (rand.Float64()*2 - 1) * 1.10
			y := (rand.Float64()*2 - 1) * 0.62

			if (x*x)/(1.10*1.10)+(y*y)/(0.62*0.62) > 1 {
				continue
			}

			x += 0.42
			local := vec.Vec{X: x, Y: y}
			dir := rotate(local, heading).Normalize()
			if dir.Length2() == 0 {
				continue
			}

			speed := helper.JitterFloat(bodySpeed, .1)
			v := dir.Scale(speed)

			ttl := time.Duration(helper.JitterInt(2700, 5./27)) * time.Millisecond
			trail := helper.JitterInt(32, 1./8)
			phases := particle.PhaseTiming{
				SpawningEnd: float64(fishSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(fishActiveDuration) / float64(ttl),
			}

			p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)

			if !yield(p) {
				return
			}
		}

		tailCount := helper.JitterInt(14, 2./7)
		tailBase := heading + math.Pi
		tailSpeed := helper.JitterFloat(80, .1)
		for range tailCount {
			δ := (rand.Float64() - 0.5) * 0.42 * math.Pi
			φ := tailBase + δ
			dir := vec.Vec{X: math.Cos(φ), Y: math.Sin(φ)}

			speed := helper.JitterFloat(tailSpeed, 0.22)
			v := dir.Scale(speed)

			ttl := time.Duration(helper.JitterInt(2450, 9./49)) * time.Millisecond
			trail := helper.JitterInt(28, 1./7)
			phases := particle.PhaseTiming{
				SpawningEnd: float64(fishSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(fishActiveDuration) / float64(ttl),
			}

			if !yield(particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)) {
				return
			}
		}
	}
}

func (s *fishSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}

func rotate(v vec.Vec, θ float64) vec.Vec {
	c := math.Cos(θ)
	s := math.Sin(θ)

	return vec.Vec{
		X: v.X*c - v.Y*s,
		Y: v.X*s + v.Y*c,
	}
}
