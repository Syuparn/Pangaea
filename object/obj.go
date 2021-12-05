package object

import (
	"bytes"
	"errors"
	"sort"
)

// ObjType is a type of PanObj.
const ObjType = "ObjType"

// NewPanObj makes new PanObj instance.
func NewPanObj(pairs *map[SymHash]Pair, proto PanObject, options ...panObjOption) *PanObj {
	publicKeys, privateKeys := keyHashes(pairs)
	obj := &PanObj{
		Pairs:       pairs,
		Keys:        &publicKeys,
		PrivateKeys: &privateKeys,
		proto:       proto,
	}

	// inherit proto's zero value as default
	if obj.proto != nil {
		obj.zero = obj.proto.Zero()
	}

	for _, op := range options {
		op(obj)
	}

	return obj
}

// FIXME: use PanObjInstancePtr instead
// PanObjInstance makes new obj literal.
func PanObjInstance(pairs *map[SymHash]Pair) PanObj {
	publicKeys, privateKeys := keyHashes(pairs)
	return PanObj{
		Pairs:       pairs,
		Keys:        &publicKeys,
		PrivateKeys: &privateKeys,
		proto:       BuiltInObjObj,
		// use Obj's zero value {}
		zero: zeroObj,
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

// FIXME: call NewPanObj inside
// ChildPanObjPtr makes new child object of proto with props in src.
func ChildPanObjPtr(proto PanObject, src *PanObj, options ...panObjOption) *PanObj {
	// share pairs with src because objects are immutable
	obj := &PanObj{
		Pairs:       src.Pairs,
		Keys:        src.Keys,
		PrivateKeys: src.PrivateKeys,
		proto:       proto,
	}

	// inherit proto's zero value as default
	if obj.proto != nil {
		obj.zero = obj.proto.Zero()
	}

	for _, op := range options {
		op(obj)
	}

	return obj
}

type panObjOption func(*PanObj)

// WithZero can set zero value to the new *PanObj.
func WithZero(zero PanObject) panObjOption {
	return func(o *PanObj) { o.zero = zero }
}

// WithZeroFromSelf can set zero value created from new object itself.
func WithZeroFromSelf(f func(*PanObj) PanObject) panObjOption {
	return func(o *PanObj) { o.zero = f(o) }
}

// PanObj is object for not only obj literal but also
// any PanObject except specific data structure literal.
type PanObj struct {
	Pairs       *map[SymHash]Pair
	Keys        *[]SymHash
	PrivateKeys *[]SymHash
	proto       PanObject
	zero        PanObject
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

// Repr returns pritty-printed string of this object.
func (o *PanObj) Repr() string {
	// if self has '_name and it is str, use it
	// NOTE: only refer _name in self (NOT PROTO)!,
	// otherwise all children are shown as _name
	if name, ok := o.extractName(); ok {
		return name.Value
	}

	// child repr
	var out bytes.Buffer
	pairs := []Pair{}

	for _, p := range *o.Pairs {
		pairs = append(pairs, p)
	}

	out.WriteString("{")
	// NOTE: sort by key order otherwise output changes randomly
	// depending on inner map structure
	out.WriteString(sortedPairsRepr(pairs))
	out.WriteString("}")

	return out.String()
}

func (o *PanObj) extractName() (*PanStr, bool) {
	pair, ok := (*o.Pairs)[GetSymHash("_name")]
	if !ok {
		return nil, false
	}

	name, ok := TraceProtoOfStr(pair.Value)
	if !ok {
		return nil, false
	}

	return name, true
}

// Proto returns proto of this object.
func (o *PanObj) Proto() PanObject {
	return o.proto
}

// Zero returns zero value of this object.
func (o *PanObj) Zero() PanObject {
	// if unset, use zero value of obj {}.
	// NOTE: since initializing zeroObj.zero requires zeroObj itself, nil is treated as zeroObj
	if o.zero == nil {
		return zeroObj
	}
	return o.zero
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
