package single

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

// CollectErr produces a single value from a sequence. It starts with an
// initial output value, and for every item in the sequence the value is
// permuted by the process function using the item from the sequence and the
// current value.
//
// The callback may return errors to stop processing, either to indicate
// failure, or in the case of ErrStopIteration to simply end processing.
func CollectErr[T, U any](s sequence.Sequence[T], initial U, process func(T, U) (U, error)) (U, error) {
	err := sequence.Each(s)(func(t T) error {
		var err error
		initial, err = process(t, initial)
		return err
	})
	return initial, err
}

// Collect produces a single value from a sequence. It starts with an initial
// output value, and for every item in the sequence the value is permuted by
// the process function using the item from the sequence and the current value.
func Collect[T, U any](s sequence.Sequence[T], initial U, process func(T, U) U) (U, error) {
	return CollectErr[T, U](s, initial, func(t T, u U) (U, error) {
		return process(t, u), nil
	})
}
