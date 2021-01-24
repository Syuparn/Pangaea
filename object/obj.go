package object

import (
	"bytes"
	"errors"
	"sort"
)

// ObjType is a type of PanObj.
const ObjType = "ObjType"

// NewPanObj makes new PanObj instance.
func NewPanObj(pairs *map[SymHash]Pair, proto PanObject) *PanObj {
	publicKeys, privateKeys := keyHashes(pairs)
	return &PanObj{
		Pairs:       pairs,
		Keys:        &publicKeys,
		PrivateKeys: &privateKeys,
		proto:       proto,
	}
}

// PanObjInstance makes new obj literal.
func PanObjInstance(pairs *map[SymHash]Pair) PanObj {
	publicKeys, privateKeys := keyHashes(pairs)
	return PanObj{
		Pairs:       pairs,
		Keys:        &publicKeys,
		PrivateKeys: &privateKeys,
		proto:       BuiltInObjObj,
	}
}

// PanObjInstancePtr makes new obj literal.
func PanObjInstancePtr(pairs *map[SymHash]Pair) PanObject {
	i := PanObjInstance(pairs)
	return &i
}

// EmptyPanObjPtr makes new empty obj literal.
func EmptyPanObjPtr() *PanObj {
	i := PanObjInstance(&map[SymHash]Pair{})
	return &i
}

// ChildPanObjPtr makes new child object of proto with props in src.
func ChildPanObjPtr(proto PanObject, src *PanObj) *PanObj {
	// share pairs with src because objects are immutable
	i := PanObj{
		Pairs:       src.Pairs,
		Keys:        src.Keys,
		PrivateKeys: src.PrivateKeys,
		proto:       proto,
	}
	return &i
}

// PanObj is object for not only obj literal but also
// any PanObject except specific data structure literal.
type PanObj struct {
	Pairs       *map[SymHash]Pair
	Keys        *[]SymHash
	PrivateKeys *[]SymHash
	proto       PanObject
}

// Type returns type of this PanObject.
func (o *PanObj) Type() PanObjType {
	return ObjType
}

// Inspect returns formatted source code of this object.
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

// Proto returns proto of this object.
func (o *PanObj) Proto() PanObject {
	return o.proto
}

func keyHashes(pairs *map[SymHash]Pair) ([]SymHash, []SymHash) {
	publicKeyStrs := []string{}
	privateKeyStrs := []string{}

	for _, pair := range *pairs {
		str, ok := pair.Key.(*PanStr)
		// must be ok (obj keys are str)
		if !ok {
			continue
		}
		if str.IsPublic {
			publicKeyStrs = append(publicKeyStrs, str.Value)
		} else {
			privateKeyStrs = append(privateKeyStrs, str.Value)
		}
	}

	sort.Strings(publicKeyStrs)
	sort.Strings(privateKeyStrs)

	publicHashes := []SymHash{}
	for _, str := range publicKeyStrs {
		publicHashes = append(publicHashes, GetSymHash(str))
	}

	privateHashes := []SymHash{}
	for _, str := range privateKeyStrs {
		privateHashes = append(privateHashes, GetSymHash(str))
	}

	return publicHashes, privateHashes
}

// AddPairs adds pairs to obj.
// NOTE: Use this method only for prop DI. Otherwise immutability gets broken.
func (o *PanObj) AddPairs(pairs *map[SymHash]Pair) error {
	if pairs == nil {
		return errors.New("pairs must not be nil")
	}

	// add new pairs
	for k, v := range *pairs {
		// set only if prop does not exist
		if _, ok := (*o.Pairs)[k]; !ok {
			(*o.Pairs)[k] = v
		}
	}

	// update keys
	publicKeys, privateKeys := keyHashes(o.Pairs)
	o.Keys = &publicKeys
	o.PrivateKeys = &privateKeys
	return nil
}
