package props

import (
	"../object"
)

// NilProps provides built-in props for Nil.
// NOTE: Some Nil props are defind by native code (not by this function).
func NilProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Nil itself! (guarantee `Nil == Nil`)
				if args[0] == object.BuiltInNilObj && args[1] == object.BuiltInNilObj {
					return object.BuiltInTrue
				}

				if _, ok := object.TraceProtoOfNil(args[0]); !ok {
					return object.BuiltInFalse
				}
				if _, ok := object.TraceProtoOfNil(args[1]); !ok {
					return object.BuiltInFalse
				}

				return object.BuiltInTrue
			},
		),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Nil#B requires at least 1 arg")
				}
				_, ok := object.TraceProtoOfNil(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be nil`)
				}
				return object.BuiltInFalse
			},
		),
	}
}
