package launcher

import (
	"slices"
	"testing"

	"github.com/erik-adelbert/firework/internal/vec"
)

func TestAllDDAZeroLength(t *testing.T) {
	start := vec.Vec{X: 2.0, Y: 3.0}
	cells := slices.Collect(allDDA(start, start))

	if len(cells) != 1 {
		t.Fatalf("expected exactly one cell for zero-length line, got %d", len(cells))
	}

	if cells[0] != (cell{c: 2, r: 3}) {
		t.Fatalf("expected only cell {2,3}, got {%d,%d}", cells[0].c, cells[0].r)
	}
}

func TestAllDDAVerticalLine(t *testing.T) {
	start := vec.Vec{X: 2.0, Y: 1.0}
	end := vec.Vec{X: 2.0, Y: 4.0}
	cells := slices.Collect(allDDA(start, end))

	if len(cells) < 2 {
		t.Fatalf("expected multiple cells for vertical line, got %d", len(cells))
	}

	if cells[0] != (cell{c: 2, r: 1}) {
		t.Fatalf("expected first cell {2,1}, got {%d,%d}", cells[0].c, cells[0].r)
	}

	last := cells[len(cells)-1]
	if last != (cell{c: 2, r: 4}) {
		t.Fatalf("expected last cell {2,4}, got {%d,%d}", last.c, last.r)
	}

	for i, c := range cells {
		if c.c != 2 {
			t.Fatalf("expected all x=2 for vertical line, got x=%d at index %d", c.c, i)
		}
	}
}

func TestAllDDAHorizontalLine(t *testing.T) {
	start := vec.Vec{X: 1.0, Y: 3.0}
	end := vec.Vec{X: 4.0, Y: 3.0}
	cells := slices.Collect(allDDA(start, end))

	if len(cells) < 2 {
		t.Fatalf("expected multiple cells for horizontal line, got %d", len(cells))
	}

	if cells[0] != (cell{c: 1, r: 3}) {
		t.Fatalf("expected first cell {1,3}, got {%d,%d}", cells[0].c, cells[0].r)
	}

	last := cells[len(cells)-1]
	if last != (cell{c: 4, r: 3}) {
		t.Fatalf("expected last cell {4,3}, got {%d,%d}", last.c, last.r)
	}

	for i, c := range cells {
		if c.r != 3 {
			t.Fatalf("expected all y=3 for horizontal line, got y=%d at index %d", c.r, i)
		}
	}
}

func TestAllDDADiagonalLine(t *testing.T) {
	start := vec.Vec{X: 0.0, Y: 0.0}
	end := vec.Vec{X: 3.0, Y: 3.0}
	cells := slices.Collect(allDDA(start, end))

	if len(cells) < 2 {
		t.Fatalf("expected multiple cells for diagonal line, got %d", len(cells))
	}

	if cells[0] != (cell{c: 0, r: 0}) {
		t.Fatalf("expected first cell {0,0}, got {%d,%d}", cells[0].c, cells[0].r)
	}

	last := cells[len(cells)-1]
	if last != (cell{c: 3, r: 3}) {
		t.Fatalf("expected last cell {3,3}, got {%d,%d}", last.c, last.r)
	}

	for i := 1; i < len(cells); i++ {
		if cells[i].c < cells[i-1].c || cells[i].r < cells[i-1].r {
			t.Fatalf("expected non-decreasing path, got {%d,%d} after {%d,%d}", cells[i].c, cells[i].r, cells[i-1].c, cells[i-1].r)
		}
	}
}
