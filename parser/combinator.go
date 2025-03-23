// Package parser provides a set of combinators for building parsers.
package parser

// Fmap applies a function to the result of a parser.
// It takes a parser p of type T and a function f that maps T to U,
// and returns a new parser that produces a result of type U.
func Fmap[T, U any](p Parser[T], f func(T) U) Parser[U] {
	return NewParser(func(s string) ParserFuncRet[U] {
		m := p.Parse(s)
		if m.IsNothing() {
			return Nothing[Tuple[U, string]]()
		}
		t := m.Get()
		return Just(NewTuple(f(t.First), t.Second))
	})
}

// Bind sequences two parsers where the second parser depends on the result of the first.
// It takes a parser p of type T and a function f that maps T to a parser of type U,
// and returns a new parser that produces a result of type U.
func Bind[T, U any](p Parser[T], f func(T) Parser[U]) Parser[U] {
	return NewParser(func(s string) ParserFuncRet[U] {
		m := p.Parse(s)
		if m.IsNothing() {
			return Nothing[Tuple[U, string]]()
		}
		t := m.Get()
		return f(t.First).Parse(t.Second)
	})
}

// OrElse tries a sequence of parsers in order and returns the result of the first successful one.
// It takes a variable number of parsers of type T and returns a new parser of type T.
func OrElse[T any](ps ...Parser[T]) Parser[T] {
	return NewParser(func(s string) ParserFuncRet[T] {
		for _, p := range ps {
			m := p.Parse(s)
			if m.IsJust() {
				return m
			}
		}
		return Nothing[Tuple[T, string]]()
	})
}

// ZeroOrMore matches zero or more occurrences of a parser.
// It takes a parser p of type T and returns a new parser that produces a slice of type T.
func ZeroOrMore[T any](p Parser[T]) Parser[[]T] {
	return NewParser(func(s string) ParserFuncRet[[]T] {
		m := p.Parse(s)
		if m.IsNothing() {
			return Just(NewTuple([]T{}, s))
		}
		t := m.Get()
		return Bind(ZeroOrMore(p), func(ts []T) Parser[[]T] {
			return Pure(append([]T{t.First}, ts...))
		}).Parse(t.Second)
	})
}

// OneOrMore matches one or more occurrences of a parser.
// It takes a parser p of type T and returns a new parser that produces a slice of type T.
func OneOrMore[T any](p Parser[T]) Parser[[]T] {
	return NewParser(func(s string) ParserFuncRet[[]T] {
		m := p.Parse(s)
		if m.IsNothing() {
			return Nothing[Tuple[[]T, string]]()
		}
		t := m.Get()
		return Bind(ZeroOrMore(p), func(ts []T) Parser[[]T] {
			return Pure(append([]T{t.First}, ts...))
		}).Parse(t.Second)
	})
}

// ZeroOrOne matches zero or one occurrence of a parser.
// It takes a parser p of type T and returns a new parser that produces a Maybe type of T.
func ZeroOrOne[T any](p Parser[T]) Parser[Maybe[T]] {
	return NewParser(func(s string) ParserFuncRet[Maybe[T]] {
		m := p.Parse(s)
		if m.IsNothing() {
			return Just(NewTuple(Nothing[T](), s))
		}
		t := m.Get()
		return Just(NewTuple(Just(t.First), t.Second))
	})
}

// OmitLeft runs two parsers in sequence and discards the result of the first.
// It takes a parser p of type T and a parser q of type U, and returns a new parser of type U.
func OmitLeft[T, U any](p Parser[T], q Parser[U]) Parser[U] {
	return Bind(p, func(_ T) Parser[U] {
		return q
	})
}

// OmitRight runs two parsers in sequence and discards the result of the second.
// It takes a parser p of type T and a parser q of type U, and returns a new parser of type T.
func OmitRight[T, U any](p Parser[T], q Parser[U]) Parser[T] {
	return Bind(p, func(t T) Parser[T] {
		return Bind(q, func(_ U) Parser[T] {
			return Pure(t)
		})
	})
}

