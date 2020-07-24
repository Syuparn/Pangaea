package evaluator

import (
	"../object"
	"../props"
)

// inject and set up builtin object props,
// which are defined in "props" package
func InjectBuiltInProps() {
	injectProps(object.BuiltInArrObj, props.ArrProps)
	injectProps(object.BuiltInFuncObj, props.FuncProps)
	injectProps(object.BuiltInIntObj, props.IntProps)
	injectProps(object.BuiltInMapObj, props.MapProps)
	injectProps(object.BuiltInObjObj, props.ObjProps)
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
	"Func_call":    &object.PanBuiltIn{Fn: evalFuncCall},
	"Map_at":       &object.PanBuiltIn{Fn: findElemInMap},
	"Obj_at":       &object.PanBuiltIn{Fn: findElemInObj},
	"Obj_callProp": &object.PanBuiltIn{Fn: builtInCallProp},
}
