package object

import (
	"fmt"
)

// BoolType is a type of PanBool.
const BoolType = "BoolType"

// PanBool is object of bool literal.
type PanBool struct {
	Value bool
	// NOTE: bool does not have proto field because true and false must be singletons.
}

// Type returns type of this PanObject.
func (b *PanBool) Type() PanObjType {
	return BoolType
}

// Inspect returns formatted source code of this object.
func (b *PanBool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Repr returns pritty-printed string of this object.
func (b *PanBool) Repr() string {
	return b.Inspect()
}

// Proto returns proto of this object.
func (b *PanBool) Proto() PanObject {
	if b.Value {
		return BuiltInOneInt
	}
	return BuiltInZeroInt
}

// Zero returns zero value of this object.
func (b *PanBool) Zero() PanObject {
	return b
}

// Hash returns hashkey of this object.
func (b *PanBool) Hash() HashKey {
	var v uint64
	if b.Value {
		v = 1
	} else {
		v = 0
	}

	return HashKey{BoolType, v}
}
