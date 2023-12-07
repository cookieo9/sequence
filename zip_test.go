package sequence

import (
	"testing"

	"github.com/cookieo9/sequence/tools"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestZip(t *testing.T) {
	seq1 := New(1, 2, 3, 4)
	seq2 := New("one", "two", "three")

	opt := cmpopts.EquateComparable(Pair[int, string]{})

	short := Zip(seq1, seq2)
	gotShort := tools.Must(short.ToSlice())
	wantShort := []Pair[int, string]{{1, "one"}, {2, "two"}, {3, "three"}}

	if diff := cmp.Diff(gotShort, wantShort, opt); diff != "" {
		t.Errorf("unexpected diff in Zip results (-got, +want):\n%s", diff)
	}

	long := ZipLongest(seq1, seq2)
	gotLong := tools.Must(long.ToSlice())
	wantLong := []Pair[int, string]{{1, "one"}, {2, "two"}, {3, "three"}, {4, ""}}

	if diff := cmp.Diff(gotLong, wantLong, opt); diff != "" {
		t.Errorf("unexpected diff in Zip results (-got, +want):\n%s", diff)
	}
}
