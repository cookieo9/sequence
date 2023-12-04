package build

import "github.com/cookieo9/sequence"

// Infinite generates a sequence where the given value is returned forever, on
// each iteration of the sequence.
func Infinite[T any](t T) sequence.Sequence[T] {
	return sequence.Generate(func(f func(T) error) error {
		for {
			if err := f(t); err != nil {
				return err
			}
		}
	})
}

// Repeat generates a sequence where the given item is returns a fixed number
// of times.
func Repeat[T any](t T, n int) sequence.Sequence[T] {
	return sequence.Generate[T](func(f func(T) error) error {
		for i := 0; i < n; i++ {
			if err := f(t); err != nil {
				return err
			}
		}
		return nil
	})
}

// Single returns a sequence where the given value is returned once.
func Single[T any](t T) sequence.Sequence[T] {
	return Repeat(t, 1)
}

// Error returns a sequence where the given error is return when iterated
// upon.
func Error[T any](err error) sequence.Sequence[T] {
	return sequence.Generate[T](func(f func(T) error) error { return err })
}
