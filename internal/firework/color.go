package firework

import (
	"image/color"
	"math/rand/v2"
)

const (
	Black = iota
	Blue
	Cyan
	Green
	Magenta
	Orange
	Red
	Yellow
	White
)

func Palette() map[string]color.RGBA {
	return map[string]color.RGBA{
		"black":   lut[Black],
		"blue":    lut[Blue],
		"cyan":    lut[Cyan],
		"green":   lut[Green],
		"magenta": lut[Magenta],
		"orange":  lut[Orange],
		"red":     lut[Red],
		"yellow":  lut[Yellow],
		"white":   lut[White],
	}
}

var lut = []color.RGBA{
	{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}, // black
	{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}, // blue
	{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}, // cyan
	{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}, // green
	{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF}, // magenta
	{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}, // orange
	{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}, // red
	{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}, // yellow
	{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}, // white
}

func pick(c int) color.RGBA {
	if c < 0 || c >= len(lut) {
		return lut[Black]
	}

	return lut[c]
}

func SolidWhite() color.RGBA {
	return pick(White)
}

func SolidBlack() color.RGBA {
	return pick(Black)
}

func SolidRed() color.RGBA {
	return pick(Red)
}

func SolidOrange() color.RGBA {
	return pick(Orange)
}

func SolidYellow() color.RGBA {
	return pick(Yellow)
}

func SolidGreen() color.RGBA {
	return pick(Green)
}

func SolidCyan() color.RGBA {
	return pick(Cyan)
}

func SolidBlue() color.RGBA {
	return pick(Blue)
}

func SolidMagenta() color.RGBA {
	return pick(Magenta)
}

func SolidColor() color.RGBA {
	return pick(rand.IntN(len(lut)))
}
