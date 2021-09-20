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
	diamondNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		diamondNativesCh <- mustReadNativeCode("Diamond", env)
	}()
	eitherNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		eitherNativesCh <- mustReadNativeCode("Either", env)
	}()
	eitherErrNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		eitherErrNativesCh <- mustReadNativeCode("EitherErr", env)
	}()
	eitherValNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		eitherValNativesCh <- mustReadNativeCode("EitherVal", env)
	}()
	floatNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		floatNativesCh <- mustReadNativeCode("Float", env)
	}()
	funcNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		funcNativesCh <- mustReadNativeCode("Func", env)
	}()
	intNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		intNativesCh <- mustReadNativeCode("Int", env)
	}()
	iterNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		iterNativesCh <- mustReadNativeCode("Iter", env)
	}()
	iterableNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		iterableNativesCh <- mustReadNativeCode("Iterable", env)
	}()
	kernelNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		kernelNativesCh <- mustReadNativeCode("Kernel", env)
	}()
	mapNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		mapNativesCh <- mustReadNativeCode("Map", env)
	}()
	objNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		objNativesCh <- mustReadNativeCode("Obj", env)
	}()
	rangeNativesCh := make(chan *map[object.SymHash]object.Pair)
	go func() {
		rangeNativesCh <- mustReadNativeCode("Range", env)
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
	diamondNatives := <-diamondNativesCh
	eitherNatives := <-eitherNativesCh
	eitherErrNatives := <-eitherErrNativesCh
	eitherValNatives := <-eitherValNativesCh
	floatNatives := <-floatNativesCh
	funcNatives := <-funcNativesCh
	intNatives := <-intNativesCh
	iterNatives := <-iterNativesCh
	iterableNatives := <-iterableNativesCh
	kernelNatives := <-kernelNativesCh
	mapNatives := <-mapNativesCh
	objNatives := <-objNativesCh
	rangeNatives := <-rangeNativesCh
	strNatives := <-strNativesCh
	wrappableNatives := <-wrappableNativesCh

	// NOTE: injection order is important! (if same-named props appear, first one remains)
	injectProps(object.BuiltInArrObj, toPairs(props.ArrProps(ctn)), arrNatives, iterableNatives)
	injectProps(object.BuiltInAssertionErr, toPairs(props.AssertionErrProps(ctn)))
	injectProps(object.BuiltInBaseObj, toPairs(props.BaseObjProps(ctn)), baseObjNatives)
	injectProps(object.BuiltInComparableObj, toPairs(props.ComparableProps(ctn)), comparableNatives)
	injectProps(object.BuiltInDiamondObj, toPairs(props.DiamondProps(ctn)), diamondNatives, iterableNatives)
	injectProps(object.BuiltInEitherObj, toPairs(props.EitherProps(ctn)), eitherNatives, wrappableNatives)
	injectProps(object.BuiltInEitherErrObj, toPairs(props.EitherErrProps(ctn)), eitherErrNatives)
	injectProps(object.BuiltInEitherValObj, toPairs(props.EitherValProps(ctn)), eitherValNatives)
	injectProps(object.BuiltInErrObj, toPairs(props.ErrProps(ctn)))
	injectProps(object.BuiltInFloatObj, toPairs(props.FloatProps(ctn)), floatNatives, comparableNatives)
	injectProps(object.BuiltInFuncObj, toPairs(props.FuncProps(ctn)), funcNatives)
	injectProps(object.BuiltInIntObj, toPairs(props.IntProps(ctn)), intNatives, iterableNatives, comparableNatives)
	injectProps(object.BuiltInIterObj, toPairs(props.IterProps(ctn)), iterNatives, iterableNatives)
	injectProps(object.BuiltInIterableObj, toPairs(props.IterableProps(ctn)), iterableNatives)
	injectProps(object.BuiltInJSONObj, toPairs(props.JSONProps(ctn)))
	injectProps(object.BuiltInKernelObj, toPairs(props.KernelProps(ctn)), kernelNatives)
	injectProps(object.BuiltInMatchObj, toPairs(props.MatchProps(ctn)))
	injectProps(object.BuiltInMapObj, toPairs(props.MapProps(ctn)), mapNatives, iterableNatives)
	injectProps(object.BuiltInNameErr, toPairs(props.NameErrProps(ctn)))
	injectProps(object.BuiltInNilObj, toPairs(props.NilProps(ctn)))
	injectProps(object.BuiltInNoPropErr, toPairs(props.NoPropErrProps(ctn)))
	injectProps(object.BuiltInNotImplementedErr, toPairs(props.NotImplementedErrProps(ctn)))
	injectProps(object.BuiltInNumObj, toPairs(props.NumProps(ctn)))
	injectProps(object.BuiltInObjObj, toPairs(props.ObjProps(ctn)), objNatives, iterableNatives)
	injectProps(object.BuiltInRangeObj, toPairs(props.RangeProps(ctn)), rangeNatives, iterableNatives)
	injectProps(object.BuiltInStopIterErr, toPairs(props.StopIterErrProps(ctn)))
	injectProps(object.BuiltInStrObj, toPairs(props.StrProps(ctn)), strNatives, iterableNatives, comparableNatives)
	injectProps(object.BuiltInSyntaxErr, toPairs(props.SyntaxErrProps(ctn)))
	injectProps(object.BuiltInTypeErr, toPairs(props.TypeErrProps(ctn)))
	injectProps(object.BuiltInValueErr, toPairs(props.ValueErrProps(ctn)))
	injectProps(object.BuiltInWrappableObj, toPairs(props.WrappableProps(ctn)), wrappableNatives)
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
