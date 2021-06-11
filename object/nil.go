package object

// NilType is a type of PanNil.
const NilType = "NilType"

// PanNil is object of nil literal.
type PanNil struct{}

// Type returns type of this PanObject.
func (n *PanNil) Type() PanObjType {
	return NilType
}

// Inspect returns formatted source code of this object.
func (n *PanNil) Inspect() string {
	return "nil"
}

// Repr returns pritty-printed string of this object.
func (n *PanNil) Repr() string {
	return n.Inspect()
}

// Proto returns proto of this object.
func (n *PanNil) Proto() PanObject {
	return BuiltInNilObj
}

// Hash returns hashkey of this object.
func (n *PanNil) Hash() HashKey {
	return HashKey{NilType, 0}
}

// NewPanNil returns new nil object.
func NewPanNil() *PanNil {
	return BuiltInNil
}
