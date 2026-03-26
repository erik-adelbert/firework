package helper

import "math/rand/v2"

func Pick[T any](s []T) (T, bool) {
	var zero T

	if len(s) == 0 {
		return zero, false
	}

	return s[rand.IntN(len(s))], true
}
