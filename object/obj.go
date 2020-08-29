package object

import (
	"bytes"
	"sort"
)

const OBJ_TYPE = "OBJ_TYPE"

func NewPanObj(pairs *map[SymHash]Pair, proto PanObject) *PanObj {
	return &PanObj{
		Pairs: pairs,
		Keys:  keyHashes(pairs),
		proto: proto,
	}
}

func PanObjInstance(pairs *map[SymHash]Pair) PanObj {
	return PanObj{
		Pairs: pairs,
		Keys:  keyHashes(pairs),
		proto: BuiltInObjObj,
	}
}

// NOTE: for making new PanObj instance pointer
// because `&(NewPanObjInstance(...))` is syntax error
func PanObjInstancePtr(pairs *map[SymHash]Pair) PanObject {
	i := PanObjInstance(pairs)
	return &i
}

func EmptyPanObjPtr() *PanObj {
	i := PanObjInstance(&map[SymHash]Pair{})
	return &i
}

func ChildPanObjPtr(proto PanObject, src *PanObj) *PanObj {
	// share pairs with src because objects are immutable
	i := PanObj{Pairs: src.Pairs, proto: proto}
	return &i
}

type PanObj struct {
	Pairs *map[SymHash]Pair
	Keys  *[]SymHash
	proto PanObject
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
	return o.proto
}

func keyHashes(pairs *map[SymHash]Pair) *[]SymHash {
	keyStrs := []string{}

	for _, pair := range *pairs {
		str, ok := pair.Key.(*PanStr)
		// must be ok (obj keys are str)
		if !ok {
			continue
		}
		if str.IsPublic {
			keyStrs = append(keyStrs, str.Value)
		}
	}

	sort.Strings(keyStrs)

	hashes := []SymHash{}
	for _, str := range keyStrs {
		hashes = append(hashes, GetSymHash(str))
	}

	return &hashes
}
