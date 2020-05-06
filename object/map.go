package object

const MAP_TYPE = "MAP_TYPE"

type PanMap struct {
	Pairs *map[HashKey]Pair
}

func (m *PanMap) Type() PanObjType {
	return ""
}

func (m *PanMap) Inspect() string {
	return ""
}

func (m *PanMap) Proto() PanObject {
	return m
}
