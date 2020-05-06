package object

import (
//"fmt"
)

const OBJ_TYPE = "OBJ_TYPE"

type PanObj struct {
	Pairs *map[SymHash]Pair
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
