package firework

import (
	"image/color"
	"testing"
)

func TestPickNegativeIndex(t *testing.T) {
	result := pick(-1)
	expected := lut[Black]
	if result != expected {
		t.Errorf("pick(-1) = %v, expected %v", result, expected)
	}
}

func TestPickOutOfBounds(t *testing.T) {
	result := pick(len(lut))
	expected := lut[Black]
	if result != expected {
		t.Errorf("pick(%d) = %v, expected %v", len(lut), result, expected)
	}

	result = pick(999)
	if result != expected {
		t.Errorf("pick(999) = %v, expected %v", result, expected)
	}
}

func TestSolidColors(t *testing.T) {
	tests := []struct {
		name     string
		fn       func() color.RGBA
		expected color.RGBA
	}{
		{"White", SolidWhite, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}},
		{"Black", SolidBlack, color.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}},
		{"Red", SolidRed, color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF}},
		{"Orange", SolidOrange, color.RGBA{R: 0xFF, G: 0xA5, B: 0x00, A: 0xFF}},
		{"Yellow", SolidYellow, color.RGBA{R: 0xFF, G: 0xFF, B: 0x00, A: 0xFF}},
		{"Green", SolidGreen, color.RGBA{R: 0x00, G: 0xFF, B: 0x00, A: 0xFF}},
		{"Cyan", SolidCyan, color.RGBA{R: 0x00, G: 0xFF, B: 0xFF, A: 0xFF}},
		{"Blue", SolidBlue, color.RGBA{R: 0x00, G: 0x00, B: 0xFF, A: 0xFF}},
		{"Magenta", SolidMagenta, color.RGBA{R: 0xFF, G: 0x00, B: 0xFF, A: 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn()
			if result != tt.expected {
				t.Errorf("%s() = %v, expected %v", tt.name, result, tt.expected)
			}
		})
	}
}
