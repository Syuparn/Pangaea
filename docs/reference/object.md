# Object

`Obj` represents an object, pairs of a property name and a value. All object literals inherit `Obj`.
objects works as the base of Pangaea object system (See [Inheritance](./inheritance.md)).

```pangaea
# use symbol for property names
{'a: 1, 'b: 2}
# bare symbol can also be used
{a: 1, b: 2}
# operator method (note: bare symbol cannot be used)
{
  '+: m{|other| self.value + other.value},
  value: 1,
}
# if you want to use variable as key, use `^`
name := 'foo
{^name: "1"} # {"foo": "1"}

# indexing
{a: 1, b: 2}['a] # 1
```

Property names are used for property call ([Calls](./calls.md)).

```
obj := {
  a: 1,
  hello: m{"hi!".p},
}
obj.a # 1
obj.hello # hi!
```

Object pairs are sorted by names.

```pangaea
{c: "C", a: "A", d: "D", b: "B"} # {"a": "A", "b": "B", "c": "C", "d": "D"}
```

## Duplicated keys

If a object literal contains duplicate keys, the first one remains.

```pangaea
{a: 1, b: 2, a: 3} # {a: 1, b: 2}
```

### Why first one (not last one)?

This specification is suitable for Mix-in design pattern ([Mix-in](./mix-in.md)).

## Unpacking

Object pairs can be unpacked by `**`.

```pangaea
{a: 1, **{b: 2}, **{c: 3, d: 4}} # {"a": 1, "b": 2, "c": 3, "d": 4}

# unpack keyword arguments
f := {|a: 1, b: 2, c: 3| a*100 + b*10 + c}
f(**{a: 1, b: 2, c: 3}) # 234
```

## Private properties

Property starting with `_` is treated as a private property.
Private properties are ignored by `Obj#keys` and `Obj#_iter`.
For that reason, the private properties are never appeared in list chains ([Chains](./chains.md)).

```pangaea
obj := {_private: "secret!", public: "you can see me"}
obj@p # ["public", "you can see me"]

obj.keys # ["public"]
# obtain private properties
obj.keys(private?: true) # ["public", "_private"]
```

:warning: `_` cannot be used for private property because it is `NotImplementedErr` constant ([Errors](./errors.md))
