# Map

`Map` represents an associative array. All map literals inherit `Map`.

```pangaea
%{"a": 1, "b": 2}
# any object can be a map key
%{1: "one", "two": 2, [1, 2, 3]: "array", {a: 1}: {a: "obj"}}
# indexing
%{1: "one", "two": 2}["two"] # 2
```

Map keeps pair order.

```pangaea
%{"a": 1, 2: 3, [4, 5]: 6}.A # [["a", 1], [2, 3], [[4, 5], 6]]
```

:memo: Due to constraints of implementation, hashable keys(int, string, float, nil) precedes the others.

```pangaea
%{1: "one", [2]: "two", "three": 3, {"fo": "ur"}: 4}.A # [[1, "one"], ["three", 3], [[2], "two"], [{"fo": "ur"}, 4]]
```

## What's the difference between Obj and Map?

Map is pairs of key and value. On the other hand, Obj is pairs of property name and property value.
All properties in object are guaranteed to be called by property call ([Calls](./calls.md)).
While map is a key-value container designed just for indexing.

## Duplicated keys

If a map literal contains duplicate keys, the first one remains (same as object literal).

```pangaea
%{"a": 1, "b": 2, "a": 3} # %{"a": 1, "b": 2}
```

## Unpacking

Map pairs can be unpacked by `**`.

```pangaea
%{"a": 1, **%{"b": 2}, **%{"c": 3, "d": 4}} # %{"a": 1, "b": 2, "c": 3, "d": 4}
```
