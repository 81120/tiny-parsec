// Package parser provides utility types for working with optional values and tuples.
package parser

// Maybe is a generic type that represents an optional value.
// It can either hold a value (Just) or be empty (Nothing).
type Maybe[T any] struct {
	// value is a pointer to the underlying value. If nil, the Maybe is Nothing.
	value *T
}

// Just creates a new Maybe instance that holds a value.
// It takes a value of type T and returns a Maybe[T] containing that value.
func Just[T any](value T) Maybe[T] {
	return Maybe[T]{value: &value}
}

// Nothing creates a new Maybe instance that represents the absence of a value.
// It returns a Maybe[T] with a nil value pointer.
func Nothing[T any]() Maybe[T] {
	return Maybe[T]{value: nil}
}

// Get retrieves the value from the Maybe instance.
// If the Maybe is Nothing, it returns the zero value of type T.
func (o Maybe[T]) Get() T {
	if o.value == nil {
		// Return the zero value of type T if the Maybe is Nothing.
		var zero T
		return zero
	}
	// Return the underlying value if the Maybe is Just.
	return *o.value
}

// IsJust checks if the Maybe instance contains a value.
// It returns true if the Maybe is Just, false otherwise.
func (o Maybe[T]) IsJust() bool {
	return o.value != nil
}

// IsNothing checks if the Maybe instance is empty.
// It returns true if the Maybe is Nothing, false otherwise.
func (o Maybe[T]) IsNothing() bool {
	return o.value == nil
}

// Tuple is a generic type that represents a pair of values.
// It holds two values of types U and V, accessible via the First and Second fields.
type Tuple[U, V any] struct {
	// First is the first value in the tuple.
	First U
	// Second is the second value in the tuple.
	Second V
}

// NewTuple creates a new Tuple instance with the given values.
// It takes two values of types U and V and returns a Tuple[U, V] containing those values.
func NewTuple[U, V any](first U, second V) Tuple[U, V] {
	return Tuple[U, V]{First: first, Second: second}
}
