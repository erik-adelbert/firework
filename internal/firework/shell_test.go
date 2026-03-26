package firework

import (
	"image/color"
	"iter"
	"testing"
	"time"

	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

type mockSpawner struct{}

func (m mockSpawner) Spawn(o vec.Vec, now time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return nil
}

func (m mockSpawner) Gradient(p *particle.Particle) color.RGBA {
	return color.RGBA{}
}

func (m mockSpawner) Reset() {}

func TestMkShell_PositiveStep(t *testing.T) {
	step := 100 * time.Millisecond
	spawner := mockSpawner{}

	shell := MkShell(step, spawner)

	if shell.step != step {
		t.Errorf("expected step %v, got %v", step, shell.step)
	}
	if shell.kind != Sustained {
		t.Errorf("expected kind Sustained, got %v", shell.kind)
	}
}

func TestMkShell_ZeroStep(t *testing.T) {
	spawner := mockSpawner{}

	shell := MkShell(0, spawner)

	if shell.step != 0 {
		t.Errorf("expected step 0, got %v", shell.step)
	}
	if shell.kind != Instant {
		t.Errorf("expected kind Instant, got %v", shell.kind)
	}
}

func TestMkShell_NegativeStep(t *testing.T) {
	step := -50 * time.Millisecond
	spawner := mockSpawner{}

	shell := MkShell(step, spawner)

	if shell.step != 0 {
		t.Errorf("expected step 0, got %v", shell.step)
	}
	if shell.kind != Instant {
		t.Errorf("expected kind Instant, got %v", shell.kind)
	}
}
