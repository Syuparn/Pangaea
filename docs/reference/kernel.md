# Kernel

`Kernel` is used for a container of top-level functions.
Properties in kernel are assigned to top-level environment variables so that you don't have to specify the receiver `Kernel`.

```pangaea
Kernel.keys(private?: true) # ["argv", "assert", "assertEq", "assertRaises", "_init", "_name"]

# you don't have to write `Kernel['argv]()`
argv() # array of command line arguments
assertEq(2 + 2, 5) # AssertionErr: 4 != 5
```
