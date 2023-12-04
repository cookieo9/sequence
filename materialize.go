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
func (s Sequence[T]) Materialize() Sequence[T] {
	var data []T
	srcErr := EachSimple[T](s)(func(item T) bool {
		data = append(data, item)
		return true
	})

	return Generate[T](func(f func(T) error) error {
		for _, x := range data {
			if err := f(x); err != nil {
				return err
			}
		}
		return srcErr
	})
}
