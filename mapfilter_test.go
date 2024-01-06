package sequence

import (
	"errors"
	"strconv"
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func TestMap(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		input := New(1, 2, 3)
		expected := New(2, 4, 6)

		result := input.Map(func(i int) int { return i * 2 })
		compareSequences(t, result, expected)
	})

	t.Run("Empty", func(t *testing.T) {
		input := New[int]()
		expected := New[int]()
		result := input.Map(func(i int) int { return i * 2 })
		compareSequences(t, result, expected)
	})

	t.Run("Nil", func(t *testing.T) {
		input := New[int](1, 2, 3)
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Map didn't panic")
			} else {
				t.Log("Map panicked as expected: ", r)
			}
		}()

		result := input.Map(nil)
		compareSequences(t, result, input)
	})
}

func TestMapErr(t *testing.T) {

	t.Run("Valid", func(t *testing.T) {
		input := New("1", "2", "3")
		expected := New(1, 2, 3)
		result := MapErr(input, strconv.Atoi)
		compareSequences(t, result, expected)
	})

	t.Run("Invalid", func(t *testing.T) {
		input := New("1", "abc", "3")
		result := MapErr(input, strconv.Atoi)
		err := checkErrorSequence(t, result, strconv.ErrSyntax)
		t.Log("Error: ", err)
	})

}

func TestFilterErr(t *testing.T) {
	t.Run("basic filter", func(t *testing.T) {
		input := New(1, 2, 3, 4)
		expected := New(2, 4)

		result := FilterErr(input, func(i int) (bool, error) {
			return i%2 == 0, nil
		})
		compareSequences(t, result, expected)
	})

	t.Run("filter error", func(t *testing.T) {
		input := New(1, 2, 3, 4)

		errTest := errors.New("test error")
		result := FilterErr(input, func(i int) (bool, error) {
			if i == 3 {
				return false, errTest
			}
			return true, nil
		})

		err := checkErrorSequence(t, result, errTest)
		t.Log("Error: ", err)
	})

	t.Run("empty", func(t *testing.T) {
		input := New[int]()
		expected := New[int]()

		result := FilterErr(input, func(i int) (bool, error) {
			return true, nil
		})

		compareSequences(t, result, expected)
	})
}

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
