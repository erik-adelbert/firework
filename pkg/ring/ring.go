// Package ring implements a fixed-size circular buffer.
package ring

import "iter"

type Ring[T any] struct {
	buf  []T
	head int
	size int
}

func NewRing[T any](cap int) *Ring[T] {
	return &Ring[T]{buf: make([]T, cap)}
}

func (r *Ring[T]) RepeatFront(x T, n int) int {
	for range n {
		r.PushFront(x)
	}

	return n
}

func (r *Ring[T]) RepeatBack(x T, n int) int {
	for range n {
		r.PushBack(x)
	}

	return n
}

func (r *Ring[T]) PushBack(v T) {
	if r.size < len(r.buf) {
		r.buf[(r.head+r.size)%len(r.buf)] = v
		r.size++
	} else {
		r.buf[r.head] = v
		r.head = (r.head + 1) % len(r.buf)
	}
}

func (r *Ring[T]) PushFront(v T) {
	if r.size < len(r.buf) {
		r.head = (r.head - 1 + len(r.buf)) % len(r.buf)
		r.buf[r.head] = v
		r.size++
	} else {
		r.head = (r.head - 1 + len(r.buf)) % len(r.buf)
		r.buf[r.head] = v
	}
}

func (r *Ring[T]) PopBack() {
	if r.size > 0 {
		r.size--
	}
}

func (r *Ring[T]) PopFront() {
	if r.size > 0 {
		r.head = (r.head + 1) % len(r.buf)
		r.size--
	}
}

func (r *Ring[T]) AllFront() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range r.size {
			if !yield(i, r.buf[(r.head+i)%len(r.buf)]) {
				return
			}
		}
	}
}

func (r *Ring[T]) AllBack() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := range r.size {
			if !yield(i, r.buf[(r.head+r.size-1-i+len(r.buf))%len(r.buf)]) {
				return
			}
		}
	}
}

func (r *Ring[T]) Clear() {
	r.head = 0
	r.size = 0
}

func (r *Ring[T]) At(i int) T {
	if i < 0 || i >= r.size {
		panic("index out of bounds")
	}

	return r.buf[(r.head+i)%len(r.buf)]
}

func (r *Ring[T]) Len() int {
	return r.size
}

func (r *Ring[T]) Cap() int {
	return len(r.buf)
}
