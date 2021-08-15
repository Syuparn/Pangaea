# Object system

Pangaea has prototype-based object system. Each object inherits one object (single inheritance).

## Inheritance tree of built-in objects

- `BaseObj`
    - `Obj`
        - `Arr`
        - `Comparable`
        - `<>` (Diamond)
        - `Either`
        - `Err`
            - `AssertionErr`
            - `NameErr`
            - `NoPropErr`
            - `NotImplementedErr`
                - `_`
            - `StopIterErr`
            - `SyntaxErr`
            - `TypeErr`
            - `ValueErr`
            - `ZeroDivisionErr`
        - `Func`
        - `Iter`
        - `Iterable`
        - `JSON`
        - `Kernel`
        - `Map`
        - `Nil`
            - `nil`
        - `Num`
            - `Float`
            - `Int`
                - `0`
                    - `false`
                - `1`
                    - `true`
        - `Range`
        - `Str`
        - `Wrappable`

Object prototypes can be referred by `ancestors` method.

```pangaea
true.ancestors # [1, Int, Num, Obj, BaseObj]
```

### Why prototype-based?

Prototype-based object system is simpler than class-based one.
You don't have to distinguish between classes and instances.
