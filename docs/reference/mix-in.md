# Mix-in

Although each object can inherit only one prototype,
you may want to implement common methods in objects.
Mix-in design pattern helps to implement properties outside of prototypes.
This is just an object unpacking ([Object](./object.md)).

```pangaea
Common := {
  log: m{"log: #{self}".p},
}

Person := {
  new: _init('name, 'age),
  hello: m{"I am #{.name}".p},
  # unpack Common to implement `log` 
  **Common,
}

person := Person.new("Taro", 20)
person.log # log: {"age": 20, "name": "Taro"}
```

## Built-in objects to mix-in

Pangaea provides some objects to be mixed-in.
You don't have to implement commonly-used methods in your objects.

|name|required method to use it|providing methods|
|-|-|-|
|`Comparable`|`<=>` returning `-1`, `0`, or `1`|comparsion methods such as `<` or `!=`|
|`Iterable`|`_iter` returning an iterator|iteration methods such as `withI` or `zip`|
|`Wrappable`|`_fmap` handling a wrapper's element|wrapper methods / mixed in `Either`|
