package tools

// Integer is a type constraint for all integral types.
type Integer interface {
	~uintptr | ~uint | ~int |
		~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~int8 | ~int16 | ~int32 | ~int64
}

// Real is a type constraint for all floating point real types.
type Real interface {
	~float32 | ~float64
}

// Complex is a type constraint for all floating point complex types.
type Complex interface {
	~complex64 | ~complex128
}

// Arithmetic is a type constraint for all types that support the basic
// arithmetic operations.
type Arithmetic interface {
	Integer | Real | Complex
}

// Add is a binary operation returning the sum of two Arithmetic values.
func Add[T Arithmetic](a, b T) T { return a + b }

// Mul is a binary operation returning the product of two Arithmetic values.
func Mul[T Arithmetic](a, b T) T { return a * b }

// Concat is a binary operation that concatenates two string values.
func Concat[T ~string](a, b T) T { return a + b }

// BitOr is a binary operation that does a bitwise or of two integral values.
func BitOr[T Integer](a, b T) T { return a | b }

// BitAnd is a binary operation that does a bitwise and of two integral values.
func BitAnd[T Integer](a, b T) T { return a & b }
