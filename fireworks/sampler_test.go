package fireworks

import (
	"math"
	"testing"
)

func TestAllUniformDiskCountAndBounds(t *testing.T) {
	const (
		ρ = 10.0
		n = 500
	)

	count := 0
	for v := range AllUniformDisk(ρ, n) {
		count++

		if v.X < -ρ || v.X > ρ || v.Y < -ρ || v.Y > ρ {
			t.Fatalf("point out of bounding box: %+v", v)
		}

		if v.Length2() > ρ*ρ {
			t.Fatalf("point outside disk: %+v", v)
		}
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}
}

func TestAllUniformFanCountAndMirroredAngles(t *testing.T) {
	const (
		ρ = 12.0
		α = 0.2
		ω = 1.0
		n = 300
	)

	count := 0
	for v := range AllUniformFan(ρ, α, ω, n) {
		count++

		if v.Length2() > ρ*ρ {
			t.Fatalf("point outside fan radius: %+v", v)
		}

		θ := v.Angle()
		if θ < -ω || θ > -α {
			t.Fatalf("point angle outside mirrored fan range: angle=%f, want [%f,%f]", θ, -ω, -α)
		}
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}
}

func TestAllUniformCircleCountAndRadius(t *testing.T) {
	const (
		ρ = 7.5
		n = 400
	)

	count := 0
	for v := range AllUniformCircle(ρ, n) {
		count++

		err := math.Abs(v.Length2() - ρ*ρ)
		if err > 1e-9 {
			t.Fatalf("point not on circle: %+v (err=%g)", v, err)
		}
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}
}

func TestAllUniformCircleBalancedCentered(t *testing.T) {
	const (
		ρ = 7.5
		n = 400
	)

	var sx, sy float64
	count := 0

	for v := range AllUniformCircleBalanced(ρ, n) {
		count++
		sx += v.X
		sy += v.Y
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}

	if math.Abs(sx) > 1e-9 || math.Abs(sy) > 1e-9 {
		t.Fatalf("expected centered balanced circle, got sum=(%g,%g)", sx, sy)
	}
}

func TestAllUniformArcCountAndRadius(t *testing.T) {
	const (
		ρ = 9.0
		α = 0.3
		ω = 1.1
		n = 250
	)

	count := 0
	for v := range AllUniformArc(ρ, α, ω, n) {
		count++

		err := math.Abs(v.Length2() - ρ*ρ)
		if err > 1e-9 {
			t.Fatalf("point not on arc radius: %+v (err=%g)", v, err)
		}

		theta := v.Angle()
		if theta < -ω || theta > -α {
			t.Fatalf("point angle outside mirrored arc range: angle=%f, want [%f,%f]", theta, -ω, -α)
		}
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}
}

func TestAllNormalDiskCountAndBounds(t *testing.T) {
	const (
		ρ = 15.0
		σ = 2.0
		n = 600
	)

	count := 0
	for v := range AllNormalDisk(ρ, σ, n) {
		count++

		if v.X < -ρ || v.X > ρ || v.Y < -ρ || v.Y > ρ {
			t.Fatalf("point out of bounding box: %+v", v)
		}

		if v.Length2() > ρ*ρ {
			t.Fatalf("point outside disk: %+v", v)
		}
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}
}

func TestAllNormalDiskBalancedCentered(t *testing.T) {
	const (
		ρ = 15.0
		σ = 2.0
		n = 600
	)

	var sx, sy float64
	count := 0

	for v := range AllNormalDiskBalanced(ρ, σ, n) {
		count++
		sx += v.X
		sy += v.Y
	}

	if count != n {
		t.Fatalf("unexpected point count: got %d, want %d", count, n)
	}

	if math.Abs(sx) > 1e-9 || math.Abs(sy) > 1e-9 {
		t.Fatalf("expected centered balanced disk, got sum=(%g,%g)", sx, sy)
	}
}
