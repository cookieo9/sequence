package sequence

import (
	"fmt"

	"github.com/cookieo9/sequence/tools"
)

// Counter creates a continuously increasing sequence of numbers starting
// from the given value and incrementing by 1. It is an infinite sequence.
func Counter[T tools.Arithmetic](start T) Sequence[T] {
	return Generate(func(f func(T) error) error {
		for i := start; ; i++ {
			if err := f(i); err != nil {
				return err
			}
		}
	})
}

// NumberSequence creates a finite sequence of numbers starting at a given
// value, incrementing by another value, and ending when the counter exceeds
// a given point analogeous to a "for i:=start; i<stop; i+=step {}" loop.
//
// - If start >= stop, a 0-item sequence is produced.
// - If step <= 0, an erroring sequnce is produced as it's likely a mistake.
func NumberSequence[T tools.Integer | tools.Real](start, stop, step T) Sequence[T] {
	if step <= 0 {
		return Error[T](fmt.Errorf("called NumberSequence with invalid step (%v <= 0)", step))
	}

	return Generate[T](func(f func(T) error) error {
		for i := start; i < stop; i += step {
			if err := f(i); err != nil {
				return err
			}
		}
		return nil
	})
}

// Count returns a Result with the number of items in the given sequence,
// or an error if one was produced while iterating the sequence.
func Count[T any](s Sequence[T]) Result[int] {
	var n int
	err := EachSimple(s)(func(t T) bool {
		n++
		return true
	})
	return MakeResult(n, err)
}

// Count is a helper to call the package function [Count] on the receiver.
func (s Sequence[T]) Count() Result[int] {
	return Count(s)
}
