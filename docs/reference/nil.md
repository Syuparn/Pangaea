# Nil

`nil` is an object representing *nothing*. `nil` inherits `Nil`.
It is returned from functions, if expressions, indexing, and so on.

```pangaea
# empty function returns nil
{||}() # nil
# if cond is false, if expression returns nil
"y" if false # nil
# if index is not found, indexing returns nil
"abc"[10] # nil
```

## Why does not Pangaea raise error instead of nil?

Although nil makes source codes fragile and unsafe,
it does not require any error handling or guards.
Pangaea values comfortability of one-liners.
