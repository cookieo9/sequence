package build

import "github.com/cookieo9/sequence"

// FromSlice returns a sequence where the elements of the slice are returned.
// The source slice is reference from the sequence so certain changes to that
// slice may affect the sequence.
func FromSlice[T any](items []T) sequence.Sequence[T] {
	return sequence.Generate(func(f func(T) error) error {
		for _, x := range items {
			if err := f(x); err != nil {
				return err
			}
		}
		return nil
	})
}
