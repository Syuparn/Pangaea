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
		return findElemInObj(env, kwargs, args...)
	}

	// TODO: duck typing for keys (allow child of arr)
	indexArr, ok := args[1].(*object.PanArr)
	if !ok {
		return findElemInObj(env, kwargs, args...)
	}

	if len(indexArr.Elems) < 1 {
		return findElemInObj(env, kwargs, args...)
	}

	switch index := indexArr.Elems[0].(type) {
	case *object.PanInt:
		// TODO: duck typing for keys (allow child of int)
		return arrIndex(index.Value, self)
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
