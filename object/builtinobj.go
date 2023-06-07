package object

// initialize built-in objects like Int, Arr, Str...
func init() {
	// set zero values of {} (refer itself)
	zeroObj.zero = zeroObj

	// NOTE:
	// 1. Props are inserted in evaluator not to make package object and evaluator circular reference
	// 2. initialized objects must be set to the pointer allocated in the statement
	//    because other objects may refer the pointer
	*BuiltInArrObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj, WithZero(zeroArr))
	*BuiltInBaseObj = *NewPanObj(&map[SymHash]Pair{}, nil)
	*BuiltInComparableObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInDiamondObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInEitherObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInEitherErrObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInEitherObj)
	*BuiltInEitherValObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInEitherObj)
	*BuiltInFloatObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInNumObj, WithZero(zeroFloat))
	*BuiltInFuncObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInIntObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInNumObj, WithZero(BuiltInZeroInt))
	*BuiltInIOObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInIterObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInIterableObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInJSONObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInKernelObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInMapObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj, WithZero(zeroMap))
	*BuiltInMatchObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInNilObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj, WithZero(BuiltInNil))
	*BuiltInNumObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInObjObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInBaseObj)
	*BuiltInRangeObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj, WithZero(zeroRange))
	*BuiltInStrObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj, WithZero(zeroStr))
	*BuiltInWrappableObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)

	*BuiltInErrObj = *NewPanObj(&map[SymHash]Pair{}, BuiltInObjObj)
	*BuiltInAssertionErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInNameErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInNoPropErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInNotImplementedErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInStopIterErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInSyntaxErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInTypeErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInValueErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
	*BuiltInZeroDivisionErr = *NewPanObj(&map[SymHash]Pair{}, BuiltInErrObj)
}

// constants

// zeroArr is a zero value of Arr []
var zeroArr = NewPanArr()

// zeroFloat is a zero value of Float 0.0
var zeroFloat = NewPanFloat(0.0)

// zeroMap is a zero value of Map %{}
var zeroMap = NewPanMap()

// zeroObj is a zero value of Obj {}
var zeroObj = &PanObj{
	Pairs:       &map[SymHash]Pair{},
	Keys:        &[]SymHash{},
	PrivateKeys: &[]SymHash{},
	proto:       BuiltInBaseObj,
	// zero value cannot be set here
	// zero: zeroObj,
}

// zeroRange is a zero value of Range (nil:nil:nil)
var zeroRange = NewPanRange(BuiltInNil, BuiltInNil, BuiltInNil)

// zeroStr is a zero value of Str ""
var zeroStr = NewPanStr("")

// BuiltInOneInt is an int object `1`.
var BuiltInOneInt = &PanInt{Value: 1, proto: BuiltInIntObj}

// BuiltInZeroInt is an int object `0`.
var BuiltInZeroInt = &PanInt{Value: 0, proto: BuiltInIntObj}

// BuiltInTrue is a bool object `true`.
var BuiltInTrue = &PanBool{true}

// BuiltInFalse is a bool object `false`.
var BuiltInFalse = &PanBool{false}

// BuiltInNil is a nil object `nil`.
var BuiltInNil = &PanNil{proto: BuiltInNilObj}

// BuiltInNotImplemented is an object of _
var BuiltInNotImplemented = NewNotImplementedErr("Not implemented")

// define built-in objects
// NOTE: built-in objects should be initialized in init() (evaluated after vars),
// otherwise initialization cycle occurs

// BuiltInIntObj is an object of Int (proto of each int).
var BuiltInIntObj = &PanObj{}

// BuiltInFloatObj is an object of Float (proto of each float).
var BuiltInFloatObj = &PanObj{}

// BuiltInNumObj is an object of Num.
var BuiltInNumObj = &PanObj{}

// BuiltInNilObj is an object of Nil (proto of nil).
var BuiltInNilObj = &PanObj{}

// BuiltInStrObj is an object of Str (proto of each str).
var BuiltInStrObj = &PanObj{}

// BuiltInArrObj is an object of Arr (proto of each arr).
var BuiltInArrObj = &PanObj{}

// BuiltInRangeObj is an object of Range (proto of each range).
var BuiltInRangeObj = &PanObj{}

// BuiltInFuncObj is an object of Func (proto of each func).
var BuiltInFuncObj = &PanObj{}

// BuiltInIterObj is an object of Iter (proto of each iter).
var BuiltInIterObj = &PanObj{}

// BuiltInMatchObj is an object of Match (proto of each match).
var BuiltInMatchObj = &PanObj{}

// BuiltInObjObj is an object of Obj (proto of each obj).
var BuiltInObjObj = &PanObj{}

// BuiltInBaseObj is an object of BaseObj (ancestor of all objects).
var BuiltInBaseObj = &PanObj{}

// BuiltInMapObj is an object of Map (proto of each map).
var BuiltInMapObj = &PanObj{}

// BuiltInIOObj is an object of IO (proto of each io).
var BuiltInIOObj = &PanObj{}

// BuiltInDiamondObj is an object of Diamond.
var BuiltInDiamondObj = &PanObj{}

// BuiltInKernelObj is an object of Kernel, whose props can be used in top-level.
var BuiltInKernelObj = &PanObj{}

// BuiltInJSONObj is an object of JSON, whose props can be used in top-level.
var BuiltInJSONObj = &PanObj{}

// BuiltInIterableObj is an object of Iterable, which is mixed-in iterable objects.
var BuiltInIterableObj = &PanObj{}

// BuiltInComparableObj is an object of Comparable, which is mixed-in comparable objects.
var BuiltInComparableObj = &PanObj{}

// BuiltInWrappableObj is an object of Wrappable, which is mixed-in wrappable objects.
var BuiltInWrappableObj = &PanObj{}

// BuiltInEitherObj is an object of Either.
var BuiltInEitherObj = &PanObj{}

// BuiltInEitherValObj is an object of EitherVal.
var BuiltInEitherValObj = &PanObj{}

// BuiltInEitherErrObj is an object of EitherErr.
var BuiltInEitherErrObj = &PanObj{}

// BuiltInErrObj is an object of Err (proto of all specific err types).
var BuiltInErrObj = &PanObj{}

// BuiltInAssertionErr is an object of AssertionErr (proto of each assertionErr).
var BuiltInAssertionErr = &PanObj{}

// BuiltInFileNotFoundErr is an object of FileNotFoundErr (proto of each fileNotFoundErr).
var BuiltInFileNotFoundErr = &PanObj{}

// BuiltInNameErr is an object of NameErr (proto of each nameErr).
var BuiltInNameErr = &PanObj{}

// BuiltInNoPropErr is an object of NoPropErr (proto of each noPropErr).
var BuiltInNoPropErr = &PanObj{}

// BuiltInNotImplementedErr is an object of NotImplemented (proto of each notImplementdErr).
var BuiltInNotImplementedErr = &PanObj{}

// BuiltInStopIterErr is an object of StopIterErr (proto of each stopIterErr).
var BuiltInStopIterErr = &PanObj{}

// BuiltInSyntaxErr is an object of SyntaxErr (proto of each syntaxErr).
var BuiltInSyntaxErr = &PanObj{}

// BuiltInTypeErr is an object of TypeErr (proto of each typeErr).
var BuiltInTypeErr = &PanObj{}

// BuiltInValueErr is an object of ValueErr (proto of each valueErr).
var BuiltInValueErr = &PanObj{}

// BuiltInZeroDivisionErr is an object of ZeroDivisionErr (proto of each zeroDivisionErr).
var BuiltInZeroDivisionErr = &PanObj{}
