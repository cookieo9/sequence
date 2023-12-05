package sequence

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
// Otherwise, it's other main use is to overcome the limitations of volatile
// sequences.
func Materialize[T any](s Sequence[T]) Sequence[T] {
	data, srcErr := ToSlice(s)
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
