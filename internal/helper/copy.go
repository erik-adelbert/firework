package helper

// Copy copies elements from src to dst. If dst is too small, it will be extended.
func Copy[T any](dst *[]T, src []T) {
	δ := len(src) - len(*dst)

	if δ > 0 {
		*dst = append(*dst, make([]T, δ)...)
	}

	copy(*dst, src)
}
