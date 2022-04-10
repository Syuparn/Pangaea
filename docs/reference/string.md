# String

`Str` represents a string. All string literals inherit `Str`.

```pangaea
"hello"
# 1-character string
?a
# raw string
`
hello
world
`
# embedded string
"2 + 3 = #{2 + 3}" # 2 + 3 = 5
# symbol (see below)
'abc
```

`""` string controls escape sequences, while raw string <code>``</code> does not.

```pangaea
"a\nb".p
# a
# b
`a\nb`.p
# a\nb
```

## Symbol

Symbol `'foo` is used as property's key.
Object's properties are called by its key ([Object](./object.md)).
For that reason, symbol literal allows only valid property names.

```pangaea
'abc
'foo?
'foo!
'_private
'+
```

Symbol can also be used as function. This is handy for methods which require function args.
This is inspired by Ruby's `Symbol#to_proc`.

```pangaea
10.exclude {|i| i.even?} # [1, 3, 5, 7, 9]

# same as above ( 'foo behaves as {.foo} )
10.exclude('even?)
```

### Why are there so many literal forms?
Because they are neccessary!

#### raw string

Raw string prevents from increasing backslashes.

```pangaea
# print `\n` literally
"\\n".p
# evaluate above by eval
"\"\\\\n\".p".eval
# with raw string
`"\\n".p`.eval
```

#### 1-character string

1-character string is handy for `split`, `sub`, or whatever you golf.

```pangaea
"abc,def,ghi" / ?, # ["abc", "def", "ghi"]
"abc".sub(?b,?d) # "adc"
```

#### symbol

While the other string literals represents string value,
symbol represents property keys (or names).
If there were no symbols, it is impossible to prevent from defining malformed properties.

```pangaea
obj := {"weird name property!!!!": m{"weird".p}}
obj.weird name property!!!! # syntax error...
```
