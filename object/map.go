package object

import (
	"bytes"
	"strings"
)

// MapType is a type of PanMap.
const MapType = "MapType"

// PanMap is object of map literal.
type PanMap struct {
	Pairs            *map[HashKey]Pair
	NonHashablePairs *[]Pair
}

// Type returns type of this PanObject.
func (m *PanMap) Type() PanObjType {
	return MapType
}

// Inspect returns formatted source code of this object.
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

// Proto returns proto of this object.
func (m *PanMap) Proto() PanObject {
	return BuiltInMapObj
}

// NewPanMap returns new map object.
func NewPanMap(pairs ...Pair) *PanMap {
	pairMap := map[HashKey]Pair{}
	nonHashablePairs := []Pair{}

	for _, pair := range pairs {
		hashable, ok := pair.Key.(PanScalar)
		if ok {
			if _, exists := pairMap[hashable.Hash()]; !exists {
				pairMap[hashable.Hash()] = pair
			}
		} else {
			// NOTE: this method DOES NOT check duplicated nonhashable keys
			// because they should be compared by '== method
			nonHashablePairs = append(nonHashablePairs, pair)
		}
	}

	return &PanMap{&pairMap, &nonHashablePairs}
}

// NewEmptyPanMap returns new empty map object.
func NewEmptyPanMap() *PanMap {
	pairMap := map[HashKey]Pair{}
	nonHashablePairs := []Pair{}

	return &PanMap{&pairMap, &nonHashablePairs}
}
