package object

// ReturnType is a type of ReturnObj.
const ReturnType = "ReturnType"

// ReturnObj is a wrapper for returned value.
type ReturnObj struct {
	PanObject
}

// Type returns type of this PanObject.
func (o *ReturnObj) Type() PanObjType {
	return ReturnType
}
