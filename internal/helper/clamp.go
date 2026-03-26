package helper

import "cmp"

func Clamp[T cmp.Ordered](x, a, b T) T {
	switch {
	case x < a:
		return a
	case x > b:
		return b
	}

	return x
}
