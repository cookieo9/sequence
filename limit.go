package sequence

// Limit returns a sequence where only a given number of items can be accessed
// from the input sequence before hitting the end of the sequence.
func Limit[T any](s Sequence[T], n int) Sequence[T] {
	return Derive[T](s, func(f func(T) error) error {
		i := 0
		return s.Each(func(t T) error {
			if i < n {
				i++
				return f(t)
			}
			return ErrStopIteration
		})
	})
}

// Limit is a helper method to call the top level function [Limit] on the
// receiver.
func (s Sequence[T]) Limit(n int) Sequence[T] {
	return Limit(s, n)
}

// While wraps a sequence to have it end early when the predicate returns
// false.
func While[T any](s Sequence[T], pred func(T) bool) Sequence[T] {
	return Derive(s, func(f func(T) error) error {
		return s.Each(func(t T) error {
			if !pred(t) {
				return ErrStopIteration
			}
			return f(t)
		})
	})
}

// While is a helper method to call the top level function [While] on the
// receiver.
func (s Sequence[T]) While(pred func(T) bool) Sequence[T] {
	return While(s, pred)
}

// Until returns a sequence that mirrors the provided sequence until the given
// predicate returns true. It's the logical counterpoint to [While].
func Until[T any](s Sequence[T], pred func(T) bool) Sequence[T] {
	return Derive[T](s, func(f func(T) error) error {
		return s.Each(func(t T) error {
			if pred(t) {
				return ErrStopIteration
			}
			return f(t)
		})
	})
}

// Until is a helper method to call the top level function [Until] on the
// receiver.
func (s Sequence[T]) Until(pred func(T) bool) Sequence[T] {
	return Until(s, pred)
}
