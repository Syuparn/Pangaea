package props

import (
	"github.com/Syuparn/pangaea/object"
)

// MapProps provides built-in props for Map.
// NOTE: Some Map props are defind by native code (not by this function).
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

				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfMap(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				return compMaps(self, other, propContainer, env)
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Map#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be map`)
				}

				return object.NewPanBuiltInIter(mapIter(self), env)
			},
		),
		"_name": object.NewPanStr("Map"),
		"at":    propContainer["Map_at"],
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Map#B requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be map`)
				}

				if len(*self.Pairs) == 0 && len(*self.NonHashablePairs) == 0 {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
		"items": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				// NOTE: order is not guaranteed!

				if len(args) < 1 {
					return object.NewTypeErr("Map#items requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.NewPanArr()
				}

				items := []object.PanObject{}
				for _, h := range *self.HashKeys {
					// NOTE: pair must be found
					pair, _ := (*self.Pairs)[h]
					items = append(items, object.NewPanArr(pair.Key, pair.Value))
				}
				for _, pair := range *self.NonHashablePairs {
					items = append(items, object.NewPanArr(pair.Key, pair.Value))
				}

				return object.NewPanArr(items...)
			},
		),
		"keys": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Map#keys requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.NewPanArr()
				}

				keys := []object.PanObject{}
				for _, h := range *self.HashKeys {
					// NOTE: pair must be found
					keys = append(keys, (*self.Pairs)[h].Key)
				}
				for _, pair := range *self.NonHashablePairs {
					keys = append(keys, pair.Key)
				}

				return object.NewPanArr(keys...)
			},
		),
		"len": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Map#len requires at least 1 arg")
				}
				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be map`)
				}

				length := len(*self.Pairs) + len(*self.NonHashablePairs)
				return object.NewPanInt(int64(length))
			},
		),
		"values": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				// NOTE: order is not guaranteed!

				if len(args) < 1 {
					return object.NewTypeErr("Map#values requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfMap(args[0])
				if !ok {
					return object.NewPanArr()
				}

				values := []object.PanObject{}
				for _, h := range *self.HashKeys {
					// NOTE: value must be found
					values = append(values, (*self.Pairs)[h].Value)
				}
				for _, pair := range *self.NonHashablePairs {
					values = append(values, pair.Value)
				}

				return object.NewPanArr(values...)
			},
		),
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

func mapIter(m *object.PanMap) object.BuiltInFunc {
	scalarYieldIdx := 0
	hashes := []object.HashKey{}
	for hash := range *m.Pairs {
		hashes = append(hashes, hash)
	}

	nonScalarYieldIdx := 0

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		if scalarYieldIdx >= len(hashes) {
			yielded := yieldNonScalar(m, nonScalarYieldIdx)
			nonScalarYieldIdx++
			return yielded
		}
		pair, ok := (*m.Pairs)[hashes[scalarYieldIdx]]
		// must be ok
		if !ok {
			return object.NewValueErr("pair data in map somehow got changed")
		}

		yielded := object.NewPanArr(pair.Key, pair.Value)

		scalarYieldIdx++
		return yielded
	}
}

func yieldNonScalar(m *object.PanMap, i int) object.PanObject {
	if i >= len(*m.NonHashablePairs) {
		return object.NewStopIterErr("iter stopped")
	}

	pair := (*m.NonHashablePairs)[i]
	return object.NewPanArr(pair.Key, pair.Value)
}
