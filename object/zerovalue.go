package object

// NOTE: zero value should be initialized in init() (evaluated after vars),
// otherwise initialization cycle occurs
// i.e.) zeroArr.proto requires BuiltInArrObj and BuiltInArrObj.zero requires zeroArr
func init() {
	// set protos
	zeroArr.proto = BuiltInArrObj
	zeroMap.proto = BuiltInMapObj
	zeroRange.proto = BuiltInRangeObj
	zeroFloat.proto = BuiltInFloatObj

	BuiltInOneInt.proto = BuiltInIntObj
	BuiltInZeroInt.proto = BuiltInIntObj
	BuiltInNil.proto = BuiltInNilObj

	// set zero values
	zeroObj.zero = zeroObj
}

// used as zero value
// NOTE: constractors cannot be used, otherwise initialization cycle occurs

var zeroArr = &PanArr{}
var zeroFloat = &PanFloat{Value: 0.0}
var zeroMap = &PanMap{
	HashKeys:         &[]HashKey{},
	Pairs:            &map[HashKey]Pair{},
	NonHashablePairs: &[]Pair{},
}

var zeroObj = &PanObj{
	Pairs:       &map[SymHash]Pair{},
	Keys:        &[]SymHash{},
	PrivateKeys: &[]SymHash{},
	proto:       BuiltInBaseObj,
	// zero value cannot be set here
	// zero: zeroObj,
}

var zeroRange = &PanRange{Start: BuiltInNil, Stop: BuiltInNil, Step: BuiltInNil}
var zeroStr = NewPanStr("")
