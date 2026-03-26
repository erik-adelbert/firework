package firework

import (
	"testing"
	"time"
)

func TestMkTiming(t *testing.T) {
	ttl := 5 * time.Second
	spawnAfter := 2 * time.Second

	timing := MkTiming(ttl, spawnAfter)

	if timing.ttl != ttl {
		t.Errorf("expected ttl %v, got %v", ttl, timing.ttl)
	}
	if timing.spawnAfter != spawnAfter {
		t.Errorf("expected spawnAfter %v, got %v", spawnAfter, timing.spawnAfter)
	}
	if !timing.t0.IsZero() {
		t.Errorf("expected t0 to be zero, got %v", timing.t0)
	}
}

func TestTimingReset(t *testing.T) {
	ttl := 3 * time.Second
	spawnAfter := 1 * time.Second

	timing := MkTiming(ttl, spawnAfter)

	// Modify t0 to non-zero
	timing.t0 = time.Now()
	if timing.t0.IsZero() {
		t.Fatal("failed to set t0 to non-zero time")
	}

	// Reset should clear t0 but preserve other fields
	timing.Reset()

	if !timing.t0.IsZero() {
		t.Errorf("expected t0 to be zero after reset, got %v", timing.t0)
	}
	if timing.ttl != ttl {
		t.Errorf("expected ttl %v after reset, got %v", ttl, timing.ttl)
	}
	if timing.spawnAfter != spawnAfter {
		t.Errorf("expected spawnAfter %v after reset, got %v", spawnAfter, timing.spawnAfter)
	}
}
