package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// EitherValProps provides built-in props for EitherVal.
// NOTE: Some Val props are defind by native code (not by this function).
func EitherValProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"A": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("EitherVal#A requires at least 1 arg")
				}

				valObj, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherVal", args[0].Inspect()))
				}

				val, ok := (*valObj.Pairs)[object.GetSymHash("_value")]
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherVal", args[0].Inspect()))
				}

				return &object.PanArr{Elems: []object.PanObject{
					val.Value,
					object.BuiltInNil,
				}}
			},
		),
		"err": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("EitherVal#err requires at least 1 arg")
				}

				return object.BuiltInNil
			},
		),
		"fmap": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("EitherVal#fmap requires at least 2 args")
				}

				valObj, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherVal", args[0].Inspect()))
				}

				val, ok := (*valObj.Pairs)[object.GetSymHash("_value")]
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherVal", args[0].Inspect()))
				}

				f := args[1]

				// eval `f(val)` (== `f.call(val)`)
				callSym := object.NewPanStr("call")
				result := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
					env, object.EmptyPanObjPtr(),
					object.EmptyPanObjPtr(), f, callSym, val.Value,
				)

				if err, ok := result.(*object.PanErr); ok {
					return toEitherErr(object.WrapErr(err))
				}

				return toEitherVal(result)
			},
		),
		"val": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("EitherVal#val requires at least 1 arg")
				}

				valObj, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherVal", args[0].Inspect()))
				}

				val, ok := (*valObj.Pairs)[object.GetSymHash("_value")]
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherVal", args[0].Inspect()))
				}

				return val.Value
			},
		),
	}
}

func toEitherErr(e *object.PanErrWrapper) object.PanObject {
	// wrap e by EitherErr
	pairMap := map[object.SymHash]object.Pair{
		object.GetSymHash("_error"): {Key: object.NewPanStr("_error"), Value: e},
	}

	return object.NewPanObj(&pairMap, object.BuiltInEitherErrObj)
}
