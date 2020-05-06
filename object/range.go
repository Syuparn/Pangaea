package object

import (
	"bytes"
	"strings"
)

const RANGE_TYPE = "RANGE_TYPE"

type PanRange struct {
	Start PanObject
	Stop  PanObject
	Step  PanObject
}

func (r *PanRange) Type() PanObjType {
	return RANGE_TYPE
}

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

func (r *PanRange) Proto() PanObject {
	return BuiltInRangeObj
}
