package object

import (
	"bytes"
	"strings"
)

// ArrType is a type of PanArr.
const ArrType = "ArrType"

// PanArr is object of arr literal.
type PanArr struct {
	Elems []PanObject
}

// Type returns type of this PanObject.
func (a *PanArr) Type() PanObjType {
	return ArrType
}

// Inspect returns formatted source code of this object.
func (a *PanArr) Inspect() string {
	var out bytes.Buffer
	elems := []string{}
	for _, e := range a.Elems {
		elems = append(elems, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}

// Proto returns proto of this object.
func (a *PanArr) Proto() PanObject {
	return BuiltInArrObj
}
