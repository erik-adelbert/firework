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

var featherColors = []color.RGBA{
	{R: 0xF2, G: 0xE9, B: 0xBE, A: 0xFF},
	{R: 0xE2, G: 0xC4, B: 0x88, A: 0xFF},
	{R: 0xFF, G: 0xF8, B: 0xFD, A: 0xFF},
}

func NewFeather(o vec.Vec) *fw.Firework {
	return NewFeatherWithGradient(o, true)
}

func NewFeatherWithGradient(o vec.Vec, enableGradient bool) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(10, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.1, 0.19, 0, nil),
		fw.MkShell(0, NewFeatherSpawner(enableGradient)),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type featherSpawner struct {
	enableGradient bool
}

func NewFeatherSpawner(enableGradient bool) *featherSpawner {
	return &featherSpawner{enableGradient: enableGradient}
}

func (s *featherSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		for v := range AllNormalDisk9(350, 35) {
			trail := helper.JitterInt(25, 1./5)
			ttl := time.Duration(helper.JitterInt(3500, .1)) * time.Millisecond

			c, _ := helper.Pick(featherColors)

			if !yield(particle.New(o, v, vec.Vec{}, c, ttl, trail)) {
				return
			}
		}
	}
}

func (s *featherSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()

	if !s.enableGradient {
		return c
	}

	g := PhaseGradient2(p.Life())

	return scaleColor(c, g)
}
