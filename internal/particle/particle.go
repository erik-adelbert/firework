package particle

import (
	"image/color"
	"iter"
	"time"

	"github.com/erik-adelbert/firework/internal/vec"
	"github.com/erik-adelbert/firework/pkg/ring"
)

type Particle struct {
	pos, vel, acc vec.Vec

	age, ttl time.Duration
	phases   PhaseTiming

	color color.RGBA
	trail *ring.Ring[vec.Vec]
}

func (p *Particle) IsAlive() bool {
	return stateWithPhaseTiming(p.age, p.ttl, p.phases) != Dead
}

func New(p0, v0, a0 vec.Vec, c0 color.RGBA, ttl time.Duration, trailLen int) *Particle {
	return NewWithPhaseTiming(p0, v0, a0, c0, ttl, trailLen, DefaultPhaseTiming())
}

func NewWithPhaseTiming(
	p0, v0, a0 vec.Vec,
	c0 color.RGBA,
	ttl time.Duration,
	trailLen int,
	phases PhaseTiming,
) *Particle {
	trail := ring.NewRing[vec.Vec](trailLen)

	for range trailLen {
		trail.PushBack(p0)
	}

	phases.normalize()

	return &Particle{
		pos: p0,
		vel: v0,
		acc: a0,

		ttl:    ttl,
		phases: phases,

		color: c0,

		trail: trail,
	}
}

func (p *Particle) Age() time.Duration {
	return p.age
}

func (p *Particle) Life() float64 {
	if p.ttl == 0 {
		return 1
	}

	return p.age.Seconds() / p.ttl.Seconds()
}

func (p *Particle) Accel() vec.Vec {
	return p.acc
}

func (p *Particle) Color() color.RGBA {
	return p.color
}

func (p *Particle) TrailLen() int {
	return p.trail.Len()
}

func (p *Particle) Pos() vec.Vec {
	return p.pos
}

func (p *Particle) Speed() float64 {
	return p.vel.Length()
}

func (p *Particle) Speed2() float64 {
	return p.vel.Length2()
}

func (p *Particle) Vel() vec.Vec {
	return p.vel
}

func (p *Particle) Update(
	dt time.Duration,
	gravity, drag, maxVelocity float64,
	applyForce func(p *Particle) vec.Vec,
) bool {
	if p.age += dt; !p.IsAlive() {
		return false
	}

	const ms = time.Millisecond

	for range dt.Milliseconds() {
		// Apply gravity and external forces
		acc := vec.Vec{X: 0, Y: 10 * gravity}.Add(p.acc)

		if applyForce != nil {
			acc = acc.Add(applyForce(p))
		}

		// Apply drag
		if drag > 0 {
			// quadratic drag: https://en.wikipedia.org/wiki/Drag_(physics)#Quadratic_drag
			if speed2 := p.Speed2(); speed2 > 0 {
				dragForce := p.vel.Normalize().Scale(-drag * speed2)
				acc = acc.Add(dragForce)
			}
		}

		// Integrate velocity
		p.vel = p.vel.Add(acc.Scale(ms.Seconds()))

		// Cap speed
		if maxVelocity > 0 && p.vel.Length2() > maxVelocity*maxVelocity {
			p.vel = p.vel.Normalize().Scale(maxVelocity)
		}

		// Integrate position
		p.pos = p.pos.Add(p.vel.Scale(ms.Seconds()))
	}

	// Update the trail
	if p.trail.Cap() > 0 {
		p.trail.PushFront(p.pos)
	}

	return true
}

func (p *Particle) AllTrailPoints() iter.Seq2[int, vec.Vec] {
	return p.trail.AllBack()
}

func (p *Particle) Trail(i int) vec.Vec {
	return p.trail.At(i)
}

func (p *Particle) State() State {
	return stateWithPhaseTiming(p.age, p.ttl, p.phases)
}

type PhaseTiming struct {
	SpawningEnd float64
	ActiveEnd   float64
}

func DefaultPhaseTiming() PhaseTiming {
	return PhaseTiming{
		SpawningEnd: 0.4,
		ActiveEnd:   0.65,
	}
}

func (pt *PhaseTiming) normalize() {
	if pt.SpawningEnd <= 0 || pt.SpawningEnd >= 1 {
		pt.SpawningEnd = DefaultPhaseTiming().SpawningEnd
	}

	if pt.ActiveEnd <= pt.SpawningEnd || pt.ActiveEnd >= 1 {
		pt.ActiveEnd = DefaultPhaseTiming().ActiveEnd
	}
}

type State int

const (
	Dead State = iota
	Fading
	Active
	Spawning
)

func state(age, ttl time.Duration) State {
	return stateWithPhaseTiming(age, ttl, DefaultPhaseTiming())
}

func stateWithPhaseTiming(age, ttl time.Duration, phases PhaseTiming) State {
	if ttl == 0 || age >= ttl {
		return Dead
	}

	phases.normalize()

	ρ := float64(age.Nanoseconds())
	ρ /= float64(ttl.Nanoseconds())

	switch {
	case ρ < phases.SpawningEnd:
		return Spawning
	case ρ < phases.ActiveEnd:
		return Active
	case ρ < 1.0:
		return Fading
	default:
		return Dead
	}
}
