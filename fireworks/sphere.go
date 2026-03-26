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
	sphereSpawningDuration = 220 * time.Millisecond
	sphereActiveDuration   = 900 * time.Millisecond
	sphereShellStep        = 20 * time.Millisecond
	sphereColorJitter      = 0.10
	sphereEndFadePortion   = 0.45
)

var sphereBaseColor = color.RGBA{R: 0x8E, G: 0xC5, B: 0xFF, A: 0xFF}

func NewSphere(o vec.Vec) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(900, 2./9)) * time.Millisecond
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.60, 0.16, 0, nil),
		fw.MkShell(sphereShellStep, NewSphereSpawner()),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type sphereSpawner struct {
	base color.RGBA
}

func NewSphereSpawner() *sphereSpawner {
	return &sphereSpawner{base: sphereBaseColor}
}

func (s *sphereSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		n := helper.JitterInt(9, 1./3)
		base := helper.JitterFloat(42, 4./21)
		burstColor := varyFromBase(s.base)

		for dir := range AllStratifiedCircle(1, n) {
			speed := base * (0.90 + 0.16*rand.Float64())
			v := dir.Scale(speed)

			ttl := time.Duration(helper.JitterInt(2200, 2./11)) * time.Millisecond
			phases := particle.PhaseTiming{
				SpawningEnd: float64(sphereSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(sphereActiveDuration) / float64(ttl),
			}

			p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, burstColor, ttl, 3, phases)

			if !yield(p) {
				return
			}
		}
	}
}

func varyFromBase(base color.RGBA) color.RGBA {
	return color.RGBA{
		R: jitterChannel(base.R, sphereColorJitter),
		G: jitterChannel(base.G, sphereColorJitter),
		B: jitterChannel(base.B, sphereColorJitter),
		A: base.A,
	}
}

func jitterChannel(v uint8, amount float64) uint8 {
	f := float64(v)
	j := 1 + (rand.Float64()*2-1)*amount
	s := f * j

	s = helper.Clamp(s, 0, 255)

	return uint8(s)
}

func (s *sphereSpawner) Gradient(p *particle.Particle) color.RGBA {
	x := p.Life()
	c := p.Color()

	g := PhaseGradient2(x)

	// Smoothly fade to black over a longer tail to avoid hard cutoff at death.
	end := (1 - x) / sphereEndFadePortion
	end = helper.Clamp(end, 0, math.Min(1, end))

	g *= end * end

	return scaleColor(c, g)
}
