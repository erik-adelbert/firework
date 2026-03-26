package launcher

import (
	"image/color"
	"iter"
	"math"
	"strings"

	"github.com/erik-adelbert/firework/internal/sym"
	"github.com/erik-adelbert/firework/internal/vec"
)

func (l *Launcher) Render(screen []sym.Symbol, h, w int) string {
	// Rasterize the fireworks by iterating over all active particles and their trails,
	// and coloring the cells they pass through according to the particle's color and age.
	if len(screen) < h*w {
		return ""
	}

	idx := func(x cell) int {
		// Convert a cell coordinate to an index in the screen slice.
		return w*x.r + x.c
	}

	clear(screen)

	for i, fwk := range l.Actives {
		for j, trail := range l.AllActiveTrails(i) {
			var (
				ii   int
				ok   bool
				α, ω vec.Vec
			)

			p := fwk.Particles[j]
			next, close := iter.Pull2(trail)
			// Consecutive pairs of points in the trail are rendered as a line segment.
			// The first point is the current position of the particle
			_, α, ok = next()
			if !ok {
				close()
				continue
			}

			pointDrawn := false

			for n := p.TrailLen(); n > 0; {
				// The end point is the previous position of the particle at
				// some time in the past.
				if ii, ω, ok = next(); !ok {
					break
				}

				// Compute the age of the particle at the previous position, as a fraction
				// of the particle's trail length. This is used to determine the brightness of
				// the line segment.

				d := float64(n-ii) / float64(n)

				allCells := allDDA(α, ω)

				// The line segment is rendered by iterating over all cells it passes
				// through using DDA, and coloring each cell according to the particle's
				// color and the distance of the cell from the particle (cells closer
				// to the particle are brighter).
				for cell := range allInside(allCells, h, w) {
					rgba := fwk.Gradient(p)
					iii := idx(cell)

					if bright(screen[iii].RGBA) <= bright(rgba) {
						screen[iii] = sym.MkSymbol(getrune(p, d), rgba)

						pointDrawn = true
					}
				}

				// The next line segment starts where the previous one ended.
				α = ω
			}

			if !pointDrawn {
				c := cell{
					r: int(math.Round(α.Y)),
					c: int(math.Round(α.X)),
				}

				if c.inside(h, w) {
					iii := idx(c)

					if screen[iii].IsZero() {
						rgba := fwk.Gradient(p)
						screen[iii] = sym.MkSymbol(getrune(p, 0), rgba)
					}
				}
			}

			close()
		}
	}

	var sb strings.Builder
	sb.Grow(2 * h * w)

	for i := range h {
		lo, hi := i*w, (i+1)*w
		row := screen[lo:hi:hi]

		for _, s := range row {
			sb.WriteString(s.String())
		}

		sb.WriteByte('\n')
	}

	return sb.String()
}

func allInside(allCells iter.Seq[cell], h, w int) iter.Seq[cell] {
	return func(yield func(cell) bool) {
		for c := range allCells {
			if c.inside(h, w) && !yield(c) {
				return
			}
		}
	}
}

// bright calculates the brightness of a color as the average
// of its RGB components.
func bright(c color.RGBA) uint8 {
	return uint8((uint16(c.R) + uint16(c.G) + uint16(c.B)) / 3)
}
