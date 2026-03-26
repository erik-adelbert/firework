package firework

import (
	"image/color"
	"iter"
	"testing"
	"time"

	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

type mockShell struct {
	spawnCount int
}

func (m *mockShell) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	m.spawnCount++
	return helper.EmptySeq[*particle.Particle]()
}

func (m *mockShell) Gradient(p *particle.Particle) color.RGBA {
	return color.RGBA{}
}

func (m *mockShell) Reset() {
	m.spawnCount = 0
}

func TestNew(t *testing.T) {
	o := vec.Vec{X: 10, Y: 20}
	f := Forces{
		g: 9.8, ar: 0.1, spd: 0.5,
		af: func(p *particle.Particle) vec.Vec {
			return vec.Vec{X: 0.2, Y: 0.3}
		},
	}
	s := MkShell(
		100*time.Millisecond,
		&mockShell{},
	)
	tm := Timing{
		ttl:        1 * time.Second,
		spawnAfter: 100 * time.Millisecond,
	}

	fw := New(o, f, s, tm)

	if fw.Center() != o {
		t.Errorf("expected center %v, got %v", o, fw.Center())
	}

	if fw.State != Ready {
		t.Errorf("expected state Ready, got %v", fw.State)
	}

	if len(fw.Particles) != 0 {
		t.Errorf("expected 0 particles, got %d", len(fw.Particles))
	}

	if cap(fw.Particles) < sizeHint {
		t.Errorf("expected capacity at least %d, got %d", sizeHint, cap(fw.Particles))
	}
}

func TestReset(t *testing.T) {
	o := vec.Vec{X: 5, Y: 15}
	f := Forces{
		g: 9.8, ar: 0.1, spd: 0.5,
		af: func(p *particle.Particle) vec.Vec {
			return vec.Vec{X: 0.2, Y: 0.3}
		},
	}
	s := MkShell(
		100*time.Millisecond,
		&mockShell{},
	)
	tm := Timing{
		ttl:        1 * time.Second,
		spawnAfter: 100 * time.Millisecond,
	}

	fw := New(o, f, s, tm)
	fw.State = Spawning
	fw.age = 500 * time.Millisecond

	fw.Reset()

	if fw.State != Ready {
		t.Errorf("expected state Ready after reset, got %v", fw.State)
	}

	if fw.age != 0 {
		t.Errorf("expected age 0 after reset, got %v", fw.age)
	}

	if len(fw.Particles) != 0 {
		t.Errorf("expected 0 particles after reset, got %d", len(fw.Particles))
	}
}

func TestTrigger(t *testing.T) {
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)

	if fw.State != Ready {
		t.Errorf("expected initial state Ready, got %v", fw.State)
	}

	now0 := time.Now()
	fw.Trigger(now0)

	if fw.State != CountingDown {
		t.Errorf("expected state CountingDown after trigger, got %v", fw.State)
	}

	now1 := time.Now()
	fw.Trigger(now1)

	if fw.State != CountingDown {
		t.Errorf("expected state to remain CountingDown, got %v", fw.State)
	}

	if fw.t0 != now0 {
		t.Errorf("expected t0 to be set to the time of the first trigger, got %v", fw.t0)
	}
}

func TestUpdateReady(t *testing.T) {
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)
	fw.State = Ready

	fw.Update(time.Now(), 16*time.Millisecond)

	if fw.age != 16*time.Millisecond {
		t.Errorf("expected age 16ms, got %v", fw.age)
	}

	if fw.State != Ready {
		t.Errorf("expected state Ready, got %v", fw.State)
	}
}

func TestUpdateCountingDown(t *testing.T) {
	s := &mockShell{}
	tm := Timing{
		ttl:        1 * time.Second,
		spawnAfter: 100 * time.Millisecond,
	}
	fw := New(vec.Vec{}, Forces{}, MkShell(100*time.Millisecond, s), tm)
	// fw.State = CountingDown

	now := time.Now()
	fw.Trigger(now)

	if fw.State != CountingDown {
		t.Errorf("expected state CountingDown, got %v", fw.State)
	}

	now = time.Now()
	fw.Update(now, 100*time.Millisecond)

	if fw.State != Dead {
		t.Errorf("expected state Dead after countdown when no particles are spawned, got %v", fw.State)
	}

	if s.spawnCount == 0 {
		t.Errorf("expected shell spawn to be attempted after countdown")
	}
}

