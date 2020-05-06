package object

import (
	"bytes"
	"strings"
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

	ps := []string{}
	for _, p := range *m.NonHashablePairs {
		ps = append(ps, p.Key.Inspect()+": "+p.Value.Inspect())
	}
	if len(ps) > 0 && len(pairs) > 0 {
		out.WriteString(", ")
	}
	out.WriteString(strings.Join(ps, ", "))

	out.WriteString("}")

	return out.String()
}

func (m *PanMap) Proto() PanObject {
	return BuiltInMapObj
}