// SepBy parses a sequence of elements separated by a separator.
// It takes a parser p of type T and a parser sep of type U, and returns a new parser that produces a slice of type T.
func SepBy[T, U any](p Parser[T], sep Parser[U]) Parser[[]T] {
	return OrElse(
		Bind(p, func(first T) Parser[[]T] {
			return Fmap(
				ZeroOrMore(Bind(sep, func(_ U) Parser[T] { return p })),
				func(rest []T) []T { return append([]T{first}, rest...) },
			)
		}),
		Pure([]T{}),
	)
}

// Satisfy parses a single rune that satisfies a given predicate.
// It takes a function f that tests a rune and returns a new parser that produces a rune.
func Satisfy(f func(rune) bool) Parser[rune] {
	return NewParser(func(s string) ParserFuncRet[rune] {
		if len(s) == 0 {
			return Nothing[Tuple[rune, string]]()
		}
		if f(rune(s[0])) {
			return Just(NewTuple(rune(s[0]), s[1:]))
		}
		return Nothing[Tuple[rune, string]]()
	})
}

// SatisfyWith applies a predicate to the result of a parser and succeeds only if the predicate is true.
// It takes a parser p of type T and a function f that tests T, and returns a new parser of type T.
func SatisfyWith[T any](p Parser[T], f func(T) bool) Parser[T] {
	return Bind(p, func(t T) Parser[T] {
		if f(t) {
			return Pure(t)
		}
		return Fail[T]()
	})
}

// TrimLeft removes leading whitespace from the result of a parser.
// It takes a parser p of type T and returns a new parser of type T.
func TrimLeft[T any](p Parser[T]) Parser[T] {
	return OmitLeft(Spaces(), p)
}

// TrimRight removes trailing whitespace from the result of a parser.
// It takes a parser p of type T and returns a new parser of type T.
func TrimRight[T any](p Parser[T]) Parser[T] {
	return OmitRight(p, Spaces())
}

// Trim removes leading and trailing whitespace from the result of a parser.
// It takes a parser p of type T and returns a new parser of type T.
func Trim[T any](p Parser[T]) Parser[T] {
	return TrimLeft(TrimRight(p))
}

// Seq parses a sequence of parsers in order and returns a slice of their results.
// It takes a variable number of parsers of type T and returns a new parser that produces a slice of type T.
func Seq[T any](ps ...Parser[T]) Parser[[]T] {
	return NewParser(func(s string) ParserFuncRet[[]T] {
		if len(ps) == 0 {
			return Just(NewTuple([]T{}, s))
		}
		return Bind(ps[0], func(t T) Parser[[]T] {
			return Bind(Seq(ps[1:]...), func(ts []T) Parser[[]T] {
				return Pure(append([]T{t}, ts...))
			})
		}).Parse(s)
	})
}

// Between parses a value between two other values and returns the middle value.
// It takes a parser p of type T, a parser q of type U, and a parser r of type V,
// and returns a new parser that produces a result of type U.
func Between[T, U, V any](p Parser[T], q Parser[U], r Parser[V]) Parser[U] {
	return Bind(p, func(_ T) Parser[U] {
		return Bind(q, func(t U) Parser[U] {
			return Bind(r, func(_ V) Parser[U] {
				return Pure(t)
			})
		})
	})
}

// Lazy defers the creation of a parser until it is needed.
// It takes a function f that returns a parser of type T and returns a new parser of type T.
func Lazy[T any](f func() Parser[T]) Parser[T] {
	return NewParser(func(s string) ParserFuncRet[T] {
		return f().Parse(s)
	})
}

// ToString converts the result of a parser to a string.
// It takes a parser p of type T and returns a new parser that produces a string.
func ToString[T rune | []rune](p Parser[T], shouldTrim bool) Parser[string] {
	if shouldTrim {
		p = Trim(p)
	}
	return Fmap(p, func(t T) string {
		return string(t)
	})
}
