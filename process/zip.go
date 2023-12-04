package process

import (
	"github.com/cookieo9/sequence"
	"github.com/cookieo9/sequence/chans"
	"github.com/cookieo9/sequence/convert/build"
	"github.com/cookieo9/sequence/tools"
)

func zipCommon[A, B any](aSeq sequence.Sequence[A], bSeq sequence.Sequence[B], shortest bool) sequence.Sequence[sequence.Pair[A, B]] {
	aCh := chans.AsChanPair(aSeq)
	bCh := chans.AsChanPair(bSeq)
	return sequence.Generate(func(f func(sequence.Pair[A, B]) error) error {
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
			if err := f(sequence.MakePair(aV, bV)); err != nil {
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
func Zip[A, B any](aSeq sequence.Sequence[A], bSeq sequence.Sequence[B]) sequence.Sequence[sequence.Pair[A, B]] {
	return zipCommon(aSeq, bSeq, false)
}

// ZipShortest is like [Zip], but stops once one sequence is empty.
func ZipShortest[A, B any](aSeq sequence.Sequence[A], bSeq sequence.Sequence[B]) sequence.Sequence[sequence.Pair[A, B]] {
	return zipCommon(aSeq, bSeq, true)
}

// FirstOnly takes a sequence of Pair values, and returns a sequence of just
// the first value from each pair.
func FirstOnly[A, B any](s sequence.Sequence[sequence.Pair[A, B]]) sequence.Sequence[A] {
	return Map(s, sequence.Pair[A, B].A)
}

// SecondOnly takes a sequence of Pair values, and returns a sequence of just
// the second value from each pair.
func SecondOnly[A, B any](s sequence.Sequence[sequence.Pair[A, B]]) sequence.Sequence[B] {
	return Map(s, sequence.Pair[A, B].B)
}

// AddFirst creates a sequence where each item is a pair of a constant value
// and a value from a given sequence. The constant value is the first item
// in the pair.
func AddFirst[A, T any](a A, s sequence.Sequence[T]) sequence.Sequence[sequence.Pair[A, T]] {
	return Zip[A, T](build.Infinite(a), s)
}

// AddSecond creates a sequence where each item is a pair of a constant value
// and a value from a given sequence. The constant value is the second item
// in the pair.
func AddSecond[T, B any](s sequence.Sequence[T], b B) sequence.Sequence[sequence.Pair[T, B]] {
	return Zip[T, B](s, build.Infinite(b))
}

// SwapPairs processes a sequence with Pair[A,B] elements to produce one where
// the elements in the pair are swapped (i.e. Pair[B,A]).
func SwapPairs[A, B any](s sequence.Sequence[sequence.Pair[A, B]]) sequence.Sequence[sequence.Pair[B, A]] {
	return Map(s, sequence.Pair[A, B].Swap)
}
