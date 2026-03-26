package helper

import "math/rand/v2"

func JitterInt(x int, amount float64) int {
	if amount <= 0 {
		return x
	}

	j := int(float64(x) * amount * (rand.Float64()*2 - 1))

	if x+j < 0 {
		return x
	}

	return x + j
}

func JitterFloat(x, amount float64) float64 {
	if amount <= 0 {
		return x
	}

	j := x * amount * (rand.Float64()*2 - 1)

	if x+j < 0 {
		return x
	}

	return x + j
}
