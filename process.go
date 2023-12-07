package sequence

import (
	"github.com/cookieo9/sequence/tools"
)

// Process provides a method of generating a destination sequence from a given
// input sequence using a BEAM style emitter.
//
// The proc callback receives each input item, one at a time, and calls the
// provided emmiter function to output zero, one, or many results from that
// single input value. It can also return an error to stop processing.
func Process[In, Out any](src Sequence[In], proc func(In, func(Out)) error) Sequence[Out] {
	return Derive(src, func(f func(Out) error) error {
		return src.Each(func(i In) error {
			var err error
			emit := func(out Out) {
				if err == nil {
					err = f(out)
				}
			}
			e2 := proc(i, emit)
			return tools.Or(err, e2)
		})
	})
}
