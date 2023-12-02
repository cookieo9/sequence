package tools

// Compose builds a composite function from the to input functions f, and g.
// The returned function computes g(f(x)).
func Compose[X, Y, Z any](f func(X) Y, g func(Y) Z) func(X) Z {
	return func(x X) Z { return g(f(x)) }
}
