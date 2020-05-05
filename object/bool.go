package object

import (
	"fmt"
)

const BOOL_TYPE = "BOOL_TYPE"

type PanBool struct {
	Value bool
}

func (b *PanBool) Type() PanObjType {
	return BOOL_TYPE
}

func (b *PanBool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *PanBool) Proto() PanObject {
	return builtInBoolObj
}

func (b *PanBool) Hash() HashKey {
	return HashKey{}
}
