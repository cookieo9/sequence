package sequence

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](items ...T) Set[T] {
	var s Set[T]
	s.AddMany(items...)
	return s
}

func (s *Set[T]) init() {
	if *s == nil {
		*s = make(Set[T])
	}
}

func (s Set[T]) Has(x T) bool {
	_, exist := s[x]
	return exist
}

func (s *Set[T]) Add(x T) bool {
	if s.Has(x) {
		return false
	}
	s.init()
	(*s)[x] = struct{}{}
	return true
}

func (s *Set[T]) AddMany(xs ...T) {
	for _, x := range xs {
		s.Add(x)
	}
}

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
	return err
}

func sequenceCompareTest[T comparable](got, want Sequence[T], opts ...cmp.Option) func(t *testing.T) {
	return func(t *testing.T) { compareSequences(t, got, want, opts...) }
}
