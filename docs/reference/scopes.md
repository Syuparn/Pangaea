# Scopes

Pangaea uses lexical scope. Closures can refer outer function's variables.

```pangaea
{|outer|
  {|inner| inner + outer}
}(2)(3) # 5
```

Note that compound assignments do not affect outer scope variables. This is just shadowing (inner scope variable `a := 2` shadows outer scope variable `a := 1`).

```pangaea
a := 1
{a += 1; a.p}() # 2
a.p # 1
```
