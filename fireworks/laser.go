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

var laserColors = []color.RGBA{
	{R: 0x98, G: 0xBA, B: 0xE3, A: 0xFF},
	{R: 0x36, G: 0x54, B: 0x75, A: 0xFF},
	{R: 0x15, G: 0x27, B: 0x3C, A: 0xFF},
}

func NewLaser(o vec.Vec) *fw.Firework {
	return NewLaserWithGradient(o, true)
}

func NewLaserWithGradient(o vec.Vec, enableGradient bool) *fw.Firework {
	ttl := time.Duration(helper.JitterInt(25, 1./5)) * time.Second
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(1.4, 0, 0, nil),
		fw.MkShell(0, NewLaserSpawner(enableGradient)),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type laserSpawner struct {
	enableGradient bool
}

func NewLaserSpawner(enableGradient bool) *laserSpawner {
	return &laserSpawner{enableGradient: enableGradient}
}

func (s *laserSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		for v := range AllNormalDisk9(450, 80) {
			trail := helper.JitterInt(38, 5./38)
			ttl := time.Duration(helper.JitterInt(3750, 1./15)) * time.Millisecond
			c, _ := helper.Pick(laserColors)

			p := particle.New(o, v, vec.Vec{}, c, ttl, trail)

			if !yield(p) {
				return
			}
		}
	}
}

func (s *laserSpawner) Gradient(p *particle.Particle) color.RGBA {
	c := p.Color()
	if !s.enableGradient {
		return c
	}

	g := PhaseGradient5(p.Life())

	return scaleColor(c, g)
}
