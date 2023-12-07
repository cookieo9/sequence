package sequence

import (
	"errors"
	"sync"
)

// ErrRepeatedUse is returned by a VolatileSequence if it's accessed more than
// once.
var ErrRepeatedUse = errors.New("volatile sequence used more than once")

// Volatile returns a new sequence that depends on the the original, but is
// marked volatile, meaning it should only be iterated over a single time. Any
// sequence, even if not actually volatile can be wrapped, and will prevent
// repeat iterations over the sequence.
//
// The exact nature of the volatile sequence is as follows:
//   - The first time it's used it will behave exactly as normal. It can stop
//     early, fail, or run to completion
//   - Any subsequent accesses will return an error immediately. The behaviour
//     is unaffected by the outcome of the first iteration
//   - Iterations concurrent with the first will also fail.
func Volatile[T any](s Sequence[T]) Sequence[T] {
	var (
		used bool
		lock sync.Mutex
	)
	v := Derive(s, func(f func(T) error) error {
		lock.Lock()
		if used {
			lock.Unlock()
			return ErrRepeatedUse
		}
		used = true
		lock.Unlock()
		return s.Each(f)
	})
	v.volatile = true
	return v
}
