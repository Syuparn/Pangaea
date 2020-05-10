package object

type BuiltInFunc func(env *Env, args ...PanObject) PanObject

const BUILTIN_TYPE = "BUILTIN_TYPE"

type PanBuiltIn struct {
	Fn BuiltInFunc
}

func (b *PanBuiltIn) Type() PanObjType {
	return BUILTIN_TYPE
}

func (b *PanBuiltIn) Inspect() string {
	return "{|| [builtin]}"
}

func (b *PanBuiltIn) Proto() PanObject {
	return BuiltInFuncObj
}
