package props

import (
	"github.com/Syuparn/pangaea/object"
)

// BaseObjProps provides built-in props for BaseObj.
// NOTE: Some BaseObj props are defind by native code (not by this function).
func BaseObjProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
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

				self, ok := args[0].(*object.PanObj)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := args[1].(*object.PanObj)
				if !ok {
					return object.BuiltInFalse
				}

				return compObjs(self, other, propContainer, env)
			},
		),
		"_name": object.NewPanStr("BaseObj"),
		"at":    propContainer["BaseObj_at"],
		"bear": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("BaseObj#bear requires at least 1 arg")
				}
				proto := args[0]

				if len(args) < 2 {
					// default src
					src := object.EmptyPanObjPtr()
					return object.ChildPanObjPtr(proto, src)
				}

				src, ok := args[1].(*object.PanObj)
				if !ok {
					return object.NewTypeErr("BaseObj#bear requires obj literal src")
				}
				return object.ChildPanObjPtr(proto, src)
			},
		),
		"proto": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("proto requires at least 1 arg")
				}

				proto := args[0].Proto()
				if proto == nil {
					return object.BuiltInNil
				}
				return proto
			},
		),
	}
}

func compObjs(
	o1 *object.PanObj,
	o2 *object.PanObj,
	propContainer map[string]object.PanObject,
	env *object.Env,
) object.PanObject {
	if len(*o1.Pairs) != len(*o2.Pairs) {
		return object.BuiltInFalse
	}

	for sym, pair1 := range *o1.Pairs {
		pair2, ok := (*o2.Pairs)[sym]
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
	return object.BuiltInTrue
}
