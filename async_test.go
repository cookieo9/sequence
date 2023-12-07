package sequence

import (
	"sync/atomic"
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func AtomicSum[T tools.Arithmetic](s Sequence[T]) (T, error) {
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
	return tools.CleanErrors(sum.Load().(T), err)
}

func AtomicIntSum(s Sequence[int64]) (int64, error) {
	var sum atomic.Int64
	err := EachSimple(s)(func(i int64) bool {
		sum.Add(i)
		return true
	})
	return tools.CleanErrors(sum.Load(), err)
}

func TestAsync(t *testing.T) {
	numbers := NumberSequence[int64](0, 1_000_000, 1)

	want := tools.Must(Sum(numbers))
	sum := tools.Must(Sum(numbers.Async().Sync()))

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
			tools.Must(Sum(numbers))
		}
	})

	b.Run("Sync+Sum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Must(Sum(numbers.Async().Sync()))
		}
	})

	b.Run("AtomicSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Must(AtomicSum(numbers))
		}
	})

	b.Run("Async+AtomicSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Must(AtomicSum(numbers.Async()))
		}
	})

	b.Run("AtomicIntSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Must(AtomicIntSum(numbers))
		}
	})

	b.Run("Async+AtomicIntSum", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tools.Must(AtomicIntSum(numbers.Async()))
		}
	})

}
