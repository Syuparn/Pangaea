package object

import (
//"fmt"
)

const FLOAT_TYPE = "FLOAT_TYPE"

type PanFloat struct {
	Value float64
}

func (f *PanFloat) Type() PanObjType {
	return FLOAT_TYPE
}

func (f *PanFloat) Inspect() string {
	return ""
}

func (f *PanFloat) Proto() PanObject {
	return builtInFloatObj
}