func TestUpdateSpawning(t *testing.T) {
	s := MkShell(100*time.Millisecond, &mockShell{})
	tm := Timing{ttl: 200 * time.Millisecond, spawnAfter: 0}
	fw := New(vec.Vec{}, Forces{}, s, tm)
	fw.State = Spawning

	now := time.Now()
	fw.Update(now, 50*time.Millisecond)

	// if fw.State != Spawning {
	// 	t.Errorf("expected state Spawning while age < ttl, got %v", fw.State)
	// }

	fw.age = 250 * time.Millisecond
	fw.Update(now, 10*time.Millisecond)

	if fw.State != Dead {
		t.Errorf("expected state Dead after ttl exceeded, got %v", fw.State)
	}

	if len(fw.Particles) != 0 {
		t.Errorf("expected 0 particles when dead, got %d", len(fw.Particles))
	}
}

func TestUpdateParticleFiltering(t *testing.T) {
	fw := New(
		vec.Vec{},
		Forces{
			g: 0.1, ar: 0.01, spd: 0.5,
			af: func(p *particle.Particle) vec.Vec {
				return vec.Vec{X: 0.05, Y: 0.05}
			},
		},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)

	// Create mock particles manually for testing
	fw.Particles = make([]*particle.Particle, 0, 10)

	dt := 16 * time.Millisecond
	fw.Update(time.Now(), dt)

	if len(fw.Particles) != 0 {
		t.Errorf("expected 0 particles, got %d", len(fw.Particles))
	}
}

func TestAge(t *testing.T) {
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)

	if fw.Age() != 0 {
		t.Errorf("expected age 0 initially, got %v", fw.Age())
	}

	fw.age = 150 * time.Millisecond
	if fw.Age() != 150*time.Millisecond {
		t.Errorf("expected age 150ms, got %v", fw.Age())
	}
}

func TestTTL(t *testing.T) {
	tm := Timing{
		ttl: 500 * time.Millisecond, spawnAfter: 100 * time.Millisecond,
	}
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		tm,
	)

	if fw.TTL() != 500*time.Millisecond {
		t.Errorf("expected TTL 500ms, got %v", fw.TTL())
	}
}

func TestSpawnAfter(t *testing.T) {
	tm := Timing{
		ttl: 500 * time.Millisecond, spawnAfter: 75 * time.Millisecond,
	}
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		tm,
	)

	if fw.SpawnAfter() != 75*time.Millisecond {
		t.Errorf("expected SpawnAfter 75ms, got %v", fw.SpawnAfter())
	}
}

func TestCenter(t *testing.T) {
	o := vec.Vec{X: 42.5, Y: 17.3}
	fw := New(
		o,
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)

	if fw.Center() != o {
		t.Errorf("expected center %v, got %v", o, fw.Center())
	}
}

func TestIsAlive(t *testing.T) {
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)

	fw.State = Ready
	if !fw.IsAlive() {
		t.Errorf("expected firework to be alive in Ready state")
	}

	fw.State = CountingDown
	if !fw.IsAlive() {
		t.Errorf("expected firework to be alive in CountingDown state")
	}

	fw.State = Spawning
	if !fw.IsAlive() {
		t.Errorf("expected firework to be alive in Spawning state")
	}

	fw.State = Dead
	if fw.IsAlive() {
		t.Errorf("expected firework to be dead in Dead state")
	}
}

func TestAllParticles(t *testing.T) {
	fw := New(
		vec.Vec{},
		Forces{},
		MkShell(100*time.Millisecond, &mockShell{}),
		Timing{},
	)

	count := 0
	for range fw.AllParticles() {
		count++
	}

	if count != 0 {
		t.Errorf("expected 0 particles, got %d", count)
	}
}
