package props

import (
	"github.com/Syuparn/pangaea/object"
)

// ErrProps provides built-in props for Err.
// NOTE: internally, these props are used for ErrWrappers
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
		"A": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Err#A requires at least 1 arg")
				}

				return &object.PanArr{Elems: []object.PanObject{
					object.BuiltInNil,
					args[0],
				}}
			},
		),
		"fmap": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("Err#fmap requires at least 2 args")
				}

				return args[0]
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
