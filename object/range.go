package object

import (
	"bytes"
	"strings"
)

// RangeType is a type of PanRange.
const RangeType = "RangeType"

// PanRange is object of range literal.
type PanRange struct {
	Start PanObject
	Stop  PanObject
	Step  PanObject
}

// Type returns type of this PanObject.
func (r *PanRange) Type() PanObjType {
	return RangeType
}

// Inspect returns formatted source code of this object.
func (r *PanRange) Inspect() string {
	var out bytes.Buffer
	elems := make([]string, 3)
	elems[0] = r.Start.Inspect()
	elems[1] = r.Stop.Inspect()
	elems[2] = r.Step.Inspect()

	out.WriteString("(")
	out.WriteString(strings.Join(elems, ":"))
	out.WriteString(")")
	return out.String()
}

// Repr returns pritty-printed string of this object.
func (r *PanRange) Repr() string {
	var out bytes.Buffer
	elems := make([]string, 3)
	elems[0] = r.Start.Repr()
	elems[1] = r.Stop.Repr()
	elems[2] = r.Step.Repr()

	out.WriteString("(")
	out.WriteString(strings.Join(elems, ":"))
	out.WriteString(")")
	return out.String()
}

// Proto returns proto of this object.
func (r *PanRange) Proto() PanObject {
	return BuiltInRangeObj
}

// Zero returns zero value of this object.
func (r *PanRange) Zero() PanObject {
	return r
}

// NewPanRange returns new range object.
func NewPanRange(start, stop, step PanObject) *PanRange {
	return &PanRange{start, stop, step}
}
