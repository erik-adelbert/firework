package firework

import (
	"image/color"
	"iter"
	"time"

	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

type Spawner interface {
	Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle]
	Gradient(p *particle.Particle) color.RGBA
}
