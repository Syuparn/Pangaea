package object

import (
	"fmt"
)

// BoolType is a type of PanBool.
const BoolType = "BoolType"

// PanBool is object of bool literal.
type PanBool struct {
	Value bool
}

// Type returns type of this PanObject.
func (b *PanBool) Type() PanObjType {
	return BoolType
}

// Inspect returns formatted source code of this object.
func (b *PanBool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Proto returns proto of this object.
func (b *PanBool) Proto() PanObject {
	if b.Value {
		return BuiltInOneInt
	}
	return BuiltInZeroInt
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
