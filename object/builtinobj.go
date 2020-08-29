package object

// initialize built-in objects like Int, Arr, Str...
// NOTE: Props are inserted in package eval not to make
// package object and package BuiltIn circular reference
var BuiltInIntObj = NewPanObj(&map[SymHash]Pair{}, BuiltInNumObj)
var BuiltInFloatObj = NewPanObj(&map[SymHash]Pair{}, BuiltInNumObj)
var BuiltInNumObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInNilObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInStrObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInArrObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInRangeObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

var BuiltInFuncObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInIterObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInMatchObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInObjObj = NewPanObj(&map[SymHash]Pair{}, BuiltInBaseObj)
var BuiltInBaseObj = NewPanObj(&map[SymHash]Pair{}, nil)
var BuiltInMapObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
var BuiltInIOObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

var BuiltInOneInt = &PanInt{1}
var BuiltInZeroInt = &PanInt{0}

var BuiltInTrue = &PanBool{true}
var BuiltInFalse = &PanBool{false}
var BuiltInNil = &PanNil{}

// error hierarchy
var BuiltInErrObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

var BuiltInAssertionErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInNameErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInNoPropErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInNotImplementedErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInStopIterErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInSyntaxErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInTypeErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInValueErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
var BuiltInZeroDivisionErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
