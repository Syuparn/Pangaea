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
	// HACK: read source codes concurrently to run faster
	// (evaluating native codes are bottleneck of injection)
	arrNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		arrNativesCh <- mustReadNativeCode("Arr", env)
	}()
	baseObjNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		baseObjNativesCh <- mustReadNativeCode("BaseObj", env)
	}()
	comparableNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		comparableNativesCh <- mustReadNativeCode("Comparable", env)
	}()
	intNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		intNativesCh <- mustReadNativeCode("Int", env)
	}()
	iterableNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		iterableNativesCh <- mustReadNativeCode("Iterable", env)
	}()
	objNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		objNativesCh <- mustReadNativeCode("Obj", env)
	}()
	strNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		strNativesCh <- mustReadNativeCode("Str", env)
	}()
	wrappableNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		wrappableNativesCh <- mustReadNativeCode("Wrappable", env)
	}()

	arrNatives := <-arrNativesCh
	baseObjNatives := <-baseObjNativesCh
	comparableNatives := <-comparableNativesCh
	intNatives := <-intNativesCh
	iterableNatives := <-iterableNativesCh
	objNatives := <-objNativesCh
	strNatives := <-strNativesCh
	wrappableNatives := <-wrappableNativesCh

	// NOTE: injection order is important! (if same-named props appear, first one remains)
	injectProps(object.BuiltInArrObj, toPairs(props.ArrProps(ctn)), arrNatives, iterableNatives)
	injectProps(object.BuiltInAssertionErr, toPairs(props.AssertionErrProps(ctn)))
	injectProps(object.BuiltInBaseObj, toPairs(props.BaseObjProps(ctn)), baseObjNatives)
	injectProps(object.BuiltInComparableObj, comparableNatives)
	injectProps(object.BuiltInDiamondObj, toPairs(props.DiamondProps(ctn)), iterableNatives)
	injectProps(object.BuiltInEitherObj, wrappableNatives)
	injectProps(object.BuiltInErrObj, toPairs(props.ErrProps(ctn)))
	injectProps(object.BuiltInFloatObj, toPairs(props.FloatProps(ctn)), comparableNatives)
	injectProps(object.BuiltInFuncObj, toPairs(props.FuncProps(ctn)))
	injectProps(object.BuiltInIntObj, toPairs(props.IntProps(ctn)), intNatives, iterableNatives, comparableNatives)
	injectProps(object.BuiltInIterObj, toPairs(props.IterProps(ctn)))
	injectProps(object.BuiltInIterableObj, iterableNatives)
	injectProps(object.BuiltInKernelObj, toPairs(props.KernelProps(ctn)))
	injectProps(object.BuiltInMapObj, toPairs(props.MapProps(ctn)), iterableNatives)
	injectProps(object.BuiltInNameErr, toPairs(props.NameErrProps(ctn)))
	injectProps(object.BuiltInNilObj, toPairs(props.NilProps(ctn)))
	injectProps(object.BuiltInNoPropErr, toPairs(props.NoPropErrProps(ctn)))
	injectProps(object.BuiltInNotImplementedErr, toPairs(props.NotImplementedErrProps(ctn)))
	injectProps(object.BuiltInNumObj, toPairs(props.NumProps(ctn)))
	injectProps(object.BuiltInObjObj, toPairs(props.ObjProps(ctn)), objNatives, iterableNatives)
	injectProps(object.BuiltInRangeObj, toPairs(props.RangeProps(ctn)), iterableNatives)
	injectProps(object.BuiltInStopIterErr, toPairs(props.StopIterErrProps(ctn)))
	injectProps(object.BuiltInStrObj, toPairs(props.StrProps(ctn)), strNatives, iterableNatives, comparableNatives)
	injectProps(object.BuiltInSyntaxErr, toPairs(props.SyntaxErrProps(ctn)))
	injectProps(object.BuiltInTypeErr, toPairs(props.TypeErrProps(ctn)))
	injectProps(object.BuiltInEitherValObj, toPairs(props.EitherValProps(ctn)))
	injectProps(object.BuiltInEitherErrObj, toPairs(props.EitherErrProps(ctn)))
	injectProps(object.BuiltInValueErr, toPairs(props.ValueErrProps(ctn)))
	injectProps(object.BuiltInWrappableObj, wrappableNatives)
	injectProps(object.BuiltInZeroDivisionErr, toPairs(props.ZeroDivisionErrProps(ctn)))
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
