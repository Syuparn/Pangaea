package object

type PanObjType string

type PanObject interface {
	Type() PanObjType
	Inspect() string
	Proto() PanObject
}
