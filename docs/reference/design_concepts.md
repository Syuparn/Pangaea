# Design concepts

Pangaea is designed for general-purpose programming language. It especially aims at comfortable one-liner script.

## How confortable?

- short coding
    - short name methods
        - e.g: convert to string `1.S # "1"`
    - collection conversion by list/reduce chain (e.g. ``)
        - e.g: names of each person in people `people@name`
- no need to break lines
    - methods designed for chains
        - e.g: `('a:'d).append('f).A.join(?&).p # a&b&c&f`
    - pipeline programming with literalcall
        - e.g: `(1:4)@{\ * 2} # [2, 4, 6]`
- easy to understand
    - all objects are immutable
    - object-oriented

## Features

- Readable one-liner
- Interpreted
- Dynamically typed
- Prototype-based object oriented
- Everything is object
- Immutable objects
- First-class functions with lexical scopes
- Method chains with context ([Chains](./chains.md) for details)
- Metaprogramming with magic methods (e.g: `_missing`, `asFor?`)
