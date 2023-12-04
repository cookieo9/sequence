package extract

import "github.com/cookieo9/sequence"

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
