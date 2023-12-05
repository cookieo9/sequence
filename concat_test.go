package sequence

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConcat(t *testing.T) {
	a := FromSlice([]int{1, 2, 3})
	b := Repeat(4, 3)
	ab := a.Append(b)
	want := []int{1, 2, 3, 4, 4, 4}
	got, err := ab.ToSlice()
	if err != nil {
		t.Errorf("unexpected error upon conversion to slice: %v", err)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected diff of joined sequences (-got, +want):\n%s", diff)
	}

	wantErr := errors.New("my wonderful error")
	eSeq := Error[int](wantErr)
	aeb := Concat(a, eSeq, b)
	result, err := ToSlice(aeb)
	if err != wantErr {
		t.Errorf("error mismatch, got %v, want %v", err, wantErr)
	}
	if result != nil {
		t.Errorf("result mismatch in error case, got %v, want %v", result, nil)
	}
}

func TestFlatten(t *testing.T) {
	seq := FromSlice([][]int{{1, 2, 3}, {6, 5, 4}})
	flat := Flatten(seq)
	want := []int{1, 2, 3, 6, 5, 4}
	got, err := flat.ToSlice()
	if err != nil {
		t.Errorf("unexpected error upon conversion to slice: %v", err)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("unexpected diff from flatten (-got, +want):\n%s", diff)
	}
}
