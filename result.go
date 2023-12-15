package sequence

// A Result represents a (Value,error) tuple where both could be present. It
// will panic if the value is accessed alone when the error is non-nil.
type Result[T any] struct {
	value T
	err   error
}

// Pair returns both the value stored and error (if any). It does not panic.
func (r Result[T]) Pair() (T, error) {
	return r.value, r.err
}

// Error returns the error stored in this result (if any) or nil.
func (r Result[T]) Error() error {
	return r.err
}

// HasError returns true if there's an error present in this Result.
func (r Result[T]) HasError() bool {
	return r.err != nil
}

// Value returns the value stored in this result (if any). If there is an error
// stored in the result, this method will panic, otherwise the value is
// returned or a zero value if none was stored.
func (r Result[T]) Value() T {
	if r.HasError() {
		panic(r.err)
	}
	return r.value
}

// DropError returns a new result copied from the current result, with the
// error explicity set to nil. This should only be used if the error can safely
// be ignored or has already been dealt with.
func (r Result[T]) DropError() Result[T] {
	r.err = nil
	return r
}

// Clean returns a new result where the presence of an error means that the
// value is the zero value.
func (r Result[T]) Clean() Result[T] {
	if r.HasError() {
		r.value = *new(T)
	}
	return r
}

// ResultValue creates a result using just the provided value, assuming a nil
// error.
func ResultValue[T any](value T) Result[T] {
	return MakeResult(value, nil)
}

// ResultError creates a result with just an error. One may need to pass a type
// argument, since it's likely the generic type inference won't detect the
// the needed value type of the Result.
func ResultError[T any](err error) Result[T] {
	return MakeResult(*new(T), err)
}

// MakeResult is the base constructor for a result where both the value and
// error are specified. Useful for wrapping a (T,error) returning function.
func MakeResult[T any](value T, err error) Result[T] {
	return Result[T]{value: value, err: err}
}

// NextResult accepts a result, and a function that will process the value
// stored within, to produce a new result value. If the input result has
// an error, no work is done other than propogating the error to the output
// result value.
func NextResult[In, Out any](in Result[In], f func(In) (Out, error)) Result[Out] {
	if in.HasError() {
		return ResultError[Out](in.Error())
	}
	return MakeResult(f(in.Value()))
}
