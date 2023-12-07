package sequence

import (
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func BenchmarkMap(b *testing.B) {
	b.Run("RawSingle", func(b *testing.B) {
		seq := NumberSequence(0, b.N, 1)
		tools.Check(Each(seq)(func(i int) error { return nil }))
	})

	b.Run("MapFilterSingle", func(b *testing.B) {
		seq := NumberSequence(0, b.N, 1)
		mf := MapFilter(seq, func(i int) (int, bool, error) { return i, true, nil })
		tools.Check(Each(mf)(func(i int) error { return nil }))
	})

	b.Run("MapErrSingle", func(b *testing.B) {
		seq := NumberSequence(0, b.N, 1)
		me := MapErr(seq, func(i int) (int, error) { return i, nil })
		tools.Check(Each(me)(func(i int) error { return nil }))
	})

	b.Run("MapSingle", func(b *testing.B) {
		seq := NumberSequence(0, b.N, 1)
		m := Map(seq, func(i int) int { return i })
		tools.Check(Each(m)(func(i int) error { return nil }))
	})

	seq := NumberSequence(0, 1_000_000, 1)

	b.Run("Raw", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Check(Each(seq)(func(i int) error { return nil }))
		}
	})

	mapFilter := MapFilter(seq,
		func(x int) (int, bool, error) { return x, true, nil })

	b.Run("MapFilter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Check(Each(mapFilter)(func(i int) error { return nil }))
		}
	})

	mapErr := MapErr(seq,
		func(x int) (int, error) { return x, nil })

	b.Run("MapErr", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Check(Each(mapErr)(func(i int) error { return nil }))
		}
	})

	mapNoErr := Map(seq,
		func(x int) int { return x })

	b.Run("Map", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Check(Each(mapNoErr)(func(i int) error { return nil }))
		}
	})
}

func BenchmarkFilter(b *testing.B) {
	seq := NumberSequence(0, 1_000_000, 1)
	b.Run("Raw", func(b *testing.B) {
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
