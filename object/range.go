package object

import ()

const RANGE_TYPE = "RANGE_TYPE"

type PanRange struct {
	Start PanObject
	Stop  PanObject
	Step  PanObject
}

func (r *PanRange) Type() PanObjType {
	return ""
}

func (r *PanRange) Inspect() string {
	return ""
}

func (r *PanRange) Proto() PanObject {
	return r
}
