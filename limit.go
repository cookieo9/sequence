package sequence

// Limit returns a sequence where only a given number of items can be accessed
// from the input sequence before hitting the end of the sequence.
func Limit[T any](s Sequence[T], n int) Sequence[T] {
	return Generate[T](func(f func(T) error) error {
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
