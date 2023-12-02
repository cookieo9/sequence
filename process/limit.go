package process

import "github.com/cookieo9/sequence"

// Limit returns a sequence where only a given number of items can be accessed
// from the input sequence before hitting the end of the sequence.
func Limit[T any](s sequence.Sequence[T], n int) sequence.Sequence[T] {
	return sequence.Generate[T](func(f func(T) error) error {
		i := 0
		return s.Each(func(t T) error {
			if i < n {
				i++
				return f(t)
			}
			return nil
		})
	})
}
