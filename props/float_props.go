package props

import (
	"../object"
	"fmt"
)

func FloatProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				self, ok := traceProtoOf(args[0], isFloat)
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := traceProtoOf(args[1], isFloat)
				if !ok {
					return object.BuiltInFalse
				}

				if self.(*object.PanFloat).Value == other.(*object.PanFloat).Value {
					return object.BuiltInTrue
				}
				return object.BuiltInFalse
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("+ requires at least 2 args")
				}

				self, ok := traceProtoOf(args[0], isFloat)
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as int", self.Inspect()))
				}
				other, ok := traceProtoOf(args[1], isFloat)
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as int", other.Inspect()))
				}

				res := self.(*object.PanFloat).Value + other.(*object.PanFloat).Value
				return &object.PanFloat{Value: res}
			},
		),
	}
}
