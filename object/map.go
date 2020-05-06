package object

import (
	"bytes"
)

const MAP_TYPE = "MAP_TYPE"

type PanMap struct {
	Pairs            *map[HashKey]Pair
	NonHashablePairs *[]Pair
}

func (m *PanMap) Type() PanObjType {
	return MAP_TYPE
}

func (m *PanMap) Inspect() string {
	var out bytes.Buffer
	pairs := []Pair{}

	// NOTE: refer map because range cannot treat map pointer
	for _, p := range *m.Pairs {
		pairs = append(pairs, p)
	}

	out.WriteString("%{")
	// NOTE: sort by key order otherwise output changes randomly
	// depending on inner map structure
	out.WriteString(sortedPairsString(pairs))
	out.WriteString("}")

	return out.String()
}

func (m *PanMap) Proto() PanObject {
	return BuiltInMapObj
}
