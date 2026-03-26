package fireworks

import (
	"iter"
	"math"
	"math/rand/v2"

	"github.com/erik-adelbert/firework/internal/vec"
)

func AllUniformDisk(ρ float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		for i := 0; i < n; {
			v := vec.Vec{
				X: 2*ρ*rand.Float64() - ρ,
				Y: 2*ρ*rand.Float64() - ρ,
			}

			switch {
			case v.Length2() > ρ*ρ:
				continue
			case !yield(v):
				return
			}

			i++
		}
	}
}

func AllUniformFan(ρ, α, ω float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		for i := 0; i < n; {
			v := vec.Vec{
				X: 2*ρ*rand.Float64() - ρ,
				Y: 2*ρ*rand.Float64() - ρ,
			}

			θ := v.Angle()
			if θ >= α && θ <= ω && v.Length2() <= ρ*ρ {
				v.Y = -v.Y

				if !yield(v) {
					return
				}

				i++
			}
		}
	}

}

func AllUniformCircle(ρ float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		for range n {
			θ := 2 * math.Pi * rand.Float64()

			v := vec.Vec{
				X: ρ * math.Cos(θ),
				Y: -ρ * math.Sin(θ),
			}

			if !yield(v) {
				return
			}
		}
	}
}

func AllUniformCircleBalanced(ρ float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		for range n / 2 {
			θ := 2 * math.Pi * rand.Float64()

			v := vec.Vec{
				X: ρ * math.Cos(θ),
				Y: -ρ * math.Sin(θ),
			}

			if !yield(v) {
				return
			}

			if !yield(v.Scale(-1)) {
				return
			}
		}

		if n%2 == 1 {
			θ := 2 * math.Pi * rand.Float64()

			v := vec.Vec{
				X: ρ * math.Cos(θ),
				Y: -ρ * math.Sin(θ),
			}

			if !yield(v) {
				return
			}
		}
	}
}

// AllStratifiedCircle divides the full circle into n equal sectors and
// samples one uniformly random angle per sector, guaranteeing even angular
// coverage regardless of particle count.
func AllStratifiedCircle(ρ float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		slice := 2 * math.Pi / float64(n)

		for i := range n {
			θ := (float64(i) + rand.Float64()) * slice

			v := vec.Vec{
				X: ρ * math.Cos(θ),
				Y: -ρ * math.Sin(θ),
			}

			if !yield(v) {
				return
			}
		}
	}
}

func AllUniformArc(ρ, α, ω float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		for range n {
			θ := α + (ω-α)*rand.Float64()

			v := vec.Vec{
				X: ρ * math.Cos(θ),
				Y: -ρ * math.Sin(θ),
			}

			if !yield(v) {
				return
			}
		}
	}
}

func AllNormalDisk(ρ, σ float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		i := 0
		for i < n {
			v := vec.Vec{
				X: normFloat64(0, σ),
				Y: normFloat64(0, σ),
			}

			if v.X > ρ || v.X < -ρ || v.Y > ρ || v.Y < -ρ {
				continue
			}

			if v.Length2() <= ρ*ρ {
				if !yield(v) {
					return
				}

				i++
			}
		}
	}
}

func AllNormalDisk9(ρ float64, n int) iter.Seq[vec.Vec] {
	return AllNormalDisk(ρ, ρ/9, n)
}

func AllNormalDiskBalanced(ρ, σ float64, n int) iter.Seq[vec.Vec] {
	return func(yield func(vec.Vec) bool) {
		for i := 0; i+1 < n; i += 2 {
			v := sampleNormalDisk(ρ, σ)

			if !yield(v) {
				return
			}

			if !yield(v.Scale(-1)) {
				return
			}
		}

		if n%2 == 1 {
			if !yield(sampleNormalDisk(ρ, σ)) {
				return
			}
		}
	}
}

func AllNormalDisk9Balanced(ρ float64, n int) iter.Seq[vec.Vec] {
	return AllNormalDiskBalanced(ρ, ρ/9, n)
}

func sampleNormalDisk(ρ, σ float64) vec.Vec {
	for {
		v := vec.Vec{
			X: normFloat64(0, σ),
			Y: normFloat64(0, σ),
		}

		if v.X > ρ || v.X < -ρ || v.Y > ρ || v.Y < -ρ {
			continue
		}

		if v.Length2() <= ρ*ρ {
			return v
		}
	}
}

func normFloat64(mean, σ float64) float64 {
	return mean + rand.NormFloat64()*σ
}
