package object

// BuiltInIterType is a type of PanBuiltInIter.
const BuiltInIterType = "BuiltInIterType"

// PanBuiltInIter is object of built-in iter literal.
// NOTE: it has env to save state
type PanBuiltInIter struct {
	Fn  BuiltInFunc
	Env *Env
}

// Type returns type of this PanObject.
func (b *PanBuiltInIter) Type() PanObjType {
	return BuiltInIterType
}

// Inspect returns formatted source code of this object.
func (b *PanBuiltInIter) Inspect() string {
	return "<{|| [builtin]}>"
}

// Proto returns proto of this object.
func (b *PanBuiltInIter) Proto() PanObject {
	return BuiltInIterObj
}
