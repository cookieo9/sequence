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

// New creates a simple sequence using the provided items.
func New[T any](items ...T) Sequence[T] {
	return FromSlice(items)
}

// Derive creates a new sequence where the properties of the input sequence are
// copied over to the created sequence. This makes sense when permuting a
// sequence as it will usually have the same properties.
//
// For example: using Map to permute a Volatile sequence results in the output
// sequence being Volatile as well.
//
// The "sequence function" passed as parameter is the same as for [Generate],
// and neither its element type nor that of the returned sequence need to match
// the input sequence. The sole purpose of the input sequence is to provide the
// sequence properties for the output.
//
// The properties that are copied include:
//   - volatile
func Derive[Out, In any](input Sequence[In], f func(func(Out) error) error) Sequence[Out] {
	out := Generate(f)
	out.volatile = input.volatile
	return out
}
