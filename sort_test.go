package sequence

import (
	"cmp"
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSorts(t *testing.T) {
	//              0  1  2  3  4  5  6  7  8  9  a
	seq := New(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5)
	want := New(1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9)

	t.Run("Reverse", func(t *testing.T) {
		want := New(5, 3, 5, 6, 2, 9, 5, 1, 4, 1, 3)
		got := seq.Reverse()
		compareSequences(t, got, want)
	})

	t.Run("SortOrdered", func(t *testing.T) {
		ordered := SortOrdered(seq)
		compareSequences(t, ordered, want)
	})
	t.Run("Sort", func(t *testing.T) {
		funcOrder := seq.Sort(cmp.Compare)
		compareSequences(t, funcOrder, want)
	})
	t.Run("SortStable", func(t *testing.T) {
		stableOrder := seq.SortStable(cmp.Compare)
		compareSequences(t, stableOrder, want)
	})

	indexed := Zip(seq, Counter(0))
	wantPair := New[Pair[int, int]]([]Pair[int, int]{
		{1, 1},
		{1, 3},
		{2, 6},
		{3, 0},
		{3, 9},
		{4, 2},
		{5, 4},
		{5, 8},
		{5, 10},
		{6, 7},
		{9, 5},
	}...)

	opt := cmpopts.EquateComparable(Pair[int, int]{})

	t.Run("PairStable", func(t *testing.T) {
		pairOrdered := indexed.SortStable(PairCompare)
		compareSequences(t, pairOrdered, wantPair, opt)
	})
}
