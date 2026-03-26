package particle

import (
	"image/color"
	"testing"
	"time"

	"github.com/erik-adelbert/firework/internal/vec"
)

func TestNew(t *testing.T) {
	p0 := vec.Vec{X: 1, Y: 2}
	v0 := vec.Vec{X: 3, Y: 4}
	a0 := vec.Vec{X: 5, Y: 6}

	c0 := color.RGBA{R: 255, G: 128, B: 64, A: 200}

	ttl := 1 * time.Second

	trailSize := 10

	p := New(p0, v0, a0, c0, ttl, trailSize)

	if p.Pos() != p0 {
		t.Errorf("expected pos %v, got %v", p0, p.Pos())
	}
	if p.Vel() != v0 {
		t.Errorf("expected vel %v, got %v", v0, p.Vel())
	}
	if p.Accel() != a0 {
		t.Errorf("expected acc %v, got %v", a0, p.Accel())
	}
	if p.Color() != c0 {
		t.Errorf("expected color %v, got %v", c0, p.Color())
	}
}

func TestIsAlive(t *testing.T) {
	p := New(vec.Vec{}, vec.Vec{}, vec.Vec{}, color.RGBA{}, 1*time.Second, 0)
	if !p.IsAlive() {
		t.Error("expected particle to be alive")
	}

	p2 := New(vec.Vec{}, vec.Vec{}, vec.Vec{}, color.RGBA{}, 0, 0)
	if p2.IsAlive() {
		t.Error("expected particle with ttl=0 to be dead")
	}
}

func TestSpeed(t *testing.T) {
	p := New(vec.Vec{}, vec.Vec{X: 3, Y: 4}, vec.Vec{}, color.RGBA{}, 1*time.Second, 0)
	expectedSpeed := 5.0 // sqrt(3^2 + 4^2) = 5

	if p.Speed() != expectedSpeed {
		t.Errorf("expected speed %v, got %v", expectedSpeed, p.Speed())
	}
}

func TestSpeed2(t *testing.T) {
	p := New(vec.Vec{}, vec.Vec{X: 3, Y: 4}, vec.Vec{}, color.RGBA{}, 1*time.Second, 0)
	expectedSpeed2 := 25.0 // 3^2 + 4^2 = 25

	if p.Speed2() != expectedSpeed2 {
		t.Errorf("expected speed2 %v, got %v", expectedSpeed2, p.Speed2())
	}
}

func TestState(t *testing.T) {
	tests := []struct {
		name     string
		age      time.Duration
		ttl      time.Duration
		expected State
	}{
		{
			name:     "spawning state",
			age:      100 * time.Millisecond,
			ttl:      1 * time.Second,
			expected: Spawning,
		},
		{
			name:     "active state",
			age:      500 * time.Millisecond,
			ttl:      1 * time.Second,
			expected: Active,
		},
		{
			name:     "fading state",
			age:      800 * time.Millisecond,
			ttl:      1 * time.Second,
			expected: Fading,
		},
		{
			name:     "dead state",
			age:      1 * time.Second,
			ttl:      1 * time.Second,
			expected: Dead,
		},
		{
			name:     "zero ttl means dead",
			age:      0,
			ttl:      0,
			expected: Dead,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := state(tt.age, tt.ttl)
			if s != tt.expected {
				t.Errorf("expected state %v, got %v", tt.expected, s)
			}
		})
	}
}

func TestStateWithPhaseTiming(t *testing.T) {
	phases := PhaseTiming{
		SpawningEnd: 0.2,
		ActiveEnd:   0.35,
	}

	tests := []struct {
		name     string
		age      time.Duration
		ttl      time.Duration
		expected State
	}{
		{
			name:     "custom spawning state",
			age:      150 * time.Millisecond,
			ttl:      1 * time.Second,
			expected: Spawning,
		},
		{
			name:     "custom active state",
			age:      300 * time.Millisecond,
			ttl:      1 * time.Second,
			expected: Active,
		},
		{
			name:     "custom fading state",
			age:      500 * time.Millisecond,
			ttl:      1 * time.Second,
			expected: Fading,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := stateWithPhaseTiming(tt.age, tt.ttl, phases)
			if s != tt.expected {
				t.Errorf("expected state %v, got %v", tt.expected, s)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		10,
	)

	// Update for 10 milliseconds with no gravity or drag
	alive := p.Update(10*time.Millisecond, 0, 0, 0, nil)

	if !alive {
		t.Error("expected particle to be alive after update")
	}
	if p.Pos().X <= 0 {
		t.Errorf("expected x position to increase, got %v", p.Pos().X)
	}
}

func TestUpdateWithGravity(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 10},
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		0,
	)

	gravity := -9.8
	p.Update(100*time.Millisecond, gravity, 0, 0, nil)

	if p.Pos().Y >= 10 {
		t.Errorf("expected y position to decrease due to gravity, got %v", p.Pos().Y)
	}
	if p.Vel().Y >= 0 {
		t.Errorf("expected y velocity to be negative due to gravity, got %v", p.Vel().Y)
	}
}

