# Tiny Parsec

A parser combinator library inspired by Haskell's Parsec, implemented in Go with strong type safety.

## Features

- Type-safe parser combinators
- JSON parser implementation
- Immutable parser state
- Comprehensive combinator library
- Error handling with parser state tracking

## Installation

```bash
go get github.com/81120/tiny-parsec
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/81120/tiny-parsec/parser"
)

func main() {
	p := parser.Trim(parser.Char('a'))
	result := p.Parse(" a ")
	fmt.Println(result) // Just('a')
}
```
