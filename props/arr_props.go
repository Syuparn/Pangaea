package props

import (
	"../object"
)

func ArrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Arr itself! (guarantee `Arr == Arr`)
				if args[0] == object.BuiltInArrObj && args[1] == object.BuiltInArrObj {
					return object.BuiltInTrue
				}

				self, ok := traceProtoOf(args[0], isArr)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isArr)
				if !ok {
					return object.BuiltInFalse
				}

				return compArrs(self.(*object.PanArr), other.(*object.PanArr),
					propContainer, env)
			},
		),
		"at": propContainer["Arr_at"],
	}
}

func compArrs(
	a1 *object.PanArr,
	a2 *object.PanArr,
	propContainer map[string]object.PanObject,
	env *object.Env,
) object.PanObject {
	if len(a1.Elems) != len(a2.Elems) {
		return object.BuiltInFalse
	}

	for i, e := range a1.Elems {
		// == comparison for both elements
		res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
			env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), e, eqSym, a2.Elems[i],
		)
		if res == object.BuiltInFalse {
			return object.BuiltInFalse
		}
	}
	return object.BuiltInTrue
}
