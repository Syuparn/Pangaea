package object

const NIL_TYPE = "NIL_TYPE"

type PanNil struct{}

func (n *PanNil) Type() PanObjType {
	return NIL_TYPE
}

func (n *PanNil) Inspect() string {
	return "nil"
}

func (n *PanNil) Proto() PanObject {
	return BuiltInNilObj
}

func (n *PanNil) Hash() HashKey {
	return HashKey{NIL_TYPE, 0}
}
