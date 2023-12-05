package sequence

// Concat performs a concatenation of multiple sequences, where each sequence
// is iterated in turn until the final one is completed.
func Concat[T any](seqs ...Sequence[T]) Sequence[T] {
	return ConcatSequences(FromSlice[Sequence[T]](seqs))
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

// ConcatSequences processes a sequence of sequences, and returns a single
// sequence where the contents of each sequence in order appear.
func ConcatSequences[T any](src Sequence[Sequence[T]]) Sequence[T] {
	return Process[Sequence[T], T](src, func(s Sequence[T], f func(T)) error {
		return s.Each(func(t T) error {
			f(t)
			return nil
		})
	})
}
