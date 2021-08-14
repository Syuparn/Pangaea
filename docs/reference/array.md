# Array

`Arr` represents an array. All array literal inherits `Arr`.

```pangaea
[1, 2, 3]
# array can contain any objects
[1, "two", %{"three": (4:5:6)}, [7]]
# indexing
["a", "b", "c"][1] # "b"
```

## Unpacking

Array elements can be unpacked by `*`.

```pangaea
[1, 2, *[3, 4], 5, *[6, 7]] # [1, 2, 3, 4, 5, 6, 7]

# unpack arguments
f := {|a, b, c| a*100 + b*10 + c}
f(*[2, 3, 4]) # 234
```

In literal call ([Calls](./calls.md)), a receiver array is unpacked automatically if arity is more than 1.

```pangaea
[2, 3, 4].{|a, b, c| a*100 + b*10 + c} # 234
[2, 3, 4].{|a| "#{a} is an array"} # "[2, 3, 4] is an array"
```
