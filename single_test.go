package sequence

import (
	"errors"
	"testing"
)

func TestMinMax(t *testing.T) {
	testErr := errors.New("foo")
	testCases := []struct {
		name     string
		seq      Sequence[int]
		min, max int
		err      error
	}{
		{name: "empty", seq: New[int]()},
		{name: "pi", seq: New(3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5), min: 1, max: 9},
		{name: "err", seq: Error[int](testErr), err: testErr},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			seq := tc.seq
			wantMin := tc.min
			wantMax := tc.max
			less := func(x, y int) bool { return x < y }

			minFunc, err := MinFunc(seq, less)
			if err != tc.err {
				t.Errorf("unexpected error difference for MinFunc; got %v, want %v", err, tc.err)
			}
			if minFunc != wantMin {
				t.Errorf("unexpected difference for MinFunc; got %v, want %v", minFunc, wantMin)
			}

			minCmp, err := Min(seq)
			if err != tc.err {
				t.Errorf("unexpected error difference for Min; got %v, want %v", err, tc.err)
			}
			if minCmp != wantMin {
				t.Errorf("unexpected difference for Min; got %v, want %v", minCmp, wantMin)
			}

			maxFunc, err := MaxFunc(seq, less)
			if err != tc.err {
				t.Errorf("unexpected error difference for MaxFunc; got %v, want %v", err, tc.err)
			}
			if maxFunc != wantMax {
				t.Errorf("unexpected difference for MaxFunc; got %v, want %v", maxFunc, wantMax)
			}

			maxCmp, err := Max(seq)
			if err != tc.err {
				t.Errorf("unexpected error difference for Max; got %v, want %v", err, tc.err)
			}
			if maxCmp != wantMax {
				t.Errorf("unexpected difference for Max; got %v, want %v", maxCmp, wantMax)
			}
		})
	}
}
