package iterutils

import "iter"

type Nexter interface {
	Next() bool
}

// FromNexter converts a Nexter to an iter.Seq.
//
// Nexter is intended for converting traditional iteration objects, such as
// [database/sql.Rows] or [bufio.Scanner], into an iter.Seq. This function takes
// a Nexter `n` and a transformation function `fn`, applying `fn` to each element
// produced by `n` and yielding the results to the sequence.
//
// The iteration continues as long as `n.Next()` returns true and the `yield` function
// returns true. If `yield` returns false, the iteration stops early.
func FromNexter[T Nexter, U any](n T, fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for n.Next() {
			if !yield(fn(n)) {
				return
			}
		}
	}
}

// FromNexter2 converts a Nexter to an iter.Seq2.
//
// This function is similar to [FromNexter], but it allows the transformation
// function `fn` to return two values (`U` and `V`). These values are then yielded
// together to the sequence. This is useful when each iteration produces a pair
// of related values.
//
// The iteration continues as long as `n.Next()` returns true and the `yield` function
// returns true. If `yield` returns false, the iteration stops early.
func FromNexter2[T Nexter, U any, V any](n T, fn func(T) (U, V)) iter.Seq2[U, V] {
	return func(yield func(U, V) bool) {
		for n.Next() {
			if !yield(fn(n)) {
				return
			}
		}
	}
}

// NexterWithT is used to convert a type T that performs traditional iteration
// into a [Nexter].
//
// This interface provides a mechanism to wrap traditional iterators that work
// with a specific type T, allowing them to be used as Nexters.
type NexterWithT[T any] interface {
	Nexter
	T() T
}

type nexterWithT[T any] struct {
	t    T
	next func(T) bool
}

// NewNexter creates a new [NexterWithT] instance.
//
// This function accepts an initial value of type T and a `next` function that
// determines the continuation of the iteration. The returned [NexterWithT] can
// be used to iterate over the value of type T, with the current value accessible
// via the T() method.
func NewNexterWithT[T any](t T, next func(T) bool) NexterWithT[T] {
	return &nexterWithT[T]{
		t:    t,
		next: next,
	}
}

func (n *nexterWithT[T]) Next() bool {
	return n.next(n.t)
}

func (n *nexterWithT[T]) T() T {
	return n.t
}

// FromTryNexter2 can be applied to types that return both a value and an error
// during iteration, such as [encoding/csv.Reader].
//
// This function converts a [TryNexter] to an [iter.Seq2], allowing you to iterate
// over a sequence of two values. The `fn` function takes the [TryNexter] `t`,
// the value `v`, and the error `err` returned by `t.Next()`, and it produces
// two values (`U` and `W`) that are yielded to the sequence.
//
// The iteration continues until the yield function returns false.
func FromTryNexter2[V any, T TryNexter[V], U any, W any](t T, fn func(T, V, error) (U, W)) iter.Seq2[U, W] {
	return func(yield func(U, W) bool) {
		for {
			v, err := t.Next()
			if !yield(fn(t, v, err)) {
				return
			}
		}
	}
}

// TryNexter is an interface for iterators that return both a value of type V
// and an error during iteration.
//
// This interface is commonly used for iterators where each iteration may
// result in an error, such as reading from a file or a network stream.
type TryNexter[V any] interface {
	Next() (V, error)
}

// TryNexterWithT is an interface that is created exclusively through the
// [NewTryNexterWithT] function.
//
// This interface is designed to wrap existing iteration types that return
// both a value and an error, allowing them to be used with [FromTryNexter2].
// In addition to the iteration, it provides access to the current state
// or context of type T via the T() method.
type TryNexterWithT[T any, V any] interface {
	TryNexter[V]
	T() T
}

type tryNexterWithT[T any, V any] struct {
	t    T
	next func(T) (V, error)
}

// NewTryNexterWithT creates a new [TryNexterWithT] instance.
//
// This function takes an initial value of type T and a `next` function that
// controls the iteration. The `next` function returns a value of type V and
// an error. The returned [TryNexterWithT] allows iteration with error handling,
// while also providing access to the current state via the T() method.
func NewTryNexterWithT[T any, V any](t T, next func(T) (V, error)) TryNexterWithT[T, V] {
	return &tryNexterWithT[T, V]{
		t:    t,
		next: next,
	}
}

func (n *tryNexterWithT[T, V]) Next() (V, error) {
	return n.next(n.t)
}

func (n *tryNexterWithT[T, V]) T() T {
	return n.t
}
