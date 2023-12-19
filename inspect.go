package sequence

// Inspect adds an "inpection" phase to the given sequence, where each value,
// and it's index is passed to the callback. The result of the callback may
// be an error to stop iteration, but otherwise, the value is returned from
// the sequence as though the inspection step never happened.
//
// Note: a reference type can be changed by the inspection function, since
// only a shallow copy is made for passing on.
func Inspect[T any](s Sequence[T], inspect func(int, T) error) Sequence[T] {
	return Derive(s, func(f func(T) error) error {
		n := 0
		return s.Each(func(t T) error {
			err := inspect(n, t)
			n++
			if err != nil {
				return err
			}
			return f(t)
		})
	})
}

// Inspect is a method helper to call the package function [Inspect] on the
// receiver.
func (s Sequence[T]) Inspect(inspect func(int, T) error) Sequence[T] {
	return Inspect(s, inspect)
}
