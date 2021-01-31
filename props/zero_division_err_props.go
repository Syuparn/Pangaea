package props

import (
	"github.com/Syuparn/pangaea/object"
)

// ZeroDivisionErrProps provides built-in props for ZeroDivisionErr.
// NOTE: internally, these props are also used for ErrWrappers
// NOTE: Some Val props are defind by native code (not by this function).
func ZeroDivisionErrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"_name": object.NewPanStr("ZeroDivisionErr"),
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				return constructErr(propContainer, env, object.NewZeroDivisionErr, args...)
			},
		),
	}
}
