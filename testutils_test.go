package sequence

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Set is a generic set implementation using a map. It provides methods for
// adding, checking membership, and adding multiple elements. The Set is not
// thread-safe, but the zero value is usable.
type Set[T comparable] map[T]struct{}

// NewSet creates a new generic Set from the provided items.
func NewSet[T comparable](items ...T) Set[T] {
	var s Set[T]
	s.AddMany(items...)
	return s
}

// init initializes the set if it is currently nil.
// It is used internally to ensure the set value is allocated before use.
func (s *Set[T]) init() {
	if *s == nil {
		*s = make(Set[T])
	}
}

// Has reports whether x is a member of the Set.
func (s Set[T]) Has(x T) bool {
	_, exist := s[x]
	return exist
}

// Add adds an element x to the set s. If x is already in s, Add returns false.
// Otherwise, Add initializes s if needed and inserts x, returning true.
func (s *Set[T]) Add(x T) bool {
	if s.Has(x) {
		return false
	}
	s.init()
	(*s)[x] = struct{}{}
	return true
}

// AddMany adds multiple elements xs to the set s. It iterates over xs and calls
// Add on each element. AddMany allows adding multiple elements in a single call
// instead of multiple individual Add calls.
func (s *Set[T]) AddMany(xs ...T) {
	for _, x := range xs {
		s.Add(x)
	}
}

// compareSequences compares two Sequence[T] values using cmp.Diff on the
// converted slice representations. It logs and errors on failure to convert
// either sequence to a slice, as well as any difference between the two
// sequences.
func compareSequences[T comparable](tb testing.TB, got, want Sequence[T], opts ...cmp.Option) {
	tb.Helper()

	gotSlice, err := ToSlice(got).Pair()
	if err != nil {
		tb.Errorf("problem converting got sequence to slice: %v", err)
	}
	tb.Logf("Got: %v (err: %v)", gotSlice, err)

	wantSlice, err := ToSlice[T](want).Pair()
	if err != nil {
		tb.Errorf("problem converting want sequence to slice: %v", err)
	}
	tb.Logf("Want: %v (err: %v)", wantSlice, err)

	if diff := cmp.Diff(gotSlice, wantSlice, opts...); diff != "" {
		tb.Errorf("unexpect diff when comparing sequences (-got, +want):\n%s", diff)
	}
}

// checkErrorSequence checks that processing the given sequence results in the
// expected error. It calls Each() on the sequence and expects it to return an
// error. If the expected error is non-nil, it verifies the returned error
// matches using errors.Is(). This is useful for testing error cases when
// processing a sequence. The error is returned so that additional checks can
// be performed by the caller.
func checkErrorSequence[T comparable](tb testing.TB, got Sequence[T], expected error) error {
	tb.Helper()
	err := got.Each(func(t T) error { return nil })
	if err == nil {
		tb.Error("expect error when processing sequence")
	}
	if expected != nil {
		if !errors.Is(err, expected) {
			tb.Errorf("expect error to be %v, but got %v", expected, err)
		}
	}
	tb.Log("Got error: ", err)
	return err
}

func sequenceCompareTest[T comparable](got, want Sequence[T], opts ...cmp.Option) func(t *testing.T) {
	return func(t *testing.T) { compareSequences(t, got, want, opts...) }
}
