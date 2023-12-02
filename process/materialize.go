package process

import (
	"sync"

	"github.com/cookieo9/sequence"
	"github.com/cookieo9/sequence/slices"
	"github.com/cookieo9/sequence/tools"
)

// Materialize wraps a sequence by storing its results to an internal slice
// and then serving that back to anyone accessing it.
//
// It serves as an easy way to wrap a VolatileSequence such that it can be
// used multiple times. Also, it can speed up processing if it is known the
// sequence is used multiple times, as after the first materialization all
// access is to the cached data itself.
//
// The two main downsides are that all values that eventually become part of
// the materialized sequence must be held in memory at all times, and the
// entire sequence passed to Materialize will be processed immediately before
// returning the new sequence.
func Materialize[T any](s sequence.Sequence[T]) sequence.Sequence[T] {
	data, err := slices.Collect(s)
	if err != nil {
		return sequence.Generate(func(f func(T) error) error { return err })
	}
	return slices.AsSequence(data)
}

// Buffer is like [Materialize] in that it keeps a cached copy of data from an
// underlying sequence, both for faster access and to support re-using a
// volatile sequence.
//
// Unlike [Materialize] a buffered sequence fills it's storage as it delivers
// values to a client, so there is no wait when calling Buffer before using
// the new sequence, and an unused Buffer sequence takes up no additional
// memory aside from the generator itself. It's slightly slower as locks are
// needed to ensure multiple readers play well with the cache, as some may be
// updating the cache at the same time. This is not true of Materialize where
// every access after it's built is a read and can easily be parallelized.
func Buffer[T any](s sequence.Sequence[T]) sequence.Sequence[T] {
	var (
		lock sync.Mutex
		data []T
		done bool
	)
	return sequence.Generate[T](func(f func(T) error) error {
		if !done {
			lock.Lock()
			defer lock.Unlock()
		}

		if len(data) > 0 {
			for _, x := range data {
				if err := f(x); err != nil {
					return err
				}
			}
		}
		if done {
			return nil
		}

		exit := false
		var errInner error

		errEnd := s.Each(func(t T) error {
			data = append(data, t)
			if exit {
				return nil
			}
			if err := f(t); err != nil {
				exit = true
				errInner = err
			}
			return nil
		})
		done = true
		return tools.Or(errInner, errEnd)
	})
}
