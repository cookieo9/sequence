package sequence

// Concat performs a concatenation of multiple sequences, where each sequence
// is iterated in turn until the final one is completed.
func Concat[T any](seqs ...Sequence[T]) Sequence[T] {
	return Generate(func(f func(T) error) error {
		for _, seq := range seqs {
			if err := Each(seq)(f); err != nil {
				return err
			}
		}
		return nil
	})
}

// Flatten processes a sequence of slices, producing a new sequence of the
// element type, and iterating across each slice as it's pulled from the input.
func Flatten[T any, S ~[]T](src Sequence[S]) Sequence[T] {
	return Process[S, T](src, func(s S, f func(T)) error {
		for _, x := range s {
			f(x)
		}
		return nil
	})
}

// Concat is a helper method which calls the top level function [Concat] to
// build a combined sequence of receiver followed by the given sequence(s).
func (s Sequence[T]) Concat(next ...Sequence[T]) Sequence[T] {
	seqs := append([]Sequence[T]{s}, next...)
	return Concat(seqs...)
}
