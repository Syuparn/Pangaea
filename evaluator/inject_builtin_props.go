package evaluator

import (
	"github.com/Syuparn/pangaea/object"
)

// NewPropContainer returns container of built-in props
// which are tightly coupled with evalXXX() and injected into built-in object.
func NewPropContainer() map[string]object.PanObject {
	return map[string]object.PanObject{
		// name format: "ObjName_propName"
		"Arr_at":       &object.PanBuiltIn{Fn: findElemInArr},
		"BaseObj_at":   &object.PanBuiltIn{Fn: findElemInObj},
		"Func_call":    &object.PanBuiltIn{Fn: evalFuncCall},
		"Int_at":       &object.PanBuiltIn{Fn: findBitInInt},
		"Iter_new":     &object.PanBuiltIn{Fn: iterNew},
		"Iter_next":    &object.PanBuiltIn{Fn: iterNext},
		"Map_at":       &object.PanBuiltIn{Fn: findElemInMap},
		"Obj_callProp": &object.PanBuiltIn{Fn: builtInCallProp},
		"Str_at":       &object.PanBuiltIn{Fn: findElemInStr},
	}
}
