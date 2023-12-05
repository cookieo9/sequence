package sequence

import (
	"github.com/cookieo9/sequence/tools"
)

func zipCommon[A, B any](aSeq Sequence[A], bSeq Sequence[B], shortest bool) Sequence[Pair[A, B]] {
	aCh := ToChanPair(aSeq)
	bCh := ToChanPair(bSeq)
	return Generate(func(f func(Pair[A, B]) error) error {
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
}

// Zip combines two input sequences into a sequence of Pairs where each pair
// contains an item from the first sequence matched up with one from the second
// sequence. When one seqence ends, the zipped sequence continues producing
// Pairs where one element is the zero value until both sequences are finished.
func Zip[A, B any](aSeq Sequence[A], bSeq Sequence[B]) Sequence[Pair[A, B]] {
	return zipCommon(aSeq, bSeq, false)
}

// ZipShortest is like [Zip], but stops once one sequence is empty.
func ZipShortest[A, B any](aSeq Sequence[A], bSeq Sequence[B]) Sequence[Pair[A, B]] {
	return zipCommon(aSeq, bSeq, true)
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
