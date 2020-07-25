package props

import (
	"../object"
)

func IterProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Iter itself! (guarantee `Iter == Iter`)
				if args[0] == object.BuiltInIterObj && args[1] == object.BuiltInIterObj {
					return object.BuiltInTrue
				}

				// same as func comparison
				self, ok := traceProtoOf(args[0], isFunc)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isFunc)
				if !ok {
					return object.BuiltInFalse
				}

				return compFuncs(self.(*object.PanFunc), other.(*object.PanFunc))
			},
		),
		"new": propContainer["Iter_new"],
	}
}
