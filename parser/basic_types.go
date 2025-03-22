// Package parser provides the basic types and functions for parsing operations.
package parser

// ParserFuncRet is an alias for Maybe[Tuple[T, string]], representing the return type of a parser function.
// It encapsulates the result of a parsing operation, which may contain a value of type T and the remaining unparsed string.
type ParserFuncRet[T any] = Maybe[Tuple[T, string]]

// ParserFunc is a function type that takes a string as input and returns a ParserFuncRet[T].
// It represents a parsing function that attempts to parse the input string and returns the result.
type ParserFunc[T any] func(string) ParserFuncRet[T]

// Parser is a generic struct that encapsulates a parsing function.
// It provides a unified interface for different parsing operations.
type Parser[T any] struct {
	// Parse is the parsing function that attempts to parse a string and returns a ParserFuncRet[T].
	Parse ParserFunc[T]
}

// NewParser creates a new Parser instance with the given parsing function.
// It takes a ParserFunc[T] as input and returns a Parser[T] instance.
func NewParser[T any](parse ParserFunc[T]) Parser[T] {
	return Parser[T]{Parse: parse}
}
