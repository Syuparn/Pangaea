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
