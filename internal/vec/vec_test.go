package vec

import (
	"math"
	"testing"
)

func floatEquals(a, b float64) bool {
	const eps = 1e-9
	return math.Abs(a-b) < eps
}

func TestAdd(t *testing.T) {
	v := Vec{1, 2}
	u := Vec{3, 4}

	res := v.Add(u)

	if res != (Vec{4, 6}) {
		t.Errorf("Add: expected {4,6}, got %+v", res)
	}
}

func TestSub(t *testing.T) {
	v := Vec{5, 7}
	u := Vec{2, 3}

	res := v.Sub(u)

	if res != (Vec{3, 4}) {
		t.Errorf("Sub: expected {3,4}, got %+v", res)
	}
}

func TestMul(t *testing.T) {
	v := Vec{2, 3}
	u := Vec{4, 5}

	res := v.Mul(u)

	if res != (Vec{8, 15}) {
		t.Errorf("Mul: expected {8,15}, got %+v", res)
	}
}

func TestScale(t *testing.T) {
	v := Vec{2, 3}

	res := v.Scale(2.5)

	if !floatEquals(res.X, 5) || !floatEquals(res.Y, 7.5) {
		t.Errorf("Scale: expected {5,7.5}, got %+v", res)
	}
}

func TestLength2(t *testing.T) {
	v := Vec{3, 4}

	if !floatEquals(v.Length2(), 25) {
		t.Errorf("Length2: expected 25, got %v", v.Length2())
	}
}

func TestLength(t *testing.T) {
	v := Vec{3, 4}

	if !floatEquals(v.Length(), 5) {
		t.Errorf("Length: expected 5, got %v", v.Length())
	}
}

func TestNormalize(t *testing.T) {
	v := Vec{3, 4}

	n := v.Normalize()

	if !floatEquals(n.Length(), 1) {
		t.Errorf("Normalize: expected length 1, got %v", n.Length())
	}

	zero := Vec{0, 0}.Normalize()

	if zero != (Vec{0, 0}) {
		t.Errorf("Normalize: expected {0,0} for zero vector, got %+v", zero)
	}
}

func TestDist2(t *testing.T) {
	v := Vec{1, 2}
	u := Vec{4, 6}

	if !floatEquals(v.Dist2(u), 25) {
		t.Errorf("Dist2: expected 25, got %v", v.Dist2(u))
	}
}

func TestDist(t *testing.T) {
	v := Vec{1, 2}
	u := Vec{4, 6}

	if !floatEquals(v.Dist(u), 5) {
		t.Errorf("Dist: expected 5, got %v", v.Dist(u))
	}
}

func TestAngle(t *testing.T) {
	v := Vec{0, 1}

	if !floatEquals(v.Angle(), math.Pi/2) {
		t.Errorf("Angle: expected %v, got %v", math.Pi/2, v.Angle())
	}

	v = Vec{1, 0}

	if !floatEquals(v.Angle(), 0) {
		t.Errorf("Angle: expected 0, got %v", v.Angle())
	}

	v = Vec{-1, 0}

	if !floatEquals(v.Angle(), math.Pi) && !floatEquals(v.Angle(), -math.Pi) {
		t.Errorf("Angle: expected pi or -pi, got %v", v.Angle())
	}
}
