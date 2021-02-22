package evaluator

import (
	"github.com/Syuparn/pangaea/object"
)

// NewPropContainer returns container of built-in props
// which are tightly coupled with evalXXX() and injected into built-in object.
func NewPropContainer() map[string]object.PanObject {
	return map[string]object.PanObject{
		// name format: "ObjName_propName"
		"Arr_at":       object.NewPanBuiltInFunc(findElemInArr),
		"BaseObj_at":   object.NewPanBuiltInFunc(findElemInObj),
		"Func_call":    object.NewPanBuiltInFunc(evalFuncCall),
		"Int_at":       object.NewPanBuiltInFunc(findBitInInt),
		"Iter_new":     object.NewPanBuiltInFunc(iterNew),
		"Iter_next":    object.NewPanBuiltInFunc(iterNext),
		"Map_at":       object.NewPanBuiltInFunc(findElemInMap),
		"Obj_callProp": object.NewPanBuiltInFunc(builtInCallProp),
		"Str_at":       object.NewPanBuiltInFunc(findElemInStr),
	}
}
