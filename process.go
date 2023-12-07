package sequence

import (
	"sync/atomic"

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
		var (
			err atomic.Value
		)

		emit := func(out Out) {
			if err.Load() == nil {
				if e := f(out); e != nil {
					err.CompareAndSwap(nil, e)
				}
			}
		}

		return src.Each(func(i In) error {
			e2 := proc(i, emit)
			if e1 := err.Load(); e1 != nil {
				return tools.Or(e1.(error), e2)
			}
			return e2
		})
	})
}
