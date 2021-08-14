# Input and output

## Input

You can read stdin lines by diamond(`<>`).
This produces each line of stdin like iterators.

```pangaea
# stdin.pangaea

# read a line from stdin
<>.S.{"1st line: #{\}"}.p
# read next line
<>.S.{"2nd line: #{\}"}.p
# any string methods can be used
<>.IA.p

# read rest of lines at once
<>@p
```

```bash
$ cat tmp.txt
abcde
fghij
12 34 56
foo
bar
hoge
$ cat tmp.txt | pangaea ./stdin.pangaea
1st line: abcde
2nd line: fghij
[12, 34, 56]
foo
bar
hoge
```

## Output

`Obj#p` (or alias `Obj#puts`) can be used to write object to stdout.

```pangaea
$ pangaea -e '"Hello, world!".p'
Hello, world!
```

In method `p`, the receiver is inplicitly converted to string as `S` method.

```pangaea
{a:1,b:2}.S # `{"a": 1, "b": 2}`
{a:1,b:2}.p # {"a": 1, "b": 2}
```

If you don't want to mess up stdout by large objects, use `repr` before `p` (see [Metaprogramming](./metaprogramming.md) for more about `repr`).

```
# pritty-printed in REPL
>>> Obj
Obj

# p shows all properties
>>> Obj.p
{"!": {|| [builtin]}, "!==": {|self, other| (!(self === other))}, "===": {|self, other| (((self == other) || .kindOf?(other)) || other.asFor?(self))}, "A": {|self| @{|| \}}, "B": {|| [builtin]}, "S": {|| [builtin]}, "_iter": {|| [builtin]}, "_name": "Obj", "acc": {|self, f, init: nil| ._iter().{|it| <{|acc| ...

# repr converts objects to string printed in REPL
>>> Obj.repr.p
Obj
nil
```
