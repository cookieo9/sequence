package sequence

import (
	"sync/atomic"
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func AtomicSum[T tools.Arithmetic](s Sequence[T]) Result[T] {
	var sum atomic.Value
	sum.Store(T(0))
	err := EachSimple(s)(func(t T) bool {
		for {
			val := sum.Load().(T)
			if sum.CompareAndSwap(val, val+t) {
				return true
			}
		}
	})
	return MakeResult(sum.Load().(T), err)
}

func AtomicIntSum(s Sequence[int64]) Result[int64] {
	var sum atomic.Int64
	err := EachSimple(s)(func(i int64) bool {
		sum.Add(i)
		return true
	})
	return MakeResult(sum.Load(), err)
}

func TestAsync(t *testing.T) {
	numbers := NumberSequence[int64](0, 1_000_000, 1)

	want := Sum(numbers)
	sum := Sum(numbers.Async().Sync())

	if sum != want {
		t.Errorf("mismatch in sum, got %d, want %d", sum, want)
	}
}

func BenchmarkAsync(b *testing.B) {
	numbers := NumberSequence[int64](0, 1_000_000, 1).Materialize()

	b.Log(Sum(numbers.Async().Sync()))
	b.Log(AtomicSum(numbers.Async()))
	b.Log(AtomicIntSum(numbers.Async()))

	b.Run("Sum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Sum(numbers).Value()
		}
	})

	b.Run("Sync+Sum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Sum(numbers.Async().Sync()).Value()
		}
	})

	b.Run("AtomicSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AtomicSum(numbers).Value()
		}
	})

	b.Run("Async+AtomicSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AtomicSum(numbers.Async()).Value()
		}
	})

	b.Run("AtomicIntSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AtomicIntSum(numbers).Value()
		}
	})

	b.Run("Async+AtomicIntSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			AtomicIntSum(numbers.Async()).Value()
		}
	})

}
