# Loop

Since Pangaea has list chain ([Chains](./chains.md)), there are no loop statements such as `for` or `while`.

```pangaea
["foo", "bar", "hoge"]@capital # ["Foo", "Bar", "Hoge"]
```

Traditional for loop `for(i=0;i<n;i++)` can be replated with int list chain.

```pangaea
3@{|i| (i * 2).p}
# 2
# 4
# 6

# loop elements with indices
["foo", "bar", "hoge"].withI@{|i, e| "#{i}: #{e}".p}
# 0: foo
# 1: bar
# 2: hoge
```
