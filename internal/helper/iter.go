package helper

import "iter"

func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range seq {
			if !yield(i, v) {
				return
			}

			i++
		}
	}
}

func EmptySeq[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}
