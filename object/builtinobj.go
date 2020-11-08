package object

// initialize built-in objects like Int, Arr, Str...
// NOTE: Props are inserted in package eval not to make
// package object and package BuiltIn circular reference

// BuiltInIntObj is an object of Int (proto of each int).
var BuiltInIntObj = NewPanObj(&map[SymHash]Pair{}, BuiltInNumObj)

// BuiltInFloatObj is an object of Float (proto of each float).
var BuiltInFloatObj = NewPanObj(&map[SymHash]Pair{}, BuiltInNumObj)

// BuiltInNumObj is an object of Num.
var BuiltInNumObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInNilObj is an object of Nil (proto of nil).
var BuiltInNilObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInStrObj is an object of Str (proto of each str).
var BuiltInStrObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInArrObj is an object of Arr (proto of each arr).
var BuiltInArrObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInRangeObj is an object of Range (proto of each range).
var BuiltInRangeObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInFuncObj is an object of Func (proto of each func).
var BuiltInFuncObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInIterObj is an object of Iter (proto of each iter).
var BuiltInIterObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInMatchObj is an object of Match (proto of each match).
var BuiltInMatchObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInObjObj is an object of Obj (proto of each obj).
var BuiltInObjObj = NewPanObj(&map[SymHash]Pair{}, BuiltInBaseObj)

// BuiltInBaseObj is an object of BaseObj (ancestor of all objects).
var BuiltInBaseObj = NewPanObj(&map[SymHash]Pair{}, nil)

// BuiltInMapObj is an object of Map (proto of each map).
var BuiltInMapObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInIOObj is an object of IO (proto of each io).
var BuiltInIOObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInDiamondObj is an object of Diamond.
var BuiltInDiamondObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInKernelObj is an object of Kernel, whose props can be used in top-level.
var BuiltInKernelObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInIterableObj is an object of Iterable, which is mixed-in iterable objects.
var BuiltInIterableObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInOneInt is an int object `1`.
var BuiltInOneInt = &PanInt{1}

// BuiltInZeroInt is an int object `0`.
var BuiltInZeroInt = &PanInt{0}

// BuiltInTrue is a bool object `true`.
var BuiltInTrue = &PanBool{true}

// BuiltInFalse is a bool object `false`.
var BuiltInFalse = &PanBool{false}

// BuiltInNil is a nil object `nil`.
var BuiltInNil = &PanNil{}

// BuiltInErrObj is an object of Err (proto of all specific err types).
var BuiltInErrObj = NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

// BuiltInAssertionErr is an object of AssertionErr (proto of each assertionErr).
var BuiltInAssertionErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInNameErr is an object of NameErr (proto of each nameErr).
var BuiltInNameErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInNoPropErr is an object of NoPropErr (proto of each noPropErr).
var BuiltInNoPropErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInNotImplementedErr is an object of NotImplemented (proto of each notImplementdErr).
var BuiltInNotImplementedErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInStopIterErr is an object of StopIterErr (proto of each stopIterErr).
var BuiltInStopIterErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInSyntaxErr is an object of SyntaxErr (proto of each syntaxErr).
var BuiltInSyntaxErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInTypeErr is an object of TypeErr (proto of each typeErr).
var BuiltInTypeErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInValueErr is an object of ValueErr (proto of each valueErr).
var BuiltInValueErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInZeroDivisionErr is an object of ZeroDivisionErr (proto of each zeroDivisionErr).
var BuiltInZeroDivisionErr = NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)

// BuiltInNotImplemented is an object of _
var BuiltInNotImplemented = NewNotImplementedErr("Not implemented")
