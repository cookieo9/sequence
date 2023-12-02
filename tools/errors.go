package tools

import "fmt"

// Must accepts a value/error pair, and returns the value as long as the error
// is nil. If there is a non-nil error Must will panic with it.
func Must[T any](value T, err error) T {
	Check(err)
	return value
}

// Check accepts an error and will panic with it if it's non-nil.
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

// Annotate adds extra context to an error message, essentially via a call to
// fmt.Errorf("%s: %w", ctx, err), but if the error was nil, no wrapped error
// is prodduced and nil is returned, allowing for clean returns from functions.
func Annotate(err error, ctx string) error {
	if err == nil {
		return err
	}
	return fmt.Errorf("%s: %w", ctx, err)
}

// Annotatef is like annotate, but accepting Printf style arguments to generate
// the context message. Like with Annotate, no work is done if the error was
// nil, and a nil is simply returned.
func Annotatef(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return Annotate(err, fmt.Sprintf(format, args...))
}

// CleanErrors wraps the return values from a function of the form (value,error)
// and ensures only one of the value or error is returned. If the error is
// non-nil, then a zero value is returned for value. This makes it easier
// to have "clean" return values where it doesn't make sense for the user
// to proceed when there is an error.
func CleanErrors[T any](value T, err error) (T, error) {
	if err != nil {
		return *new(T), err
	}
	return value, nil
}
