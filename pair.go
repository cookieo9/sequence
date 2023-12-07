package sequence

import (
	"cmp"
	"fmt"
)

// A Pair is a 2 element tuple containing values of independant types.
type Pair[A, B any] struct {
	a A
	b B
}

// MakePair creates a pair using its two arguments
func MakePair[A, B any](a A, b B) Pair[A, B] { return Pair[A, B]{a: a, b: b} }

// A returns the first value of the pair.
func (p Pair[A, B]) A() A { return p.a }

// B returns the second value of the pair.
func (p Pair[A, B]) B() B { return p.b }

// AB returns both values from the pair in order.
func (p Pair[A, B]) AB() (A, B) { return p.a, p.b }

// BA returns both values from the pair in reverse order.
func (p Pair[A, B]) BA() (B, A) { return p.b, p.a }

// Swap retuns a new pair where the elements of the current pair are swapped.
func (p Pair[A, B]) Swap() Pair[B, A] { return MakePair(p.BA()) }

// String returns a string representation of this pair.
func (p Pair[A, B]) String() string {
	return fmt.Sprintf("(%v, %v)", p.a, p.b)
}

// PairCompare compares two Pair values with the same element types provided
// both items in the pair are a cmp.Ordered type.  The comparison first looks
// at the first (A) value, and if equal, moves on to the second (B) value.
func PairCompare[A, B cmp.Ordered](a, b Pair[A, B]) int {
	if c := cmp.Compare(a.a, b.a); c != 0 {
		return c
	}
	return cmp.Compare(a.b, b.b)
}

// PairCompareFirst compares 2 pairs by only their first elements using
// cmp.Compare. The second element of the Pair is completely ignored for the
// purpose of the comparison.
func PairCompareFirst[A cmp.Ordered, B any](a, b Pair[A, B]) int {
	return cmp.Compare(a.a, b.a)
}

// PairCompareSecond compares 2 pairs by only their second elements using
// cmp.Compare. The first element of the Pair is completely ignored for the
// purpose of the comparison.
func PairCompareSecond[A, B cmp.Ordered](a, b Pair[A, B]) int {
	return cmp.Compare(a.b, b.b)
}
