package props

import (
	"../object"
)

var eqSym = object.NewPanStr("==")

func traceProtoOf(
	obj object.PanObject,
	cond func(object.PanObject) bool,
) (object.PanObject, bool) {
	for o := obj; o.Proto() != nil; o = o.Proto() {
		if cond(o) {
			return o, true
		}
	}
	return nil, false
}

func isArr(o object.PanObject) bool {
	return o.Type() == object.ARR_TYPE
}

func isBuiltInFunc(o object.PanObject) bool {
	return o.Type() == object.BUILTIN_TYPE
}

func isFloat(o object.PanObject) bool {
	return o.Type() == object.FLOAT_TYPE
}

func isFunc(o object.PanObject) bool {
	return o.Type() == object.FUNC_TYPE
}

func isInt(o object.PanObject) bool {
	return o.Type() == object.INT_TYPE
}

func isMap(o object.PanObject) bool {
	return o.Type() == object.MAP_TYPE
}

func isNil(o object.PanObject) bool {
	return o.Type() == object.NIL_TYPE
}

func isObj(o object.PanObject) bool {
	return o.Type() == object.OBJ_TYPE
}

func isRange(o object.PanObject) bool {
	return o.Type() == object.RANGE_TYPE
}

func isStr(o object.PanObject) bool {
	return o.Type() == object.STR_TYPE
}

func propIn(obj *object.PanObj, propName string) (object.Pair, bool) {
	propSym := object.GetSymHash(propName)
	pair, ok := (*obj.Pairs)[propSym]
	return pair, ok
}
