package helper

import (
	"reflect"
	"testing"
)

func TestEnumerate(t *testing.T) {
	seq := func(yield func(string) bool) {
		for _, v := range []string{"a", "b", "c"} {
			if !yield(v) {
				return
			}
		}
	}

	var got [][2]any
	for i, v := range Enumerate(seq) {
		got = append(got, [2]any{i, v})
	}

	want := [][2]any{{0, "a"}, {1, "b"}, {2, "c"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected enumerate result: got %v, want %v", got, want)
	}
}

func TestEnumerateEarlyStop(t *testing.T) {
	seq := func(yield func(int) bool) {
		for _, v := range []int{10, 20, 30, 40} {
			if !yield(v) {
				return
			}
		}
	}

	count := 0
	for i, v := range Enumerate(seq) {
		_, _ = i, v

		if count++; count == 2 {
			break
		}
	}

	if count != 2 {
		t.Fatalf("expected early stop after 2 elements, got %d", count)
	}
}
