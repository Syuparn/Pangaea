package object

import (
	"bytes"
	"strings"
)

const ARR_TYPE = "ARR_TYPE"

type PanArr struct {
	Elems []PanObject
}

func (a *PanArr) Type() PanObjType {
	return ARR_TYPE
}

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

func (a *PanArr) Proto() PanObject {
	return BuiltInArrObj
}
