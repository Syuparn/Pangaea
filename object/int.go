package object

const INT_TYPE = "INT_TYPE"

type PanInt struct {
	value int64
}

func (i *PanInt) Type() PanObjType {
	return INT_TYPE
}

func (i *PanInt) Inspect() string {
	return ""
}
