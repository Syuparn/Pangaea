package object

import ()

const ARR_TYPE = "ARR_TYPE"

type PanArr struct {
	Elems []PanObject
}

func (a *PanArr) Type() PanObjType {
	return ""
}

func (a *PanArr) Inspect() string {
	return ""
}

func (a *PanArr) Proto() PanObject {
	return a
}
