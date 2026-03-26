package helper

import "testing"

type myInt int

type myFloat float64

func TestClampInt(t *testing.T) {
	if got := Clamp(-5, 0, 10); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}

	if got := Clamp(7, 0, 10); got != 7 {
		t.Fatalf("expected 7, got %d", got)
	}

	if got := Clamp(15, 0, 10); got != 10 {
		t.Fatalf("expected 10, got %d", got)
	}
}

func TestClampUint(t *testing.T) {
	if got := Clamp[uint](1, 2, 8); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}

	if got := Clamp[uint](5, 2, 8); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}

	if got := Clamp[uint](20, 2, 8); got != 8 {
		t.Fatalf("expected 8, got %d", got)
	}
}

func TestClampFloat(t *testing.T) {
	if got := Clamp(-0.25, 0.0, 1.0); got != 0.0 {
		t.Fatalf("expected 0.0, got %f", got)
	}

	if got := Clamp(0.4, 0.0, 1.0); got != 0.4 {
		t.Fatalf("expected 0.4, got %f", got)
	}

	if got := Clamp(1.25, 0.0, 1.0); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestClampNamedTypes(t *testing.T) {
	if got := Clamp[myInt](3, 1, 9); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}

	if got := Clamp[myFloat](0.1, 0.9, 0.95); got != 0.9 {
		t.Fatalf("expected 0.9, got %f", got)
	}
}
