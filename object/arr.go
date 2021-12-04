package object

import (
	"bytes"
	"strings"
)

// ArrType is a type of PanArr.
const ArrType = "ArrType"

// used as zero value
var zeroArr = NewPanArr()

// PanArr is object of arr literal.
type PanArr struct {
	Elems []PanObject
	proto PanObject
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

// Repr returns pritty-printed string of this object.
func (a *PanArr) Repr() string {
	var out bytes.Buffer
	elems := []string{}
	for _, e := range a.Elems {
		elems = append(elems, e.Repr())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")
	return out.String()
}

// Proto returns proto of this object.
func (a *PanArr) Proto() PanObject {
	return a.proto
}

// NewPanArr returns new arr object.
func NewPanArr(elems ...PanObject) *PanArr {
	return &PanArr{
		Elems: elems,
		proto: BuiltInArrObj,
	}
}

// NewInheritedArr returns new arr object born of proto.
func NewInheritedArr(proto PanObject, elems ...PanObject) *PanArr {
	return &PanArr{
		Elems: elems,
		proto: proto,
	}
}
