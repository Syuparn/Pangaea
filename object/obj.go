package object

import (
	"bytes"
)

const OBJ_TYPE = "OBJ_TYPE"

type PanObj struct {
	Pairs *map[SymHash]Pair
}

func (o *PanObj) Type() PanObjType {
	return OBJ_TYPE
}

func (o *PanObj) Inspect() string {
	var out bytes.Buffer
	pairs := []Pair{}

	// NOTE: refer map because range cannot treat map pointer
	for _, p := range *o.Pairs {
		pairs = append(pairs, p)
	}

	out.WriteString("{")
	// NOTE: sort by key order otherwise output changes randomly
	// depending on inner map structure
	out.WriteString(sortedPairsString(pairs))
	out.WriteString("}")

	return out.String()
}

func (o *PanObj) Proto() PanObject {
	return builtInObjObj
}
