package props

import (
	"github.com/Syuparn/pangaea/object"
)

// EitherProps provides built-in props for Either.
// NOTE: Some Val props are defind by native code (not by this function).
func EitherProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"A":    object.BuiltInNotImplemented,
		"err":  object.BuiltInNotImplemented,
		"fmap": object.BuiltInNotImplemented,
		"or":   object.BuiltInNotImplemented,
		"val":  object.BuiltInNotImplemented,
	}
}
