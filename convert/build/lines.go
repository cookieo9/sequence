package build

import (
	"bufio"
	"io"

	"github.com/cookieo9/sequence"
)

// Lines returns a sequence that produces one at a time, lines of text from the
// provided reader. Errors while reading will interupt processing.
//
// This version is Volatile, meaning that you can only iterate over the input a
// single time, as it can't "go back" and return old lines, since none are
// stored.
//
// If you need to use the lines multiple times, consider calling Materialize on
// the returned sequence, or wrap it with [process.Buffer].
func Lines(r io.Reader) sequence.Sequence[string] {
	scn := bufio.NewScanner(r)
	return sequence.GenerateVolatile(func(f func(string) error) error {
		for scn.Scan() {
			if err := f(scn.Text()); err != nil {
				return err
			}
		}
		return scn.Err()
	})
}
