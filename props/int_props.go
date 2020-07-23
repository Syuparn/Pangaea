package props

import (
	"../object"
)

func IntProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				isInt := func(o object.PanObject) bool {
					return o.Type() == object.INT_TYPE
				}

				self, ok := traceProtoOf(args[0], isInt)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isInt)
				if !ok {
					return object.BuiltInFalse
				}

				if self.(*object.PanInt).Value == other.(*object.PanInt).Value {
					return object.BuiltInTrue
				}
				return object.BuiltInFalse
			},
		),
	}
}
