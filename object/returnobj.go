package object

const RETURN_TYPE = "RETURN_TYPE"

// wrapper for returned value
type ReturnObj struct {
	PanObject
}

func (o *ReturnObj) Type() PanObjType {
	return RETURN_TYPE
}
