package props

import (
	"fmt"

	"github.com/Syuparn/pangaea/object"
)

// IterProps provides built-in props for Iter.
// NOTE: Some Iter props are defind by native code (not by this function).
func IterProps(propContainer map[string]object.PanObject) map[string]object.PanObject {
	// NOTE: inject some built-in functions which relate to parser or evaluator
	return map[string]object.PanObject{
		"==": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 2 {
					return object.NewTypeErr("== requires at least 2 args")
				}

				// necessary for Iter itself! (guarantee `Iter == Iter`)
				if args[0] == object.BuiltInIterObj && args[1] == object.BuiltInIterObj {
					return object.BuiltInTrue
				}

				// same as func comparison
				self, ok := object.TraceProtoOfFunc(args[0])
				if !ok {
					return object.BuiltInFalse
				}
				other, ok := object.TraceProtoOfFunc(args[1])
				if !ok {
					return object.BuiltInFalse
				}

				return compFuncs(self, other)
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Iter#_iter requires at least 1 arg")
				}

				self, ok := object.TraceProtoOfFunc(args[0])
				if !ok || self.FuncKind != object.IterFunc {
					return object.NewTypeErr(
						fmt.Sprintf("%s cannot be treated as iter", args[0].Inspect()))
				}

				// return copied new iter not to share state(=env)
				newEnv := object.NewCopiedEnv(self.Env)

				return &object.PanFunc{
					FuncWrapper: self.FuncWrapper,
					FuncKind:    object.IterFunc,
					Env:         newEnv,
				}
			},
		),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Iter#B requires at least 1 arg")
				}
				_, ok := object.TraceProtoOfFunc(args[0])
				if !ok {
					return object.NewTypeErr(`\1 must be func`)
				}
				return object.BuiltInTrue
			},
		),
		"new":  propContainer["Iter_new"],
		"next": propContainer["Iter_next"],
	}
}
