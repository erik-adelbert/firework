package fireworks

import (
	"image/color"
	"iter"
	"math"
	"math/rand/v2"
	"time"

	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/helper"
	"github.com/erik-adelbert/firework/internal/particle"
	"github.com/erik-adelbert/firework/internal/vec"
)

const (
	digitSpawningDuration = 220 * time.Millisecond
	digitActiveDuration   = 760 * time.Millisecond
)

type segment struct {
	a vec.Vec
	b vec.Vec
}

func NewDigit(o vec.Vec, digit int) *fw.Firework {
	if digit < 0 || digit > 3 {
		digit = 0
	}

	ttl := time.Duration(helper.JitterInt(2400, 1./8)) * time.Millisecond
	spawnAfter := time.Duration(helper.JitterInt(25, 1.)) * time.Millisecond

	return fw.New(
		o,
		fw.MkForces(0.42, 0.12, 0, nil),
		fw.MkShell(0, NewDigitSpawner(digit)),
		fw.MkTiming(ttl, spawnAfter),
	)
}

type digitSpawner struct {
	LUT
	digit int
}

func NewDigitSpawner(digit int) *digitSpawner {
	return &digitSpawner{LUT: PickLut(), digit: digit}
}

func (s *digitSpawner) Spawn(o vec.Vec, t time.Time, dt time.Duration) iter.Seq[*particle.Particle] {
	return func(yield func(*particle.Particle) bool) {
		points := sampleDigitPoints(s.digit, 22)
		if len(points) == 0 {
			return
		}

		scale := helper.JitterFloat(58, 0.08)
		travelSeconds := helper.JitterFloat(0.78, 0.1)

		for _, p0 := range points {
			target := p0.Scale(scale)
			v := target.Scale(1 / travelSeconds)
			v = v.Scale(helper.JitterFloat(1, 0.12))

			ttl := time.Duration(helper.JitterInt(2300, 3./23)) * time.Millisecond
			trail := helper.JitterInt(24, 1./6)
			phases := particle.PhaseTiming{
				SpawningEnd: float64(digitSpawningDuration) / float64(ttl),
				ActiveEnd:   float64(digitActiveDuration) / float64(ttl),
			}

			p := particle.NewWithPhaseTiming(o, v, vec.Vec{}, s.PickColor(), ttl, trail, phases)
			if !yield(p) {
				return
			}
		}
	}
}

func (s *digitSpawner) Gradient(p *particle.Particle) color.RGBA {
	g := PhaseGradient2(p.Life())
	c := p.Color()

	return scaleColor(c, g)
}

func sampleDigitPoints(digit int, density float64) []vec.Vec {
	segments := digitSegments(digit)
	points := make([]vec.Vec, 0, 96)

	for _, seg := range segments {
		delta := seg.b.Sub(seg.a)
		length := delta.Length()
		if length == 0 {
			continue
		}

		count := int(math.Ceil(length * density))
		count = max(count, 2)

		normal := vec.Vec{X: -delta.Y, Y: delta.X}.Normalize()
		for i := range count {
			t := 0.0
			if count > 1 {
				t = float64(i) / float64(count-1)
			}

			p := seg.a.Add(delta.Scale(t))

			lineJitter := rand.NormFloat64() * 0.045
			tangentJitter := rand.NormFloat64() * 0.015
			p = p.Add(normal.Scale(lineJitter)).Add(delta.Normalize().Scale(tangentJitter))

			points = append(points, p)
		}
	}

	return points
}

func digitSegments(digit int) []segment {
	switch digit {
	case 1:
		return []segment{
			{a: vec.Vec{X: -0.24, Y: -0.76}, b: vec.Vec{X: 0.00, Y: -1.00}},
			{a: vec.Vec{X: 0.00, Y: -1.00}, b: vec.Vec{X: 0.00, Y: 1.00}},
			{a: vec.Vec{X: -0.33, Y: 1.00}, b: vec.Vec{X: 0.33, Y: 1.00}},
		}
	case 2:
		return []segment{
			{a: vec.Vec{X: -0.58, Y: -0.95}, b: vec.Vec{X: 0.54, Y: -0.95}},
			{a: vec.Vec{X: 0.54, Y: -0.95}, b: vec.Vec{X: 0.66, Y: -0.58}},
			{a: vec.Vec{X: 0.66, Y: -0.58}, b: vec.Vec{X: -0.50, Y: 0.28}},
			{a: vec.Vec{X: -0.50, Y: 0.28}, b: vec.Vec{X: -0.65, Y: 1.00}},
			{a: vec.Vec{X: -0.65, Y: 1.00}, b: vec.Vec{X: 0.68, Y: 1.00}},
		}
	case 3:
		return []segment{
			{a: vec.Vec{X: -0.56, Y: -0.92}, b: vec.Vec{X: 0.44, Y: -0.92}},
			{a: vec.Vec{X: 0.44, Y: -0.92}, b: vec.Vec{X: 0.66, Y: -0.60}},
			{a: vec.Vec{X: 0.66, Y: -0.60}, b: vec.Vec{X: 0.20, Y: -0.14}},
			{a: vec.Vec{X: 0.20, Y: -0.14}, b: vec.Vec{X: 0.68, Y: 0.30}},
			{a: vec.Vec{X: 0.68, Y: 0.30}, b: vec.Vec{X: 0.44, Y: 0.90}},
			{a: vec.Vec{X: 0.44, Y: 0.90}, b: vec.Vec{X: -0.56, Y: 0.90}},
		}
	default:
		return []segment{
			{a: vec.Vec{X: -0.46, Y: -1.00}, b: vec.Vec{X: 0.46, Y: -1.00}},
			{a: vec.Vec{X: 0.46, Y: -1.00}, b: vec.Vec{X: 0.68, Y: -0.74}},
			{a: vec.Vec{X: 0.68, Y: -0.74}, b: vec.Vec{X: 0.68, Y: 0.74}},
			{a: vec.Vec{X: 0.68, Y: 0.74}, b: vec.Vec{X: 0.46, Y: 1.00}},
			{a: vec.Vec{X: 0.46, Y: 1.00}, b: vec.Vec{X: -0.46, Y: 1.00}},
			{a: vec.Vec{X: -0.46, Y: 1.00}, b: vec.Vec{X: -0.68, Y: 0.74}},
			{a: vec.Vec{X: -0.68, Y: 0.74}, b: vec.Vec{X: -0.68, Y: -0.74}},
			{a: vec.Vec{X: -0.68, Y: -0.74}, b: vec.Vec{X: -0.46, Y: -1.00}},
		}
	}
}
