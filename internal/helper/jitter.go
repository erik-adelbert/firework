package helper

import "math/rand/v2"

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

func JitterInt(x int, amount float64) int {
	return int(JitterFloat(float64(x), amount))
}
