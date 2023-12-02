//go:build go1.22

package tools

import "cmp"

// Or is an alias for cmp.Or from go 1.22 and beyond.
func Or[T comparable](items ...T) T {
	return cmp.Or(items...)
}
