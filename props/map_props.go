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
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Map#_iter requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isMap)
				if !ok {
					return object.NewTypeErr(`\1 must be map`)
				}

				return &object.PanBuiltInIter{
					Fn:  mapIter(self.(*object.PanMap)),
					Env: env, // not used
				}
			},
		),
		"at": propContainer["Map_at"],
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Map#B requires at least 1 arg")
				}
				self, ok := traceProtoOf(args[0], isMap)
				if !ok {
					return object.NewTypeErr(`\1 must be map`)
				}

				m, _ := self.(*object.PanMap)
				if len(*m.Pairs) == 0 && len(*m.NonHashablePairs) == 0 {
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

				self, ok := traceProtoOf(args[0], isMap)
				if !ok {
					return &object.PanArr{Elems: []object.PanObject{}}
				}
				map_, _ := self.(*object.PanMap)

				items := []object.PanObject{}
				for _, pair := range *map_.Pairs {
					items = append(items, &object.PanArr{Elems: []object.PanObject{
						pair.Key,
						pair.Value,
					}})
				}
				for _, pair := range *map_.NonHashablePairs {
					items = append(items, &object.PanArr{Elems: []object.PanObject{
						pair.Key,
						pair.Value,
					}})
				}

				return &object.PanArr{Elems: items}
			},
		),
		"keys": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				// NOTE: order is not guaranteed!

				if len(args) < 1 {
					return object.NewTypeErr("Map#keys requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isMap)
				if !ok {
					return &object.PanArr{Elems: []object.PanObject{}}
				}
				map_, _ := self.(*object.PanMap)

				keys := []object.PanObject{}
				for _, pair := range *map_.Pairs {
					keys = append(keys, pair.Key)
				}
				for _, pair := range *map_.NonHashablePairs {
					keys = append(keys, pair.Key)
				}

				return &object.PanArr{Elems: keys}
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

				self, ok := traceProtoOf(args[0], isMap)
				if !ok {
					return &object.PanArr{Elems: []object.PanObject{}}
				}
				map_, _ := self.(*object.PanMap)

				values := []object.PanObject{}
				for _, pair := range *map_.Pairs {
					values = append(values, pair.Value)
				}
				for _, pair := range *map_.NonHashablePairs {
					values = append(values, pair.Value)
				}

				return &object.PanArr{Elems: values}
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
	for hash, _ := range *m.Pairs {
		hashes = append(hashes, hash)
	}

	nonScalarYieldIdx := 0

	return func(
		env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
	) object.PanObject {
		if scalarYieldIdx >= len(hashes) {
			yielded := yieldNonScalar(m, nonScalarYieldIdx)
			nonScalarYieldIdx += 1
			return yielded
		}
		pair, ok := (*m.Pairs)[hashes[scalarYieldIdx]]
		// must be ok
		if !ok {
			return object.NewValueErr("pair data in map somehow got changed")
		}

		yielded := &object.PanArr{Elems: []object.PanObject{
			pair.Key,
			pair.Value,
		}}
		scalarYieldIdx += 1
		return yielded
	}
}

func yieldNonScalar(m *object.PanMap, i int) object.PanObject {
	if i >= len(*m.NonHashablePairs) {
		return object.NewStopIterErr("iter stopped")
	}

	pair := (*m.NonHashablePairs)[i]
	return &object.PanArr{Elems: []object.PanObject{
		pair.Key,
		pair.Value,
	}}
}
