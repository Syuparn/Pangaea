# Error handling

## Catch the error

Once an error is raised, the evaluation terminates.
If you want to catch and handle the error, you should use `Obj#try`.

```pangaea
5.try.{|n| n / 0}.val # nil (error occurred)
5.try.{|n| n / 0}.err # [ZeroDivisionErr: cannot be divided by 0]
```

`Obj#try` wraps the receiver by `Either`. As the name implies, it contains either of an evaluated value or an error.

```pangaea
#                      value, error
5.try.{|n| n / 2}.A # [2.500000, nil]
5.try.{|n| n / 0}.A # [nil, [ZeroDivisionErr: cannot be divided by 0]]
```

If either contains an error, it ignores any calls.
You can handle all errors at the end of the whole chain.

```pangaea
# any call to eitherErr is no operation

# error not raised
20.try.{|n| n + 5}.{|n| 100 // n}.sqrt.A.{|v, err| return err.p if err; v.p} # 2.0
#            err
"a".try.{|n| n + 5}.{|n| 100 // n}.sqrt.A.{|v, err| return err.p if err; v.p} # [TypeErr: 5 cannot be treated as str]
#                       err
-5.try.{|n| n + 5}.{|n| 100 // n}.sqrt.A.{|v, err| return err.p if err; v.p} # [ZeroDivisionErr: cannot be divided by 0]
#                                 err
-6.try.{|n| n + 5}.{|n| 100 // n}.sqrt.A.{|v, err| return err.p if err; v.p} # [ValueErr: sqrt of -100 is not a real number]
```

Error objects generated by either can be treated as normal objects (no errors are raised).

```pangaea
6.try.{|n| n / 0}.err.kindOf?(ZeroDivisionErr) # true
```

### Why is either introduced?

Ordinal try-catch statements is not suitable for Pangaea one-liners.
Also, either keeps chain flow even though error is raised inside.

### Thoughtful chain vs Either

While thoughtful chain `~` just ignores the raised error, either can handle caught errors.

## Raise caught error object again

Error objects generated by either is just an object.
If you want to raise it again, use `raise` statement.

```pangaea
err := 6.try.{|n| n / 0}.err
raise err
```