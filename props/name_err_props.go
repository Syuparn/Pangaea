package props

import (
	"github.com/Syuparn/pangaea/object"
)

// NameErrProps provides built-in props for NameErr.
// NOTE: internally, these props are also used for ErrWrappers
// NOTE: Some Val props are defind by native code (not by this function).
func NameErrProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"new": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				return constructErr(propContainer, env, object.NewNameErr, args...)
			},
		),
	}
}
