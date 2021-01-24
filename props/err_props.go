package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// ErrProps provides built-in props for Err.
// NOTE: internally, these props are also used for ErrWrappers
// NOTE: Some Val props are defind by native code (not by this function).
func ErrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Arr itself! (guarantee `Err == Err`)
				if args[0] == object.BuiltInErrObj && args[1] == object.BuiltInErrObj {
					return object.BuiltInTrue
				}

				self, ok := object.TraceProtoOfErrWrapper(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfErrWrapper(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				return compErrWrappers(self, other, propContainer, env)
			},
		),
		"msg": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Err#msg requires at least 1 arg")
				}

				err, ok := object.TraceProtoOfErrWrapper(args[0])

				if !ok {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as err", args[0].Inspect()))
				}

				return object.NewPanStr(err.PanErr.Msg)
			},
		),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				return constructErr(propContainer, env, object.NewPanErr, args...)
			},
		),
		"type": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Err#type requires at least 1 arg")
				}

				errObj, ok := object.TraceProtoOfObj(args[0])

				if !ok {
					// return default error type (never occurred)
					return object.BuiltInErrObj
				}

				return errObj
			},
		),
	}
}

func compErrWrappers(
	e1 *object.PanErrWrapper,
	e2 *object.PanErrWrapper,
	propContainer map[string]object.PanObject,
	env *object.Env,
) object.PanObject {
	if e1.Kind() != e2.Kind() {
		return object.BuiltInFalse
	}

	if e1.Message() != e2.Message() {
		return object.BuiltInFalse
	}

	return object.BuiltInTrue
}

func constructErr(
	propContainer map[string]object.PanObject,
	env *object.Env,
	newErr func(string) *object.PanErr,
	args ...object.PanObject,
) object.PanObject {
	if len(args) < 2 {
		// NOTE: not raising typeErr because constructor should not raise other type err
		return newErr("nil")
	}

	str, ok := object.TraceProtoOfStr(args[1])
	if ok {
		return newErr(str.Value)
	}

	// get args[1].S instead
	s := propContainer["Obj_callProp"].(*object.PanBuiltIn).Fn(
		env, object.EmptyPanObjPtr(),
		object.EmptyPanObjPtr(), args[1], sSym,
	)
	convertedStr, ok := object.TraceProtoOfStr(s)
	if ok {
		return newErr(convertedStr.Value)
	}

	return newErr(args[1].Inspect())
}
