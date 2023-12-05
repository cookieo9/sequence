package sequence

import "github.com/cookieo9/sequence/tools"

// Infinite generates a sequence where the given value is returned forever, on
// each iteration of the sequence.
func Infinite[T any](t T) Sequence[T] {
	return Generate(func(f func(T) error) error {
		for {
			if err := f(t); err != nil {
				return err
			}
		}
	})
}

// Repeat generates a sequence where the given item is returns a fixed number
// of times.
func Repeat[T any](t T, n int) Sequence[T] {
	return Generate[T](func(f func(T) error) error {
		for i := 0; i < n; i++ {
			if err := f(t); err != nil {
				return err
			}
		}
		return nil
	})
}

// Single returns a sequence where the given value is returned once.
func Single[T any](t T) Sequence[T] {
	return Repeat(t, 1)
}

// Error returns a sequence where the given error is return when iterated
// upon.
func Error[T any](err error) Sequence[T] {
	return Generate[T](func(f func(T) error) error { return err })
}

// CollectErr produces a single value from a sequence. It starts with an
// initial output value, and for every item in the sequence the value is
// permuted by the process function using the item from the sequence and the
// current value.
//
// The callback may return errors to stop processing, either to indicate
// failure, or in the case of ErrStopIteration to simply end processing.
func CollectErr[T, U any](s Sequence[T], initial U, process func(T, U) (U, error)) (U, error) {
	err := Each(s)(func(t T) error {
		var err error
		initial, err = process(t, initial)
		return err
	})
	return initial, err
}

// Collect produces a single value from a sequence. It starts with an initial
// output value, and for every item in the sequence the value is permuted by
// the process function using the item from the sequence and the current value.
func Collect[T, U any](s Sequence[T], initial U, process func(T, U) U) (U, error) {
	return CollectErr[T, U](s, initial, func(t T, u U) (U, error) {
		return process(t, u), nil
	})
}

// First returns the first value from the given sequence, or an error if even
// that isn't possible.
func First[T any](s Sequence[T]) (T, error) {
	var value T
	err := EachSimple(s)(func(t T) bool { value = t; return false })
	return tools.CleanErrors(value, err)
}

// Last returns the final value from the given sequence and any error should
// the sequence end with an error. Unlike most functions both value and error
// will be returns as it's really upto the caller to determine if an error
// should stop downstream processing in this case. Note: ErrStopIteration
// will still be suppressed as normal as it indicates early exit without error.
func Last[T any](s Sequence[T]) (T, error) {
	var value T
	err := EachSimple(s)(func(t T) bool { value = t; return true })
	return value, err
}