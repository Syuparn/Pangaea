package object

// BuiltInFunc is a function for PanBuiltIn.
type BuiltInFunc func(env *Env, kwargs *PanObj, args ...PanObject) PanObject

// BuiltInType is a type of PanBuiltIn.
const BuiltInType = "BuiltInType"

// PanBuiltIn is object of built-in func literal.
type PanBuiltIn struct {
	Fn BuiltInFunc
}

// Type returns type of this PanObject.
func (b *PanBuiltIn) Type() PanObjType {
	return BuiltInType
}

// Inspect returns formatted source code of this object.
func (b *PanBuiltIn) Inspect() string {
	return "{|| [builtin]}"
}

// Proto returns proto of this object.
func (b *PanBuiltIn) Proto() PanObject {
	return BuiltInFuncObj
}

// NewPanBuiltInFunc returns new BuiltInFunc object.
func NewPanBuiltInFunc(f BuiltInFunc) *PanBuiltIn {
	return &PanBuiltIn{f}
}
