package sequence

// MapFilter performs both a [Map] and [Filter] operation on the input sequence
// where the convert function returns both the new value, and a boolean to
// indicate if it should be added at all.
//
// MapFilter allows the callback to stop iteration with a generic error, or use
// ErrStopIteration to simply indicate that no more values should be proceed.
func MapFilter[In, Out any](s Sequence[In], convert func(In) (Out, bool, error)) Sequence[Out] {
	return Derive(s, func(f func(Out) error) error {
		return s.Each(func(i In) error {
			out, ok, err := convert(i)
			if err != nil {
				return err
			}
			if ok {
				return f(out)
			}
			return nil
		})
	})
}

// Map takes in an input sequence and returns a sequence where every input
// value is permuted by the provided convert function. The new sequence may
// be of a different type due to the conversion, but will have the same number
// of items.
func Map[In, Out any](s Sequence[In], convert func(In) Out) Sequence[Out] {
	return MapFilter(s, func(in In) (Out, bool, error) {
		return convert(in), true, nil
	})
}

// MapErr takes in an input sequence and returns a sequence where every input
// value is permuted by the provided convert function. The new sequence may
// be of a different type due to the conversion, but will have the same number
// of items.
//
// MapErr allows the callback to stop iteration with a generic error, or use
// ErrStopIteration to simply indicate that no more values should be proceed.
func MapErr[In, Out any](s Sequence[In], convert func(In) (Out, error)) Sequence[Out] {
	return MapFilter(s, func(in In) (Out, bool, error) {
		out, err := convert(in)
		return out, true, err
	})
}

// Filter takes an input sequence and creates a sequence where only the values
// that pass the provided predicate function will be emitted. The ouput
// sequence will be the same type as the input sequence.
func Filter[T any](s Sequence[T], pred func(T) bool) Sequence[T] {
	return Derive[T](s, func(f func(T) error) error {
		return s.Each(func(t T) error {
			if pred(t) {
				return f(t)
			}
			return nil
		})
	})
}

// FilterErr takes an input sequence and creates a sequence where only the
// values that pass the provided predicate function will be emitted. The ouput
// sequence will have the same element type as the input sequence.
//
// FilterErr allows the callback to stop iteration with a generic error, or use
// ErrStopIteration to simply indicate that no more values should be proceed.
func FilterErr[T any](s Sequence[T], pred func(T) (bool, error)) Sequence[T] {
	return MapFilter(s, func(t T) (T, bool, error) { ok, err := pred(t); return t, ok, err })
}

// Filter is a helper method that creates a new sequence via the top level
// function [Filter] using the receiver and the given predicate function.
func (s Sequence[T]) Filter(pred func(T) bool) Sequence[T] {
	return Filter(s, pred)
}
