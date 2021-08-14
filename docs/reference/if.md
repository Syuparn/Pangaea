# If

## If expression

Pangaea does not have if statements. Instead, you can use if expression.

```pangaea
"hi".p if 2 > 0 # hi
"hi".p if -1 > 0 # (nothing printed)
("even" if 3.even? else "odd").p # "odd"
```

If expression without `else` is evaluated as `nil` if condition is falsy.
This can be used as filter for list chains (note that `nil` elements are eliminated by list chains! ([Chains](./chains.md))).

```pangaea
"even" if 3.even? # nil

# filter elements by if
(1:10)@{|i| i if .even?} # [2, 4, 6, 8]
```

### Why not if statements?

If statements requires line breaks. That's annoying for one-liners.

### Why not traditional ternary `a ? b : c`?

`a ? b : c` is much shorter than `b if a else c`.
But this conflicts with range literal `(a:b:c)` and 1-charcter string literal `?a`.

Also, `b if a` syntax is same as jump statement guards ([Statements](./statements.md)).

## Truthy and falsy

Condition clause in if expression allows not only booleans but also any type of objects.
The rule of truthy/falsy value is simple.

- `o` is truthy if `o.B == true`
- `o` is falsy otherwise

As far as built-in objects, zero values are treated as falsy values.

- `0`
- `0.0`
- `""`
- `{}`
- `%{}`
- `[]`
- `nil`

