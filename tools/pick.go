package tools

// Pick returns trueValue if condition is true, falseValue otherwise. This
// allows for simple selection of a value based on a condition that would
// normally result in a 4 line (or more) code block due to the lack of a
// ternary operator in Go.
//
// Unlike the ?: operator in C, since this is still a function, any computation
// needed to create the values passed as trueValue and falseValue are likely
// still performed as it's not guaranteed to be a short-circuit. Inlining may
// still produce benefits, but any side effects of producing the values will
// always happen, regardless of which value is selected by the condition.
func Pick[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}
