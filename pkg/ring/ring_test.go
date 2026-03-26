package ring

import (
	"testing"
)

func TestNewRing(t *testing.T) {
	r := NewRing[int](5)

	if r.Len() != 0 {
		t.Errorf("expected length 0, got %d", r.Len())
	}

	if cap(r.buf) != 5 {
		t.Errorf("expected capacity 5, got %d", cap(r.buf))
	}
}

func TestNewRingZeroCapacity(t *testing.T) {
	r := NewRing[int](0)

	if r.Cap() != 0 {
		t.Errorf("expected capacity 0 for zero capacity, got %v", r.Cap())
	}
}

func TestPushBack(t *testing.T) {
	r := NewRing[int](3)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	if r.Len() != 3 {
		t.Errorf("expected length 3, got %d", r.Len())
	}

	if r.At(0) != 1 || r.At(1) != 2 || r.At(2) != 3 {
		t.Errorf("unexpected values: %d, %d, %d", r.At(0), r.At(1), r.At(2))
	}

	// Test overflow
	r.PushBack(4)
	if r.Len() != 3 {
		t.Errorf("expected length 3, got %d", r.Len())
	}

	if r.At(0) != 2 || r.At(1) != 3 || r.At(2) != 4 {
		t.Errorf("unexpected values after overflow: %d, %d, %d", r.At(0), r.At(1), r.At(2))
	}
}

func TestPushFront(t *testing.T) {
	r := NewRing[int](3)

	r.PushFront(1)
	r.PushFront(2)
	r.PushFront(3)

	if r.Len() != 3 {
		t.Errorf("expected length 3, got %d", r.Len())
	}

	if r.At(0) != 3 || r.At(1) != 2 || r.At(2) != 1 {
		t.Errorf("unexpected values: %d, %d, %d", r.At(0), r.At(1), r.At(2))
	}

	// Test overflow
	r.PushFront(4)
	if r.At(0) != 4 || r.At(1) != 3 || r.At(2) != 2 {
		t.Errorf("unexpected values after overflow: %d, %d, %d", r.At(0), r.At(1), r.At(2))
	}
}

func TestPopBack(t *testing.T) {
	r := NewRing[int](3)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	r.PopBack()
	if r.Len() != 2 {
		t.Errorf("expected length 2, got %d", r.Len())
	}

	if r.At(0) != 1 || r.At(1) != 2 {
		t.Errorf("unexpected values: %d, %d", r.At(0), r.At(1))
	}

	// Pop from empty
	r.PopBack()
	r.PopBack()
	r.PopBack()

	if r.Len() != 0 {
		t.Errorf("expected length 0, got %d", r.Len())
	}
}

func TestPopFront(t *testing.T) {
	r := NewRing[int](3)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	r.PopFront()
	if r.Len() != 2 {
		t.Errorf("expected length 2, got %d", r.Len())
	}

	if r.At(0) != 2 || r.At(1) != 3 {
		t.Errorf("unexpected values: %d, %d", r.At(0), r.At(1))
	}
}

func TestRepeatFront(t *testing.T) {
	r := NewRing[int](5)
	n := r.RepeatFront(42, 3)

	if n != 3 {
		t.Errorf("expected 3 pushes, got %d", n)
	}

	if r.Len() != 3 {
		t.Errorf("expected length 3, got %d", r.Len())
	}

	for i := range 3 {
		if r.At(i) != 42 {
			t.Errorf("expected 42 at index %d, got %d", i, r.At(i))
		}
	}
}

func TestRepeatBack(t *testing.T) {
	r := NewRing[int](5)
	n := r.RepeatBack(42, 3)

	if n != 3 {
		t.Errorf("expected 3 pushes, got %d", n)
	}

	if r.Len() != 3 {
		t.Errorf("expected length 3, got %d", r.Len())
	}

	for i := range 3 {
		if r.At(i) != 42 {
			t.Errorf("expected 42 at index %d, got %d", i, r.At(i))
		}
	}
}

func TestAllFront(t *testing.T) {
	r := NewRing[int](5)
	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	idx, expected := 0, []int{1, 2, 3}
	for i, v := range r.AllFront() {
		if i != idx || v != expected[idx] {
			t.Errorf("expected index %d value %d, got index %d value %d", idx, expected[idx], i, v)
		}
		idx++
	}
}

func TestAllBack(t *testing.T) {
	r := NewRing[int](5)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	idx, expected := 0, []int{3, 2, 1}
	for i, v := range r.AllBack() {
		if i != idx || v != expected[idx] {
			t.Errorf("expected index %d value %d, got index %d value %d", idx, expected[idx], i, v)
		}
		idx++
	}
}

func TestClear(t *testing.T) {
	r := NewRing[int](5)

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)

	r.Clear()
	if r.Len() != 0 {
		t.Errorf("expected length 0 after clear, got %d", r.Len())
	}

	if r.head != 0 {
		t.Errorf("expected head 0 after clear, got %d", r.head)
	}
}

func TestAt(t *testing.T) {
	r := NewRing[string](3)

	r.PushBack("a")
	r.PushBack("b")
	r.PushBack("c")

	if r.At(0) != "a" || r.At(1) != "b" || r.At(2) != "c" {
		t.Errorf("unexpected values: %s, %s, %s", r.At(0), r.At(1), r.At(2))
	}
}

func BenchmarkPushBack(b *testing.B) {
	r := NewRing[int](1000)

	for b.Loop() {
		r.PushBack(42)
	}
}

func BenchmarkPushFront(b *testing.B) {
	r := NewRing[int](1000)

	for b.Loop() {
		r.PushFront(42)
	}
}

func BenchmarkPopBack(b *testing.B) {
	r := NewRing[int](1000)

	for i := range 1000 {
		r.PushBack(i)
	}

	for b.Loop() {
		r.PopBack()

		if r.Len() == 0 {
			r.size = 1000
		}
	}
}

func BenchmarkPopFront(b *testing.B) {
	r := NewRing[int](1000)

	for i := range 1000 {
		r.PushBack(i)
	}

	for b.Loop() {
		r.PopFront()

		if r.Len() == 0 {
			r.size = 1000
		}
	}
}

var sink int

func BenchmarkAt(b *testing.B) {
	r := NewRing[int](100)

	for i := range 100 {
		r.PushBack(i)
	}

	b.ResetTimer()
	for i := range b.N {
		sink = r.At(i % 100)
	}
}

func BenchmarkLen(b *testing.B) {
	r := NewRing[int](100)
	for i := range 100 {
		r.PushBack(i)
	}

	for b.Loop() {
		sink = r.Len()
	}
}

func BenchmarkClear(b *testing.B) {
	r := NewRing[int](1000)

	for range 1000 {
		r.PushBack(42)
	}

	for b.Loop() {
		r.Clear()
		r.size = 1000 // fake refill
	}
}
