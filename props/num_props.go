package props

import (
	"fmt"
	"math"

	"github.com/Syuparn/pangaea/object"
)

// NumProps provides built-in props for Num.
// NOTE: Some Num props are defind by native code (not by this function).
func NumProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("Num"),
		"ceil": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Num#ceil requires at least 1 arg")
				}

				if i, ok := object.TraceProtoOfInt(args[0]); ok {
					return i
				}

				if f, ok := object.TraceProtoOfFloat(args[0]); ok {
					return object.NewPanInt(int64(math.Ceil(f.Value)))
				}

				return object.NewTypeErr(fmt.Sprintf("%s cannot be treated as num",
					args[0].Repr()))
			},
		),
		"F": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Num#F requires at least 1 arg")
				}

				if f, ok := object.TraceProtoOfFloat(args[0]); ok {
					return f
				}

				if i, ok := object.TraceProtoOfInt(args[0]); ok {
					return object.NewPanFloat(float64(i.Value))
				}

				return object.NewTypeErr(fmt.Sprintf("%s cannot be treated as num",
					args[0].Repr()))
			},
		),
		"floor": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Num#floor requires at least 1 arg")
				}

				if i, ok := object.TraceProtoOfInt(args[0]); ok {
					return i
				}

				if f, ok := object.TraceProtoOfFloat(args[0]); ok {
					return object.NewPanInt(int64(math.Floor(f.Value)))
				}

				return object.NewTypeErr(fmt.Sprintf("%s cannot be treated as num",
					args[0].Repr()))
			},
		),
		"round": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Num#round requires at least 1 arg")
				}

				if i, ok := object.TraceProtoOfInt(args[0]); ok {
					return i
				}

				if f, ok := object.TraceProtoOfFloat(args[0]); ok {
					return object.NewPanInt(int64(math.Round(f.Value)))
				}

				return object.NewTypeErr(fmt.Sprintf("%s cannot be treated as num",
					args[0].Repr()))
			},
		),
	}
}
