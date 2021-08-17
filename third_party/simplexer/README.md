simplexer
=========

[![Build Status](https://travis-ci.org/macrat/simplexer.svg?branch=master)](https://travis-ci.org/macrat/simplexer)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2208d3ac4fbcdcd2b78a/test_coverage)](https://codeclimate.com/github/macrat/simplexer/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/2208d3ac4fbcdcd2b78a/maintainability)](https://codeclimate.com/github/macrat/simplexer/maintainability)
[![GoDoc](https://godoc.org/github.com/macrat/simplexer?status.svg)](https://godoc.org/github.com/macrat/simplexer)

A simple lexical analyzser for Go.

## example
### simplest usage
``` go
package main

import (
	"fmt"
	"strings"

	"github.com/macrat/simplexer"
)

func Example() {
	input := "hello_world = \"hello world\"\nnumber = 1"
	lexer := simplexer.NewLexer(strings.NewReader(input))

	fmt.Println(input)
	fmt.Println("==========")

	for {
		token, err := lexer.Scan()
		if err != nil {
			panic(err.Error())
		}
		if token == nil {
			fmt.Println("==========")
			return
		}

		fmt.Printf("line %2d, column %2d: %s: %s\n",
			token.Position.Line,
			token.Position.Column,
			token.Type,
			token.Literal)
	}
}
```

It is output as follow.

``` text
hello_world = "hello world"
number = 1
==========
line  0, column  0: IDENT: hello_world
line  0, column 12: OTHER: =
line  0, column 14: STRING: "hello world"
line  1, column  0: IDENT: number
line  1, column  7: OTHER: =
line  1, column  9: NUMBER: 1
==========
```

### more examples
Please see [godoc](https://godoc.org/github.com/macrat/simplexer).
