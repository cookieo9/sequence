package sequence

import (
	"errors"
)

var (
	// ErrStopIteration is used by iteration callbacks to indicate that
	// the iteration should be stopped early, but there was no other issue.
	ErrStopIteration = errors.New("iteration stopped")

	// ErrEmptySequence is returned when an operation expects a non-empty
	// sequence, but doesn't get one.
	ErrEmptySequence = errors.New("empty sequence")
)

// A Sequence represents a functionally immutable series of values. The
// sequence can only be used to fetch the values stored within, but the
// sequence itself can't be used to make changes.
type Sequence[T any] struct {
	source   func(func(T) error) error
	volatile bool
	async    bool
}

// IsVolatile returns true if the sequence is volatile, meaning that it should
// only be iterated over once. It is an error to call Each/EachSimple/Iterator
// more than once. Volatile sequences can be created with the [Volatlie]
// function.
func (s Sequence[T]) IsVolatile() bool { return s.volatile }

// IsAsync will return true if the sequence will asynchronously generate values
// during iteration.
func (s Sequence[T]) IsAsync() bool { return s.async }

// Each iterates over every item in the sequence calling the passed callback
// with each item. An error returned from the callback, or one that arose from
// the processing of the sequence will be returned if they arise and iteration
// will stop. Unlike the top-level functions that iterate over sequences, this
// method will not on it's own handle [ErrStopIteration] and return it instead.
func (s Sequence[T]) Each(f func(T) error) error {
	return s.source(f)
}
