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

				if _, ok := traceProtoOf(args[0], isNil); !ok {
					return object.BuiltInFalse
				}
				if _, ok := traceProtoOf(args[1], isNil); !ok {
					return object.BuiltInFalse
				}

				return object.BuiltInTrue
			},
		),
	}
}
