package fireworks

import (
	fw "github.com/erik-adelbert/firework/internal/firework"
	"github.com/erik-adelbert/firework/internal/vec"
)

func Catalog() map[string]fw.NewFirework {
	return map[string]fw.NewFirework{
		"zero":    NewZero,
		"one":     NewOne,
		"two":     NewTwo,
		"three":   NewThree,
		"chrys":   NewChrysanthemum,
		"comet":   NewComet,
		"feather": NewFeather,
		"fish":    NewFish,
		"glitter": NewGlitter,
		"kamuro":  NewKamuro,
		"laser":   NewLaser,
		"palm":    NewPalm,
		"peony":   NewPeony,
		"saturn":  NewSaturn,
		"sphere":  NewSphere,
		"sun":     NewSun,
		"willow":  NewWillow,
	}
}

func NewZero(o vec.Vec) *fw.Firework {
	return NewDigit(o, 0)
}

func NewOne(o vec.Vec) *fw.Firework {
	return NewDigit(o, 1)
}

func NewTwo(o vec.Vec) *fw.Firework {
	return NewDigit(o, 2)
}

func NewThree(o vec.Vec) *fw.Firework {
	return NewDigit(o, 3)
}
