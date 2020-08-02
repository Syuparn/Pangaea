package props

import (
	"../object"
)

func NilProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Nil itself! (guarantee `Nil == Nil`)
				if args[0] == object.BuiltInNilObj && args[1] == object.BuiltInNilObj {
					return object.BuiltInTrue
				}

				if _, ok := traceProtoOf(args[0], isNil); !ok {
					return object.BuiltInFalse
				}
				if _, ok := traceProtoOf(args[1], isNil); !ok {
					return object.BuiltInFalse
				}

				return object.BuiltInTrue
			},
		),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Nil#B requires at least 1 arg")
				}
				_, ok := traceProtoOf(args[0], isNil)
				if !ok {
					return object.NewTypeErr(`\1 must be nil`)
				}
				return object.BuiltInFalse
			},
		),
	}
}
