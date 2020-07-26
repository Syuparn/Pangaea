package object

const YIELD_TYPE = "YIELD_TYPE"

// wrapper for yielded value
type YieldObj struct {
	PanObject
}

func (o *YieldObj) Type() PanObjType {
	return YIELD_TYPE
}
