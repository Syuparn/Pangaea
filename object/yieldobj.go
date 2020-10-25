package object

// YieldType is a type of YieldObj.
const YieldType = "YieldType"

// YieldObj is a wrapper for yielded value.
type YieldObj struct {
	PanObject
}

// Type returns type of this PanObject.
func (o *YieldObj) Type() PanObjType {
	return YieldType
}
