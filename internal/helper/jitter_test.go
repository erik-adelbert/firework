package helper

import (
	"math"
	"testing"
)

func TestJitterFloat_NoAmount(t *testing.T) {
	x := 10.0
	amount := 0.0
	result := JitterFloat(x, amount)
	if result != x {
		t.Errorf("Expected %v, got %v", x, result)
	}
}

func TestJitterFloat_NegativeAmount(t *testing.T) {
	x := 10.0
	amount := -1.0
	result := JitterFloat(x, amount)
	if result != x {
		t.Errorf("Expected %v, got %v", x, result)
	}
}

func TestJitterFloat_NonNegativeResult(t *testing.T) {
	x := 1.0
	amount := 2.0
	for i := 0; i < 100; i++ {
		result := JitterFloat(x, amount)
		if result < 0 {
			t.Errorf("Result should not be negative, got %v", result)
		}
	}
}

func TestJitterFloat_Range(t *testing.T) {
	x := 100.0
	amount := 0.5
	min := x - x*amount
	max := x + x*amount
	for i := 0; i < 100; i++ {
		result := JitterFloat(x, amount)
		if result < min || result > max {
			t.Errorf("Result %v out of expected range [%v, %v]", result, min, max)
		}
	}
}

func TestJitterInt_Basic(t *testing.T) {
	x := 100
	amount := 0.2
	for i := 0; i < 100; i++ {
		result := JitterInt(x, amount)
		// JitterInt returns x + int(jitter), so result can be less than x or more than x
		// But should not be negative
		if result < 0 {
			t.Errorf("Result should not be negative, got %v", result)
		}
	}
}

func TestJitterInt_ZeroAmount(t *testing.T) {
	x := 42
	amount := 0.0
	result := JitterInt(x, amount)
	if result != x {
		t.Errorf("Expected %v, got %v", x, result)
	}
}

func TestJitterInt_NegativeAmount(t *testing.T) {
	x := 42
	amount := -1.0
	result := JitterInt(x, amount)
	if result != x {
		t.Errorf("Expected %v, got %v", x, result)
	}
}

func TestJitterFloat_ZeroInput(t *testing.T) {
	x := 0.0
	amount := 1.0
	result := JitterFloat(x, amount)
	if result != 0.0 {
		t.Errorf("Expected 0.0, got %v", result)
	}
}

func TestJitterInt_ZeroInput(t *testing.T) {
	x := 0
	amount := 1.0
	result := JitterInt(x, amount)
	if result != 0 {
		t.Errorf("Expected 0, got %v", result)
	}
}

func TestJitterFloat_LargeAmount(t *testing.T) {
	x := 10.0
	amount := 10.0
	for i := 0; i < 100; i++ {
		result := JitterFloat(x, amount)
		if result < 0 {
			t.Errorf("Result should not be negative, got %v", result)
		}
		if math.IsNaN(result) || math.IsInf(result, 0) {
			t.Errorf("Result should not be NaN or Inf, got %v", result)
		}
	}
}
