# Iterator

`Iter` represents an iterator. All iterator literals inherit `Iter`.
Iterator generates a sequence of values one by one. This is used mainly for list chains and reduce chains ([Chains](./chains.md)).

Iterators look like functions, but it can keep next call arguments inside (iterator is the only mutable object in Pangaea).
This enables to generate next value from current evaluated value.

```pangaea
# generate 1 forever
i1 := <{1}>
# values can be obtained by `next` method
i1.next # 1
i1.next # 1
i1.next # 1

# generate 1,2,3,...
i2 := <{|i|
  # return i
  yield i
  # set arguments for next call
  recur(i + 1)
}>.new(1) # initial arguments can be set by `new`
i2.next # 1
i2.next # 2
i2.next # 3
i2.next # 4

# generate 1 to 3
i3 := <{|i|
  # make the iterator finite by if clause 
  yield i if i <= 3
  recur(i + 1)
}>.new(1)
i3.next # 1
i3.next # 2
i3.next # 3
i3.next # StopIterErr: iter stopped

# use iterator for a list chain
i4 := <{|i|
  yield i if i < 10
  recur(i * 2)
}>.new(2)
i4@p
# 2
# 4
# 8
```

(See [Statements](./statements.md) for details about `yield`)

`new` method generates a new independent iterator.
You don't have to copy & paste iterator definitions.

```pangaea
gen := <{|i|
  yield i
  recur(i + 1)
}>
it1 := gen.new(1)
it2 := gen.new(1)
it1.next # 1
it1.next # 2
it2.next # 1
it2.next # 2
```

If you got a `NameErr` from `Iter#next`, you may forget to set an initial value by `new`.

```pangaea
it := <{|i| yield i; recur i+1}>
it.next # NameErr: name `i` is not defined
```

## Why is iterator designed as a stateful function?

There were 2 more ideas to realize iterators, but they required more syntactic elements than the current design.

### 1. bidirectional jump statements

If jump statement `yield` brought back the evaluation control to the called function,
iterator literal would be replaced with a ordinary function.
But it has 2 problems below.

- requires `for` statement
- makes flow control more complicated

```
# REJECTED SYNTAX
it := {
  # NOTE: list chain `@` cannot be used because it requires iterator!
  # (list chain uses a iterator -> the iterator uses another list chain ->...)
  for(i := 0; i < 3; i++) {
    yield i
  }
}
it.next # 1 (then bring control back to the next line of the yield)
it.next # 2

it2 := {
  for(i := 0; i < 3; i++) {
    # what happens if return appeared before yield?
    return "a" if i == 2
    yield i
  }
}
```

### 2. Channels

If channel operator is introduced,
iterator literal would be replaced with a ordinary function.
Problems are same as idea 1.

```
# REJECTED SYNTAX
it := {|ch|
  # NOTE: list chain `@` cannot be used because it requires iterator!
  # (list chain uses a iterator -> the iterator uses another list chain ->...)
  for(i := 0; i < 3; i++) {
    ch <- i
  }
}

nums := Channel.new
it(nums)
n := <-nums
n # 1
```
