# Range

`Range` represents a range. All range literal inherits `Range`.
A triplet in a range means *start*, *stop*, and *step* respectively.
Range is usually used for indexing.

```pangaea
(1:10:2)
('e:'a:-1)
(1:10) # (1:10:nil)
(::-1) # (nil:nil:-1)

# () can be omitted in () or []
"hello, world!"[1:10:2] # "el,wr"
```

## Why is `()` neccessary?

If `()` can be omitted, nested range gets ambiguous.

```pangaea
# REJECTED SYNTAX

1:2:3:4 # what is it?
# possible results
# (1:2:(3:4))
# ((1:2):(3:4))
# (1:(2:3:4))
# (1:(2:(3:4)))
# ...
```
