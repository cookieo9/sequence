package sequence

import "fmt"

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
