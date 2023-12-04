package extract

import (
	"github.com/cookieo9/sequence"
)

// Append processes the given sequence such that each item returned is
// appended to the provided destination slice. The resulting final
// slice is returned.
func Append[S ~[]T, T any](dst S, s sequence.Sequence[T]) (S, error) {
	return Collect(s, dst, func(t T, s S) S {
		return append(s, t)
	})
}

// ToSlice returns a new slice containing every item received from the given
// sequence. If the sequence has an exact size, it will be used to construct
// a slice of the correct size.
func ToSlice[T any](s sequence.Sequence[T]) ([]T, error) {
	return Append([]T(nil), s)
}
