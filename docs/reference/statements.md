# Statements

## Expressions and statements
Almost all syntactic elements in Pangaea are expressions, which do not require line breaks. The only exceptions are jump statements.

## Jump statements

Jump statements relates to the end of function evaluation.

### return

`return` ends function evaluation immidiately and returns value to the caller.

```pangaea
{|i|
  return i * 2
  i.p # never evaluated
}(3).p # 6
```

Combined to `if`, you can write guard.

```pangaea
fact := {|n|
  return 1 if n == 0
  n * fact(n - 1)
}

fact(5).p # 120
```

### yield

`yield` is used in iterators.

It sets returns value to the caller but keeps evaluation. Also, `yield if` raises `StopIterErr` if the condition is false.

```pangaea
<{|n|
  yield n if n < 4 # set n to return value but keeps evaluation
  recur(n+1)
}>.new(1)@p
# 1
# 2
# 3
```

### raise

`raise` is used to raise an error.
It works similar to `return` but it unwraps error wrapper and raises it again (see [Error handling](./error_handling.md) for details).

```pangaea
neg := {|n|
  raise TypeErr.new("n must be int") if !(n.kindOf?(Int))
  -n
}
```

### defer

`defer` does not return anything. Instead, it sets expressions which will be evaluated **after** the end of function evaluation.

```pangaea
{|n|
  # defer statements are evaluated after return
  defer "--logging--".p
  defer "n: #{n}".p
  return (n *= 2)
}(3).p
# --logging--
# n: 6
# 6
```

```pangaea
{
  defer "--logging--".p # evaluated after raise
  raise Err.new("bang!")
}()
```
