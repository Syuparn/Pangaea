# Metaprogramming

## Special methods

To support polymorphism, some methods can control Pangaea syntactic elements.

|method|usage|
|-|-|
|`_literalProxy`|proxies literalcall|
|`_incBy`|used for range indexing (e.g: `(1:10:2)._iter` generates next value by `._incBy(2)`)|
|`_iter`|used to generate iterator of the receiver in list/reduce chain|
|`_missing`|called when no properties found on the prototype chain|
|`_name`|used in REPL for pretty-printing receivers (e.g: `Obj._name == "Obj"`)|
|`at`|used for indexing (`a[b]` is `a.at([b])`)|
|`B`|used to convert condition in if expression to boolean|
|`call`|used for function call (`f()` is `f.call()`)|
|`digest`|used in a list chain argument to convert evaluated value|
|`S`|used in `p` method to convert the receiver to a string|

## Writing codes dynamically

`eval` and `evalEnv` evaluate source codes dynamically.

- `eval`: returns evaluated value
- `evalEnv`: returns the environment after evaluation as an object

```pangaea
`a := 3; b := 2; a + b`.eval # 5
`a := 3; b := 2; a + b`.evalEnv # {"a": 3, "b": 2}
```

