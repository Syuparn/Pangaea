package props

import (
	"../object"
)

func MapProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				// NOTE: BaseObj#== comparison is NOTHING TO DO WITH proto hierarchy!
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Map itself! (guarantee `Map == Map`)
				if args[0] == object.BuiltInMapObj && args[1] == object.BuiltInMapObj {
					return object.BuiltInTrue
				}

				self, ok := traceProtoOf(args[0], isMap)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isMap)
				if !ok {
					return object.BuiltInFalse
				}

				return compMaps(self.(*object.PanMap), other.(*object.PanMap),
					propContainer, env)
			},
		),
		"at": propContainer["Map_at"],
	}
}

func compMaps(
	m1 *object.PanMap,
	m2 *object.PanMap,
	propContainer map[string]object.PanObject,
	env *object.Env,
) object.PanObject {
	if len(*m1.Pairs) != len(*m2.Pairs) {
		return object.BuiltInFalse
	}

	for hash, pair1 := range *m1.Pairs {
		pair2, ok := (*m2.Pairs)[hash]
		if !ok {
			return object.BuiltInFalse
		}

		// == comparison for both elements
		res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
			env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), pair1.Value, eqSym, pair2.Value,
		)
		if res == object.BuiltInFalse {
			return object.BuiltInFalse
		}
	}

	// FIXME: this algorithm takes O(len(NonHashablePairs)**2)...
	if len(*m1.NonHashablePairs) != len(*m2.NonHashablePairs) {
		return object.BuiltInFalse
	}

	for _, pair1 := range *m1.NonHashablePairs {
		// NOTE: no warry about situation that 2 different keys in pairs1
		//       match same key in pairs2
		// 		 because keys are unique about == comparison

		idx, ok := containsKey(pair1.Key, *m2.NonHashablePairs,
			propContainer, env)
		if !ok {
			return object.BuiltInFalse
		}

		pair2 := (*m2.NonHashablePairs)[idx]
		res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
			env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), pair1.Value, eqSym, pair2.Value,
		)
		if res == object.BuiltInFalse {
			return object.BuiltInFalse
		}
	}

	return object.BuiltInTrue
}

func containsKey(
	k object.PanObject,
	pairs []object.Pair,
	propContainer map[string]object.PanObject,
	env *object.Env,
) (int, bool) {
	for i, pair := range pairs {
		// == comparison
		res := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
			env, object.EmptyPanObjPtr(),
			object.EmptyPanObjPtr(), k, eqSym, pair.Key,
		)

		if res == object.BuiltInTrue {
			return i, true
		}
	}
	return -1, false
}
