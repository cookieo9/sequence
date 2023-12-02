package sequence

import (
	"errors"
	"sync"
)

// ErrRepeatedUse is returned by a VolatileSequence if it's accessed more than
// once.
var ErrRepeatedUse = errors.New("volatile sequence used more than once")

// A VolatileSequence represents a sequence where the process of iterating over
// the sequence damages or changes it in a fundamental way. Since sequences are
// meant to be immutable, some cases can't uphold this requirement.
//
// VolatileSequences are OK to be accessed once, but any further access will
// return ErrRepeatedUse immediately. They serve as an indication that the
// user may wish to wrap them using [process.Materialize] or [process.Buffer]
// if they think they may be used more than once.
type VolatileSequence[T any] interface {
	Sequence[T]
}

// Volatile wraps the provided sequence in the VolatileSequence type indicating
// it should only be iterated over a single time. Any sequence, even if not
// actually volatile can be wrapped, and will prevent repeat calls to Each.
//
// The exact nature of the wrapped sequence is as follows:
//   - The first time it's used it will behave exactly as normal
//   - it can stop early, fail, or run to completion
//   - Any subsequent accesses will return an error immediately
//   - the behaviour is unaffected by the outcome of the first iteration
//   - Iterations concurrent with the first will also fail in the same way
func Volatile[T any](s Sequence[T]) VolatileSequence[T] {
	var (
		used bool
		lock sync.Mutex
	)
	return Generate(func(f func(T) error) error {
		lock.Lock()
		if used {
			lock.Unlock()
			return ErrRepeatedUse
		}
		used = true
		lock.Unlock()
		return s.Each(f)
	})
}
