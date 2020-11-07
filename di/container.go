package di

import (
	"github.com/Syuparn/pangaea/evaluator"
	"github.com/Syuparn/pangaea/object"
	"github.com/Syuparn/pangaea/props"
)

// InjectBuiltInProps injects and sets up builtin object props,
// which are defined in "props" package.
func InjectBuiltInProps(env *object.Env) {
	propContainer := mergePropContainers(
		NewPropContainer(),
		evaluator.NewPropContainer(),
	)
	injectBuiltInProps(env, propContainer)
}

// NewPropContainer returns container of built-in props
// which are injected into built-in object.
func NewPropContainer() map[string]object.PanObject {
	return map[string]object.PanObject{
		// name format: "ObjName_propName"
		"Str_eval":    &object.PanBuiltIn{Fn: strEval},
		"Str_evalEnv": &object.PanBuiltIn{Fn: strEvalEnv},
	}
}

func injectBuiltInProps(
	env *object.Env,
	ctn map[string]object.PanObject,
) {
	injectProps(object.BuiltInArrObj, toPairs(props.ArrProps(ctn)), mustReadNativeCode("Arr", env), mustReadNativeCode("Iterable", env))
	injectProps(object.BuiltInBaseObj, toPairs(props.BaseObjProps(ctn)))
	injectProps(object.BuiltInFloatObj, toPairs(props.FloatProps(ctn)))
	injectProps(object.BuiltInFuncObj, toPairs(props.FuncProps(ctn)))
	injectProps(object.BuiltInIntObj, toPairs(props.IntProps(ctn)), mustReadNativeCode("Int", env), mustReadNativeCode("Iterable", env))
	injectProps(object.BuiltInIterObj, toPairs(props.IterProps(ctn)))
	injectProps(object.BuiltInIterableObj, mustReadNativeCode("Iterable", env))
	injectProps(object.BuiltInKernelObj, toPairs(props.KernelProps(ctn)))
	injectProps(object.BuiltInMapObj, toPairs(props.MapProps(ctn)), mustReadNativeCode("Iterable", env))
	injectProps(object.BuiltInNilObj, toPairs(props.NilProps(ctn)))
	injectProps(object.BuiltInNumObj, toPairs(props.NumProps(ctn)))
	injectProps(object.BuiltInObjObj, toPairs(props.ObjProps(ctn)), mustReadNativeCode("Obj", env), mustReadNativeCode("Iterable", env))
	injectProps(object.BuiltInRangeObj, toPairs(props.RangeProps(ctn)), mustReadNativeCode("Iterable", env))
	injectProps(object.BuiltInStrObj, toPairs(props.StrProps(ctn)), mustReadNativeCode("Str", env), mustReadNativeCode("Iterable", env))
}

func injectProps(
	obj *object.PanObj,
	pairsList ...*map[object.SymHash]object.Pair,
) {
	for _, pairs := range pairsList {
		obj.AddPairs(pairs)
	}
}

func toPairs(
	propContainer map[string]object.PanObject,
) *map[object.SymHash]object.Pair {
	pairs := map[object.SymHash]object.Pair{}
	for k, v := range propContainer {
		pair := object.Pair{
			Key:   object.NewPanStr(k),
			Value: v,
		}
		pairs[object.GetSymHash(k)] = pair
	}

	return &pairs
}

func mergePropContainers(
	containers ...map[string]object.PanObject,
) map[string]object.PanObject {
	mergedCtn := map[string]object.PanObject{}
	for _, ctn := range containers {
		for k, v := range ctn {
			// NOTE: if same key is found, first value is remained
			if _, ok := mergedCtn[k]; !ok {
				mergedCtn[k] = v
			}
		}
	}

	return mergedCtn
}
