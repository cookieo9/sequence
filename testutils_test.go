package sequence

import (
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

func compareSequences[T comparable](t *testing.T, got, want Sequence[T], opts ...cmp.Option) {
	t.Helper()

	gotSlice, err := ToSlice(got)
	if err != nil {
		t.Errorf("problem converting got sequence to slice: %v", err)
	}

	wantSlice, err := ToSlice[T](want)
	if err != nil {
		t.Errorf("problem converting want sequence to slice: %v", err)
	}

	if diff := cmp.Diff(gotSlice, wantSlice); diff != "" {
		t.Errorf("unexpect diff when comparing sequences (-got, +want):\n%s", diff)
	}
}