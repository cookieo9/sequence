package slices

import (
	"github.com/cookieo9/sequence"
	"github.com/cookieo9/sequence/single"
)

// AsSequence returns a sequence where the elements of the slice are returned.
// The source slice is reference from the sequence so certain changes to the
// slice may affect the sequence.
func AsSequence[T any](items []T) sequence.Sequence[T] {
	return sequence.Generate(func(f func(T) error) error {
		for _, x := range items {
			if err := f(x); err != nil {
				return err
			}
		}
		return nil
	})
}

// Append processes the given sequence such that each item returned is
// appended to the provided destination slice. The resulting final
// slice is returned.
func Append[S ~[]T, T any](dst S, s sequence.Sequence[T]) (S, error) {
	return single.Collect(s, dst, func(t T, s S) S {
		return append(s, t)
	})
}

// Collect returns a new slice containing every item received from the given
// sequence. If the sequence has an exact size, it will be used to construct
// a slice of the correct size.
func Collect[T any](s sequence.Sequence[T]) ([]T, error) {
	return Append([]T(nil), s)
}
