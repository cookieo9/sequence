package sequence

import (
	"context"
	"runtime"
	"sync"

	"github.com/cookieo9/sequence/tools"
	"golang.org/x/sync/errgroup"
)

// AsyncPool create a derived sequence that processes the input sequence using
// parallel goroutines. At most n goroutines may be in-flight at once, and a
// negative value results in no limit.
//
// All Aync* sequences provide no guarantee on the order of results in
// iteration, and downstream errors / early exits may not prevent the the
// processing of additional elements, although the errors will be processed
// and one will be returned.
//
// Due to the unstable nature of results, async sequences are also marked as
// volatile, meaning that they can only be iterated over once, and need to
// be generated again from the input sequence, or the results need to be
// materialized with [Materialize] or [extra.Buffer].
//
// Note: any state that is altered by a downstream operation, e.g. [Collect]
// may produce unexpected results since the order is non-deterministic, and
// the access is unsynchronized.
func AsyncPool[T any](n int, s Sequence[T]) Sequence[T] {
	if s.IsAsync() {
		return s
	}

	s2 := Derive(s, func(f func(T) error) error {
		grp, _ := errgroup.WithContext(context.Background())
		grp.SetLimit(n)

		err := s.Each(func(t T) error {
			grp.Go(func() error {
				return f(t)
			})
			return nil
		})

		return tools.Or(grp.Wait(), err)
	})
	s2.async = true
	return Volatile(s2)
}

// AsyncPool is a helper method to call the top-level function [AsyncPool]
// using the receiver.
func (s Sequence[T]) AsyncPool(n int) Sequence[T] {
	return AsyncPool(n, s)
}

// AsyncProcs is a helper that calls [AsyncPool] with the pool size set using
// [runtime.GOMAXPROCS].
func AsyncProcs[T any](s Sequence[T]) Sequence[T] {
	return AsyncPool(runtime.GOMAXPROCS(-1), s)
}

// AsyncCPUs is a helper that calls [AsyncPool] with the pool size set using
// [runtime.NumCPU].
func AsyncCPUs[T any](s Sequence[T]) Sequence[T] {
	return AsyncPool[T](runtime.NumCPU(), s)
}

// Async is an alias for [AsyncProcs].
func Async[T any](s Sequence[T]) Sequence[T] {
	return AsyncProcs(s)
}

// Async is a helper method that calls the top level function [Async] with the
// receiver.
func (s Sequence[T]) Async() Sequence[T] {
	return Async(s)
}

// Sync creates a new sequence where the iteration is guaranteed to produce
// values synchronously for downstream use. It's a no-op if the sequence is
// not asynchronous.
func Sync[T any](s Sequence[T]) Sequence[T] {
	if !s.IsAsync() {
		return s
	}
	s2 := Derive[T](s, func(f func(T) error) error {
		var lock sync.Mutex
		return s.Each(func(t T) error {
			lock.Lock()
			defer lock.Unlock()
			return f(t)
		})
	})
	s2.async = false
	return s2
}

// Sync is a helper that calls the top level function [Sync] with the reciever.
func (s Sequence[T]) Sync() Sequence[T] {
	return Sync(s)
}
