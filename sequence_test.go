package sequence

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSequence(t *testing.T) {
	input := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	s := FromSlice(input)
	full, err := s.ToSlice().Pair()
	if err != nil {
		t.Errorf("unexpected error in copy: %v", err)
	} else if diff := cmp.Diff(full, input); diff != "" {
		t.Errorf("unexpected diff in copy (-got, +want): \n%s", diff)
	}

	l := Limit(s, 5)
	want := []int{3, 1, 4, 1, 5}
	got, err := l.ToSlice().Pair()
	if err != nil {
		t.Errorf("unexpected error in limited copy: %v", err)
	} else if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected diff in limited copy (-got, +want): \n%s", diff)
	}

	if full, err = s.ToSlice().Pair(); err != nil {
		t.Errorf("unexpected error in second copy: %v", err)
	} else if diff := cmp.Diff(full, input); diff != "" {
		t.Errorf("unexpected diff in copy copy (-got, +want): \n%s", diff)
	}
}
