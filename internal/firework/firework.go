package firework

import (
	"iter"
	"slices"
	"time"

	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

const sizeHint = 50

type NewFirework func(vec.Vec) *Firework

type Firework struct {
	State
	Forces
	Shell
	Timing

	o   vec.Vec
	age time.Duration

	Particles []*particle.Particle
}

func New(o vec.Vec, f Forces, s Shell, t Timing) *Firework {
	return &Firework{
		Forces: f,
		Shell:  s,
		Timing: t,

		o: o,

		Particles: make([]*particle.Particle, 0, sizeHint),
	}
}

func (f *Firework) Reset() {
	f.Particles = f.Particles[:0]

	f.Timing.Reset()

	f.Shell.Reset()

	f.age = 0

	f.State = Ready
}

func (f *Firework) Update(now time.Time, dt time.Duration) {
	f.age += dt

	switch f.State {
	case Ready:
		return // do nothing until triggered

	case CountingDown:
		if f.t0.IsZero() {
			// If t0 is zero, it means the firework was triggered without
			// a specific time (e.g., via Trigger(time.Time{})) instead of
			// Trigger(t0)). In this case, we set t0 to now to start the
			// countdown from the current time.
			f.t0 = now
		}

		if f.age < f.spawnAfter {
			break
		}

		f.State = Spawning

		fallthrough
	case Spawning:
		if f.age < f.ttl {
			f.Particles = slices.AppendSeq(f.Particles, f.Spawn(f.o, now, dt))
			break
		}

		f.State = Dead

		fallthrough
	case Dead:
		f.Particles = f.Particles[:0]
	}

	alives := f.Particles[:0] // reuse the same slice to retain alive particles

	for _, p := range f.Particles {
		if p.Update(dt, f.g, f.ar, f.spd, f.af) {
			alives = append(alives, p)
		}
	}

	// log.Printf("Updated firework at age %v, %d particles alive\n", f.age, len(alives))

	f.Particles = alives

	// Only declare the firework dead via empty particles if it has already
	// started spawning. During CountingDown there are legitimately 0 particles,
	// and we must not kill the firework before it has had a chance to explode.
	if f.State == Spawning && len(f.Particles) == 0 {
		f.State = Dead
	}
}

func (f *Firework) Trigger(t0 time.Time) {
	if f.State == Ready {
		f.State = CountingDown
		f.t0 = t0
	}
}

func (f *Firework) Age() time.Duration {
	return f.age
}

func (f *Firework) TTL() time.Duration {
	return f.ttl
}

func (f *Firework) Life() float64 {
	if f.ttl == 0 {
		return 1
	}

	return f.age.Seconds() / f.ttl.Seconds()
}

func (f *Firework) SpawnAfter() time.Duration {
	return f.spawnAfter
}

func (f *Firework) SetCenter(o vec.Vec) {
	f.o = o
}

func (f *Firework) Center() vec.Vec {
	return f.o
}

func (f *Firework) AllParticles() iter.Seq2[int, *particle.Particle] {
	return slices.All(f.Particles)
}

func (f *Firework) IsAlive() bool {
	return f.State != Dead
}

type State int

const (
	Ready State = iota
	CountingDown
	Spawning
	Dead
)
