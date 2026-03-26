package launcher

import (
	"iter"
	"math"

	"github.com/erik-adelbert/firework/internal/vec"
)

type cell struct {
	r, c int
}

func (c cell) inside(h, w int) bool {
	return 0 <= c.r && c.r < h && 0 <= c.c && c.c < w
}

func allDDA(α, ω vec.Vec) iter.Seq[cell] {
	return func(yield func(cell) bool) {
		const Step = 0.2

		δ := ω.Sub(α)
		slope := math.Abs(δ.Y / δ.X)

		var dx, dy float64

		switch {
		case α.X < ω.X:
			dx = 1
		case α.X > ω.X:
			dx = -1
		}

		switch {
		case α.Y < ω.Y:
			dy = 1
		case α.Y > ω.Y:
			dy = -1
		}

		cur := cell{
			r: int(math.Round(α.Y)),
			c: int(math.Round(α.X)),
		}

		if !yield(cur) {
			return
		}

		ds := α.Dist2(ω)

		for p := α; p.Dist2(ω) <= ds; {
			nxt := cell{
				r: int(math.Round(p.Y)),
				c: int(math.Round(p.X)),
			}

			if cur != nxt {
				cur = nxt

				if !yield(cur) {
					return
				}

				ds = p.Dist2(ω)
			}

			if slope <= 1 {
				p.X += dx * Step
				p.Y += dy * math.Abs(Step*slope)
			} else {
				p.X += dx * Step / slope
				p.Y += dy * Step
			}
		}
	}
}
