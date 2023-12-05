package sequence

import "github.com/cookieo9/sequence/tools"

// ToSlice is a utility method that calls the top level function version of
// [ToSlice].
func (s Sequence[T]) ToSlice() ([]T, error) {
	return ToSlice(s)
}

// Append processes the given sequence such that each item returned is
// appended to the provided destination slice. The resulting final
// slice is returned.
func Append[S ~[]T, T any](dst S, s Sequence[T]) (S, error) {
	err := EachSimple(s)(func(t T) bool {
		dst = append(dst, t)
		return true
	})
	return tools.CleanErrors(dst, err)
}

// ToSlice returns a new slice containing every item received from the given
// sequence, or any error that occured while doing so.
func ToSlice[T any](s Sequence[T]) ([]T, error) {
	return Append([]T(nil), s)
}

// FromSlice returns a sequence where the elements of the slice are returned.
// The source slice is reference from the sequence so certain changes to that
// slice may affect the sequence.
func FromSlice[T any](items []T) Sequence[T] {
	return Generate(func(f func(T) error) error {
		for _, x := range items {
			if err := f(x); err != nil {
				return err
			}
		}
		return nil
	})
}
