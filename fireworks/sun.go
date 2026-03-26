package fireworks

import (
	"image/color"
	"iter"
	"math/rand/v2"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

var sunColor = color.RGBA{R: 250, G: 216, B: 68, A: 255}

const (
	sunSpeed       = 100
	sunParticles   = 300
	sunTTLMinMs    = 3000
	sunTTLRangeMs  = 2500
	sunTrailMin    = 30
	sunTrailSpread = 20
)

func NewSun(o vec.Vec) *fw.Firework {
	return NewSunWithGradient(o, true)
}

func NewSunWithGradient(o vec.Vec, enableGradient bool) *fw.Firework {
	const spawnAfter = 700 * time.Millisecond
	// TTL must exceed spawnAfter + max particle lifetime so particles die naturally.
	ttl := spawnAfter + 10*time.Second

	return fw.New(
		o,
		fw.MkForces(0, 0.15, 0, nil),
		fw.MkShell(0, NewSunSpawner(enableGradient)),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type sunSpawner struct {
	enableGradient bool
}

func NewSunSpawner(enableGradient bool) *sunSpawner {
	return &sunSpawner{enableGradient: enableGradient}
}

func (s *sunSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		for v := range AllUniformCircle(sunSpeed, sunParticles) {
			ttl := time.Duration(sunTTLMinMs+rand.IntN(sunTTLRangeMs)) * time.Millisecond
			trail := rand.IntN(sunTrailSpread) + sunTrailMin

			if !yield(particle.New(o, v, vec.Vec{}, sunColor, ttl, trail)) {
				return
			}
		}
	}
}

func (s *sunSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	if !s.enableGradient {
		return c
	}

	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}
