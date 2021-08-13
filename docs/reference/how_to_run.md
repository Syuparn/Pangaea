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

### One-liner

With `-e` command, you can execute one-liner script.

```bash
$ pangaea -e '"Hello, world!".p'
Hello, world!
```

See `pangaea -h` for details.
