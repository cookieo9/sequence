package sequence

import "sync"

// Materialize returns a new sequence where all the data (including any error)
// from the current sequence is cached and played back on each iteration.
//
// This has 3 main benefits:
//   - It is likely faster than the original sequence
//   - It has no dependencies, not even on the original sequence
//   - It can be iterated over multiple times
//
// The speed benefit only makes sense if you are going to iterate more than
// once, as Materialize otherwise requires the time to iterate the sequence to
// be created.
//
// The resulting sequence is neither volatile, nor asynchronous, although an
// asynchronous input will mean the values stored for playback will be in a
// non-deterministic order.
func Materialize[T any](s Sequence[T]) Sequence[T] {
	data, srcErr := ToSlice(s.Sync())
	return Generate[T](func(f func(T) error) error {
		for _, x := range data {
			if err := f(x); err != nil {
				return err
			}
		}
		return srcErr
	})
}

// Materialize is a utility method that calls the package level function
// [Materialize] on this sequence.
func (s Sequence[T]) Materialize() Sequence[T] {
	return Materialize(s)
}

// Buffer is an alternative to the Materialize method that waits until the
// sequence is accessed before creating the cached copy. This means that there
// is basically no cost in computation or memory if the sequence hasn't been
// used yet.
func Buffer[T any](s Sequence[T]) Sequence[T] {
	seqGen := sync.OnceValue(func() Sequence[T] {
		return s.Materialize()
	})

	return Generate[T](func(f func(T) error) error {
		return seqGen().Each(f)
	})
}

// Buffer is a helper method to call the package function [Buffer] on the
// receiver.
func (s Sequence[T]) Buffer() Sequence[T] {
	return Buffer(s)
}
