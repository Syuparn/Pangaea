package object

// NOTE: zero value should be initialized in init() (evaluated after vars),
// otherwise initialization cycle occurs
// i.e.) zeroArr.proto requires BuiltInArrObj and BuiltInArrObj.zero requires zeroArr
func init() {
	// set protos
	zeroArr.proto = BuiltInArrObj

	// set zero values
	zeroObj.zero = zeroObj
}

// used as zero value
var zeroArr = &PanArr{}
var zeroFloat = NewPanFloat(0.0)
var zeroMap = NewEmptyPanMap()

// NOTE: EmptyPanObjPtr cannot be used, otherwise initialization cycle occurs
var zeroObj = &PanObj{
	Pairs:       &map[SymHash]Pair{},
	Keys:        &[]SymHash{},
	PrivateKeys: &[]SymHash{},
	proto:       BuiltInBaseObj,
	// zero value cannot be set here
	// zero: zeroObj,
}

var zeroRange = NewPanRange(BuiltInNil, BuiltInNil, BuiltInNil)
var zeroStr = NewPanStr("")