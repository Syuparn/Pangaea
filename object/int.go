package object

import (
	"fmt"
)

// IntType is a type of PanInt.
const IntType = "IntType"

// PanInt is object of int literal.
type PanInt struct {
	Value int64
}

// Type returns type of this PanObject.
func (i *PanInt) Type() PanObjType {
	return IntType
}

// Inspect returns formatted source code of this object.
func (i *PanInt) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Repr returns pritty-printed string of this object.
func (i *PanInt) Repr() string {
	return i.Inspect()
}

// Proto returns proto of this object.
func (i *PanInt) Proto() PanObject {
	return BuiltInIntObj
}

// Zero returns zero value of this object.
func (i *PanInt) Zero() PanObject {
	return i
}

// Hash returns hashkey of this object.
func (i *PanInt) Hash() HashKey {
	return HashKey{IntType, uint64(i.Value)}
}

// NewPanInt returns new int object.
// NOTE: `0` and `1` are cached and always same instance are returned.
func NewPanInt(i int64) *PanInt {
	switch i {
	case 0:
		return BuiltInZeroInt
	case 1:
		return BuiltInOneInt
	default:
		return &PanInt{Value: i}
	}
}
