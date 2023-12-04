package build_test

import (
	"errors"
	"io"
	"strings"
	"testing"
	"testing/quick"

	"github.com/cookieo9/sequence"
	"github.com/cookieo9/sequence/convert/build"
	"github.com/cookieo9/sequence/convert/extract"
	"github.com/cookieo9/sequence/process"
	"github.com/cookieo9/sequence/tools"
)

const text = `
hello, world!
nothing much to see here.
I'm gonna head out now.
`

func quickFunc(gen func(io.Reader) sequence.Sequence[string]) func([]string) (string, int) {
	return func(s []string) (string, int) {
		full := strings.Join(s, "\n")
		seq := gen(strings.NewReader(full))
		num := 0
		result, err := extract.Collect(seq, "", func(line, merged string) string {
			num++
			return merged + line
		})
		tools.Check(err)
		return result, num
	}
}

func TestLines(t *testing.T) {
	linesWrap := func(r io.Reader) sequence.Sequence[string] {
		return build.Lines(r)
	}
	base := linesWrap
	testCases := []struct {
		name string
		gen  func(io.Reader) sequence.Sequence[string]
	}{
		{"Lines", linesWrap},
		{"Materialized", tools.Compose(linesWrap, sequence.Sequence[string].Materialize)},
		{"Buffered", tools.Compose(linesWrap, process.Buffer[string])},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			p := quickFunc(base)
			q := quickFunc(tc.gen)
			if err := quick.CheckEqual(p, q, nil); err != nil {
				t.Errorf("mismatch found: %v", err)
			}
		})
	}
}

func TestLinesBad(t *testing.T) {
	seq := build.Lines(strings.NewReader(text))
	coolErr := errors.New("my cool error")

	err := sequence.Each(seq)(func(s string) error {
		if strings.HasPrefix(s, "nothing") {
			return coolErr
		}
		t.Logf("Got a line!: %q", s)
		return nil
	})
	if !errors.Is(err, coolErr) {
		t.Errorf("didn't receive expected error; got %v, want %v", err, coolErr)
	}

	err = sequence.Each(seq)(func(s string) error { return nil })
	if err == nil {
		t.Error("no error occured on re-run of volatile sequence")
	}
	if !errors.Is(err, sequence.ErrRepeatedUse) {
		t.Errorf("got unexpected error from repeat use, got %q, want %q", err, sequence.ErrRepeatedUse)
	}
}
