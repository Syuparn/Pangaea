package object

import (
	"bytes"
	"strings"
)

const OBJ_TYPE = "OBJ_TYPE"

type PanObj struct {
	Pairs *map[SymHash]Pair
}

func (o *PanObj) Type() PanObjType {
	return OBJ_TYPE
}

func (o *PanObj) Inspect() string {
	var out bytes.Buffer
	elems := []string{}

	// NOTE: refer map because range cannot treat map pointer
	for _, p := range *o.Pairs {
		// NOTE: unwrap double quotation
		keyStr := p.Key.Inspect()
		keyStr = keyStr[1 : len(keyStr)-1]

		elemStr := keyStr + ": " + p.Value.Inspect()
		elems = append(elems, elemStr)
	}

	out.WriteString("{")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("}")

	return out.String()
}

func (o *PanObj) Proto() PanObject {
	return builtInObjObj
}
