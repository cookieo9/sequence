package sequence

import (
	"testing"

	"github.com/cookieo9/sequence/tools"
	"github.com/google/go-cmp/cmp"
)

func TestScanAndSum(t *testing.T) {
	seq := New(1, 2, 3, 4, 5)
	scn := seq.Scan(1, tools.Mul[int])

	scnWant := []int{1, 2, 6, 24, 120}
	scnGot, err := scn.ToSlice().Pair()
	if err != nil {
		t.Errorf("unexpected error converting scan to slice: %v", err)
	}
	if diff := cmp.Diff(scnGot, scnWant); diff != "" {
		t.Errorf("unexpected mismatch in scan result (-got, +want):\n%s", diff)
	}

	sumWant := 15
	sumGot, err := Sum(seq).Pair()
	if err != nil {
		t.Errorf("unexpected error computing sum: %v", err)
	}
	if sumGot != sumWant {
		t.Errorf("unexpected difference in sum result; got %v, want %v", sumGot, sumWant)
	}
}
