# How to run

## Setup

You only need to download binary from [Releases](https://github.com/Syuparn/Pangaea/releases). You can also try it online at [Pangaea Playground](https://syuparn.github.io/Pangaea/).

## Binary usage
### Run Script

```bash
$ pangaea ./hello.pangaea
Hello, world!
```

### REPL

```
$ pangaea
Pangaea 0.6.0
multi : multi-line mode
single: single-line mode (default)

>>> "Hello, world!"
"Hello, world!"
>>> 3 * 5
15
```

In multi-line mode, you can paste multi-line script.

```
$ pangaea
Pangaea 0.6.0
multi : multi-line mode
single: single-line mode (default)

>>> multi
nil
<< multi-line mode (read lines until empty line is found) >>
# convert multi-line string
`
Lorem
ipsum
dolor
sit
amet
`.split(sep:"\n")@capital@p

Lorem
Ipsum
Dolor
Sit
Amet
[]
```

#### Prettify evaluated values

If REPL's output is long and verbose, consider using `_name` prop.

```pangaea
# too long...
> asciis := (?a:?{)@({}){[\,.ord]}
{"a": 97, "b": 98, "c": 99, ..., "z": 122}

# add _name
> asciis := {_name: "asciis", **asciis}
asciis
```

### One-liner

With `-e` command, you can execute one-liner script.

```bash
$ pangaea -e '"Hello, world!".p'
Hello, world!
```

See `pangaea -h` for details.

### Jargon File

If you write the same scripts frequently, *jargon* file will help you.
Objects written in jargon file is evaluated beforehand and can be referred in REPL or one-liner.

With `-j` command, a script file defined by environment variable `$PANGAEA_JARGON_FILE` is evaluated by jargon (`~/.jargon.pangaea` by default). 

```.jargon.pangaea
colorcode := {|s| \[1:].chunk(2)@join("")@I(base:16).{"rgb".zip(\)}.A}
```

```bash
$ pangaea -j
>>> colorcode("#2e4d55")
[["r", 46], ["g", 77], ["b", 85]]
```

See `pangaea -h` for details.
