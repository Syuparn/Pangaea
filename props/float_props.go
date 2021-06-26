package props

import (
	"fmt"
	"math"

	"github.com/Syuparn/pangaea/object"
)

// FloatProps provides built-in props for Float.
// NOTE: Some Float props are defind by native code (not by this function).
func FloatProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"<=>": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkFloatInfixArgs(args, "<=>", object.NewPanFloat(0.0))
				if err != nil {
					return err
				}

				selfVal := self.Value
				otherVal := other.Value
				var res int64

				if selfVal > otherVal {
					res = 1
				} else if selfVal == otherVal {
					res = 0
				} else {
					res = -1
				}

				return object.NewPanInt(res)
			},
		),
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
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
				return object.NewPanFloat(res)
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
				return object.NewPanFloat(res)
			},
		),
		"+": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkFloatInfixArgs(args, "+", object.NewPanFloat(0.0))
				if err != nil {
					return err
				}

				res := self.Value + other.Value
				return object.NewPanFloat(res)
			},
		),
		"-": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				self, other, err := checkFloatInfixArgs(args, "-", object.NewPanFloat(0.0))
				if err != nil {
					return err
				}

				res := self.Value - other.Value
				return object.NewPanFloat(res)
			},
		),
		"_name": object.NewPanStr("Float"),
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
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Float#new requires at least 2 args")
				}
				f, ok := object.TraceProtoOfFloat(args[1])
				if ok {
					return f
				}

				i, ok := object.TraceProtoOfInt(args[1])
				if ok {
					return object.NewPanFloat(float64(i.Value))
				}

				return object.NewTypeErr(
					fmt.Sprintf("%s cannot be treated as float", args[1].Repr()))
			},
		),
	}
}

func checkFloatInfixArgs(
	args []object.PanObject,
	propName string,
	nilAs *object.PanFloat,
) (*object.PanFloat, *object.PanFloat, *object.PanErr) {
	if len(args) < 2 {
		return nil, nil, object.NewTypeErr(propName + " requires at least 2 args")
	}

	self, ok := toPanFloat(args[0])
	if !ok {
		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("%s cannot be treated as float", args[0].Repr()))
	}

	other, ok := toPanFloat(args[1])
	if !ok {
		// NOTE: nil is treated as nilAs (0 in `+` and 1 in `*` for example)
		_, ok := object.TraceProtoOfNil(args[1])
		if ok {
			return self, nilAs, nil
		}

		return nil, nil, object.NewTypeErr(
			fmt.Sprintf("%s cannot be treated as float", args[1].Repr()))
	}

	return self, other, nil
}

func toPanFloat(obj object.PanObject) (*object.PanFloat, bool) {
	if f, ok := object.TraceProtoOfFloat(obj); ok {
		return f, true
	}

	// cast int to float
	if i, ok := object.TraceProtoOfInt(obj); ok {
		return object.NewPanFloat(float64(i.Value)), true
	}

	return nil, false
}
