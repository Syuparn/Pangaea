package props

import (
	"../object"
	"fmt"
	"math"
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

				// necessary for Float itself! (guarantee `Float == Float`)
				if args[0] == object.BuiltInFloatObj && args[1] == object.BuiltInFloatObj {
					return object.BuiltInTrue
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
		"-%": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("\\- requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isFloat)
				if !ok {
					return object.NewTypeErr("\\1 must be float")
				}

				res := -self.(*object.PanFloat).Value
				return &object.PanFloat{Value: res}
			},
		),
		"/~": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("/~ requires at least 1 arg")
				}

				self, ok := traceProtoOf(args[0], isFloat)
				if !ok {
					return object.NewTypeErr("\\1 must be float")
				}

				v := self.(*object.PanFloat).Value
				// NOTE: go cannot invert float bits directly
				res := math.Float64frombits(^math.Float64bits(v))
				return &object.PanFloat{Value: res}
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
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Float#B requires at least 1 arg")
				}
				self, ok := traceProtoOf(args[0], isFloat)
				if !ok {
					return object.NewTypeErr(`\1 must be float`)
				}

				if self.(*object.PanFloat).Value == 0.0 {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
	}
}
