# Case

`Obj#case` is used to select one value matching to the receiver.

```pangaea
2.case(%{
  1: "one",
  2: "two",
  3: "three",
}) # two
```

Ranges, arrays, and prototypes can be also used as keys (See [Operators](./operators.md) for details).

```pangaea
[1, 4, 100, "a"]@case(%{
  [1, 2, 3]: "small", # match if receiver === [1, 2, 3]
  (5:10): "medium",
  Int: "large",
  Obj: "others",
}).p # ["small", "large", "large", "others"]
```
