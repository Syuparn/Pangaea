package evaluator

import (
	"../object"
)

func findElemInObj(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 2 {
		return object.NewTypeErr("Obj#at requires at least 2 args")
	}
	self := args[0]

	// TODO: duck typing for keys (allow child of arr)
	indexArr, ok := args[1].(*object.PanArr)
	if !ok {
		return object.BuiltInNil
	}

	if len(indexArr.Elems) < 1 {
		return object.BuiltInNil
	}

	// TODO: duck typing for keys (allow child of str)
	propName, ok := indexArr.Elems[0].(*object.PanStr)
	if !ok {
		return object.BuiltInNil
	}

	ret, ok := callProp(self, object.GetSymHash(propName.Value))
	if !ok {
		return object.BuiltInNil
	}

	return ret
}

func findElemInArr(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	// NOTE: if index is not found, Obj#at is called

	if len(args) < 2 {
		return object.NewTypeErr("Arr#at requires at least 2 args")
	}
	// TODO: duck typing for keys (allow child of arr)
	self, ok := args[0].(*object.PanArr)
	if !ok {
		return object.BuiltInNil
	}

	// TODO: duck typing for keys (allow child of arr)
	indexArr, ok := args[1].(*object.PanArr)
	if !ok {
		return object.BuiltInNil
	}

	if len(indexArr.Elems) < 1 {
		return object.BuiltInNil
	}

	switch index := indexArr.Elems[0].(type) {
	case *object.PanInt:
		// TODO: duck typing for keys (allow child of int)
		return arrIndex(index.Value, self)
	case *object.PanRange:
		return arrRange(index, self)
	default:
		return findElemInObj(env, kwargs, args...)
	}
}

func arrIndex(index int64, arr *object.PanArr) object.PanObject {
	length := int64(len(arr.Elems))
	if index >= length || index < -length {
		return object.BuiltInNil
	}

	if index < 0 {
		return arr.Elems[index+length]
	}

	return arr.Elems[index]
}

func arrRange(r *object.PanRange, arr *object.PanArr) object.PanObject {
	ok := canBeUsedForArrRange(r.Start) &&
		canBeUsedForArrRange(r.Stop) &&
		canBeUsedForArrRange(r.Step)
	if !ok {
		// empty array
		return &object.PanArr{Elems: []object.PanObject{}}
	}

	// default step
	step := int64(1)
	if i, ok := r.Step.(*object.PanInt); ok {
		step = i.Value
	}

	if step == 0 {
		return object.NewValueErr("cannot use 0 for range step")
	}

	start, stop := fixRange(r, int64(len(arr.Elems)), step)

	hasNext := func(i int64, stop int64) bool {
		if step < 0 {
			return i > stop
		}
		return i < stop
	}

	elems := []object.PanObject{}
	for i := start; hasNext(i, stop); i += step {
		elems = append(elems, arrIndex(i, arr))
	}

	return &object.PanArr{Elems: elems}
}

func canBeUsedForArrRange(o object.PanObject) bool {
	return o.Type() == object.INT_TYPE || o.Type() == object.NIL_TYPE
}

func fixRange(r *object.PanRange, length int64, step int64) (int64, int64) {
	fix := func(i int64) int64 {
		if i < -length {
			return 0
		}
		if i > length {
			return length
		}
		if i < 0 {
			return i + length
		}
		return i
	}

	var start, stop int64

	// default values
	if step > 0 {
		start = 0
		stop = length
	} else {
		start = length - 1
		stop = -1
	}

	// update by range value
	if i, ok := r.Start.(*object.PanInt); ok {
		start = fix(i.Value)
	}

	if i, ok := r.Stop.(*object.PanInt); ok {
		stop = fix(i.Value)
	}

	return start, stop
}

func findElemInMap(
	env *object.Env,
	kwargs *object.PanObj,
	args ...object.PanObject,
) object.PanObject {
	// NOTE: if index is not found, Obj#at is called

	if len(args) < 2 {
		return object.NewTypeErr("Obj#at requires at least 2 args")
	}

	// TODO: duck typing for keys (allow child of map)
	self, ok := args[0].(*object.PanMap)
	if !ok {
		return object.BuiltInNil
	}

	// TODO: duck typing for keys (allow child of arr)
	indexArr, ok := args[1].(*object.PanArr)
	if !ok {
		return object.BuiltInNil
	}

	if len(indexArr.Elems) < 1 {
		return object.BuiltInNil
	}

	index := indexArr.Elems[0]
	hashableIndex, ok := index.(object.PanScalar)

	if ok {
		pair, ok := (*self.Pairs)[hashableIndex.Hash()]
		if !ok {
			return findElemInObj(env, kwargs, args...)
		}
		return pair.Value
	}

	for _, pair := range *self.NonHashablePairs {
		// find key by == method
		eqSym := &object.PanStr{Value: "=="}

		// equivalent to src `key == index`
		ret := builtInCallProp(env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), index, eqSym, pair.Key)

		if ret.Type() == object.ERR_TYPE {
			return ret
		}

		if ret == object.BuiltInTrue {
			return pair.Value
		}
	}

	return object.BuiltInNil
}
