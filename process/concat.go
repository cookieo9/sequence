package process

import "github.com/cookieo9/sequence"

// Concat performs a concatenation of multiple sequences, where each sequence
// is iterated in turn until the final one is completed.
func Concat[T any](seqs ...sequence.Sequence[T]) sequence.Sequence[T] {
	return sequence.Generate[T](func(f func(T) error) error {
		var stop bool
		for _, seq := range seqs {
			err := sequence.Each(seq)(func(t T) error {
				if e := f(t); e != nil {
					stop = true
					return e
				}
				return nil
			})
			if err != nil || stop {
				return err
			}
		}
		return nil
	})
}

// Flatten processes a sequence of slices, producing a new sequence of the
// element type, and iterating across each slice as it's pulled from the input.
func Flatten[T any, S ~[]T](src sequence.Sequence[S]) sequence.Sequence[T] {
	return Process[S, T](src, func(s S, f func(T)) error {
		for _, x := range s {
			f(x)
		}
		return nil
	})
}
