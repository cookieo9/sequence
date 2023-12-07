package sequence

import (
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func BenchmarkFilter(b *testing.B) {
	seq := NumberSequence(0, 1_000_000, 1)
	b.Run("NoFilter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Check(Each(seq)(func(i int) error { return nil }))
		}
	})
	filtered := seq.Filter(func(i int) bool { return true })
	b.Run("Filter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Check(Each(filtered)(func(i int) error { return nil }))
		}
	})
}
