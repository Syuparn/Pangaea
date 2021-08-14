# Errors

Error raises when something goes wrong. Once an error is raised, the evaluation terminates unless it is caught by `Obj#try`.

```
$ pangaea -e '(0 / 0).p; "finish to calculate".p'
ZeroDivisionErr: cannot be divided by 0
line: 1, col: 7
(0 / 0).p; "finish to calculate".p
line: 1, col: 10
(0 / 0).p; "finish to calculate".p
```

## Raise an error

You can raise an error with `raise` statement.

```pangaea
twice := {|n|
  raise TypeErr.new("n must be Int") if !(n.kindOf?(Int))
  n * 2
}
twice(3).p # 6
twice("a") # TypeErr: n must be Int
```

## Error list

|Name|Meaning|
|-|-|
|`AssertionErr`|assertion failed|
|`NameErr`|variable is not defined|
|`NoPropErr`|object does not have the specified property|
|`NotImplementedErr`|the method/property is has not been implemented yet|
|`StopIterErr`|iteration stopped(used for iterables)|
|`SyntaxErr`|source code syntax is wrong|
|`TypeErr`|argument type is not supported|
|`ValueErr`|argument value is not supported|
|`ZeroDivisionErr`|value is divided by zero|

## Error constant `_`

`_` is a child of `NotImplementErr`. This is useful to notify a method is not available now but will be.

```pangaea
obj := {greatMethod: _}
obj.greatMethod # NotImplementedErr: Not implemented
```
