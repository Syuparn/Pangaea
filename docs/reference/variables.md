# Varibles

## Available characters

Regular expression: `[a-zA-Z_][a-zA-Z0-9_]*[!?]?`

## Naming rule used in built-in objects

- both variables and properties are camelCase
    - shorter than snake_case!
    - `NoPropErr`
    - `myObj`
    - `"".evalEnv`
- boolean variables and properties which return boolean end at `?`
    - `appeared?`
    - `2.even?`
- dengerous methods end at `!`
    - `foo.bang!`
- private properties start at `_`
    - properties starting at `_` are ignored by list/reduce chain
    - `Obj._name`
    - `Arr._iter`
- receiver(1st parameter) of method is `self`
    - just use syntax sugar `m{}`!

## Special variables

Nth parameter of function can be referred by `\{n}`. It is handy for single-use functions.
Similarly, keyword parameters can be referred by `\{var}`.

```pangaea
{\1 + \2}(3, 4) # 7
{\a + \b}(a: 1, b: 2) # 3
```

Since `\1` is used very frequently, this can be also referred by `\`.

```pangaea
{\ * 2}(3) # 6
```

If you want to use varargs, use `\0` and `\_`. They contain all arguments/keyword arguments respectively.

```pangaea
{\0}(1, 2, 3) # [1, 2, 3]
{\_}(a: 1, b: 2, c: 3) # {"a": 1, "b": 2, "c": 3}
```

### Why?

This helps one-line short coding! :smile:

## Reserved words

These words cannot be used as variables.

- `if`
- `else`
- `return`
- `raise`
- `yield`
- `defer`
