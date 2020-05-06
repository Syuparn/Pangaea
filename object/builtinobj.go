package object

// initialize built-in objects like Int, Arr, Str...
// NOTE: Props are inserted in package eval not to make
// package object and package builtin circular reference
var builtInIntObj = &PanObj{&map[SymHash]Pair{}, builtInNumObj}
var builtInFloatObj = &PanObj{&map[SymHash]Pair{}, builtInNumObj}
var builtInNumObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInBoolObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInNilObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInStrObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInArrObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInRangeObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}

// TODO: consider Proto of Iter!
var builtInFuncObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInMatchObj = &PanObj{&map[SymHash]Pair{}, builtInObjObj}
var builtInObjObj = &PanObj{&map[SymHash]Pair{}, builtInBaseObj}
var builtInBaseObj = &PanObj{&map[SymHash]Pair{}, nil}
