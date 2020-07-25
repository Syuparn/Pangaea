package object

// initialize built-in objects like Int, Arr, Str...
// NOTE: Props are inserted in package eval not to make
// package object and package BuiltIn circular reference
var BuiltInIntObj = &PanObj{&map[SymHash]Pair{}, BuiltInNumObj}
var BuiltInFloatObj = &PanObj{&map[SymHash]Pair{}, BuiltInNumObj}
var BuiltInNumObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInNilObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInStrObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInArrObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInRangeObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}

var BuiltInFuncObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInIterObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInMatchObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInObjObj = &PanObj{&map[SymHash]Pair{}, BuiltInBaseObj}
var BuiltInBaseObj = &PanObj{&map[SymHash]Pair{}, nil}
var BuiltInMapObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}
var BuiltInIOObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}

var BuiltInOneInt = &PanInt{1}
var BuiltInZeroInt = &PanInt{0}

var BuiltInTrue = &PanBool{true}
var BuiltInFalse = &PanBool{false}
var BuiltInNil = &PanNil{}

// error hierarchy
var BuiltInErrObj = &PanObj{&map[SymHash]Pair{}, BuiltInObjObj}

var BuiltInAssertionErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInNameErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInNoPropErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInNotImplementedErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInStopIterErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInSyntaxErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInTypeErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInValueErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
var BuiltInZeroDivisionErr = &PanObj{&map[SymHash]Pair{}, BuiltInErrObj}
