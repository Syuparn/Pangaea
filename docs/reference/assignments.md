# Assignment

Assignment expression assigns the right-side value to the left-side variable.

```pangaea
a := 2 * 3
a.p # 6

# right-side assignment
"hello" + "world" => b
b.p # helloworld
```

### Why `:=` (not `=`)?

Although `=` is shorter, it would conflict with strict chain.

```
# REJECTED SYNTAX
a = 1
# assignment of anonymous chain (a = (.b)) or strict chain property call (a=.b)?
a = .b
```

### Why `:=` is not an operator (method)?

If `:=` were an operator, it would be a mutable method that updates environment.
It can be abused to use environment as *mutable object*, which is not desirable for Pangaea concept. Instead, the idea "environment as mutable object" is encapsulated as iterators ([Iterator](./iterator.md)).

## Reassignment

You can also reassign values to variables already defined.

```pangaea
a := 2
a.p # 2
a := 3
a.p # 3
```

You can also use compound assignments. Note that this is reassignment, not mutable change.

```pangaea
a := 1
a.p # 1
a += 1
a.p # 2
```

### Why is reassignment permitted?

Reassignment is permitted for compound assignment and shadowing.
Also, reassignment is suitable for REPL.
