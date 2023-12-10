package sequence

import (
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func TestSimulate(t *testing.T) {
	seq := Simulate(2, func(i int) (int, bool) {
		i *= i
		return i, i < 1_000_000
	})
	want := New(2, 4, 16, 256, 65536)

	compareSequences(t, seq, want)
}

func BenchmarkSimulate(b *testing.B) {
	seq := Simulate(2, func(i int) (int, bool) {
		i *= i
		return i, i < 1_000_000
	})

	b.Run("Raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = tools.Must(seq.Last())
		}
	})

	mSeq := seq.Materialize()
	b.Run("Materialized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = tools.Must(mSeq.Last())
		}
	})
}
