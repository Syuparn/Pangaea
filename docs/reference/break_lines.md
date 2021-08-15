# Break lines

Each statement are separated by line breaks.

```pangaea
(2 + 3).p
"Hello".p
```

Also, you can break lines by semicolons. This is useful for one-liner iterators.

```pangaea
<{|i| yield i if i < 100; recur(i*2)}>.new(2).A # [2, 4, 8, 16, 32, 64]
```

## Join multiple lines

You can join multiple lines together by `|` before each chain.

```pangaea
100
  |@{\ if .prime?}
  |.len
  |.even? # false
```

### Why?

The syntax aligns each chain so that you can read it easily.

#### Rejected syntax

##### Join lines by special characters at the end of each line

It is impossible because Pangaea syntax already uses all ascii special characters.

```
# REJECTED SYNTAX
# "\" is already used for 1st-argument variable
1 + \
2 + \
3
```

#### Break lines by chains

If line breaks could be inserted before chains, the syntax would collide with anonymous chain.

```
# REJECTED SYNTAX
# `hoge.fuga` or `hoge; .fuga` ?
hoge
.fuga
```
