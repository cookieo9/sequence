package sequence

// Scan converts a sequence into a scaned version of the sequence where every
// element returned is the result of performing a binary operation between the
// element from the input sequence and an accumulator, with the result also
// going back in the accumulator. The initial value of the accumulator must
// also be provided. It is recommended that op be commutative and associative
// for consistent results.
func Scan[T, U any](s Sequence[T], initial U, op func(T, U) U) Sequence[U] {
	return Map[T, U](s, func(t T) U {
		initial = op(t, initial)
		return initial
	})
}

// Scan is a helper method for calling the top level function [Scan] with the
// current sequence, as long as the accumulator and binary operation all use
// and return the same type as the sequence's element type.
func (s Sequence[T]) Scan(initial T, op func(T, T) T) Sequence[T] {
	return Scan(s, initial, op)
}
