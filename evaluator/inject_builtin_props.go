package evaluator

import (
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/props"
)

// InjectBuiltInProps injects and sets up builtin object props,
// which are defined in "props" package.
func InjectBuiltInProps(ctn map[string]object.PanObject) {
	injectProps(object.BuiltInArrObj, props.ArrProps, ctn)
	injectProps(object.BuiltInBaseObj, props.BaseObjProps, ctn)
	injectProps(object.BuiltInFloatObj, props.FloatProps, ctn)
	injectProps(object.BuiltInFuncObj, props.FuncProps, ctn)
	injectProps(object.BuiltInIntObj, props.IntProps, ctn)
	injectProps(object.BuiltInIterObj, props.IterProps, ctn)
	injectProps(object.BuiltInMapObj, props.MapProps, ctn)
	injectProps(object.BuiltInNilObj, props.NilProps, ctn)
	injectProps(object.BuiltInObjObj, props.ObjProps, ctn)
	injectProps(object.BuiltInRangeObj, props.RangeProps, ctn)
	injectProps(object.BuiltInStrObj, props.StrProps, ctn)
}

func injectProps(
	obj *object.PanObj,
	props func(map[string]object.PanObject) map[string]object.PanObject,
	propContainer map[string]object.PanObject,
) {
	for propName, propVal := range props(propContainer) {
		propHash := object.GetSymHash(propName)
		(*obj.Pairs)[propHash] = object.Pair{
			Key:   object.NewPanStr(propName),
			Value: propVal,
		}
	}
}

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
