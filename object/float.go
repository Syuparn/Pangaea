package object

import (
	"fmt"
	"math"
)

const FLOAT_TYPE = "FLOAT_TYPE"

type PanFloat struct {
	Value float64
}

func (f *PanFloat) Type() PanObjType {
	return FLOAT_TYPE
}

func (f *PanFloat) Inspect() string {
	return fmt.Sprintf("%.6f", f.Value)
}

func (f *PanFloat) Proto() PanObject {
	return BuiltInFloatObj
}

func (f *PanFloat) Hash() HashKey {
	// Float64bits convert float64 to uint64 with same bit pattern
	return HashKey{FLOAT_TYPE, math.Float64bits(f.Value)}
}
