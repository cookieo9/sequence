//go:build !go1.22

package tools

// Or returns the first item from its list of arguments that isn't a zero value
// and returns zero if none are found.
//
// Note: this is likely to be added to the stdlib as cmp.Or in go 1.22.
func Or[T comparable](items ...T) T {
	var zero T
	for _, item := range items {
		if item != zero {
			return item
		}
	}
	return zero
}
