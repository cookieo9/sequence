package sequence

import (
	"errors"

	"github.com/cookieo9/sequence/tools"
)

var (
	// ErrStopIteration is used by iteration callbacks to indicate that
	// the iteration should be stopped early, but there was no other issue.
	ErrStopIteration = errors.New("iteration stopped")
)

// A Sequence represents a functionally immutable series of values. The
// sequence can only be used to fetch the values stored within, but the
// sequence itself can't be used to make changes.
type Sequence[T any] interface {
	Each(func(T) error) error
}

// Each returns a function to iterate over the sequence. The callback to
// the produced function is for users to receive and handle each value from
// the sequence, and return an error if processing should stop. The first
// error produced by the sequence, either an error from the callback or
// an error producing values will be returned. The [ErrStopIteration] error
// will not be returned from the iteration function.
func Each[T any](s Sequence[T]) func(func(T) error) error {
	return func(f func(T) error) error {
		if err := s.Each(f); !errors.Is(err, ErrStopIteration) {
			return err
		}
		return nil
	}
}

// EachSimple is like [Each], where a function is returned to iterate over
// the contents of the sequence, but the callback is in a simpler true/false
// form. An error generating a value will be returned, and iteration stopped.
// When the callback returns false, an ErrStopIteration error is passed
// through the sequence to anyone who cares about errors, but will be dropped
// before returning from EachSimple as it doesn't represent a real error.
func EachSimple[T any](s Sequence[T]) func(func(T) bool) error {
	fn := Each(s)
	return func(f func(T) bool) error {
		err := fn(tools.Compose(f, b2Err))
		if !errors.Is(err, ErrStopIteration) {
			return err
		}
		return nil
	}
}

// Iterator is a wrapper on [EachSimple] and is used to work with the upcoming
// language extension that allow functions in a for-range statement. Thus, when
// avaliable, one may write "for x := range Iterator(seq) {}" to use a sequence
// directly in a for-loop.
//
// Errors when generating the contents of the sequence will result in a panic.
func Iterator[T any](s Sequence[T]) func(func(T) bool) {
	fn := EachSimple(s)
	return func(f func(T) bool) {
		if err := fn(f); err != nil {
			panic(err)
		}
	}
}

func b2Err(b bool) error {
	if b {
		return nil
	}
	return ErrStopIteration
}
