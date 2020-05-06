package object

// initialize built-in objects like Int, Arr, Str...
// NOTE: Props are inserted in package eval not to make
// package object and package BuiltIn circular reference
var BuiltInIntObj = &PanObj{&map[SymHash]Pair{}, BuiltInNumObj}
var BuiltInFloatObj = &PanObj{&map[SymHash]Pair{}, BuiltInNumObj}
var BuiltInNumObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInBoolObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInNilObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInStrObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInArrObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInRangeObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}

// TODO: consider Proto of Iter!
var BuiltInFuncObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInMatchObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInObjObj = &PanObj{&map[SymHash]Pair{}, BuiltInBaseObj}
var BuiltInBaseObj = &PanObj{&map[SymHash]Pair{}, nil}
