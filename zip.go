package sequence

import (
	"github.com/cookieo9/sequence/tools"
)

func zipCommon[A, B any](aSeq Sequence[A], bSeq Sequence[B], shortest bool) Sequence[Pair[A, B]] {
	s := Generate(func(f func(Pair[A, B]) error) error {
		aCh := make(chan Pair[A, error], 1)
		go IntoChanPair(aCh, aSeq)
		bCh := make(chan Pair[B, error], 1)
		go IntoChanPair(bCh, bSeq)

		for {
			aP, aOk := <-aCh
			bP, bOk := <-bCh

			stopShortest := shortest && (!aOk || !bOk)
			stopBoth := !aOk && !bOk
			if stopShortest || stopBoth {
				break
			}

			aV, aErr := aP.AB()
			bV, bErr := bP.AB()
			if err := tools.Or(aErr, bErr); err != nil {
				return err
			}
			if err := f(MakePair(aV, bV)); err != nil {
				return err
			}
		}
		return nil
	})
	if aSeq.IsVolatile() || bSeq.IsVolatile() {
		return Volatile(s)
	}
	return s
}

// Zip combines two input sequences into a sequence of Pairs where each pair
// contains an item from the first sequence matched up with one from the second
// sequence. When one seqence ends, the zipped sequence stops.
//
// If at least one sequence is volatile then the output sequence will be as
// well.
func Zip[A, B any](aSeq Sequence[A], bSeq Sequence[B]) Sequence[Pair[A, B]] {
	return zipCommon(aSeq, bSeq, true)
}

// ZipLongest is like [Zip], but when the shorter sequence ends the output
// uses a zero value in it's place while the longer sequence continues.
func ZipLongest[A, B any](aSeq Sequence[A], bSeq Sequence[B]) Sequence[Pair[A, B]] {
	return zipCommon(aSeq, bSeq, false)
}

// Unzip takes a sequence of pairs and returns two sequences, one containing
// the first item of each pair, the other containing the second of each pair.
func Unzip[A, B any](s Sequence[Pair[A, B]]) (Sequence[A], Sequence[B]) {
	aS := Map(s, Pair[A, B].A)
	bS := Map(s, Pair[A, B].B)
	return aS, bS
}

// PairSelectA takes a sequence of Pair values, and returns a sequence of just
// the first value from each pair.
func PairSelectA[A, B any](s Sequence[Pair[A, B]]) Sequence[A] {
	return Map(s, Pair[A, B].A)
}

// PairSelectB takes a sequence of Pair values, and returns a sequence of just
// the second value from each pair.
func PairSelectB[A, B any](s Sequence[Pair[A, B]]) Sequence[B] {
	return Map(s, Pair[A, B].B)
}

// PairSwap processes a sequence with Pair[A,B] elements to produce one where
// the elements in the pair are swapped (i.e. Pair[B,A]).
func PairSwap[A, B any](s Sequence[Pair[A, B]]) Sequence[Pair[B, A]] {
	return Map(s, Pair[A, B].Swap)
}
