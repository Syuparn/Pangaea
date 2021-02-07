# Pangaea programming language
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)
![](https://github.com/Syuparn/Pangaea/workflows/Test/badge.svg?branch=master)
[![codecov](https://codecov.io/gh/Syuparn/Pangaea/branch/master/graph/badge.svg)](https://codecov.io/gh/Syuparn/Pangaea)

A programming language for one-liner method chain lovers! (Under construction...)

# Run

```bash:
# 1. Build (dependent tools are installed automatically)
$ go generate
$ go build

# 2. Run

# Run REPL
# (Linux, Mac)
$ ./pangaea
# (Windows)
$ ./pangaea.exe

# Run script file
# (Linux, Mac)
$ ./pangaea ./example/hello.pangaea
# (Windows)
$ ./pangaea.exe ./example/hello.pangaea

# 4. Enjoy!
```

# Requirements
## Host language
- Golang (1.15+)

## Packages

- [goyacc](https://godoc.org/golang.org/x/tools/cmd/goyacc)
- [simplexer](github.com/macrat/simplexer)
- [dtoa](github.com/tanaton/dtoa)
- [statik](github.com/rakyll/statik)

# Progress

- [x] Lexer
- [x] Parser
- [x] Evaluator
- [ ] Methods/Properties (about 20%)

# language features (Let's run your REPL!)

## One-way!
This language is tuned for a one-liner method chain!
You don't have to go back to beginning of line!

```
"Hello, world!".puts # Hello, world!
(1:5).A.sum.puts # 10
```

Looks similar to other language though?
But Chains in Pangaea has more power...

## Chain context
Dot chain is "one of" the method chain in Pangaea.
There are some kinds of chain styles, and each one shows different "context".
(The concept is from Perl :) )

There are 3 kinds of chain context(`.`, `@`, `$`).

### Scalar Chain
The receiver is left-side value, which is ordinary method chain.

```
10.puts # 10
```

### List Chain
The receiver is **each element of** left-side value.
This can be used as "map" or "filter" in other languages.

```
[1, 2, 3]@{|i| i * 2}.puts # [2, 4, 6]
["foo", "var", "hoge"]@capital.puts # ["Foo", "Var", "Hoge"]
# select only evens because nils are ignored
(1:10)@{|i| i if i.even?}.puts # [2, 4, 6, 8]
```

### Reduce Chain
The receiver is **each element of** left-side value.
Also, returned value of previous call is passed to 2nd argument.
(In short, it's reduce!)

```
# reduce chain can hold initial value.
[1, 2, 3]$(0){|acc, i| acc+i} # 6
# same as above
[1, 2, 3]$(0)+ # 6
```

### Additional context
Additional context can be prepended by main chain context.
There are 3 kinds of additional chain context(`&`, `=`, `~`).
Thus, there are 9 kinds (3 additional * 3 main) of context.

#### Lonely Chain
This chain ignores call and return `nil` if its receiver is `nil` (what a "lonely" object!),
which works same as "lonely operator" in Ruby.

```
# nil.capital.puts # NoPropErr: property `capital` is not defined.
nil&.capital.puts # nil

[1, 2, nil, 4]&@F.puts # [1.000000, 2.000000, 4.000000]
```

#### Thoughtful Chain
This chain returns receiver instead if returned value is `nil`
(it "thoughtfully" repairs failed call).

```
(1:16)~@{|i| ['fizz][i%3] + ['buzz][i%5]}.puts # [1, 2, "fizz", 4, "buzz", ..., "fizzbuzz"]

(3:20)~$([2]){|acc, n| [*acc, n] if acc.all? {|p| n % p}}.puts # [2, 3, 5, ..., 19]

# (Of course you can use built-in prime function)
20.select {.prime?}.puts # [2, 3, 5, ..., 19]
```

#### Strict Chain
This chain keeps returned `nil` value ("strictly" returns the calclation result).
This is useful only in list context, which removes returned `nil`.

```
(1:10)@{|i| i if i.even?}.puts # [2, 4, 6, 8]
(1:10)=@{|i| i if i.even?}.puts # [nil, 2, nil, 4, nil, 6, nil, 8, nil]
```
