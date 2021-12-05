package object

import (
	"fmt"
	"math"
)

// FloatType is a type of PanFloat.
const FloatType = "FloatType"

// PanFloat is object of float literal.
type PanFloat struct {
	Value float64
}

// Type returns type of this PanObject.
func (f *PanFloat) Type() PanObjType {
	return FloatType
}

// Inspect returns formatted source code of this object.
func (f *PanFloat) Inspect() string {
	return fmt.Sprintf("%.6f", f.Value)
}

// Repr returns pritty-printed string of this object.
func (f *PanFloat) Repr() string {
	return f.Inspect()
}

// Proto returns proto of this object.
func (f *PanFloat) Proto() PanObject {
	return BuiltInFloatObj
}

// Zero returns zero value of this object.
func (f *PanFloat) Zero() PanObject {
	return f
}

// Hash returns hashkey of this object.
func (f *PanFloat) Hash() HashKey {
	// Float64bits convert float64 to uint64 with same bit pattern
	return HashKey{FloatType, math.Float64bits(f.Value)}
}

// NewPanFloat returns new float object.
func NewPanFloat(f float64) *PanFloat {
	return &PanFloat{f}
}
