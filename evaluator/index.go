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
