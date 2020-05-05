package object

import ()

const NIL_TYPE = "NIL_TYPE"

type PanNil struct{}

func (n *PanNil) Type() PanObjType {
	return NIL_TYPE
}

func (n *PanNil) Inspect() string {
	return ""
}

func (n *PanNil) Proto() PanObject {
	return n
}
