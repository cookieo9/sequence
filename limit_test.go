package sequence

import "testing"

func TestWhile(t *testing.T) {
	seq := New(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5)
	var seen Set[int]
	got := seq.While(seen.Add)
	want := New[int](3, 1, 4)

	compareSequences(t, got, want)
}

func TestUntil(t *testing.T) {
	seq := New(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5)
	bad := NewSet(9)
	got := seq.Until(bad.Has)
	want := New[int](3, 1, 4, 1, 5)

	compareSequences(t, got, want)
}
