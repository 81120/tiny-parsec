// Package parser provides a set of combinators for building parsers in Go.
package parser

import (
	"strconv"
	"strings"
)

// Pure creates a parser that always succeeds without consuming any input and returns the given value.
//
// Parameters:
// - val: The value to be returned by the parser.
//
// Returns:
// - A parser that always succeeds and returns the given value.
func Pure[T any](val T) Parser[T] {
	return NewParser(func(s string) Maybe[Tuple[T, string]] {
		return Just(NewTuple(val, s))
	})
}

// Fail creates a parser that always fails without consuming any input.
//
// Returns:
// - A parser that always fails.
func Fail[T any]() Parser[T] {
	return NewParser(func(s string) Maybe[Tuple[T, string]] {
		return Nothing[Tuple[T, string]]()
	})
}

// Char creates a parser that matches a single character if it is equal to the given character.
//
// Parameters:
// - c: The character to match.
//
// Returns:
// - A parser that matches the given character.
func Char(c rune) Parser[rune] {
	return Satisfy(func(r rune) bool {
		return r == c
	})
}

// Str creates a parser that matches a given string at the beginning of the input.
//
// Parameters:
// - str: The string to match.
//
// Returns:
// - A parser that matches the given string.
func Str(str string) Parser[string] {
	return NewParser(func(s string) Maybe[Tuple[string, string]] {
		if strings.HasPrefix(s, str) {
			return Just(NewTuple(str, strings.TrimPrefix(s, str)))
		}
		return Nothing[Tuple[string, string]]()
	})
}

// Digit creates a parser that matches a single digit character.
//
// Returns:
// - A parser that matches a single digit character.
func Digit() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return r >= '0' && r <= '9'
	})
}

// Digits creates a parser that matches one or more digit characters and returns them as a string.
//
// Returns:
// - A parser that matches one or more digit characters.
func Digits() Parser[string] {
	return Fmap(OneOrMore(Digit()), func(rs []rune) string {
		return string(rs)
	})
}

// Alpha creates a parser that matches a single alphabetic character (either uppercase or lowercase).
//
// Returns:
// - A parser that matches a single alphabetic character.
func Alpha() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
	})
}

// Alphas creates a parser that matches one or more alphabetic characters and returns them as a string.
//
// Returns:
// - A parser that matches one or more alphabetic characters.
func Alphas() Parser[string] {
	return Fmap(OneOrMore(Alpha()), func(rs []rune) string {
		return string(rs)
	})
}

// Space creates a parser that matches a single whitespace character (space, tab, or newline).
//
// Returns:
// - A parser that matches a single whitespace character.
func Space() Parser[rune] {
	return Satisfy(func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n'
	})
}

// Spaces creates a parser that matches zero or more whitespace characters and returns them as a string.
//
// Returns:
// - A parser that matches zero or more whitespace characters.
func Spaces() Parser[string] {
	return Fmap(ZeroOrMore(Space()), func(rs []rune) string {
		return string(rs)
	})
}

// Symbol creates a parser that matches a given string surrounded by optional whitespace.
//
// Parameters:
// - str: The string to match.
//
// Returns:
// - A parser that matches the given string surrounded by optional whitespace.
func Symbol(str string) Parser[string] {
	return Trim(Str(str))
}

// Sign creates a parser that matches an optional sign character ('+' or '-') and returns it.
// If no sign is present, it returns '+'.
//
// Returns:
// - A parser that matches an optional sign character.
func Sign() Parser[rune] {
	return Bind(
		ZeroOrOne(OrElse(Char('-'), Char('+'))),
		func(r Maybe[rune]) Parser[rune] {
			if r.IsNothing() {
				return Pure('+')
			}
			return Pure(r.Get())
		})
}

// IntegerWithoutSign creates a parser that matches one or more digits and returns them as an integer.
//
// Returns:
// - A parser that matches one or more digits and returns them as an integer.
func IntegerWithoutSign() Parser[int64] {
	return Fmap(Digits(), func(digits string) int64 {
		i, _ := strconv.ParseInt(digits, 10, 64)
		return i
	})
}

// Integer creates a parser that matches an optional sign followed by one or more digits and returns the resulting integer.
//
// Returns:
// - A parser that matches an optional sign followed by one or more digits.
func Integer() Parser[int64] {
	return Bind(Sign(), func(sign rune) Parser[int64] {
		return Fmap(IntegerWithoutSign(), func(i int64) int64 {
			if sign == '-' {
				return -i
			}
			return i
		})
	})
}

// FloatWithoutSign creates a parser that matches a floating-point number without a sign.
//
// Returns:
// - A parser that matches a floating-point number without a sign.
func FloatWithoutSign() Parser[float64] {
	return Fmap(
		Seq(Digits(), Fmap(Char('.'), func(c rune) string { return string(c) }), Digits()),
		func(strs []string) float64 {
			f, _ := strconv.ParseFloat(strings.Join(strs, ""), 64)
			return f
		})
}

// Float creates a parser that matches an optional sign followed by a floating-point number and returns the resulting float.
//
// Returns:
// - A parser that matches an optional sign followed by a floating-point number.
func Float() Parser[float64] {
	return Bind(Sign(), func(sign rune) Parser[float64] {
		return Fmap(FloatWithoutSign(), func(f float64) float64 {
			if sign == '-' {
				return -f
			}
			return f
		})
	})
}

// String creates a parser that matches a double-quoted string, handling escape sequences.
//
// Returns:
// - A parser that matches a double-quoted string.
func String() Parser[string] {
	return NewParser(func(s string) ParserFuncRet[string] {
		// Check if the input starts with a double quote
		if len(s) == 0 || s[0] != '"' {
			return Nothing[Tuple[string, string]]()
		}
		// Skip the opening double quote
		s = s[1:]

		var b []byte
		escaped := false

		for i := 0; i < len(s); i++ {
			c := s[i]
			if escaped {
				// Append the escaped character
				b = append(b, c)
				escaped = false
			} else if c == '\\' {
				// Mark the next character as escaped
				escaped = true
			} else if c == '"' {
				// Return the parsed string and the remaining input
				return Just(NewTuple(string(b), s[i+1:]))
			} else {
				// Append the character
				b = append(b, c)
			}
		}

		// If no closing double quote is found, return Nothing
		return Nothing[Tuple[string, string]]()
	})
}
