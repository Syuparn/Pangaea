# Inheritance

Each Pangaea object inherits one object (*prototype*).

## Create a new child object

New inherited objects can be created by `bear` method.

```pangaea
parent := {a: 1}
# inherit parent
child := parent.bear({b: 2})

# get prototype by proto method
child.proto == parent # true
```

Object literal is a syntax sugar of `Obj.bear`.

```pangaea
# they are equivalent
parent := {a: 1}
parent := Obj.bear({a: 1})
```

## Prototype chains

If an object does not have the same property, prototype's one is used instead ([Calls](./calls.md)).

```pangaea
parent := {a: 1, b: 2}
child := parent.bear({b: 10, c: 20})

# call child's property
child.c # 20
# if child does not have the property, parent's one is used
child.a # 1
# if both child and parent have the property, child's one is used
child.b # 10
# if parent does not have the property, parent's prototype's one is used
# (`keys` is defined in Obj)
child.keys # ["b", "c"]
```

You can find the specific property's owner by `Obj#which`.

```pangaea
child.which('a) == parent # true
child.which('c) == child # true
child.which('keys) == Obj # true
```

## Constructors

`Obj#bear` is powerful but buggy to create children with specific properties.
You can use `Kernel._init` to create object constructors instead.

```pangaea
Person := {
  # make a constructor of Person
  new: _init('name, 'age),
  canDrink?: m{.age >= 20},
  hello: m{"I am #{.name}".p},
}

taro := Person.new("Taro", 20)
jiro := Person.new("Jiro", 18)

taro.hello # I am Taro
jiro.hello # I am Jiro
taro.canDrink? # true
jiro.canDrink? # false
```

## Update objects

Since Pangaea objects are immutable, state changes are described as new objects.
`Obj#bro` generates updated new object (`bro` means *brother*, because self and the new value have the same prototype).

```pangaea
Account := {
  new: _init('name, 'deposit),
  withdraw: m{|amount| .bro({deposit: .deposit - amount})},
}

account := Account.new("Taro", 10000)
account.deposit # 10000
account := account.withdraw(3000)
account.deposit # 7000
```

### Why `bro`(not `bear` or constructor)?

#### using `bear`

`bear` generates a self's child. If `deposit` is called twice, returned value is a prototype of a prototype of `account`. This makes prototype chains verbose.

#### using constructor

Using `Account.new` in `deposit` seems good. But it does not work in `Account`'s children.

```pangaea
Account := {
  new: _init('name, '_deposit),
  # use new instead of bro
  # NOTE: if you use self.new, the problem of bear comes up again
  withdraw: m{|amount| Account.new(.name, ._deposit - amount)},
}

# inherit Account
GreatAccount := Account.bear({
  deposit: m{|amount| GreatAccount.new(.name, ._deposit - amount)},
})

account := GreatAccount.new("Taro", 10000)
account := account.deposit(2000)
account._deposit # 12000
account := account.withdraw(3000) # withdraw returns Account!
account := account.deposit(2000) # NoPropErr: property `deposit` is not defined.
```

## Zero values

Since Pangaea does not have classes, prototypes should be able to use as values.
For that reason, prototypes of literal objects can be used as zero values.

```pangaea
# Int is the prototype of integers. So Int itself should be an integer!
1.proto # Int
# Int behaves as integer zero value 0
Int + 3 # 3
# Str behaves as string zero value ""
"abc" + Str # "abc"
```
