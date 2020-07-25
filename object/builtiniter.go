package object

const BUILTIN_ITER_TYPE = "BUILTIN_ITER_TYPE"

// PanBuiltInIter has env to save state
type PanBuiltInIter struct {
	Fn  BuiltInFunc
	Env *Env
}

func (b *PanBuiltInIter) Type() PanObjType {
	return BUILTIN_ITER_TYPE
}

func (b *PanBuiltInIter) Inspect() string {
	return "<{|| [builtin]}>"
}

func (b *PanBuiltInIter) Proto() PanObject {
	return BuiltInIterObj
}
