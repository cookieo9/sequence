package sequence

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestSimulate(t *testing.T) {
	start := MakePair(0, 1)
	sim := func(p Pair[int, int]) (Pair[int, int], bool) {
		a, b := p.B(), p.A()+p.B()
		p = MakePair(a, b)
		return p, p.B() < 1_000
	}
	seq := Simulation(start, sim)
	seq2 := Single(start).Simulate(sim)

	want := New[Pair[int, int]]([]Pair[int, int]{
		{0, 1}, {1, 1}, {1, 2}, {2, 3}, {3, 5}, {5, 8}, {8, 13}, {13, 21}, {21, 34},
		{34, 55}, {55, 89}, {89, 144}, {144, 233}, {233, 377}, {377, 610}, {610, 987},
	}...)

	opts := []cmp.Option{cmpopts.EquateComparable(Pair[int, int]{})}

	seqL4 := seq.Limit(4)
	seq2L4 := seq2.Limit(4)
	wantL4 := want.Limit(4)

	var testCases = []struct {
		name string
		got  Sequence[Pair[int, int]]
		want Sequence[Pair[int, int]]
		opts []cmp.Option
	}{
		{name: "Simulation", got: seq, want: want, opts: opts},
		{name: "Single+Simulate", got: seq2, want: want, opts: opts},
		{name: "Compare-Sims", got: seq, want: seq2, opts: opts},

		{name: "Simulation+Limit(4)", got: seqL4, want: wantL4, opts: opts},
		{name: "Single+Simulate+Limit(4)", got: seq2L4, want: wantL4, opts: opts},
		{name: "Compare-Sims+Limit(4)", got: seqL4, want: seq2L4, opts: opts},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, sequenceCompareTest(tc.got, tc.want, tc.opts...))
	}
}

func BenchmarkSimulate(b *testing.B) {
	start := 2
	sim := func(i int) (int, bool) {
		i *= i
		return i, i < 1_000_000
	}

	seq := Simulation(start, sim)
	b.Run("Simulation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = seq.Last().Value()
		}
	})

	seq2 := Single(start).Simulate(sim)
	b.Run("Single+Simulate", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = seq2.Last().Value()
		}
	})

	mSeq := seq.Materialize()
	b.Run("Simulation+Materialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = mSeq.Last().Value()
		}
	})

	mSeq2 := seq2.Materialize()
	b.Run("Single+Simulate+Materialize", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = mSeq2.Last().Value()
		}
	})
}
