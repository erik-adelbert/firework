package firework

import (
	"reflect"
	"testing"

	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

func TestMkForces(t *testing.T) {
	apply := func(p *particle.Particle) vec.Vec {
		if p == nil {
			return vec.Vec{}
		}
		return vec.Vec{X: 1, Y: -1}
	}

	got := MkForces(9.81, 0.25, 42.0, apply)

	if got.g != 9.81 {
		t.Fatalf("g: got %v, want %v", got.g, 9.81)
	}
	if got.ar != 0.25 {
		t.Fatalf("ar: got %v, want %v", got.ar, 0.25)
	}
	if got.spd != 42.0 {
		t.Fatalf("spd: got %v, want %v", got.spd, 42.0)
	}

	if reflect.ValueOf(got.af).Pointer() != reflect.ValueOf(apply).Pointer() {
		t.Fatalf("af: function mismatch")
	}

	// Sanity check: function behaves as expected.
	if v := got.af(&particle.Particle{}); v != (vec.Vec{X: 1, Y: -1}) {
		t.Fatalf("af result: got %v, want %v", v, vec.Vec{X: 1, Y: -1})
	}
}
