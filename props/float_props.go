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

				self, ok := object.TraceProtoOfFloat(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfFloat(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				if self.Value == other.Value {
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

				self, ok := object.TraceProtoOfFloat(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be float")
				}

				res := -self.Value
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

				self, ok := object.TraceProtoOfFloat(args[0])
				if !ok {
					return object.NewTypeErr("\\1 must be float")
				}

				v := self.Value
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

				self, ok := object.TraceProtoOfFloat(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as float", self.Inspect()))
				}
				other, ok := object.TraceProtoOfFloat(args[1])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as float", other.Inspect()))
				}

				res := self.Value + other.Value
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
				self, ok := object.TraceProtoOfFloat(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be float`)
				}

				if self.Value == 0.0 {
					return object.BuiltInFalse
				}
				return object.BuiltInTrue
			},
		),
	}
}
