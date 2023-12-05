package sequence

import (
	"context"
)

// FromChan produces a sequence from the provided channel. The new sequence is
// volatile since the channel itself can only be iterated over once.
func FromChan[T any](ch <-chan T) Sequence[T] {
	return GenerateVolatile(func(f func(T) error) error {
		for x := range ch {
			if err := f(x); err != nil {
				return err
			}
		}
		return nil
	})
}

// IntoChan sends the contents of the sequence into the channel. If an error
// occurs processing the sequence it will be returned.
func IntoChan[T any](ch chan<- T, s Sequence[T]) error {
	return IntoChanCtx(context.Background(), ch, s)
}

// IntoChanCtx sends the contents of the sequence into the channel. If an error
// occurs processing the sequence it will be returned.
// The provided context will be consulted for an alternate reason to stop
// iteration.
func IntoChanCtx[T any](ctx context.Context, ch chan<- T, s Sequence[T]) error {
	defer close(ch)
	return Each(s)(func(t T) error {
		select {
		case ch <- t:
			return nil
		case <-ctx.Done():
			if err := ctx.Err(); err != nil {
				return err
			}
			return ErrStopIteration
		}
	})
}

// IntoChanPair sends the items from a sequence into a channel of type
// Pair[T,error] for the purpose of passing on the error from the sequence
// should one occur. If such an error happens it will be in a <zero,error>
// element on its own.
func IntoChanPair[T any](ch chan<- Pair[T, error], s Sequence[T]) {
	defer close(ch)
	err := Each(s)(func(t T) error {
		ch <- MakePair[T, error](t, nil)
		return nil
	})
	if err != nil {
		ch <- MakePair(*new(T), err)
	}
}

// ToChan returns a channel that will receive the values from the sequence. If
// the sequence produces an error, then this code will panic.
func ToChan[T any](s Sequence[T]) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		if err := IntoChan(ch, s); err != nil {
			panic(err)
		}
	}()
	return ch
}

// ToChanErr returns 2 channels, the first containing items from the sequence,
// and the second will receive an error if the sequence failed with that error.
func ToChanErr[T any](s Sequence[T]) (<-chan T, <-chan error) {
	ch := make(chan T)
	eCh := make(chan error)
	go func() {
		defer close(ch)
		defer close(eCh)
		if err := IntoChan(ch, s); err != nil {
			eCh <- err
		}
	}()
	return ch, eCh
}

// ToChanPair returns a channel whose elements are Pair[T,error] such that
// errors that arise while processing the sequence will return a <zero,err>
// Pair.
func ToChanPair[T any](s Sequence[T]) <-chan Pair[T, error] {
	ch := make(chan Pair[T, error])
	go IntoChanPair(ch, s)
	return ch
}

// ToChanCtx returns a channel and a context where the channel returns values
// emitted from the input sequence, and the context will be cancelled if the
// sequence returns an error. An input context is used as the basis of the
// generated context, and can be used to cancel processing externally.
func ToChanCtx[T any](ctx context.Context, s Sequence[T]) (<-chan T, context.Context) {
	ch := make(chan T)
	ctx, cncl := context.WithCancelCause(ctx)
	go func() {
		defer close(ch)
		cncl(IntoChanCtx(ctx, ch, s))
	}()
	return ch, ctx
}
