# Chains

Chains are an expanded concept of `.` in OOP language method call syntax `foo.bar()`.
Pangaea can describe how to deal with receiver by chains (*chain context*).

## Chain context
Dot chain `foo.bar` is just *one of* the method chains in Pangaea.
There are some kinds of chain styles, and each one shows different "context".

|chain|name|receiver of the call|
|-|-|-|
|`foo.bar`|scalar chain|left-side value|
|`foo@bar`|list chain|each element of left-side value iterator|
|`foo$bar`|reduce chain|an accumulated value and each element of left-side value iterator|

### Scalar Chain
The receiver is left-side value, which is ordinary method chain.

```pangaea
10.p # 10
```

### List Chain
The receiver is **each element of** left-side value iterator.
This can be used as "map" or "filter" in other languages.

```pangaea
[1, 2, 3]@{|i| i * 2}.p # [2, 4, 6]
["foo", "var", "hoge"]@capital.p # ["Foo", "Var", "Hoge"]
```

In list chains, evaluated `nil` elements are ignored.
Combining to if expression, you can filter generated elements ([If](./if.md)).

```pangaea
# nil elements are ignored (because `i if i.even? == nil` if `i.even?` is false)
(1:10)@{|i| i if i.even?}.p # [2, 4, 6, 8]
```

### Reduce Chain
The receiver is **each element of** left-side value iterator.
Also, returned value of previous call is passed to 2nd argument
(In short, it's reduce!).

```pangaea
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
This chain ignores call and return `nil` if its receiver is `nil`.

```pangaea
# nil.capital.puts # NoPropErr: property `capital` is not defined.
nil&.capital.puts # nil

[1, 2, nil, 4]&@F.puts # [1.000000, 2.000000, 4.000000]
```

#### Thoughtful Chain
This chain returns receiver instead if returned value is `nil`.

```pangaea
(1:16)~@{|i| ['fizz][i%3] + ['buzz][i%5]}.puts # [1, 2, "fizz", 4, "buzz", ..., "fizzbuzz"]

(3:20)~$([2]){|acc, n| [*acc, n] if acc.all? {|p| n % p}}.puts # [2, 3, 5, ..., 19]

# (Of course you can use built-in prime function)
20.select {.prime?}.puts # [2, 3, 5, ..., 19]
```

#### Strict Chain
This chain keeps returned `nil` value.
This is useful only in list context, which removes returned `nil`.

```pangaea
(1:10)@{|i| i if i.even?}.puts # [2, 4, 6, 8]
(1:10)=@{|i| i if i.even?}.puts # [nil, 2, nil, 4, nil, 6, nil, 8, nil]
```

### Why is method call `.` expanded?

Because that's why Pangaea was made!
A main purpose of Pangaea design is "Can chain context shorten one-liner effectively?".
(In some cases, `@` and `$` is shorter than `map`, `filter`, and `reduce`)

## Chain argument

Chains can take 1 argument for specific use.

### Chain argument in list chain

List chain can only generate an array by default.
Chain argument enables to convert generated array into specific types.

```pangaea
# arr by default
(?a:?d)@{|c| [c, .uc]} # [["a", "A"], ["b", "B"], ["c", "C"]]
# convert to obj
(?a:?d)@({}){|c| [c, .uc]} # {"a": "A", "b": "B", "c": "C"}
# convert to map
(?a:?d)@(%{}){|c| [c, .uc]} # %{"a": "A", "b": "B", "c": "C"}
```

Technically, the chain argument's `digest` method is called to convert the evaluated array ([Metaprogramming](./metaprogramming.md)).

### Chain argument in reduce chain

Chain argument is used for an initial accumulator.
If no arguments are passed, the initinal value is `nil`.

```pangaea
# initial accumulator: "initial"
(?a:?e)$("initial")+ # "initialabcd"
# initial accumulator: nil
(?a:?e)$+ # "abcd" (note that nil + "a" == "a")
```
