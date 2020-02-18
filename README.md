# Pangaea programming language
A programming language for one-liner method chain lovers! (Under construction...)

# Requirements
## Host language
- Golang

## Packages
- [goyacc](https://godoc.org/golang.org/x/tools/cmd/goyacc)
- [simplexer](github.com/macrat/simplexer)

# Plans for language features
(These are not implemented yet...)

## One-way!
This language is tuned for a one-liner method chain!
You don't have to go back to beginning of line!

```
"Hello, world!".puts # Hello, world!
1:5.to_a.sum.puts # 15
```

Looks similar to other language though?
But Chains in Pangaea has more power...

## Chain context
Dot chain is "one of" the method chain in Pangaea.
There are some kinds of chain styles, and each one shows different "context".
(The concept is from Perl :) )

There are 3 kinds of chain context(`.`, `@`, `$`).

### Scalar Chain
A receiver is left-side value, which is ordinary method chain.

```
10.puts # 10
```

### List Chain
A receiver is **each element of** left-side value.
This can be used as "map" or "filter" in other languages.

```
[1, 2, 3]@{|i| i * 2}.puts # [2, 4, 6]
["foo", "var", "hoge"]@capital.puts # ["Foo", "Var", "Hoge"]
1:10@{|i| i if i.even?}.puts # [2, 4, 6, 8]
```

### Reduce Chain
A receiver is **each element of** left-side value.
Also, returned value of previous call is passed to 2nd argument.
(In short, it's reduce!)

```
# reduce chain can hold initial value.
[1, 2, 3]$(0){|acc, i| acc+i}.puts # 6
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
# nil.capital.puts # NoPropErr: nil does not have property "capital"
nil&.capital.puts # nil

[1, 2, nil, 4]&@to_f.puts # [1.0, 2.0, nil, 4.0]
```

#### Thoughtful Chain
This chain returns receiver instead if returned value is `nil`
(it "thoughtfully" repairs failed call).

```
1:16~@{|i| [`fizz][i%3] + [`buzz][i%5]}.puts # [1, 2, `fizz, 4, `buzz, ..., `fizzbuzz]

3:20~$([2]){|acc, n| [*acc, n] if acc.all? {|p| n % p}}.puts # [2, 3, 5, ..., 19]

# (Of course there is also an embedded prime function)
20@filter(prime?).puts # [2, 3, 5, ..., 19]
```

#### Strict Chain
This chain keeps returned `nil` value ("strictly" returns the calclation result).
This is useful only in list context, which removes returned `nil`.

```
1:10@{|i| i if i.even?}.puts # [2, 4, 6, 8]
1:10=@{|i| i if i.even?}.puts # [nil, 2, nil, 4, nil, 6, nil, 8, nil]
```
