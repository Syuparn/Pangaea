package evaluator

import (
	"../object"
	"../props"
)

// inject and set up builtin object props,
// which are defined in "props" package
func InjectBuiltInProps() {
	injectProps(object.BuiltInArrObj, props.ArrProps)
	injectProps(object.BuiltInBaseObj, props.BaseObjProps)
	injectProps(object.BuiltInFloatObj, props.FloatProps)
	injectProps(object.BuiltInFuncObj, props.FuncProps)
	injectProps(object.BuiltInIntObj, props.IntProps)
	injectProps(object.BuiltInIterObj, props.IterProps)
	injectProps(object.BuiltInMapObj, props.MapProps)
	injectProps(object.BuiltInNilObj, props.NilProps)
	injectProps(object.BuiltInObjObj, props.ObjProps)
	injectProps(object.BuiltInRangeObj, props.RangeProps)
	injectProps(object.BuiltInStrObj, props.StrProps)
}

func injectProps(
	obj *object.PanObj,
	props func(map[string]object.PanObject) map[string]object.PanObject,
) {
	for propName, propVal := range props(propContainer) {
		propHash := object.GetSymHash(propName)
		(*obj.Pairs)[propHash] = object.Pair{
			Key:   &object.PanStr{Value: propName},
			Value: propVal,
		}
	}
}

var propContainer = map[string]object.PanObject{
	// name format: "ObjName_propName"
	"Arr_at":       &object.PanBuiltIn{Fn: findElemInArr},
	"BaseObj_at":   &object.PanBuiltIn{Fn: findElemInObj},
	"Func_call":    &object.PanBuiltIn{Fn: evalFuncCall},
	"Iter_new":     &object.PanBuiltIn{Fn: iterNew},
	"Iter_next":    &object.PanBuiltIn{Fn: iterNext},
	"Map_at":       &object.PanBuiltIn{Fn: findElemInMap},
	"Obj_callProp": &object.PanBuiltIn{Fn: builtInCallProp},
}
