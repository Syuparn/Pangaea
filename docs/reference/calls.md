# Calls

Pangaea has 3 ways of call expression. All calls can use any of chain context ([Chains](./chains.md)).

|name|expression|description|
|-|-|-|
|property call|`foo.bar(arg)`|calls a property of receiver with arguments|
|literal call|<code>foo.{&#124;arg&#124; #...}</code>|evaluates literal function|
|variable call|`foo.^f`|works similar to literal call but uses a variable|

## Property call

Property call works similarly to the ordinary method call and property reference.
It calls a property of receiver with arguments.

```pangaea
obj := {
  foo: 1,
  hello: m{|o| "Hello, #{o}. I am obj!".p}
}
obj.foo.p # 1
obj.hello("bar") # Hello, bar. I am obj!
```

If arguments are passed to a property which is not a method (callable objects), they are just ignored.

```pangaea
obj.foo("verbose", "arg") # 1
```

If method requires only one argument (self), `()` can be omitted.

```pangaea
chocolate := {sweet?: m{true}}

# same as chocolate.sweet?()
chocolate.sweet? # true
```

If you want to obtain a method property itself, use indexing instead.

```pangaea
f := chocolate['sweet?]
f() # true
```

### Why can `()` be omitted only when arity is 1?

The syntax sugar is introduced to identify accessor method calls with property values.

```pangaea
# even? is actually a method, but it seems 2 has boolean property `even?`
2.even?
```

But `()` cannot be ommitted if 2 or more arguments are passed, otherwise one-liner may not work as you expect.

```pangaea
# REJECTED SYNTAX

# ambiguous expression
"abc,def,ghi".split sep: ","@capital
# "abc,def,ghi".split(sep: ",")@capital
# "abc,def,ghi".split(sep: ","@capital)`
```

### Why do method call and property reference use the same mechanism?

To allow an accessor method's `()` to be omitted (see above).
Accessor method and property encapsulation are another option to realize that, but it was a bit over-engineering for Pangaea immutable objects.

```pangaea
# REJECTED SYNTAX

# only methods can be called (foo cannot be referred from outside)
# instead, accessor method foo (returns property foo) is defined implicitly
obj := {
  foo: 1,
  hello: m{|o| "Hello, #{o}. I am obj!".p}
}

# call method
obj.hello("bar")
# call accessor method
obj.foo

# problem 1: no benefits to encapsulate properties because they cannot be modified
# problem 2: property and accessor method conflicts each other
obj['foo] # 1 or accesssor?
```

### Prototype chains

If the receiver does not have the specified property(method), its prototype's property is called.

If `BaseObj`(ancestor of all objects) does not have the specified property, `_missing` method is called instead.

(See [Inheritance](./inheritance.md) and [Object System](./object_system.md) for details about object systems)

```
order of property search along the prototype chain

e.g: `obj.foo`

- receiver's `foo` property
- receiver's prototype's `foo` property
- receiver's prototype's prototype's `foo` property
...
- BaseObj's `foo` property
- receiver's `_missing` property
- receiver's prototype's `_missing` property
...
- BaseObj's `_missing` property
- raise an NoPropErr
```

```pangaea
# example
# method `bear` is not defined in "a"
"a".which('bear) # BaseObj
# but bear can be called because of prototype chain
"a".bear({prop: "foo"}) # {"prop": "foo"}
```

## Literal call

Literal call evaluates a literal function with one argument, the receiver.
This idea is originated from *pipelines* in functional programming languages.

```pangaea
3.{|n| n * 2} # 6
```

Special variables helps to shorten literal call functions (see [Variables](./variables.md)).

```pangaea
3.{\ * 2} # 6
```

If the receiver is an array and the arity of literal call function is more than one,
each element of the array is assigned to the each parameter.

```pangaea
[1, 2, 3].{|a, b, c| a * 100 + b * 10 + c}.p # 123
```

## Variable call

Variable call works similar to literal call, but it uses a variable.

```pangaea
f := {|n| n * 2}
3.^f # 6
```
