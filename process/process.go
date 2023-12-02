package process

import (
	"github.com/cookieo9/sequence"
	"github.com/cookieo9/sequence/tools"
)

// Process provides a method of generating a destination sequence from a given
// input sequence using a BEAM style emitter.
//
// The proc callback receives each input item, one at a time, and calls the
// provided emmiter function to output zero, one, or many results from that
// single input value. It can also return an error to stop processing.
func Process[In, Out any](src sequence.Sequence[In], proc func(In, func(Out)) error) sequence.Sequence[Out] {
	return sequence.Generate[Out](func(f func(Out) error) error {
		return sequence.Each(src)(func(i In) error {
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
