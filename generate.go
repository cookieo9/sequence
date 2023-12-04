package sequence

// Generate is used to create a sequence manually from a "sequence function".
// A sequence function represents the for loop that will produce all the values
// of the sequence, passing them to a callback. The callback returns an error
// either on a problem, or because the user of the data wishes to stop early
// (by sending ErrStopIteration). The generator is expected to stop immediately
// on a non-nill error from the callback, perform cleanup, and then return
// a non-nil error (preferably wrapping the error it was given) or the exact
// error it received.
//
// The generator function should not handle ErrStopIteration itself, both for
// simplicity (unless it intends to wrap all errors), and to allow the
// top-level of the sequence processing to see it (or a wrapped version).
func Generate[T any](f func(func(T) error) error) Sequence[T] {
	return Sequence[T]{source: f}
}

// GenerateVolatile is a helper to do the following: Volatile(Generate(f)).
// It creates a volatile sequence from the provided sequence function.
func GenerateVolatile[T any](f func(func(T) error) error) Sequence[T] {
	return Volatile(Generate(f))
}
