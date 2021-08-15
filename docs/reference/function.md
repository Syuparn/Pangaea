# Function

`Func` represents a function. All function literals inherit `Func`.
Pangaea does not have function statements.

```pangaea
f := {|a, b|
  return a + b
}
f(2, 3) # 5

# return can be omitted (function returns the last evaluated object)
g := {|n| n * 2}
g(4) # 8

# if arguments are more than parameters, they are ignored
h := {|a, b| [a, b]}
h(1, 2, 3) # [1, 2]
# if arguments are less than parameters, `nil`s are passed instead
h(1) # [1, nil]

# empty function returns nil
{||}() # nil
```

Functions can also be used for literal calls ([Calls](./calls.md)).

```pangaea
2.{|n| n * 3} # 6
```

## Why not function statement?

Function statement have complicated syntax stuctures and it is unneccessary for one-liners!

## Keyword parameters

Keyword parameters have default values.

```pangaea
greet := {|name, to: "Pangaea"|
    "Hi, #{to}! I am #{name}.".p
}
greet("Taro") # Hi, Pangaea! I am Taro.
greet("Hanako", to: "John") # Hi, John! I am Hanako.
```

keyword arguments can be inserted anywhere in ordinary positional arguments.

```pangaea
f := {|p1, p2, k1: "k1", k2: "k2"| [p1, p2, k1, k2].p}
f(1, 2, k1: 3, k2: 4) # [1, 2, 3, 4]
f(k1: 3, k2: 4, 1, 2) # [1, 2, 3, 4]
f(k1: 3, 1, k2: 4, 2) # [1, 2, 3, 4]
f(k2: 4, 1, 2, k1: 3) # [1, 2, 3, 4]
```

## Argument unpacking

Arguments can be unpacked by `*` and `**` (See [Array](./array.md) and [Object](./object.md)).

```pangaea
f := {|a, b| a + b}
arr := [2, 3]
f(*arr) # 5
```

## Scopes and closures

Functions have lexical scopes so that they can be nested (See [Scopes](./scopes.md) for details).

```pangaea
{|outer|
  {|inner| inner + outer}
}(2)(3) # 5
```

## Specical variables

Nth argument of function can be referred as special variable `\{n}` without parameter definition (see [Variables](./variables.md) for details).

```pangaea
f := {\1 + \2}
f(3, 4) # 7
```

## Methods

Methods `m{|| }` are just function properties.
Property call uses the receiver as 1st argument and appends the rest specified arguments.

```pangaea
person := {
  name: "John",
  # syntax sugar of {|self, to| ...}
  sayHi: m{|to| "Hello, #{to}. I am #{.name}!"},
}

person.sayHi("Tom") # "Hello, Tom. I am John!"
# call method as ordinary function
person['sayHi](person, "Tom") # "Hello, Tom. I am John!"
```

Since 1st parameter `self` is prepended, anonymous chains refers 1st argument `self`, which is the property call receiver (`.name` above is same as `self.name`. see [Chains](./chains.md) for more details).

### Why is a method a function?

Pangaea keeps callable system simple. Methods can be extracted and used as standalone functions. Functions can be inserted into some objects as methods.
Thanks to the feature, Pangaea identifies method call (property call) with pipeline programming (literal call) ([Calls](./calls.md)).

These were other options to implement methods.

- handle special variable `self` as the receiver
    - method cannot be used as a standalone function (What `self` refers?)
- make methods and function different types
    - type conversion between them may make one-liners longer 
