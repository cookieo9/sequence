package sequence

import (
	"testing"

	"github.com/cookieo9/sequence/tools"
)

func euler[T tools.Integer](a, b T) T {
	return (a + b) * (b - a + 1) / 2
}

func TestCounting(t *testing.T) {
	t.Parallel()

	limited := Limit(Counter(5), 1_000_000)
	numSeq := NumberSequence(5, 1_000_005, 1)

	compareSequences(t, limited, numSeq)

	first := tools.Must(First(numSeq))
	wantFirst := 5
	if first != wantFirst {
		t.Errorf("unexpected difference for First(); got: %v, want: %v", first, wantFirst)
	}

	last := tools.Must(Last(numSeq))
	wantLast := 1_000_004
	if last != wantLast {
		t.Errorf("unexpected difference for Last(); got: %v, want: %v", last, wantLast)
	}

	count := tools.Must(Count(numSeq))
	wantCount := 1_000_000
	if count != wantCount {
		t.Errorf("unexpected difference for Count(); got: %v, want: %v", count, wantCount)
	}

	sum := tools.Must(Sum(numSeq))
	wantSum := euler(5, 1_000_004)
	if sum != wantSum {
		t.Errorf("unexpected difference for Sum(); got: %v, want: %v", sum, wantSum)
	}
}

func TestBadNumberSequences(t *testing.T) {
	t.Parallel()

	empty := NumberSequence(10, 0, 1)
	compareSequences(t, empty, New[int]())

	empty = NumberSequence(0, -10, 2)
	compareSequences(t, empty, New[int]())

	inf := NumberSequence(0, 10, 0)
	sum, err := Sum(inf)
	t.Logf("Sum(inf): %v %v", sum, err)
	if err == nil {
		t.Error("expected error, but none provided from Sum(inf)")
	}

	rev := NumberSequence(10, 0, -1)
	sum, err = Sum(rev)
	t.Logf("Sum(rev): %v %v", sum, err)
	if err == nil {
		t.Error("expected error, but none provided from Sum(rev)")
	}
}
