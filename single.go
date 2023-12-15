package sequence

import (
	"cmp"

	"github.com/cookieo9/sequence/tools"
)

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
	return Generate(func(f func(T) error) error {
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
// current value. CollectErr uses [Sync] to synchronize updating the value
// but there is no guarantee that the processing step will happen in the same
// order.
//
// The callback may return errors to stop processing, either to indicate
// failure, or in the case of ErrStopIteration to simply end processing.
func CollectErr[T, U any](s Sequence[T], initial U, process func(T, U) (U, error)) Result[U] {
	err := Each(s.Sync())(func(t T) error {
		var err error
		initial, err = process(t, initial)
		return err
	})
	return MakeResult(initial, err)
}

// Collect produces a single value from a sequence. It starts with an initial
// output value, and for every item in the sequence the value is permuted by
// the process function using the item from the sequence and the current value.
// Collect uses [Sync] to synchronize updating the accumulated value, but the
// order of updates is non-deterministic.
func Collect[T, U any](s Sequence[T], initial U, process func(T, U) U) Result[U] {
	return CollectErr[T, U](s, initial, func(t T, u U) (U, error) {
		return process(t, u), nil
	})
}

// First returns a Result containing the first value from the given sequence,
// or an error if even that isn't possible. The results of first are not
// deterministic for an asynchronous sequence.
func First[T any](s Sequence[T]) Result[T] {
	var value T
	err := EachSimple(s.Sync())(func(t T) bool { value = t; return false })
	return MakeResult(value, err)
}

// Last returns a Result containing the final value from the given sequence
// and any error should the sequence end with an error.
func Last[T any](s Sequence[T]) Result[T] {
	var value T
	err := EachSimple(s.Sync())(func(t T) bool { value = t; return true })
	return MakeResult(value, err)
}

// First is a helper method for the top level function [First].
func (s Sequence[T]) First() Result[T] {
	return First(s)
}

// Last is a helper method for the top level function [Last].
func (s Sequence[T]) Last() Result[T] {
	return Last(s)
}

// Sum is a helper function for a sequence of arithmetic values that produces
// the sum of the entire sequence.
func Sum[T tools.Arithmetic](s Sequence[T]) Result[T] {
	return Collect(s, T(0), tools.Add)
}

// Product is a helper function for a sequence of arithmetic values that
// returns the product of all values in the sequence.
func Product[T tools.Arithmetic](s Sequence[T]) Result[T] {
	return Collect(s, T(1), tools.Mul)
}

// Max returns a Result containing the smallest item from the sequence of
// cmp.Ordered items. The result will contain only an error if one stopped
// processing early.
func Max[T cmp.Ordered](s Sequence[T]) Result[T] {
	var (
		empty    = true
		maxValue T
	)
	err := EachSimple(s.Sync())(func(t T) bool {
		if empty {
			empty = false
			maxValue = t
			return true
		}
		maxValue = max(maxValue, t)
		return true
	})
	err = tools.Pick(err == nil && empty, ErrEmptySequence, err)
	return MakeResult(maxValue, err).Clean()
}

// MaxFunc returns the largest item from the sequence as determined by the
// provided comparision function. The comparison must return true if
// the first value is strictly less than the second. The returned result like
// with [Max], will only contain the error if one occurred.
func MaxFunc[T any](s Sequence[T], less func(x, y T) bool) Result[T] {
	var (
		empty    = true
		maxValue T
	)
	err := EachSimple(s.Sync())(func(t T) bool {
		if empty {
			empty = false
			maxValue = t
			return true
		}
		maxValue = tools.Pick(less(maxValue, t), t, maxValue)
		return true
	})
	err = tools.Pick(err == nil && empty, ErrEmptySequence, err)
	return MakeResult(maxValue, err).Clean()
}

// Min returns a Result containing the smallest item from the sequence of
// cmp.Ordered items. Any error during processing will stop the collection
// early, and the Result generate will only contain the error.
func Min[T cmp.Ordered](s Sequence[T]) Result[T] {
	var (
		empty    = true
		minValue T
	)
	err := EachSimple(s.Sync())(func(t T) bool {
		if empty {
			empty = false
			minValue = t
			return true
		}
		minValue = min(minValue, t)
		return true
	})
	err = tools.Pick(err == nil && empty, ErrEmptySequence, err)
	return MakeResult(minValue, err).Clean()
}

// MinFunc returns the smallest item from the sequence as determined by the
// provided comparision function. The comparison must return true if
// the first value is strictly less than the second. Otherwise it returns a
// Result with the same restrictions as [Min].
func MinFunc[T any](s Sequence[T], less func(x, y T) bool) Result[T] {
	var (
		empty    = true
		minValue T
	)
	err := EachSimple(s.Sync())(func(t T) bool {
		if empty {
			empty = false
			minValue = t
			return true
		}
		minValue = tools.Pick(less(t, minValue), t, minValue)
		return true
	})
	err = tools.Pick(err == nil && empty, ErrEmptySequence, err)
	return MakeResult(minValue, err).Clean()
}
