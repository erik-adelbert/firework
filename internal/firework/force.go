package firework

import (
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

type Forces struct {
	g   float64 // gravity
	ar  float64 // air resistance
	spd float64 // max speed

	af func(p *particle.Particle) vec.Vec // additional force
}

func MkForces(gravity, drag, maxSpeed float64, applyForce func(p *particle.Particle) vec.Vec) Forces {
	return Forces{
		g:   gravity,
		ar:  drag,
		spd: maxSpeed,

		af: applyForce,
	}
}

func DefaultForces() Forces {
	return Forces{
		g:   1.0,
		ar:  0.28,
		spd: 0,

		af: nil,
	}
}
