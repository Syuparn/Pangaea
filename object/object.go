package object

type PanObjType int

type PanObject interface {
	Type() PanObjType
	Inspect() string
	Proto() PanObject
}
