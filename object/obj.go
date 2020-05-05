package object

import (
//"fmt"
)

type PanObj struct {
	Props []PanObj
}

func (o *PanObj) Type() PanObjType {
	return ""
}

func (o *PanObj) Inspect() string {
	return ""
}

func (o *PanObj) Proto() PanObject {
	return o
}
