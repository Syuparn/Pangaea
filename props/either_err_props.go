package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// EitherErrProps provides built-in props for EitherErr.
// NOTE: Some Val props are defind by native code (not by this function).
func EitherErrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("EitherErr"),
		"A": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("EitherErr#A requires at least 1 arg")
				}

				errObj, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherErr", object.ReprStr(args[0])))
				}

				err, ok := (*errObj.Pairs)[object.GetSymHash("_error")]
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherErr", object.ReprStr(args[0])))
				}

				return &object.PanArr{Elems: []object.PanObject{
					object.BuiltInNil,
					err.Value,
				}}
			},
		),
		"err": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("EitherErr#err requires at least 1 arg")
				}

				errObj, ok := object.TraceProtoOfObj(args[0])
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherErr", object.ReprStr(args[0])))
				}

				err, ok := (*errObj.Pairs)[object.GetSymHash("_error")]
				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("`%s` cannot be treated as EitherErr", object.ReprStr(args[0])))
				}

				return err.Value
			},
		),
		"fmap": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("EitherErr#fmap requires at least 2 args")
				}

				return args[0]
			},
		),
		"or": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("EitherErr#or requires at least 2 args")
				}

				return args[1]
			},
		),
		"val": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("EitherErr#val requires at least 1 arg")
				}

				return object.BuiltInNil
			},
		),
	}
}
