package object

const STR_TYPE = "STR_TYPE"

type PanStr struct {
	Value string
}

func (s *PanStr) Type() PanObjType {
	return STR_TYPE
}

func (s *PanStr) Inspect() string {
	return s.Value
}

func (s *PanStr) Proto() PanObject {
	return s
}
