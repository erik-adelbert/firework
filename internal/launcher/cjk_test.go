package launcher

import (
	"strings"
	"testing"

	"github.com/erik-adelbert/firework/internal/particle"
)

func TestRandomRuneEmptyString(t *testing.T) {
	result := randomRune([]rune(""))

	if result != ' ' {
		t.Fatalf("expected space for empty string, got %q", result)
	}
}

func TestRandomRune(t *testing.T) {
	tests := []struct {
		name       string
		input      []rune
		iterations int
	}{
		{"ASCII", []rune("hello"), 20},
		{"CJK Characters", []rune("龖龠龜"), 12},
		{"Single Rune", []rune("a"), 4},
		{"Multibye UTF-8", []rune("你好世界"), 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for range tt.iterations {
				result := randomRune(tt.input)

				if !strings.ContainsRune(string(tt.input), result) {
					t.Fatalf("expected rune from input %q, got %q", tt.input, result)
				}
			}
		})
	}
}

func TestRandomRuneRandomness(t *testing.T) {
	input := []rune("abc")

	counts := map[rune]int{
		'a': 0,
		'b': 0,
		'c': 0,
	}

	const draws = 3000

	for range draws {
		counts[randomRune(input)]++
	}

	const minPerRune = 600
	const maxPerRune = 1400

	for r, c := range counts {
		switch {
		case c == 0:
			t.Fatalf("rune %q was never selected in %d draws", r, draws)
		case c < minPerRune || c > maxPerRune:
			t.Fatalf("rune %q count out of expected range: got %d, want [%d,%d]", r, c, minPerRune, maxPerRune)
		}
	}
}

func TestPaletteStructure(t *testing.T) {
	expectedStates := map[particle.State]bool{
		particle.Spawning: true,
		particle.Active:   true,
		particle.Fading:   true,
		particle.Dead:     true,
	}

	if len(palettes) != len(expectedStates) {
		t.Fatalf("expected palette length %d, got %d", len(expectedStates), len(palettes))
	}

	for state := range expectedStates {
		if state >= particle.State(len(palettes)) {
			t.Fatalf("palette missing entry for state %d", state)
		}

		clutList := palettes[state]
		if len(clutList) == 0 {
			t.Fatalf("palette entry for state %d is empty", state)
		}

		for i, c := range clutList {
			if c.thresh <= 0 || c.thresh > 1 {
				t.Fatalf("state %d, clut %d: invalid threshold %f (should be 0 < thresh <= 1)", state, i, c.thresh)
			}

			if len(c.runes) == 0 {
				t.Fatalf("state %d, clut %d: empty chars string", state, i)
			}

			if i > 0 {
				prevThresh := clutList[i-1].thresh
				if c.thresh <= prevThresh {
					t.Fatalf("state %d: thresholds should be increasing, got %f after %f", state, c.thresh, prevThresh)
				}
			}
		}
	}
}

func TestPaletteThresholdMonotonicity(t *testing.T) {
	for state, clutList := range palettes {
		for i := 1; i < len(clutList); i++ {
			if clutList[i].thresh <= clutList[i-1].thresh {
				t.Fatalf(
					"state %d: threshold at index %d (%f) not greater than index %d (%f)",
					state, i, clutList[i].thresh, i-1, clutList[i-1].thresh,
				)
			}
		}
	}
}

func BenchmarkRandomRune(b *testing.B) {
	input := []rune("龖龠龜时中自字木月日目火田左右点以")

	for b.Loop() {
		randomRune(input)
	}
}