func TestUpdateDeadParticle(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 1},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		0,
		0,
	)

	alive := p.Update(10*time.Millisecond, 0, 0, 0, nil)

	if alive {
		t.Error("expected particle to be dead after update")
	}
}

func TestUpdateWithMaxSpeed(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 100, Y: 100},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		0,
	)

	maxSpeed := 50.0
	p.Update(10*time.Millisecond, 0, 0, maxSpeed, nil)

	speed := p.Speed()
	if speed > maxSpeed+0.01 { // small epsilon for floating point errors
		t.Errorf("expected speed to be capped at %v, got %v", maxSpeed, speed)
	}
}

func TestTrailLen(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		5,
	)

	p.Update(10*time.Millisecond, 0, 0, 0, nil)

	if p.TrailLen() == 0 {
		t.Error("expected trail to have items after update")
	}
}

func TestTrailEmpty(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		0,
	)

	p.Update(10*time.Millisecond, 0, 0, 0, nil)

	if p.TrailLen() > 0 {
		t.Error("expected empty trail when trailSize is 0")
	}
}

func TestApplyForce(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		0,
	)

	customForce := func(p *Particle) vec.Vec {
		return vec.Vec{X: 10, Y: 0}
	}

	p.Update(10*time.Millisecond, 0, 0, 0, customForce)

	if p.Vel().X <= 0 {
		t.Errorf("expected velocity to be affected by custom force, got %v", p.Vel())
	}
}

func TestAllTrailPoints(t *testing.T) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		50,
	)

	p.Update(50*time.Millisecond, 0, 0, 0, nil)

	count := 0

	for _, pos := range p.AllTrailPoints() {
		count++
		_ = pos
	}

	if count != 50 {
		t.Error("expected AllTrailPoints to return 50 positions, got", count)
	}
}

func BenchmarkNew(b *testing.B) {
	p0 := vec.Vec{X: 1, Y: 2}
	v0 := vec.Vec{X: 3, Y: 4}
	a0 := vec.Vec{X: 5, Y: 6}
	c0 := color.RGBA{R: 255, G: 128, B: 64, A: 200}

	for b.Loop() {
		New(p0, v0, a0, c0, 1*time.Second, 10)
	}
}

func BenchmarkSpeed(b *testing.B) {
	p := New(vec.Vec{}, vec.Vec{X: 3, Y: 4}, vec.Vec{}, color.RGBA{}, 1*time.Second, 0)

	for b.Loop() {
		p.Speed()
	}
}

func BenchmarkSpeed2(b *testing.B) {
	p := New(vec.Vec{}, vec.Vec{X: 3, Y: 4}, vec.Vec{}, color.RGBA{}, 1*time.Second, 0)

	for b.Loop() {
		p.Speed2()
	}
}

func BenchmarkUpdate(b *testing.B) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		50,
	)

	for b.Loop() {
		p.Update(50*time.Millisecond, 0, 0, 0, nil)
	}
}

func BenchmarkUpdateWithGravity(b *testing.B) {
	p := New(
		vec.Vec{X: 0, Y: 10},
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		0,
	)

	for b.Loop() {
		p.Update(100*time.Millisecond, -9.8, 0, 0, nil)
	}
}

func BenchmarkIsAlive(b *testing.B) {
	p := New(vec.Vec{}, vec.Vec{}, vec.Vec{}, color.RGBA{}, 1*time.Second, 0)

	for b.Loop() {
		p.IsAlive()
	}
}

func BenchmarkAllTrailPoints(b *testing.B) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		50,
	)

	p.Update(50*time.Millisecond, 0, 0, 0, nil)

	for b.Loop() {
		for range p.AllTrailPoints() {
		}
	}
}

func BenchmarkAllTrailPointsSeq(b *testing.B) {
	p := New(
		vec.Vec{X: 0, Y: 0},
		vec.Vec{X: 1, Y: 0},
		vec.Vec{X: 0, Y: 0},
		color.RGBA{},
		1*time.Second,
		50,
	)

	p.Update(50*time.Millisecond, 0, 0, 0, nil)

	foreachTrail := p.AllTrailPoints()

	for b.Loop() {
		foreachTrail(func(_ int, _ vec.Vec) bool { return true })
	}
}
