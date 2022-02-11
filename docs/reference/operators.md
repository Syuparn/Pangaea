# Operators

## Operator list

### Prefix

|operator|meaning|
|-|-|
|`!`|not|
|`+`|plus|
|`-`|minus|
|`~`|bit invert|

### Infix

|operator|meaning|
|-|-|
|`+`|addition|
|`-`|subtraction|
|`*`|multiplication|
|`/`|division|
|`//`|floor division|
|`%`|modulus|
|`**`|power|
|`==`|equal|
|`!=`|not equal|
|`===`|topic-equal|
|`!==`|not topic-equal|
|`<`|less than|
|`>`|greater than|
|`<=`|less or equal|
|`>=`|greater or equal|
|`<=>`|spaceship|
|`<<`|left shift|
|`>>`|right shift|
|`/&`|bit and|
|<code>/&#124;</code>|bit or|
|`/^`|bit xor|
|`&&`|and|
|<code>&#124;&#124;</code>|or|

## Precedence

(from highest to lowest)

|Precedence|
|-|
|indexing `a[b]`|
|grouping `()`|
|calling `a.b`|
|prefix operators|
|chain|
|`**`|
|`*`, `/`, `//`, `%`|
|`+`, `-`|
|`<<`, `>>`|
|`/&`|
|<code>/&#124;</code>, `/^`|
|`<=>`, `==`, `!=`, `<=`, `>=`, `<`, `>`, `===`, `!==`|
|`&&`|
|<code>&#124;&#124;</code>|
|left assign `a := 1`, compound assign `a += 1`|
|right assign `1 => a`|
|jump statements `return x`|
|ternary `a if b else c`|

## Operator is a method

Operators (except `&&` and `||`) are just syntax sugar of oprarator methods.

```pangaea
QuusInt := Int.bear({
  new: _init('n),
  # operator method (bare symbol cannot be used!)
  '+: m{|other| 5 if [self, other].any? {.n >= 57} else .n + other.n},
})

m := QuusInt.new(60)
n := QuusInt.new(40)
(m + n).p # 5
```

Prefix operators are defined as `(operator)%` method.

```pangaea
a := {'-%: "minus"}
-a.p # "minus"
```

### NOTE: why `&&` and `||` are not methods?

If they are methods, short-circuit evaluation does not work because property call is eager evaluation.

```
# REJECTED SYNTAX
false || "never evaluated".p
# above is equivalent to below
false.||("never evaluated".p) # "never evaluated"
```

## Equal `==` vs Topic-Equal `===`

`==` is intended to be used as equivalence. On the other hand, `===` refers wider relationships.

`===` is used in `Obj#case`, which works similar to switch statement ([Case](case.md) for more details).

```pangaea
[1, 4, 100, "a"]@case(%{
  [1, 2, 3]: "small", # match if receiver === [1, 2, 3]
  (5:10): "medium",
  Int: "large",
  Obj: "others",
}).p # ["small", "large", "large", "others"]
```

Technically, `a === b` returns true if any of below is true.

- `a == b`: a is equivalent to b
- `a.kindOf?(b)`: a is a descendant of b (b appears in a's prototype chain).
- `b.asFor?(a)`: predicate b is true as for a.

### examples

```pangaea
# 2 + 3 is equivalent to 5
2 + 3 === 5
# true is a kind of Int (prototype chain: Int -> 1 -> true)
true === Int # true
# as for "pangaea", "^pan.*$" is true (regex "^pan.*$" is matched against "pangaea")
"pangaea" === "^pan.*$" # true
# as for 4, [1, 2, 3] is false (4 is not in [1, 2, 3])
4 === [1, 2, 3] # false
# as for 6, {|n| n % 3 == 0} is true ({|n| n % 3 == 0}(6) == true)
6 === {|n| n % 3 == 0} # true
# as for 7, (1:10:2) is true (4 is in (1:10:2))
7 === (1:10:2) # true
```

:warning: If you like Ruby, be careful that Pangaea's `===` works opposite way!

```ruby
# Ruby
case x
# match if expr === x
when expr
  # ...
end
```

```pangaea
# Pangaea
x.case(%{
  # match if x === expr
  expr: #...
})
```

### Why?

Topic-equal is inspired by Japanese particle `は(wa)`, which works as a *topic marker* in a sentence. Thanks to the particle, not only subjects but also any topics can construct a sentence. If you speak Japanese, translate `a === b` to `a は b だ。` in your mind. 
