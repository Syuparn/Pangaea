package object

import (
	"fmt"
)

const INT_TYPE = "INT_TYPE"

type PanInt struct {
	Value int64
}

func (i *PanInt) Type() PanObjType {
	return INT_TYPE
}

func (i *PanInt) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *PanInt) Proto() PanObject {
	return BuiltInIntObj
}

func (i *PanInt) Hash() HashKey {
	return HashKey{INT_TYPE, uint64(i.Value)}
}
