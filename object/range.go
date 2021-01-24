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

// Proto returns proto of this object.
func (r *PanRange) Proto() PanObject {
	return BuiltInRangeObj
}
