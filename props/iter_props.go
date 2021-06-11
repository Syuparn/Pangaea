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

				if args[0].Type() == object.FuncType && args[1].Type() == object.FuncType {
					return compFuncs(args[0].(*object.PanFunc), args[1].(*object.PanFunc))
				}

				if args[0].Type() == object.BuiltInIterType && args[1].Type() == object.BuiltInIterType {
					if args[0] == args[1] {
						return object.BuiltInTrue
					}
					return object.BuiltInFalse
				}

				return object.BuiltInFalse
			},
		),
		"_iter": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Iter#_iter requires at least 1 arg")
				}

				// return copied new iter not to share state(=env)
				if it, ok := object.TraceProtoOfFunc(args[0]); ok {
					return copiedIterFromIter(it)
				}

				if builtInIt, ok := object.TraceProtoOfBuiltInIter(args[0]); ok {
					return copiedIterFromBuiltInIter(builtInIt)
				}

				return object.NewTypeErr(
					fmt.Sprintf("%s cannot be treated as iter", args[0].Repr()))
			},
		),
		"_name": object.NewPanStr("Iter"),
		"B": f(
			func(
				env *object.Env, kwargs *object.PanObj, args ...object.PanObject,
			) object.PanObject {
				if len(args) < 1 {
					return object.NewTypeErr("Iter#B requires at least 1 arg")
				}
				return object.BuiltInTrue
			},
		),
		"new":  propContainer["Iter_new"],
		"next": propContainer["Iter_next"],
	}
}

func copiedIterFromIter(self *object.PanFunc) object.PanObject {
	if self.FuncKind != object.IterFunc {
		return object.NewTypeErr(
			fmt.Sprintf("%s cannot be treated as iter", self.Repr()))
	}

	return object.NewPanIter(self.FuncWrapper, object.NewCopiedEnv(self.Env))
}

func copiedIterFromBuiltInIter(self *object.PanBuiltInIter) object.PanObject {
	return object.NewPanBuiltInIter(
		self.Fn,
		object.NewCopiedEnv(self.Env),
	)
}
