package build

import "github.com/cookieo9/sequence"

// FromChan produces a sequence from the provided channel. The new sequence is
// volatile since the channel itself can only be iterated over once.
func FromChan[T any](ch <-chan T) sequence.Sequence[T] {
	return sequence.GenerateVolatile(func(f func(T) error) error {
		for x := range ch {
			if err := f(x); err != nil {
				return err
			}
		}
		return nil
	})
}
