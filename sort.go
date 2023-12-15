package sequence

import (
	"cmp"
	"slices"
)

// SortOrdered creates a sequence that provides the items from a sequence of
// cmp.Ordered items in sorted order.
//
// Note: This function must read and store the entire sequence prior to sorting
// it, which may be an issue for large/infinite sequences.
func SortOrdered[T cmp.Ordered](s Sequence[T]) Sequence[T] {
	data, err := s.ToSlice().Pair()
	if err != nil {
		return Error[T](err)
	}
	slices.Sort(data)
	return FromSlice(data)
}

// Sort creates a sequence using the comparision function that returns the
// items of the original sequence in sorted order. The comparison function is
// the same used for [slices.Sort] and can be [cmp.Compare] for cmp.Ordered
// types.
//
// Note: This function must read and store the entire sequence prior to sorting
// it, which may be an issue for large/infinite sequences.
func Sort[T any](s Sequence[T], cmp func(a, b T) int) Sequence[T] {
	data, err := s.ToSlice().Pair()
	if err != nil {
		return Error[T](err)
	}
	slices.SortFunc(data, cmp)
	return FromSlice(data)
}

// Sort is a helper method to call the package function [Sort] on the receiver.
func (s Sequence[T]) Sort(cmp func(T, T) int) Sequence[T] {
	return Sort(s, cmp)
}

// SortStable is like [Sort] and takes the same parameters, but the resulting
// sort/order will be stable.
//
// Note: This function must read and store the entire sequence prior to sorting
// it, which may be an issue for large/infinite sequences.
func SortStable[T any](s Sequence[T], cmp func(a, b T) int) Sequence[T] {
	data, err := s.ToSlice().Pair()
	if err != nil {
		return Error[T](err)
	}
	slices.SortStableFunc(data, cmp)
	return FromSlice(data)
}

// SortStable is a helper method to call the package function [SortStable] on
// the receiver.
func (s Sequence[T]) SortStable(cmp func(T, T) int) Sequence[T] {
	return SortStable(s, cmp)
}

// Reverse returns a sequence that produces the items from the input sequence
// in reverse order.
//
// Note: This function must read and store the entire sequence prior to
// reversing it, which may be an issue for large/infinite sequences.
func Reverse[T any](s Sequence[T]) Sequence[T] {
	data, err := s.ToSlice().Pair()
	if err != nil {
		return Error[T](err)
	}
	slices.Reverse(data)
	return FromSlice(data)
}

// Reverse is a helper method to call the package function [Reverse] on the
// receiver.
func (s Sequence[T]) Reverse() Sequence[T] {
	return Reverse(s)
}
