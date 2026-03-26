package vec

import "math"

type Vec struct {
	X float64
	Y float64
}

func (v Vec) Add(u Vec) Vec {
	return Vec{v.X + u.X, v.Y + u.Y}
}

func (v Vec) Sub(u Vec) Vec {
	return Vec{v.X - u.X, v.Y - u.Y}
}

func (v Vec) Mul(u Vec) Vec {
	return Vec{v.X * u.X, v.Y * u.Y}
}

func (v Vec) Scale(s float64) Vec {
	return Vec{v.X * s, v.Y * s}
}

func (v Vec) Length2() float64 {
	return v.X*v.X + v.Y*v.Y
}

func (v Vec) Length() float64 {
	return math.Sqrt(v.Length2())
}

func (v Vec) Normalize() Vec {
	n := v.Length()

	if n == 0 {
		return Vec{0, 0}
	}

	return v.Scale(1 / n)
}

func (v Vec) Dist2(u Vec) float64 {
	return v.Sub(u).Length2()
}

func (v Vec) Dist(u Vec) float64 {
	return math.Sqrt(v.Dist2(u))
}

func (v Vec) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}
